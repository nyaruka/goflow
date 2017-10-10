package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakify(t *testing.T) {
	var snakeTests = []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello_world"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello_world"},
		{"hi😀😃😄😁there", "hi_there"},
		{"昨夜のコ", "昨夜のコ"},
		{"this@isn't@email", "this_isn_t_email"},
	}

	for _, test := range snakeTests {
		value := Snakify(test.input)

		if value != test.expected {
			t.Errorf("Expected: '%s' Got: '%s' for input: '%s'", test.expected, value, test.input)
		}
	}
}

type testResolvableChild struct {
	foo int
}

type testResolvableParent struct {
	child *testResolvableChild
	bar   string
}

func (p *testResolvableParent) Resolve(key string) interface{} {
	switch key {
	case "child":
		return p.child
	case "bar":
		return p.bar
	}
	return fmt.Errorf("no such field on parent")
}

func (p *testResolvableParent) Default() interface{} {
	return p.child
}

func (p *testResolvableParent) String() string {
	return "parentstring"
}

func (c *testResolvableChild) Resolve(key string) interface{} {
	switch key {
	case "foo":
		return c.foo
	}
	return fmt.Errorf("no such field on child")
}

func (c *testResolvableChild) Default() interface{} {
	return c
}

func (c *testResolvableChild) String() string {
	return "childstring"
}

func TestResolving(t *testing.T) {
	env := NewDefaultEnvironment()
	child := &testResolvableChild{foo: 1234}
	parent := &testResolvableParent{child: child, bar: "hello"}

	assert.Equal(t, "hello", ResolveVariable(env, parent, "bar"))
	assert.Equal(t, child, ResolveVariable(env, parent, "child"))
	assert.Equal(t, 1234, ResolveVariable(env, parent, "child.foo"))

	parentAsString, err := ToString(env, parent)
	assert.NoError(t, err)
	assert.Equal(t, "childstring", parentAsString) // because it resolves to its child by default

	childAsString, err := ToString(env, child)
	assert.NoError(t, err)
	assert.Equal(t, "childstring", childAsString)
}
