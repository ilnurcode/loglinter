# loglinter

Линтер для проверки лог-записей в Go-проектах. Совместим с golangci-lint.

## Описание

`loglinter` — это статический анализатор для Go, который проверяет лог-сообщения на соответствие установленным правилам кодирования. Инструмент разработан в рамках тестового задания на позицию Backend Developer (Golang).

## Возможности

### Проверяемые правила

1. **Лог-сообщения должны начинаться со строчной буквы**

   ```go
   // ❌ Неправильно
   log.Info("Starting server on port 8080")
   slog.Error("Failed to connect to database")

   // ✅ Правильно
   log.Info("starting server on port 8080")
   slog.Error("failed to connect to database")
   ```

2. **Лог-сообщения должны быть только на английском языке**

   ```go
   // ❌ Неправильно
   log.Info("запуск сервера")
   log.Error("ошибка подключения к базе данных")

   // ✅ Правильно
   log.Info("starting server")
   log.Error("failed to connect to database")
   ```

3. **Лог-сообщения не должны содержать спецсимволы или эмодзи**

   ```go
   // ❌ Неправильно
   log.Info("server started!🚀")
   log.Error("connection failed!!!")
   log.Warn("warning: something went wrong...")

   // ✅ Правильно
   log.Info("server started")
   log.Error("connection failed")
   log.Warn("something went wrong")
   ```

4. **Лог-сообщения не должны содержать потенциально чувствительные данные**

   ```go
   // ❌ Неправильно
   log.Info("user password: secret123")
   log.Debug("api_key=" + apiKey)
   log.Info("token: " + token)

   // ✅ Правильно
   log.Info("user authenticated successfully")
   log.Debug("api request completed")
   log.Info("token validated")
   ```

## Поддерживаемые логгеры

| Логгер | Методы |
|--------|--------|
| `log` | Print, Println, Printf, Info, Error, Warn, Debug |
| `slog` | Info, Error, Warn, Debug, Log |
| `zap` | Info, Error, Warn, Debug, Panic, Fatal |
| `logger` | Info, Error, Warn, Debug, Print, Println, Printf |

## Установка

### Вариант 1: Установка через go install

```bash
go install github.com/ilnurcode/loglinter@latest
```

После установки линтер доступен как команда:

```bash
loglinter ./...
```

### Вариант 2: Запуск без установки

```bash
go run github.com/ilnurcode/loglinter@latest ./...
```

### Вариант 3: Локальная сборка

```bash
git clone https://github.com/ilnurcode/loglinter.git
cd loglinter
go build -o loglinter ./cmd/main.go
./loglinter ./...
```

---

## Интеграция с golangci-lint

### Использование как Module Plugin

Добавьте конфигурацию в `.golangci.yml` вашего проекта:

```yaml
linters:
  enable:
    - loglinter

linters-settings:
  custom:
    loglinter:
      type: module
      path: github.com/ilnurcode/loglinter/cmd/main.go@v1.0.0
      description: "Проверка лог-сообщений"
      original-url: github.com/ilnurcode/loglinter
```

Запустите проверку:

```bash
golangci-lint run ./...
```

### Последовательный запуск

```bash
golangci-lint run ./... && loglinter ./...
```

## Конфигурация

Линтер поддерживает настройку через YAML-файл. Пример `.loglinter.yaml`:

```yaml
sensitive_patterns:
  - password
  - passwd
  - secret
  - token
  - api_key
  - apikey
  - credential
  - private_key
  - access_token

allowed_loggers:
  log:
    - Print
    - Println
    - Printf
    - Info
    - Error
    - Warn
    - Debug
  slog:
    - Info
    - Error
    - Warn
    - Debug
    - Log
  zap:
    - Info
    - Error
    - Warn
    - Debug
    - Panic
    - Fatal
```

## Технические детали

### Архитектура

Линтер реализован с использованием `golang.org/x/tools/go/analysis` и состоит из следующих компонентов:

- **analyzer.go** — ядро анализатора, регистрация и обход AST
- **rules.go** — реализация правил проверки и функции авто-исправления
- **loglinter_test.go** — unit-тесты для всех правил

### Авто-исправления

Для правил 1 (строчная буква) и 3 (спецсимволы) реализованы `SuggestedFixes`, которые позволяют автоматически исправить найденные проблемы при запуске с флагом `-fix`.

### Требования

- Go 1.22+
- golangci-lint v1.59+ (для интеграции как Module Plugin)

## Структура проекта

```
loglinter/
├── cmd/main.go                 # CLI утилита
├── pkg/loglinter/
│   ├── analyzer.go             # Анализатор
│   ├── rules.go                # Правила проверки
│   └── loglinter_test.go       # Unit-тесты
├── internal/analysisutil/
│   └── util.go                 # Вспомогательные функции
├── .github/workflows/ci.yml    # CI/CD конфигурация
├── .loglinter.yaml             # Пример конфигурации
├── go.mod
├── go.sum
└── test_example.go             # Пример для проверки
```

## Тестирование

Запуск unit-тестов:

```bash
go test -v ./...
```

Тесты покрывают все 4 правила проверки и включают 12 тестовых случаев.

## CI/CD

Проект использует GitHub Actions для автоматического тестирования и сборки при каждом push и pull request.

## Лицензия

MIT
