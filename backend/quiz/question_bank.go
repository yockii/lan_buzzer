package quiz

import (
	"math/rand"
	"os"
	"time"
)

// QuestionBank manages a collection of questions
type QuestionBank struct {
	questions []*Question
	askedIDs  map[string]bool
}

// NewQuestionBank creates a new question bank
func NewQuestionBank(questions []*Question) *QuestionBank {
	return &QuestionBank{
		questions: questions,
		askedIDs:  make(map[string]bool),
	}
}

// HasQuestions returns true if there are questions available
func (qb *QuestionBank) HasQuestions() bool {
	return len(qb.questions) > 0
}

// GetRandomQuestion returns a random unasked question
func (qb *QuestionBank) GetRandomQuestion() *Question {
	// Filter unasked questions
	var unasked []*Question
	for _, q := range qb.questions {
		if !qb.askedIDs[q.ID] {
			unasked = append(unasked, q)
		}
	}

	// Still no questions (empty bank or all asked)
	if len(unasked) == 0 {
		return nil
	}

	// Return random question
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(unasked))
	return unasked[idx]
}

// MarkAsked marks a question as asked
func (qb *QuestionBank) MarkAsked(id string) {
	qb.askedIDs[id] = true
}

// ResetAsked clears the asked list
func (qb *QuestionBank) ResetAsked() {
	qb.askedIDs = make(map[string]bool)
}

// LoadFromFile loads questions from a file
func LoadQuestionBank(filePath string) (*QuestionBank, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	questions, err := ParseQuestions(file)
	if err != nil {
		return nil, err
	}

	return NewQuestionBank(questions), nil
}