package ai

import "fmt"

func langFullName(code string) string {
	switch code {
	case "uz":
		return "Uzbek"
	case "ru":
		return "Russian"
	default:
		return "English"
	}
}

var langInstruction = map[string]string{
	"uz": "IMPORTANT: Always reply and explain everything in Uzbek. The user is an Uzbek speaker learning English.",
	"ru": "IMPORTANT: Always reply and explain everything in Russian. The user is a Russian speaker learning English.",
}

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

func BuildPremiumChatPromptWithHistory(level, lang, text string, recentErrors []string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Reply naturally in English, keeping your response concise (2-4 sentences)."
	}

	historySection := ""
	if len(recentErrors) > 0 {
		historySection = "\n\nThe user has been making these errors recently. Pay special attention to them:\n"
		for _, e := range recentErrors {
			historySection += "- " + e + "\n"
		}
		historySection += "\nProvide extra detailed explanations when these errors occur."
	}

	return fmt.Sprintf(`You are an AI English tutor. The user's English level is %s.

%s

You are in PREMIUM mode — provide deeper analysis.
%s

Rules:
1. Only correct ACTUAL mistakes. Do NOT invent rules that don't exist.
2. If unsure about a correction, skip it. False positives are worse than missed corrections.
3. Always return JSON only (no markdown, no code fences):
{
  "reply": "your reply in %s",
  "corrections": [
    {
      "original": "incorrect phrase",
      "corrected": "corrected phrase",
      "explanation_uz": "if lang is uz, explain in Uzbek",
      "explanation_ru": "if lang is ru, explain in Russian",
      "type": "grammar|vocabulary|spelling|word_order|punctuation",
      "severity": "critical|major|minor",
      "category": "grammar|vocabulary|spelling|word_order|punctuation",
      "learning_tip": "short tip to remember this rule",
      "rule_violated": "English grammar rule name"
    }
  ],
  "premium_analysis": {
    "overall_grade": "A|B|C|D",
    "strengths": ["strength 1", "strength 2"],
    "areas_for_improvement": ["area 1", "area 2"],
    "suggested_topic": "next topic to practice"
  }
}
4. If no mistakes, return "corrections": [].
5. corrections array may contain 1-5 items.
6. premium_analysis must be specific to the user's message, not generic.
7. areas_for_improvement should list 2-3 concrete points from the message.

User message: %s`, level, instr, historySection, langFullName(lang), text)
}
