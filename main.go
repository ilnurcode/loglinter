// Package main предоставляет CLI для loglinter.
//
// Использование:
//
//	# Запуск линтера
//	loglinter ./...
//
//	# Запуск с конфигурацией
//	loglinter -config .loglinter.yaml ./...
//
//	# Интеграция с golangci-lint
//	# Добавьте в .golangci.yml:
//	# linters:
//	#   enable:
//	#     - loglinter
//	# linters-settings:
//	#   custom:
//	#     loglinter:
//	#       type: module
//	#       path: github.com/yourusername/loglinter/cmd/main.go@v1.0.0
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilnurcode/loglinter/pkg/loglinter"

	"golang.org/x/tools/go/analysis/singlechecker"
	"gopkg.in/yaml.v3"
)

// Version версия линтера (для semver)
const Version = "v1.0.0"

// configPath флаг для пути к конфигурационному файлу
var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "путь к конфигурационному файлу (yaml)")
	// Флаг 'c' уже определён в singlechecker, не переопределяем
}

func main() {
	flag.Parse()

	// Загрузка конфигурации если указан файл
	if configPath != "" {
		cfg, err := loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка загрузки конфигурации: %v\n", err)
			os.Exit(1)
		}
		loglinter.SetConfig(cfg)
	}

	// Запуск анализатора
	singlechecker.Main(loglinter.Analyzer)
}

// loadConfig загружает конфигурацию из YAML файла
func loadConfig(path string) (*loglinter.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
	}

	cfg := loglinter.DefaultConfig()

	if len(data) > 0 {
		// Парсим YAML во временную структуру
		var rawCfg map[string]interface{}
		if err := yaml.Unmarshal(data, &rawCfg); err != nil {
			return nil, fmt.Errorf("не удалось разобрать YAML: %w", err)
		}

		// Обработка sensitive_patterns
		if patterns, ok := rawCfg["sensitive_patterns"].([]interface{}); ok {
			cfg.SensitivePatterns = make([]string, len(patterns))
			for i, p := range patterns {
				if s, ok := p.(string); ok {
					cfg.SensitivePatterns[i] = s
				}
			}
		}

		// Обработка allowed_loggers
		if loggers, ok := rawCfg["allowed_loggers"].(map[string]interface{}); ok {
			cfg.AllowedLoggers = make(map[string][]string)
			for k, v := range loggers {
				if methods, ok := v.([]interface{}); ok {
					cfg.AllowedLoggers[k] = make([]string, len(methods))
					for i, m := range methods {
						if s, ok := m.(string); ok {
							cfg.AllowedLoggers[k][i] = s
						}
					}
				}
			}
		}
	}

	return cfg, nil
}
