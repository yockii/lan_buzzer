package quiz

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
    return q.Answer == answer
}

// NeedsManualJudgment returns true if this question type requires manual judging
func (q *Question) NeedsManualJudgment() bool {
    return q.Type == OpenEnded
}