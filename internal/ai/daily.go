package ai

type DailyExercise struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Options  []string `json:"options"`
}

type DailyVocab struct {
	Word          string `json:"word"`
	TranslationUz string `json:"translation_uz"`
	TranslationRu string `json:"translation_ru"`
}

type DailyLessonResponse struct {
	Topic        string          `json:"topic"`
	ExplanationUz string         `json:"explanation_uz"`
	ExplanationRu string         `json:"explanation_ru"`
	Examples     []string        `json:"examples"`
	Exercises    []DailyExercise `json:"exercises"`
	Vocabulary   []DailyVocab    `json:"vocabulary"`
}
