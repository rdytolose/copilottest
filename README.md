# Devin Router (Go)

Создает Devin-сессии для анализа PR и возвращает ссылку на полноценное окно Devin.

## Требования
- Go 1.22+
- 4 Devin API ключа (v1)

## Настройка
```bash
export DEVIN_API_KEYS="apk_1,apk_2,apk_3,apk_4"
export GITHUB_TOKEN="ghp_..." # если приватные репо или для лимитов
```

## Запуск
```bash
go run ./cmd/devin-router --repo owner/repo --listen :8080
```

## Использование
```bash
curl -X POST http://localhost:8080/api/analyze/pr \
  -H 'Content-Type: application/json' \
  -d '{"pr_number":123}'
```

Ответ:
```json
{"session_id":"devin-xxx","url":"https://app.devin.ai/..." }
```
Открывайте `url` в браузере — это официальный Devin UI.
