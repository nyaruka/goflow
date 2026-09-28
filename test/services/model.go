package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/nyaruka/gocommon/i18n"
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
	} else if strings.HasPrefix(instructions, "Translate") { // "Translate..." leetifies the input
		tv, err := translate(input)
		if err != nil {
			return nil, err
		}
		output = tv
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &core.ModelResponse{Output: output, Tokens: core.ModelTokens{Input: 45, Output: 78}}, nil
}

func (s *ModelService) Classify(ctx context.Context, input string, options []*core.ClassifierOption) (*core.Classification, error) {
	// the last option is chosen, like a categorize prompt
	option := options[len(options)-1].Name
	probs := make(map[string]float64, len(options))
	for _, o := range options {
		if o.Name == option {
			probs[o.Name] = 0.9
		} else {
			probs[o.Name] = 0.1
		}
	}

	return &core.Classification{Option: option, Confidence: 0.8, Probabilities: probs, Tokens: core.ModelTokens{Input: 34, Output: 5}}, nil
}

func (s *ModelService) Translate(ctx context.Context, source, target i18n.Language, items map[string][]string) (*core.Translation, error) {
	translated := make(map[string][]string, len(items))

	for key, vals := range items {
		tvals := make([]string, len(vals))
		for i, v := range vals {
			tv, err := translate(v)
			if err != nil {
				return nil, err
			}
			tvals[i] = tv
		}
		if !slices.Contains(tvals, "<CANT>") {
			translated[key] = tvals
		}
	}

	return &core.Translation{Items: translated, Tokens: core.ModelTokens{Input: 56, Output: 67}}, nil
}

// MockModelResult is a canned result for a call to a MockModel. A call to Response uses Output, a call to Classify uses
// Option, Confidence and Probabilities, a call to Translate uses Items, and any returns Error instead if it's set.
type MockModelResult struct {
	Output        string              `json:"output,omitempty"`
	Option        string              `json:"option,omitempty"`
	Confidence    float64             `json:"confidence,omitempty"`
	Probabilities map[string]float64  `json:"probabilities,omitempty"`
	Items         map[string][]string `json:"items,omitempty"`
	Tokens        core.ModelTokens    `json:"tokens,omitzero"`
	Error         string              `json:"error,omitempty"`
}

// ModelCall is a call made to a MockModel
type ModelCall struct {
	Instructions string // set for Response calls
	Input        string
	MaxTokens    int                      // set for Response calls
	Options      []*core.ClassifierOption // set for Classify calls
	Source       i18n.Language            // set for Translate calls
	Target       i18n.Language            // set for Translate calls
	Items        map[string][]string      // set for Translate calls
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
	if r.Option != "" || r.Items != nil {
		panic("mock model result for another call type used for a response call")
	}

	return &core.ModelResponse{Output: r.Output, Tokens: r.Tokens}, nil
}

func (m *MockModel) Classify(ctx context.Context, input string, options []*core.ClassifierOption) (*core.Classification, error) {
	r := m.next(&ModelCall{Input: input, Options: options})
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}
	if !slices.ContainsFunc(options, func(o *core.ClassifierOption) bool { return o.Name == r.Option }) {
		panic(fmt.Sprintf("mock model result option '%s' isn't one of the classify call's options", r.Option))
	}

	return &core.Classification{Option: r.Option, Confidence: r.Confidence, Probabilities: r.Probabilities, Tokens: r.Tokens}, nil
}

func (m *MockModel) Translate(ctx context.Context, source, target i18n.Language, items map[string][]string) (*core.Translation, error) {
	r := m.next(&ModelCall{Source: source, Target: target, Items: items})
	if r.Error != "" {
		return nil, errors.New(r.Error)
	}
	if r.Option != "" || r.Output != "" {
		panic("mock model result for another call type used for a translate call")
	}

	return &core.Translation{Items: r.Items, Tokens: r.Tokens}, nil
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
