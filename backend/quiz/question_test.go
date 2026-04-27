package quiz

import "testing"

func TestQuestionTypes(t *testing.T) {
    mcQ := &Question{
        ID:       "q1",
        Type:     SingleChoice,
        Question: "What is CPU?",
        Options:  []string{"A. CPU", "B. GPU", "C. RAM", "D. Disk"},
        Answer:   "A",
    }

    if mcQ.Type != SingleChoice {
        t.Errorf("Expected SingleChoice, got %s", mcQ.Type)
    }

    tfQ := &Question{
        ID:       "q2",
        Type:     TrueFalse,
        Question: "Earth is round?",
        Options:  []string{"对", "错"},
        Answer:   "对",
    }

    if tfQ.Type != TrueFalse {
        t.Errorf("Expected TrueFalse, got %s", tfQ.Type)
    }

    oeQ := &Question{
        ID:       "q3",
        Type:     OpenEnded,
        Question: "Capital of China?",
        Answer:   "北京",
    }

    if oeQ.Type != OpenEnded {
        t.Errorf("Expected OpenEnded, got %s", oeQ.Type)
    }
}

func TestQuestionMethods(t *testing.T) {
    // Test IsCorrect
    q := &Question{
        ID:       "q1",
        Type:     SingleChoice,
        Question: "Test?",
        Answer:   "A",
    }

    if !q.IsCorrect("A") {
        t.Error("Expected IsCorrect to return true for matching answer")
    }

    if q.IsCorrect("B") {
        t.Error("Expected IsCorrect to return false for non-matching answer")
    }

    // Test NeedsManualJudgment
    mcQ := &Question{Type: SingleChoice}
    tfQ := &Question{Type: TrueFalse}
    oeQ := &Question{Type: OpenEnded}

    if mcQ.NeedsManualJudgment() {
        t.Error("Single choice should not need manual judgment")
    }

    if tfQ.NeedsManualJudgment() {
        t.Error("True/false should not need manual judgment")
    }

    if !oeQ.NeedsManualJudgment() {
        t.Error("Open-ended should need manual judgment")
    }
}