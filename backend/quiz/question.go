package quiz

import "strings"

// QuestionType represents the type of question
type QuestionType string

const (
	SingleChoice QuestionType = "single_choice"
	TrueFalse    QuestionType = "true_false"
	OpenEnded   QuestionType = "open_ended"
)

// Question represents a quiz question
type Question struct {
	ID       string
	Type     QuestionType
	Question string
	Options  []string // For single_choice and true_false
	Answer   string   // Correct answer
}

// IsCorrect checks if the given answer is correct
func (q *Question) IsCorrect(answer string) bool {
	// Trim space and punctuation for more lenient matching
	// This handles cases like "北京" vs "北京市" or "对" vs "对的"
	normalizedCorrect := normalizeAnswer(q.Answer)
	normalizedAnswer := normalizeAnswer(answer)
	println("[DEBUG] IsCorrect: question.Answer='", q.Answer, "' answer='", answer, "' normalizedCorrect='", normalizedCorrect, "' normalizedAnswer='", normalizedAnswer, "' result='", normalizedCorrect == normalizedAnswer, "'")
	return normalizedCorrect == normalizedAnswer
}

// normalizeAnswer removes common punctuation and extra whitespace
func normalizeAnswer(s string) string {
	// Remove common Chinese and English punctuation
	punctuation := []string{"。", "！", "？", "，", "、", "；", "：", "\"", "'", "（", "）", "的", "了", "着", "啊",
		".", ",", "!", "?", ";", ":", "\"", "'", "(", ")", " ", "\t", "\n", "\r"}
	result := s
	for _, p := range punctuation {
		result = strings.ReplaceAll(result, p, "")
	}
	return strings.TrimSpace(result)
}

// NeedsManualJudgment returns true if this question type requires manual judging
func (q *Question) NeedsManualJudgment() bool {
	return q.Type == OpenEnded
}
