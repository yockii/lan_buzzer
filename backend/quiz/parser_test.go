package quiz

import (
    "os"
    "strings"
    "testing"
)

func TestParseQuestions(t *testing.T) {
    file, err := os.Open("test_questions.txt")
    if err != nil {
        t.Fatalf("Failed to open test file: %v", err)
    }
    defer file.Close()

    questions, err := ParseQuestions(file)
    if err != nil {
        t.Fatalf("Failed to parse questions: %v", err)
    }

    if len(questions) != 3 {
        t.Errorf("Expected 3 questions, got %d", len(questions))
    }

    // Check first question (single choice)
    q1 := questions[0]
    if q1.Type != SingleChoice {
        t.Errorf("Expected single_choice, got %s", q1.Type)
    }
    if q1.Question != "什么是CPU?" {
        t.Errorf("Expected '什么是CPU?', got '%s'", q1.Question)
    }
    if len(q1.Options) != 4 {
        t.Errorf("Expected 4 options, got %d", len(q1.Options))
    }
    if q1.Answer != "A" {
        t.Errorf("Expected answer 'A', got '%s'", q1.Answer)
    }

    // Check second question (true/false)
    q2 := questions[1]
    if q2.Type != TrueFalse {
        t.Errorf("Expected true_false, got %s", q2.Type)
    }
    if len(q2.Options) != 2 {
        t.Errorf("Expected 2 options, got %d", len(q2.Options))
    }

    // Check third question (open-ended)
    q3 := questions[2]
    if q3.Type != OpenEnded {
        t.Errorf("Expected open_ended, got %s", q3.Type)
    }
    if len(q3.Options) != 0 {
        t.Errorf("Expected 0 options for open-ended, got %d", len(q3.Options))
    }
}

func TestParseQuestionsInvalidFormat(t *testing.T) {
    input := "invalid question without separator"
    reader := strings.NewReader(input)

    _, err := ParseQuestions(reader)
    if err == nil {
        t.Error("Expected error for invalid format, got nil")
    }
}

func TestParseQuestionsMissingTypeHeader(t *testing.T) {
    input := "什么是CPU?|A.中央处理器|B.显卡|A"
    reader := strings.NewReader(input)

    _, err := ParseQuestions(reader)
    if err == nil {
        t.Error("Expected error for question without type header")
    }
}