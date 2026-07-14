package excellent_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// a mix of templates that exercise every segment kind - literal text, @@ escapes, @identifier and @(...) forms,
// valid and invalid identifiers, malformed expressions, and evaluation errors
var templateTestCases = []string{
	``,
	`hello world`,
	`@string1`,
	`@string1 world`,
	`Hello @(title(string1))`,
	`@(string1 & string2)`,
	`@("hello\nworld")`,
	`@("\"hello\nworld\"")`,
	`@( "quoted )pa)ren" )`,
	`@(  1  +  2  )`,
	`@(title("hello"))`,
	`@(title(hello))`,
	`@string1.xxx`,
	`@hello`,
	`@hello.bar`,
	`My email is foo@bar.com`,
	`@`,
	`@@`,
	`@@string1`,
	`@@@string1`,
	`@@@@string1`,
	`Hello @@string1 world`,
	`@string1@string2`,
	`@string1.@string2`,
	`@string1.@string2.@string3`,
	`@(string1`,
	`@ (string1`,
	`@ (string1)`,
	`@((`,
	`@('x')`,
	`@(0 / )`,
	`@(0 / )@('x')`,
	`@(1.1.0)`,
	`@array1`,
	`@array1[0]`,
	`@array1.0`,
	`@(array1[0])`,
	`@array2d.0.2`,
	`@_special`,
	`@汉字 chinese`,
	`@(汉字)`,
	`@😀`,
	`@legacy and @legacy again`,
	`@(1 / 0)`,
	`@(err)`,
	`1 + 2`,
	"line1\nline2 @string1\n@@line3",
}

func makeTemplateContext() *types.XObject {
	legacy := types.NewXText("old")
	legacy.SetDeprecated("use something else")

	return types.NewXObject(map[string]types.XValue{
		"string1":  types.NewXText("foo"),
		"string2":  types.NewXText("bar"),
		"_special": types.NewXText("🐒"),
		"汉字":       types.NewXText("simplified chinese"),
		"int1":     types.NewXNumberFromInt(1),
		"dec1":     types.RequireXNumberFromString("1.5"),
		"words":    types.NewXText("one two three"),
		"array1":   types.NewXArray(types.NewXText("one"), types.NewXText("two"), types.NewXText("three")),
		"legacy":   legacy,
		"err":      types.NewXError(fmt.Errorf("an error")),
		"func":     functions.Lookup("upper"),
		"thing": types.NewXObject(map[string]types.XValue{
			"foo": types.NewXText("bar"),
			"zed": types.NewXNumberFromInt(123),
		}),
	})
}

// asserts that parsing the given source round-trips exactly and that evaluating the parsed template produces
// the same output, warnings and error as the existing evaluator
func assertTemplateEquivalent(t *testing.T, env envs.Environment, ctx *types.XObject, src string, escaping excellent.Escaping) {
	t.Helper()

	template := excellent.ParseTemplate(src)
	require.Equal(t, src, template.String(), "round trip mismatch for template %q", src)

	expectedVal, expectedWarns, expectedErr := excellent.NewEvaluator().Template(env, ctx, src, escaping)
	actualVal, actualWarns, actualErr := template.Evaluate(env, ctx, escaping)

	assert.Equal(t, expectedVal, actualVal, "output mismatch for template %q", src)
	assert.Equal(t, expectedWarns, actualWarns, "warnings mismatch for template %q", src)

	if expectedErr != nil {
		require.Error(t, actualErr, "expected error for template %q", src)
		assert.Equal(t, expectedErr.Error(), actualErr.Error(), "error mismatch for template %q", src)
	} else {
		assert.NoError(t, actualErr, "unexpected error for template %q", src)
	}
}

func TestParseTemplate(t *testing.T) {
	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()

	escaping := func(s string) string { return strings.ReplaceAll(s, `"`, `\"`) }

	for _, src := range templateTestCases {
		assertTemplateEquivalent(t, env, ctx, src, nil)
		assertTemplateEquivalent(t, env, ctx, src, escaping)
	}

	// context with no matching properties means identifiers are treated as literal text
	empty := types.NewXObject(map[string]types.XValue{})
	assertTemplateEquivalent(t, env, empty, `Hello @string1 and @(string2)`, nil)
}

