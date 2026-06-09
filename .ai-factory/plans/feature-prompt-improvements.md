# Улучшение промптов (Prompt Engineering)

**Branch:** `feature/prompt-improvements`
**Created:** 2026-06-07 15:30

## Settings

- **Testing:** yes — unit tests for prompt output parsing
- **Logging:** verbose
- **Docs:** no
- **Roadmap Linkage:** Оптимизация → Улучшение промптов (меньше галлюцинаций)

## Анализ

**Текущие проблемы в промптах:**

| Проблема | Промпт | Описание |
|----------|--------|----------|
| Галлюцинации | Chat, Premium, Grammar | AI придумывает несуществующие грамматические правила |
| False positives | Chat, Grammar | AI «исправляет» правильно написанные предложения |
| Однотипные упражнения | Daily Lesson | Только fill-in-the-blank |
| Высокая температура | gemini.go | 0.4 для всех задач — для грамматики нужно 0.2 |

**План исправлений:**
1. Добавить anti-hallucination guard во все промпты
2. Добавить «не исправлять то, что правильно»
3. Разнообразить типы упражнений в Daily Lesson
4. Снизить температуру для точных задач (grammar, corrections)
5. Улучшить фильтрацию мата в Vocab

## Задачи

### Task 1: BuildChatPrompt — anti-hallucination + false positive guard

**Файлы:** `internal/ai/prompts.go`

**Добавить в существующий промпт:**
```
- Only correct ACTUAL mistakes. If the sentence is correct, return "corrections": [].
- Do NOT invent grammar rules that don't exist in standard English.
- If you are unsure about a correction, do NOT include it.
- False positives frustrate the user. It's better to miss a correction than to make one up.
```

**Логирование:** `[prompt] chat prompt updated for level=..., lang=...`

### Task 2: BuildGrammarCheckPrompt — добавить reply + false positive guard

**Файлы:** `internal/ai/prompts.go`

**Изменить структуру ответа — добавить поле `reply`:**
```json
{
  "reply": "brief summary of what was checked",
  "corrections": [...]
}
```

**Добавить guards:**
```
- Only flag REAL errors. Do not correct informal but acceptable English.
- If the text has no errors, return "corrections": [].
```

### Task 3: BuildPremiumChatPrompt — усилить анализ + guards

**Файлы:** `internal/ai/prompts.go`

**Добавить те же guards что и в Chat, плюс:**
```
- premium_analysis must be specific to the user's message, not generic.
- areas_for_improvement should list 2-3 concrete points from the message.
```

### Task 4: BuildDailyLessonPrompt — разнообразие упражнений

**Файлы:** `internal/ai/prompts.go`

**Обновить инструкцию:**
```
- Vary exercise types within one lesson (mix fill-in-blank, multiple choice, matching).
- At least 2 different exercise types per lesson.
```

### Task 5: BuildVocabPrompt — усилить фильтрацию мата

**Файлы:** `internal/ai/prompts.go`

**Улучшить инструкцию:**
```
- Reject any offensive, obscene, or vulgar words with {"error": "inappropriate_word"}.
- Also reject words that are misspellings or variants of offensive words.
```

### Task 6: gemini.go — оптимизация температуры

**Файлы:** `internal/ai/gemini.go`

- Chat/Grammar: `0.4` → `0.2` (меньше креатива = точнее)
- Premium Chat: `0.4` → `0.3` (баланс точности и глубины)
- Daily Lesson: `0.4` — оставить
- Vocab/Quiz: `0.4` → `0.3`

### Task 7: Тесты для промптов

**Файлы:** `internal/ai/prompts_test.go`

- `TestBuildChatPrompt_IncludesAntiHallucination` — проверяет наличие guards
- `TestBuildGrammarCheckPrompt_HasReplyField` — проверяет reply в JSON-шаблоне
- `TestBuildDailyLessonPrompt_HasVarietyInstruction` — проверяет инструкцию о разнообразии
- `TestBuildVocabPrompt_RejectsProfanity` — проверяет фильтрацию мата

## Commit Plan

1. `refactor(prompts): add anti-hallucination guards to all prompts`
2. `refactor(prompts): add exercise variety instruction to daily lesson`
3. `refactor(ai): optimize temperature for accuracy-sensitive tasks`
4. `test(ai): add prompt structure tests`
