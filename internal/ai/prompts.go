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

func BuildIeltsWritingPrompt(taskType, lang, userText, taskDescription string) string {
	format := "report describing a chart/graph/data"
	if taskType == "task2" {
		format = "essay discussing a topic"
	}
	return fmt.Sprintf(`You are an IELTS Writing examiner. Evaluate this %s according to official IELTS criteria.

Task: %s

User's response: %s

Return JSON only:
{
  "band_score": 6.5,
  "criteria": {
    "task_achievement": 6.0,
    "coherence_cohesion": 6.5,
    "lexical_resource": 7.0,
    "grammatical_range": 6.0
  },
  "feedback": "detailed feedback in %s",
  "corrections": [
    {
      "original": "wrong phrase",
      "corrected": "fixed phrase",
      "explanation_uz": "if lang=uz, explanation in Uzbek",
      "explanation_ru": "if lang=ru, explanation in Russian",
      "type": "grammar|vocabulary|spelling"
    }
  ],
  "improvement_tips": ["tip 1", "tip 2", "tip 3"]
}`, format, taskDescription, userText, langFullName(lang))
}

func BuildIeltsSpeakingPrompt(part int, lang string) string {
	var partDesc string
	switch part {
	case 1:
		partDesc = "Part 1 — Introduction and general questions about the candidate's life, work, studies, hobbies. Generate 3-4 simple questions."
	case 2:
		partDesc = `Part 2 — Cue Card. Generate ONE topic card in this format:
"The candidate should speak for 1-2 minutes on the following topic. Describe/Lets talk about [topic]. Include the following points: [3 bullet points]".`
	case 3:
		partDesc = "Part 3 — Discussion. Generate 2-3 abstract/discussion questions related to the Part 2 topic. These should be more complex, requiring opinion, comparison, prediction."
	}

	return fmt.Sprintf(`You are an IELTS Speaking examiner. Generate questions for %s.

Return JSON only:
{
  "part": %d,
  "questions": ["question 1", "question 2", "question 3"],
  "cue_card": "full cue card text for Part 2 (optional)"
}

Explain exam instructions in %s.`, partDesc, part, langFullName(lang))
}

func BuildIeltsSpeakingEvaluatePrompt(part int, lang, question, userResponse string) string {
	return fmt.Sprintf(`You are an IELTS Speaking examiner. Evaluate this candidate's response for Part %d.

Question: %s

Candidate's response: %s

Return JSON only:
{
  "band_score": 6.0,
  "criteria": {
    "fluency_coherence": 6.0,
    "lexical_resource": 6.5,
    "grammatical_range": 6.0,
    "pronunciation": 5.5
  },
  "feedback": "detailed feedback in %s",
  "improvement_tips": ["tip 1", "tip 2", "tip 3"]
}`, part, question, userResponse, langFullName(lang))
}

func BuildIeltsReadingPrompt(lang string) string {
	return fmt.Sprintf(`You are an IELTS Reading examiner. Generate a reading passage with 10 questions.

Requirements:
- Passage length: 400-600 words
- Mix of question types: multiple_choice, true_false, gap_fill, matching
- The passage should be academic/general interest
- Difficulty: IELTS Academic level

Return JSON only:
{
  "title": "passage title",
  "passage": "full reading text here...",
  "word_count": 450,
  "questions": [
    {
      "type": "multiple_choice|true_false|gap_fill|matching",
      "question": "question text",
      "options": ["A", "B", "C", "D", "True", "False", "Not Given"],
      "correct": 0
    }
  ]
}

Explain instructions in %s.`, langFullName(lang))
}

func BuildIeltsReadingEvaluatePrompt(lang, passage string, questionsJSON string, userAnswersJSON string) string {
	return fmt.Sprintf(`You are an IELTS Reading examiner. Check these answers against the passage.

Passage: %s

Questions: %s

User's answers (indices): %s

Return JSON only:
{
  "correct_answers": 7,
  "total_questions": 10,
  "band_score": 6.0,
  "results": [
    {
      "question_index": 0,
      "user_answer": 0,
      "correct_answer": 1,
      "is_correct": false,
      "explanation_uz": "if lang=uz, explanation in Uzbek",
      "explanation_ru": "if lang=ru, explanation in Russian"
    }
  ],
  "feedback": "overall feedback in %s"
}`, passage, questionsJSON, userAnswersJSON, langFullName(lang))
}
