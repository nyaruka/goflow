package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/flows"
)

// ModelService is a deterministic model service for testing which derives its output from its input. Tests which need
// specific results should use MockModel instead.
type ModelService struct{}

func NewModel() *ModelService {
	return &ModelService{}
}

var leetify = strings.NewReplacer(
	"a", "4", "A", "4",
	"b", "8", "B", "8",
	"e", "3", "E", "3",
	"g", "9", "G", "9",
	"i", "1", "I", "1",
	"l", "1", "L", "1",
	"o", "0", "O", "0",
	"s", "5", "S", "5",
	"t", "7", "T", "7",
).Replace

func translate(s string) (string, error) {
	if s == "error" {
		return "", errors.New("simulated model error")
	}
	if s == "untranslatable" {
		return "<CANT>", nil
	}
	return leetify(s), nil
}

func (s *ModelService) Response(ctx context.Context, instructions, input string, maxTokens int) (*core.ModelResponse, error) {
	var output string
	if strings.HasPrefix(instructions, "Categorize") { // instructions like "Categorize... Category2, Category3]" will return "Category3"
		words := strings.Fields(instructions)
		output = strings.TrimSuffix(words[len(words)-1], "]")
	} else if strings.HasPrefix(instructions, "Translate") { // "Translate..." leetifies the input; if "JSON" is mentioned, values of a string->[]string object
		if strings.Contains(instructions, "JSON") {
			obj := map[string][]string{}
			if err := json.Unmarshal([]byte(input), &obj); err != nil {
				return nil, fmt.Errorf("invalid JSON object input: %w", err)
			}
			for k, vs := range obj {
				for i, v := range vs {
					tv, err := translate(v)
					if err != nil {
						return nil, err
					}
					vs[i] = tv
				}
				obj[k] = vs
			}
			b, _ := json.Marshal(obj)
			output = string(b)
		} else {
			tv, err := translate(input)
			if err != nil {
				return nil, err
			}
			output = tv
		}
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &core.ModelResponse{Output: output, TokensInput: 45, TokensOutput: 78}, nil
}

func (s *ModelService) Classify(ctx context.Context, input string, categories []string) (*core.ModelClassification, error) {
	// the last category is chosen, like a categorize prompt
	category := categories[len(categories)-1]
	probs := make(map[string]float64, len(categories))
	for _, c := range categories {
		if c == category {
			probs[c] = 0.9
		} else {
			probs[c] = 0.1
		}
	}

	return &core.ModelClassification{Category: category, Confidence: 0.8, Probabilities: probs, TokensInput: 34, TokensOutput: 5}, nil
}

// MockModelResult is a canned result for a call to a MockModel. A call to Response uses Output, a call to Classify uses
// Category, Confidence and Probabilities, and either returns Error instead if it's set.
type MockModelResult struct {
	Output        string             `json:"output,omitempty"`
	Category      string             `json:"category,omitempty"`
	Confidence    float64            `json:"confidence,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
	TokensInput   int64              `json:"tokens_input,omitempty"`
	TokensOutput  int64              `json:"tokens_output,omitempty"`
	Error         string             `json:"error,omitempty"`
}

// ModelCall is a call made to a MockModel
type ModelCall struct {
	Instructions string // set for Response calls
	Input        string
	MaxTokens    int      // set for Response calls
	Categories   []string // set for Classify calls
}

// MockModel is a model service for testing which answers each call with the next of its given results
type MockModel struct {
	mutex   sync.Mutex
	results []*MockModelResult
	calls   []*ModelCall
}

// NewMockModel creates a new mock model service which will return the given results in order
func NewMockModel(results ...*MockModelResult) *MockModel {
	return &MockModel{results: slices.Clone(results)}
}

func (m *MockModel) Response(ctx context.Context, instructions, input string, maxTokens int) (*core.ModelResponse, error) {
	r := m.next(&ModelCall{Instructions: instructions, Input: input, MaxTokens: maxTokens})
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}
	if r.Category != "" {
		panic("mock model result with category used for a response call")
	}

	return &core.ModelResponse{Output: r.Output, TokensInput: r.TokensInput, TokensOutput: r.TokensOutput}, nil
}

func (m *MockModel) Classify(ctx context.Context, input string, categories []string) (*core.ModelClassification, error) {
	r := m.next(&ModelCall{Input: input, Categories: categories})
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}
	if !slices.Contains(categories, r.Category) {
		panic(fmt.Sprintf("mock model result category '%s' isn't one of the classify call's categories", r.Category))
	}

	return &core.ModelClassification{Category: r.Category, Confidence: r.Confidence, Probabilities: r.Probabilities, TokensInput: r.TokensInput, TokensOutput: r.TokensOutput}, nil
}

func (m *MockModel) next(call *ModelCall) *MockModelResult {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.calls = append(m.calls, call)

	if len(m.results) == 0 {
		panic(fmt.Sprintf("missing mock model result for call with input '%s'", call.Input))
	}
	r := m.results[0]
	m.results = m.results[1:]
	return r
}

// Calls returns the calls made to this service so far
func (m *MockModel) Calls() []*ModelCall {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return slices.Clone(m.calls)
}

// HasUnused returns whether there are results which haven't been used
func (m *MockModel) HasUnused() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return len(m.results) > 0
}

var _ flows.ModelService = (*ModelService)(nil)
var _ flows.ModelService = (*MockModel)(nil)
