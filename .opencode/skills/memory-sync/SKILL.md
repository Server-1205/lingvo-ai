---
name: memory-sync
description: Use after completing any significant feature, API change, architecture decision, or UI component creation — to update the memory/context files for Lingvo AI. Keeps project knowledge persistent across sessions. Run automatically after major implementations.
---

# Memory Sync

После завершения значимых изменений обновляй файлы в `.opencode/memory/`:

## Что и когда обновлять

| Файл | Когда обновлять |
|---|---|
| `memory/context.md` | Новый компонент, изменение архитектуры, смена фокуса |
| `memory/architecture-decisions.md` | Принято новое архитектурное решение |
| `memory/project-state.md` | Реализован компонент, начат новый спринт |

## Правила

1. Не удаляй старые записи из project-state — только меняй статус (❌ → ⏳ → ✅)
2. ADR добавляй сверху
3. context.md обновляй если изменился стек, структура или архитектура
