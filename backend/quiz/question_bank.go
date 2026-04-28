package quiz

import (
	"math/rand"
	"os"
	"sync"
)

// QuestionBank manages a collection of questions
type QuestionBank struct {
	questions    []*Question
	askedIDs     map[string]bool
	availableIDs []string
	mu           sync.RWMutex
}

// NewQuestionBank creates a new question bank
func NewQuestionBank(questions []*Question) *QuestionBank {
	availableIDs := make([]string, len(questions))
	for i, q := range questions {
		availableIDs[i] = q.ID
	}

	return &QuestionBank{
		questions:    questions,
		askedIDs:     make(map[string]bool),
		availableIDs: availableIDs,
	}
}

// HasQuestions returns true if there are questions available
func (qb *QuestionBank) HasQuestions() bool {
	qb.mu.RLock()
	defer qb.mu.RUnlock()
	return len(qb.availableIDs) > 0
}

// GetRandomQuestion returns a random unasked question
func (qb *QuestionBank) GetRandomQuestion() *Question {
	qb.mu.RLock()
	defer qb.mu.RUnlock()

	// No available questions
	if len(qb.availableIDs) == 0 {
		return nil
	}

	// Return random question from available ones
	idx := rand.Intn(len(qb.availableIDs))
	questionID := qb.availableIDs[idx]

	// Find the question object
	for _, q := range qb.questions {
		if q.ID == questionID {
			return q
		}
	}

	return nil
}

// MarkAsked marks a question as asked
func (qb *QuestionBank) MarkAsked(id string) {
	qb.mu.Lock()
	defer qb.mu.Unlock()

	qb.askedIDs[id] = true

	// Remove from availableIDs
	for i, availID := range qb.availableIDs {
		if availID == id {
			qb.availableIDs = append(qb.availableIDs[:i], qb.availableIDs[i+1:]...)
			break
		}
	}
}

// ResetAsked clears the asked list
func (qb *QuestionBank) ResetAsked() {
	qb.mu.Lock()
	defer qb.mu.Unlock()

	qb.askedIDs = make(map[string]bool)

	// Reset availableIDs to include all question IDs
	qb.availableIDs = make([]string, len(qb.questions))
	for i, q := range qb.questions {
		qb.availableIDs[i] = q.ID
	}
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