// tests that a template parsed once can be safely evaluated concurrently (run with -race).. note that contexts
// are per-goroutine because XObject lazily initializes itself and so isn't itself safe for concurrent use
func TestTemplateConcurrentUse(t *testing.T) {
	env := envs.NewBuilder().Build()

	template := excellent.ParseTemplate(`Hello @string1, @thing.foo is @(int1 + 2) @@here`)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := makeTemplateContext()
			for range 100 {
				val, _, err := template.Evaluate(env, ctx, nil)
				assert.NoError(t, err)
				assert.Equal(t, `Hello foo, bar is 3 @here`, val)
			}
		}()
	}
	wg.Wait()
}

// loads a corpus of template strings from the flow definitions in the repo's test fixtures
func loadFixtureTemplates(tb testing.TB) []string {
	paths, err := filepath.Glob("../test/testdata/runner/*.json")
	require.NoError(tb, err)
	require.NotEmpty(tb, paths)

	seen := make(map[string]bool)
	var templates []string

	var collect func(v any)
	collect = func(v any) {
		switch typed := v.(type) {
		case string:
			if !seen[typed] {
				seen[typed] = true
				templates = append(templates, typed)
			}
		case map[string]any:
			for _, item := range typed {
				collect(item)
			}
		case []any:
			for _, item := range typed {
				collect(item)
			}
		}
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		require.NoError(tb, err)

		var parsed any
		require.NoError(tb, json.Unmarshal(data, &parsed))
		collect(parsed)
	}

	return templates
}

func TestParseTemplateWithFixtures(t *testing.T) {
	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()

	templates := loadFixtureTemplates(t)
	require.Greater(t, len(templates), 1000) // sanity check that we loaded a real corpus

	for _, src := range templates {
		assertTemplateEquivalent(t, env, ctx, src, nil)
	}
}

func FuzzParseTemplate(f *testing.F) {
	for _, src := range templateTestCases {
		f.Add(src)
	}

	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()

	f.Fuzz(func(t *testing.T, src string) {
		// the scanner reads runes so it only round-trips valid UTF-8, and treats NUL as end of input
		if !utf8.ValidString(src) || strings.ContainsRune(src, 0) {
			t.Skip()
		}

		template := excellent.ParseTemplate(src)
		if template.String() != src {
			t.Fatalf("round trip mismatch: %q != %q", template.String(), src)
		}

		// skip evaluation of very large inputs to keep fuzzing fast
		if len(src) > 1000 {
			return
		}

		expectedVal, expectedWarns, expectedErr := excellent.NewEvaluator().Template(env, ctx, src, nil)
		actualVal, actualWarns, actualErr := template.Evaluate(env, ctx, nil)

		if actualVal != expectedVal {
			t.Fatalf("output mismatch for %q: %q != %q", src, actualVal, expectedVal)
		}
		if !slices.Equal(actualWarns, expectedWarns) {
			t.Fatalf("warnings mismatch for %q: %v != %v", src, actualWarns, expectedWarns)
		}
		if (actualErr != nil) != (expectedErr != nil) || (actualErr != nil && actualErr.Error() != expectedErr.Error()) {
			t.Fatalf("error mismatch for %q: %v != %v", src, actualErr, expectedErr)
		}
	})
}

const benchmarkTemplate = `Hello @string1, @thing.foo is @(int1 + 2) @@here with @(upper(words))`

func BenchmarkEvaluatorTemplate(b *testing.B) {
	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()
	eval := excellent.NewEvaluator()

	b.ResetTimer()
	for range b.N {
		eval.Template(env, ctx, benchmarkTemplate, nil)
	}
}

func BenchmarkTemplateEvaluate(b *testing.B) {
	env := envs.NewBuilder().Build()
	ctx := makeTemplateContext()
	template := excellent.ParseTemplate(benchmarkTemplate)

	b.ResetTimer()
	for range b.N {
		template.Evaluate(env, ctx, nil)
	}
}
