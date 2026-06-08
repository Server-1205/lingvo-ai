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
	"uz": "IMPORTANT: Always reply and explain everything in Uzbek. The user is an Uzbek speaker learning English. All explanations, corrections, and the reply itself must be in Uzbek.",
	"ru": "IMPORTANT: Always reply and explain everything in Russian. The user is a Russian speaker learning English. All explanations, corrections, and the reply itself must be in Russian.",
}

func BuildChatPrompt(level, lang, text string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Reply naturally in English, keeping your response concise (2-4 sentences)."
	}
	return fmt.Sprintf(`You are an AI English tutor. The user's English level is %s.

%s

Rules:
1. Only correct ACTUAL mistakes. If the sentence is correct, return "corrections": [].
2. Do NOT invent grammar rules that don't exist in standard English.
3. If you are unsure about a correction, do NOT include it. False positives frustrate the user.
4. Always return JSON only (no markdown, no code fences):
{
  "reply": "your reply in %s",
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
5. If no mistakes, return "corrections": [].
6. corrections array must be empty or contain 1-3 items max.

User message: %s`, level, instr, langFullName(lang), text)
}

func BuildGrammarCheckPrompt(level, lang, text string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Explain in English."
	}
	return fmt.Sprintf(`Check this text for grammar errors. User level: %s.

%s

Rules:
1. Only flag REAL errors. Do not correct informal but acceptable English.
2. If the text has no errors, return "corrections": [].
3. Always return JSON only (no markdown, no code fences):

{
  "reply": "brief summary of what was checked",
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

Text: %s`, level, instr, text)
}

func BuildPremiumChatPrompt(level, lang, text string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Reply naturally in English, keeping your response concise (2-4 sentences)."
	}
	return fmt.Sprintf(`You are an AI English tutor. The user's English level is %s.

%s

You are in PREMIUM mode — provide deeper analysis.

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

User message: %s`, level, instr, langFullName(lang), text)
}

func BuildVocabPrompt(lang, word string) string {
	return fmt.Sprintf(`You are a bilingual English-Uzbek/Russian dictionary. The user enters a word in any language (English, Uzbek, or Russian). You must detect the language, find the English word, and return translations to both Uzbek and Russian.

Rules:
- If the input is not English, first identify the correct English word, then translate it.
- Reject any offensive, obscene, or vulgar words with {"error": "inappropriate_word"}.
- Also reject words that are misspellings or variants of offensive words.
- Always return JSON only:

{
  "word_en": "the English word",
  "translation_uz": "translation to Uzbek",
  "translation_ru": "translation to Russian",
  "examples_uz": ["example sentence in Uzbek 1", "example sentence in Uzbek 2"],
  "examples_ru": ["example sentence in Russian 1", "example sentence in Russian 2"],
  "level": "a1|a2|b1|b2|c1"
}

Input: %s. User language: %s`, word, lang)
}

func BuildQuizPrompt(topic string, count int, lang string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Explain in English."
	}
	return fmt.Sprintf(`Generate %d English quiz questions about "%s".

%s

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
}`, count, topic, instr)
}

func BuildDailyLessonPrompt(level, lang string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Explain in English."
	}
	return fmt.Sprintf(`You are an AI English tutor. Generate a 5-minute personalized lesson for a student at level %s.

%s

Rules:
1. Choose ONE grammar topic appropriate for the student's level.
2. Explain the rule simply with 2-3 examples.
3. Provide 3-4 exercises with at least 2 different exercise types (mix fill-in-blank, multiple choice, matching).
4. Include 2 new vocabulary words related to the topic with translations.
5. Always return JSON only (no markdown, no code fences):

{
  "topic": "Grammar topic name",
  "explanation_uz": "explanation in Uzbek",
  "explanation_ru": "explanation in Russian",
  "examples": ["Example sentence 1", "Example sentence 2", "Example sentence 3"],
  "exercises": [
    {
      "question": "She ___ (go) to school yesterday.",
      "answer": "went",
      "options": ["go", "went", "gone"]
    }
  ],
  "vocabulary": [
    {
      "word": "yesterday",
      "translation_uz": "kecha",
      "translation_ru": "vchera"
    }
  ]
}

Student level: %s`, level, instr, level)
}

func BuildLevelTestPrompt(lang string) string {
	instr := langInstruction[lang]
	if instr == "" {
		instr = "Explain in English."
	}
	return fmt.Sprintf(`Generate 10 English level test questions (A1 to C1).

%s

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
}`, instr)
}
