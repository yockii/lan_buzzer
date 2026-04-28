package quiz

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseQuestions parses questions from a text file
func ParseQuestions(file io.Reader) ([]*Question, error) {
	scanner := bufio.NewScanner(file)
	var questions []*Question
	var currentType QuestionType
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for type header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			typeStr := strings.Trim(line, "[]")
			switch typeStr {
			case "单选":
				currentType = SingleChoice
			case "判断":
				currentType = TrueFalse
			case "问答":
				currentType = OpenEnded
			default:
				return nil, fmt.Errorf("line %d: unknown question type '%s'", lineNum, typeStr)
			}
			continue
		}

		// Parse question
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected at least 2 fields separated by |", lineNum)
		}

		// Check if type header has been set
		if currentType == "" {
			return nil, fmt.Errorf("line %d: question appears before type header", lineNum)
		}

		question := &Question{
			ID:       fmt.Sprintf("q%d", lineNum),
			Type:     currentType,
			Question: parts[0],
			Answer:   parts[len(parts)-1],
		}

		// Add options for single_choice and true_false
		if currentType == SingleChoice || currentType == TrueFalse {
			// For true/false questions, we need to generate the "对" and "错" options
			if currentType == TrueFalse {
				question.Options = []string{"对", "错"}
			} else {
				// Parts[1:-1] are options for single_choice
				question.Options = parts[1 : len(parts)-1]
			}
		}

		questions = append(questions, question)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return questions, nil
}
