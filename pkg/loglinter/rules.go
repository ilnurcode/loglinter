// Package loglinter предоставляет правила проверки лог-сообщений.
package loglinter

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var loggers = map[string][]string{
	"log":    {"Print", "Println", "Printf", "Info", "Error", "Warn", "Debug"},
	"slog":   {"Info", "Error", "Warn", "Debug", "Log"},
	"zap":    {"Info", "Error", "Warn", "Debug", "Panic", "Fatal"},
	"logger": {"Info", "Error", "Warn", "Debug", "Print", "Println", "Printf"},
}

// Config представляет конфигурацию линтера.
type Config struct {
	SensitivePatterns []string            `yaml:"sensitive_patterns"`
	AllowedLoggers    map[string][]string `yaml:"allowed_loggers"`
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() *Config {
	return &Config{
		SensitivePatterns: []string{"password", "passwd", "secret", "token", "api_key", "apikey", "credential"},
		AllowedLoggers:    loggers,
	}
}

var currentConfig = DefaultConfig()

// SetConfig устанавливает пользовательскую конфигурацию.
func SetConfig(cfg *Config) {
	if cfg != nil {
		currentConfig = cfg
	}
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}
	pkgName := pkgIdent.Name
	funcName := sel.Sel.Name

	if !isSupportedLogger(pkgName, funcName) {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	var msgLit *ast.BasicLit
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			msgLit = lit
			break
		}
	}
	if msgLit == nil {
		return
	}

	message := strings.Trim(msgLit.Value, `"`)
	runRules(pass, msgLit, message)
}

func isSupportedLogger(pkg, method string) bool {
	if methods, ok := currentConfig.AllowedLoggers[pkg]; ok {
		for _, m := range methods {
			if m == method {
				return true
			}
		}
	}
	return false
}

func runRules(pass *analysis.Pass, lit *ast.BasicLit, msg string) {
	if msg == "" {
		return
	}

	// Правило 1: начинается со строчной буквы
	if !startsWithLowercase(msg) {
		fix := createCapitalizationFix(lit, msg)
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message should start with a lowercase letter",
			SuggestedFixes: []analysis.SuggestedFix{fix},
		})
	}

	// Правило 2: только английский язык
	if !isEnglishOnly(msg) {
		pass.Reportf(lit.Pos(), `log message should contain only English letters: %q`, msg)
	}

	// Правило 3: нет спецсимволов и эмодзи
	if hasSpecialCharsOrEmoji(msg) {
		fix := createSpecialCharsFix(lit, msg)
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message should not contain special characters or emojis",
			SuggestedFixes: []analysis.SuggestedFix{fix},
		})
	}

	// Правило 4: нет чувствительных данных
	if containsSensitiveData(msg) {
		pass.Reportf(lit.Pos(), `log message should not contain sensitive data keywords: %q`, msg)
	}
}

func startsWithLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
		if !unicode.IsSpace(r) {
			// Если первый не пробельный символ не буква, считаем OK
			return true
		}
	}
	return true
}

func isEnglishOnly(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && r > 127 {
			return false
		}
	}
	return true
}

func hasSpecialCharsOrEmoji(s string) bool {
	// Проверка на эмодзи и специальные символы (за пределами ASCII)
	for _, r := range s {
		if r > 127 {
			// Эмодзи находятся в ranges So (Symbol, Other) и Sm (Symbol, Math)
			if unicode.Is(unicode.So, r) || unicode.Is(unicode.Sm, r) {
				return true
			}
			// Любые не-ASCII символы (включая кириллицу) уже отлавливаются в isEnglishOnly
		}
	}

	// Проверка на последовательности спецсимволов
	if strings.HasSuffix(s, "...") || strings.HasSuffix(s, "!!!") {
		return true
	}

	return false
}

func containsSensitiveData(s string) bool {
	lowerMsg := strings.ToLower(s)

	// Проверка паттернов по умолчанию
	for _, key := range currentConfig.SensitivePatterns {
		if strings.Contains(lowerMsg, key) {
			return true
		}
	}

	return false
}

// createCapitalizationFix создает исправление для правила капитализации.
func createCapitalizationFix(lit *ast.BasicLit, msg string) analysis.SuggestedFix {
	if len(msg) == 0 {
		return analysis.SuggestedFix{}
	}

	// Находим первую букву и делаем её строчной
	runes := []rune(msg)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToLower(r)
			break
		}
	}
	fixedMsg := string(runes)

	return analysis.SuggestedFix{
		Message: "Capitalize first letter to lowercase",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos() + 1, // +1 чтобы не затрагивать открывающую кавычку
				End:     lit.End() - 1, // -1 чтобы не затрагивать закрывающую кавычку
				NewText: []byte(fixedMsg),
			},
		},
	}
}

// createSpecialCharsFix создает исправление для удаления спецсимволов.
func createSpecialCharsFix(lit *ast.BasicLit, msg string) analysis.SuggestedFix {
	fixedMsg := msg

	// Удаляем последовательности "..." и "!!!" без проверки (unconditional)
	fixedMsg = strings.TrimSuffix(fixedMsg, "...")
	fixedMsg = strings.TrimSuffix(fixedMsg, "!!!")

	// Удаляем эмодзи и специальные символы
	var cleaned strings.Builder
	for _, r := range fixedMsg {
		if r <= 127 || !unicode.Is(unicode.So, r) && !unicode.Is(unicode.Sm, r) {
			cleaned.WriteRune(r)
		}
	}
	fixedMsg = strings.TrimSpace(cleaned.String())

	return analysis.SuggestedFix{
		Message: "Remove special characters and emojis",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos() + 1,
				End:     lit.End() - 1,
				NewText: []byte(fixedMsg),
			},
		},
	}
}

// CompileSensitivePatterns компилирует пользовательские паттерны в regexp.
func CompileSensitivePatterns(patterns []string) ([]*regexp.Regexp, error) {
	regexps := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		regexps = append(regexps, re)
	}
	return regexps, nil
}
