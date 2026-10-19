package tests

import (
	"context"
	"testing"

	"github.com/DPLakmal/guidenow-platform/go/internal/ai"
)

func TestAISearchFilterParsing(t *testing.T) {
	parser := ai.NewParser("") // Use offline rule-based heuristic parser
	ctx := context.Background()

	tests := []struct {
		name               string
		query              string
		expectedLocation   string
		expectedLanguages  []string
		expectedMaxRate    float64
		expectedImmediate  bool
	}{
		{
			name:               "Galle history tour with language and budget",
			query:              "Looking for an English or German speaking guide in Galle for history and architecture under $25/hr",
			expectedLocation:   "Galle",
			expectedLanguages:  []string{"English", "German"},
			expectedMaxRate:    25.0,
			expectedImmediate:  false,
		},
		{
			name:               "Weligama surfing instructor with immediate availability",
			query:              "Need a surfing coach in Weligama under $35 available right now",
			expectedLocation:   "Weligama",
			expectedLanguages:  nil,
			expectedMaxRate:    35.0,
			expectedImmediate:  true,
		},
		{
			name:               "Sigiriya French guide",
			query:              "Find a French guide for Sigiriya Lion rock sunrise climb",
			expectedLocation:   "Sigiriya",
			expectedLanguages:  []string{"French"},
			expectedMaxRate:    0.0,
			expectedImmediate:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parser.Parse(ctx, tc.query)
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			if res.Location != tc.expectedLocation {
				t.Errorf("Expected location %s, got %s", tc.expectedLocation, res.Location)
			}

			if tc.expectedMaxRate > 0 && res.MaxRate != tc.expectedMaxRate {
				t.Errorf("Expected max rate %.2f, got %.2f", tc.expectedMaxRate, res.MaxRate)
			}

			if tc.expectedImmediate != res.ImmediateOnly {
				t.Errorf("Expected immediate %v, got %v", tc.expectedImmediate, res.ImmediateOnly)
			}

			for _, expectedLang := range tc.expectedLanguages {
				found := false
				for _, actualLang := range res.Languages {
					if actualLang == expectedLang {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected language %s not found in result %v", expectedLang, res.Languages)
				}
			}
		})
	}
}
