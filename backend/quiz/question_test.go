package quiz

import "testing"

func TestQuestionTypes(t *testing.T) {
    mcQ := &Question{
        ID:       "q1",
        Type:     "single_choice",
        Question: "What is CPU?",
        Options:  []string{"A. CPU", "B. GPU", "C. RAM", "D. Disk"},
        Answer:   "A",
    }

    if mcQ.Type != "single_choice" {
        t.Errorf("Expected single_choice, got %s", mcQ.Type)
    }

    tfQ := &Question{
        ID:       "q2",
        Type:     "true_false",
        Question: "Earth is round?",
        Options:  []string{"对", "错"},
        Answer:   "对",
    }

    if tfQ.Type != "true_false" {
        t.Errorf("Expected true_false, got %s", tfQ.Type)
    }

    oeQ := &Question{
        ID:       "q3",
        Type:     "open_ended",
        Question: "Capital of China?",
        Answer:   "北京",
    }

    if oeQ.Type != "open_ended" {
        t.Errorf("Expected open_ended, got %s", oeQ.Type)
    }
}
