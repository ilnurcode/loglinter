// Package loglinter предоставляет анализатор для проверки лог-сообщений.
// Он проверяет логи на соответствие правилам:
//   - начало со строчной буквы
//   - только английский язык
//   - отсутствие спецсимволов и эмодзи
//   - отсутствие чувствительных данных
package loglinter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `loglinter проверяет лог-сообщения на:
	- начало со строчной буквы
	- содержание только английских букв
	- отсутствие спецсимволов и эмодзи
	- отсутствие ключевых слов чувствительных данных

Поддерживаемые логгеры:
	- log (Print, Println, Printf, Info, Error, Warn, Debug)
	- slog (Info, Error, Warn, Debug, Log)
	- zap (Info, Error, Warn, Debug, Panic, Fatal)
`

var Analyzer = &analysis.Analyzer{
	Name:             "loglinter",
	Doc:              Doc,
	Run:              run,
	RunDespiteErrors: false,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := inspector.New(pass.Files)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkLogCall(pass, call)
	})
	return nil, nil
}
