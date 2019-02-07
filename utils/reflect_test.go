package utils_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

type testSubType struct {
	text string `subtag:"hi"`
}

type testReflectType struct {
	text        string `tag1:"weee"`
	number      int
	textSlice   []string
	sub         testSubType `tag2:""`
	subPtr      *testSubType
	subSlice    []testSubType
	subPtrSlice []*testSubType
	numberMap   map[string]int
	subMap      map[string]testSubType
}

func TestVisitFields(t *testing.T) {
	s := &testReflectType{
		text:      "Hello",
		number:    123,
		textSlice: []string{"a", "b", "c"},
		sub:       testSubType{text: "abc"},
		subPtr:    &testSubType{text: "def"},
		subSlice: []testSubType{
			{text: "ghi"},
			{text: "jkl"},
		},
		subPtrSlice: []*testSubType{
			{text: "mno"},
			{text: "pqr"},
		},
		numberMap: map[string]int{"x": 24, "y": 25},
		subMap: map[string]testSubType{
			"eng": {text: "stu"},
			"fra": {text: "vwx"},
		},
	}

	visits := make([]string, 0)
	utils.VisitFields(s, func(v reflect.Value, tag reflect.StructTag) {
		fmt.Printf("type=%s tag=%s\n", v.Type(), tag)
		visits = append(visits, fmt.Sprintf("type=%s tag=%s", v.Type(), tag))
	})

	assert.Equal(t, []string{
		`type=string tag=tag1:"weee"`,
		`type=int tag=`,
		`type=[]string tag=`,
		`type=utils_test.testSubType tag=tag2:""`,
		`type=string tag=subtag:"hi"`,
		`type=*utils_test.testSubType tag=`,
		`type=string tag=subtag:"hi"`,
		`type=[]utils_test.testSubType tag=`,
		`type=string tag=subtag:"hi"`,
		`type=string tag=subtag:"hi"`,
		`type=[]*utils_test.testSubType tag=`,
		`type=string tag=subtag:"hi"`,
		`type=string tag=subtag:"hi"`,
		`type=map[string]int tag=`,
		`type=map[string]utils_test.testSubType tag=`,
		`type=string tag=subtag:"hi"`,
		`type=string tag=subtag:"hi"`,
	}, visits)
}
