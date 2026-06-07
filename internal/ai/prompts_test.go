package ai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPremiumChatPrompt_ContainsPremiumSections(t *testing.T) {
	prompt := BuildPremiumChatPrompt("intermediate", "uz", "I has a apple")

	assert.Contains(t, prompt, "severity")
	assert.Contains(t, prompt, "category")
	assert.Contains(t, prompt, "learning_tip")
	assert.Contains(t, prompt, "rule_violated")
	assert.Contains(t, prompt, "overall_grade")
	assert.Contains(t, prompt, "strengths")
	assert.Contains(t, prompt, "areas_for_improvement")
	assert.Contains(t, prompt, "suggested_topic")
}

func TestBuildPremiumChatPrompt_IncludesUserInput(t *testing.T) {
	prompt := BuildPremiumChatPrompt("beginner", "ru", "He go to school")
	assert.Contains(t, prompt, "He go to school")
}

func TestBuildPremiumChatPrompt_HasPremiumFlag(t *testing.T) {
	prompt := BuildPremiumChatPrompt("advanced", "en", "test")
	assert.Contains(t, prompt, "premium")
}

func TestBuildPremiumChatPrompt_DifferentLevels(t *testing.T) {
	promptBeg := BuildPremiumChatPrompt("beginner", "uz", "test")
	promptAdv := BuildPremiumChatPrompt("advanced", "uz", "test")

	assert.Contains(t, promptBeg, "beginner")
	assert.Contains(t, promptAdv, "advanced")
	assert.NotEqual(t, promptBeg, promptAdv)
}

func TestBuildChatPrompt_DoesNotContainPremiumSections(t *testing.T) {
	prompt := BuildChatPrompt("intermediate", "uz", "I has a apple")

	assert.NotContains(t, prompt, "severity")
	assert.NotContains(t, prompt, "overall_grade")
	assert.NotContains(t, prompt, "premium_analysis")
}

func TestBuildPremiumChatPrompt_IsLongerThanNormal(t *testing.T) {
	input := "I has a apple"
	normal := BuildChatPrompt("intermediate", "uz", input)
	premium := BuildPremiumChatPrompt("intermediate", "uz", input)

	assert.Greater(t, len(premium), len(normal))
}

func TestBuildPremiumChatPrompt_LevelAndLang(t *testing.T) {
	prompt := BuildPremiumChatPrompt("beginner", "uz", "hello")
	assert.True(t, strings.Contains(prompt, "beginner"), "should contain level")
	assert.True(t, strings.Contains(prompt, "uz"), "should contain language code")
}

func TestBuildPremiumChatPrompt_EmptyInput(t *testing.T) {
	prompt := BuildPremiumChatPrompt("beginner", "uz", "")
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "User message:")
	assert.Contains(t, prompt, "premium_analysis")
}

func TestBuildChatPrompt_NormalPromptDoesntLeakPremium(t *testing.T) {
	input := "I has a apple"
	premium := BuildPremiumChatPrompt("intermediate", "uz", input)
	normal := BuildChatPrompt("intermediate", "uz", input)

	premiumKeywords := []string{"severity", "category", "learning_tip", "overall_grade", "strengths", "areas_to_improve", "premium_analysis", "suggested_topic"}
	for _, kw := range premiumKeywords {
		assert.NotContains(t, normal, kw, "normal prompt should not contain premium keyword: %s", kw)
	}

	assert.Greater(t, len(premium), len(normal))
}
