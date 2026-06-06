package ai

import (
	"math"
	"time"
)

type SM2Card struct {
	Repetitions int
	Interval    int
	EaseFactor  float64
}

func NewSM2Card() SM2Card {
	return SM2Card{
		Repetitions: 0,
		Interval:    0,
		EaseFactor:  2.5,
	}
}

type SM2Result struct {
	Repetitions   int
	Interval      int
	EaseFactor    float64
	NextReview    time.Time
	LastReviewed  time.Time
}

func ProcessReview(card SM2Card, quality int) SM2Result {
	now := time.Now().Truncate(24 * time.Hour)

	result := SM2Result{
		LastReviewed: now,
	}

	if quality < 3 {
		result.Repetitions = 0
		result.Interval = 1
		result.EaseFactor = card.EaseFactor
	} else {
		result.Repetitions = card.Repetitions + 1

		switch result.Repetitions {
		case 1:
			result.Interval = 1
		case 2:
			result.Interval = 6
		default:
			result.Interval = int(math.Round(float64(card.Interval) * card.EaseFactor))
		}

		newEF := card.EaseFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
		if newEF < 1.3 {
			newEF = 1.3
		}
		result.EaseFactor = math.Round(newEF*100) / 100
	}

	result.NextReview = now.AddDate(0, 0, result.Interval)

	return result
}
