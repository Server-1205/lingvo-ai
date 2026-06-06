# Lingvo AI — Project State

> Отслеживает что сделано и что в работе.

## Компоненты

### Backend
| Компонент | Статус | Примечание |
|---|---|---|
| Go module init | ✅ Done | go.mod, go.sum |
| DB schema + migration | ✅ Done | 5 tables, indexes |
| Model structs | ✅ Done | User, Correction, ChatResponse, Usage |
| main.go wiring | ✅ Done | DB, router, bot |
| Bot handlers | ⏳ Stub | только лог |
| API endpoints | ⏳ Stub | все 501 |
| AI integration | ❌ Not started | |
| Auth middleware | ❌ Not started | |
| Rate limiter | ❌ Not started | |

### Frontend
| Компонент | Статус | Примечание |
|---|---|---|
| Vite + React setup | ✅ Done | шаблон |
| i18n translations | ✅ Done | uz.json + ru.json |
| Chat UI | ❌ Not started | |
| Vocabulary UI | ❌ Not started | |
| Progress UI | ❌ Not started | |
| Subscription UI | ❌ Not started | |
| API client | ❌ Not started | |

### Infrastructure
| Компонент | Статус | Примечание |
|---|---|---|
| Dockerfile | ❌ Not started | |
| Railway deploy | ❌ Not started | |
| CI/CD | ❌ Not started | |

## Текущий спринт
1. Реализовать API /api/chat с Gemini
2. Написать auth middleware
3. Написать rate limiter middleware
4. Реализовать bot handlers (start, help, payment)
