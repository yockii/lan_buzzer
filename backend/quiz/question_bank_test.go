package quiz

import (
	"testing"
)

func TestQuestionBank(t *testing.T) {
	questions := []*Question{
		{ID: "q1", Type: SingleChoice, Question: "Q1", Answer: "A"},
		{ID: "q2", Type: SingleChoice, Question: "Q2", Answer: "B"},
		{ID: "q3", Type: SingleChoice, Question: "Q3", Answer: "C"},
	}

	qb := NewQuestionBank(questions)

	// Test HasQuestions
	if !qb.HasQuestions() {
		t.Error("Expected HasQuestions to return true")
	}

	// Test GetRandomQuestion
	q := qb.GetRandomQuestion()
	if q == nil {
		t.Error("Expected question, got nil")
	}

	// Test MarkAsked
	qb.MarkAsked(q.ID)
	q2 := qb.GetRandomQuestion()
	if q2 == nil {
		t.Error("Expected question, got nil")
	}

	// Eventually all questions should be asked
	for i := 0; i < 10; i++ {
		q := qb.GetRandomQuestion()
		if q != nil {
			qb.MarkAsked(q.ID)
		}
	}

	// After resetting, should get questions again
	qb.ResetAsked()
	q3 := qb.GetRandomQuestion()
	if q3 == nil {
		t.Error("Expected question after reset, got nil")
	}
}

func TestQuestionBankDeduplication(t *testing.T) {
	questions := []*Question{
		{ID: "q1", Type: SingleChoice, Question: "Q1", Answer: "A"},
		{ID: "q2", Type: SingleChoice, Question: "Q2", Answer: "B"},
	}

	qb := NewQuestionBank(questions)

	// Get first question
	q1 := qb.GetRandomQuestion()
	if q1 == nil {
		t.Error("Expected question")
	}

	// Mark it as asked
	qb.MarkAsked(q1.ID)

	// Should get the other question now
	q2 := qb.GetRandomQuestion()
	if q2 == nil {
		t.Error("Expected question after marking one as asked")
	}
	if q2.ID == q1.ID {
		t.Errorf("Expected different question, got same: %s", q2.ID)
	}

	// Mark the second as asked too
	qb.MarkAsked(q2.ID)

	// Now should get nil (no unasked questions)
	q3 := qb.GetRandomQuestion()
	if q3 != nil {
		t.Errorf("Expected nil when all questions asked, got %s", q3.ID)
	}

	// After reset, should get questions again
	qb.ResetAsked()
	q4 := qb.GetRandomQuestion()
	if q4 == nil {
		t.Error("Expected question after reset")
	}
}
