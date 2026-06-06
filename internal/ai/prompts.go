package ai

import "fmt"

func BuildChatPrompt(level, lang, text string) string {
	return fmt.Sprintf(`You are an AI English tutor. The user's English level is %s. Explain in %s.

Rules:
1. Reply naturally in English, keeping your response concise (2-4 sentences).
2. If the user makes mistakes, correct them gently.
3. Always return JSON only (no markdown, no code fences):
{
  "reply": "your reply in English",
  "corrections": [
    {
      "original": "incorrect phrase",
      "corrected": "corrected phrase",
      "explanation_uz": "if lang is uz, explain in Uzbek",
      "explanation_ru": "if lang is ru, explain in Russian",
      "type": "grammar|vocabulary|spelling"
    }
  ]
}
4. If no mistakes, return "corrections": [].
5. corrections array must be empty or contain 1-3 items max.

User message: %s`, level, lang, text)
}

func BuildGrammarCheckPrompt(level, lang, text string) string {
	return fmt.Sprintf(`Check this text for grammar errors. User level: %s. Explain in %s.

Return JSON:
{
  "corrections": [
    {
      "original": "wrong part",
      "corrected": "corrected part",
      "explanation_uz": "explanation in Uzbek",
      "explanation_ru": "explanation in Russian",
      "type": "grammar|vocabulary|spelling"
    }
  ]
}

Text: %s`, level, lang, text)
}

func BuildVocabPrompt(lang, word string) string {
	return fmt.Sprintf(`You are an English-Uzbek/Russian dictionary. Return JSON only:
{
  "translation_uz": "translation to Uzbek",
  "translation_ru": "translation to Russian",
  "examples": ["example sentence 1", "example sentence 2", "example sentence 3"],
  "level": "a1|a2|b1|b2|c1"
}

Word: %s. Language: %s`, word, lang)
}

func BuildQuizPrompt(topic string, count int, lang string) string {
	return fmt.Sprintf(`Generate %d English quiz questions about "%s". Explain in %s.

Return JSON:
{
  "questions": [
    {
      "question": "question text",
      "options": ["A", "B", "C", "D"],
      "correct": 0,
      "explanation_uz": "explanation in Uzbek",
      "explanation_ru": "explanation in Russian"
    }
  ]
}`, count, topic, lang)
}

func BuildLevelTestPrompt(lang string) string {
	return fmt.Sprintf(`Generate 10 English level test questions (A1 to C1). Explain in %s.

Return JSON:
{
  "questions": [
    {
      "question": "question text",
      "options": ["A", "B", "C", "D"],
      "correct": 0,
      "level": "a1|a2|b1|b2|c1",
      "explanation_uz": "explanation in Uzbek",
      "explanation_ru": "explanation in Russian"
    }
  ]
}`, lang)
}
