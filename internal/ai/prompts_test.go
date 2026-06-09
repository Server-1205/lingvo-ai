package ai

import (
	"strings"
	"testing"
)

func TestBuildChatPrompt_IncludesAntiHallucination(t *testing.T) {
	prompt := BuildChatPrompt("a1", "uz", "hello")
	if !strings.Contains(prompt, "Never invent grammar rules") && !strings.Contains(prompt, "Do NOT invent grammar rules") {
		t.Error("Chat prompt missing anti-hallucination guard")
	}
	if !strings.Contains(prompt, "only correct ACTUAL mistakes") && !strings.Contains(prompt, "Only correct ACTUAL mistakes") {
		t.Error("Chat prompt missing false positive guard")
	}
}

func TestBuildChatPrompt_ReturnsJSONOnly(t *testing.T) {
	prompt := BuildChatPrompt("b1", "ru", "test")
	if !strings.Contains(prompt, "No markdown") && !strings.Contains(prompt, "no markdown") {
		t.Error("Chat prompt missing JSON-only instruction")
	}
}

func TestBuildGrammarCheckPrompt_HasReplyField(t *testing.T) {
	prompt := BuildGrammarCheckPrompt("a2", "uz", "test")
	if !strings.Contains(prompt, `"reply"`) {
		t.Error("Grammar prompt missing reply field")
	}
	if !strings.Contains(prompt, "Only flag REAL errors") {
		t.Error("Grammar prompt missing false positive guard")
	}
}

func TestBuildPremiumChatPrompt_HasSpecificAnalysis(t *testing.T) {
	prompt := BuildPremiumChatPrompt("b2", "ru", "test")
	if !strings.Contains(prompt, "specific to the user's message") {
		t.Error("Premium prompt missing specificity instruction")
	}
	if !strings.Contains(prompt, "Do NOT invent") {
		t.Error("Premium prompt missing anti-hallucination guard")
	}
}

func TestBuildDailyLessonPrompt_HasVarietyInstruction(t *testing.T) {
	prompt := BuildDailyLessonPrompt("a1", "uz")
	if !strings.Contains(prompt, "different exercise types") {
		t.Error("Daily lesson prompt missing exercise variety instruction")
	}
}

func TestBuildDailyLessonPrompt_ContainsVocabulary(t *testing.T) {
	prompt := BuildDailyLessonPrompt("b1", "ru")
	if !strings.Contains(prompt, "vocabulary") {
		t.Error("Daily lesson prompt missing vocabulary section")
	}
}

func TestBuildVocabPrompt_RejectsProfanity(t *testing.T) {
	prompt := BuildVocabPrompt("uz", "test")
	if !strings.Contains(prompt, "offensive, obscene, or vulgar") {
		t.Error("Vocab prompt missing profanity filter")
	}
	if !strings.Contains(prompt, "misspellings or variants") {
		t.Error("Vocab prompt missing variant detection")
	}
}

func TestBuildQuizPrompt_ContainsLanguage(t *testing.T) {
	prompt := BuildQuizPrompt("grammar", 5, "ru")
	if !strings.Contains(prompt, "explanation_ru") {
		t.Error("Quiz prompt missing Russian explanation field")
	}
}

func TestBuildLevelTestPrompt_ContainsLanguages(t *testing.T) {
	prompt := BuildLevelTestPrompt("uz")
	if !strings.Contains(prompt, "explanation_uz") && !strings.Contains(prompt, "explanation_ru") {
		t.Error("Level test prompt missing explanation fields")
	}
}

func TestLangFullName(t *testing.T) {
	if langFullName("uz") != "Uzbek" {
		t.Error("uz should map to Uzbek")
	}
	if langFullName("ru") != "Russian" {
		t.Error("ru should map to Russian")
	}
	if langFullName("en") != "English" {
		t.Error("en should map to English")
	}
	if langFullName("fr") != "English" {
		t.Error("unknown should default to English")
	}
}

func TestLangInstruction(t *testing.T) {
	if _, ok := langInstruction["uz"]; !ok {
		t.Error("Missing uz language instruction")
	}
	if _, ok := langInstruction["ru"]; !ok {
		t.Error("Missing ru language instruction")
	}
}
