package invalid

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRange(t *testing.T) {
	testRange1(t)
	testRangeExpend(t)
	testRangeCross(t)
	testSingleLineRange1(t)
	testSingleLineRange2(t)
}

var testNode *yaml.Node

func init() {
	testNode = &yaml.Node{
		Kind:   yaml.ScalarNode,
		Style:  yaml.LiteralStyle,
		Tag:    "!!str",
		Value:  "who are you",
		Line:   2,
		Column: 3,
	}
}

func testRange1(t *testing.T) {
	l1, err := NewLineByYAMLNode(testNode)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, l1.Line)
	assert.EqualValues(t, 3, l1.ColumnStart)
	assert.EqualValues(t, 14, l1.ColumnEnd)

	l2 := Line{
		Line:        5,
		ColumnStart: 10,
		ColumnEnd:   15,
	}

	r1 := NewRange(l1, &l2)
	assert.EqualValues(t, 2, r1.Start.Line)
	assert.EqualValues(t, 3, r1.Start.ColumnStart)
	assert.EqualValues(t, 14, r1.Start.ColumnEnd)
	assert.EqualValues(t, 5, r1.End.Line)
	assert.EqualValues(t, 10, r1.End.ColumnStart)
	assert.EqualValues(t, 15, r1.End.ColumnEnd)
}

func testRangeExpend(t *testing.T) {
	l1, err := NewLineByYAMLNode(testNode)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, l1.Line)
	assert.EqualValues(t, 3, l1.ColumnStart)
	assert.EqualValues(t, 14, l1.ColumnEnd)

	l2 := Line{
		Line:        5,
		ColumnStart: 10,
		ColumnEnd:   15,
	}

	r1 := NewRange(l1, &l2)
	assert.EqualValues(t, 2, r1.Start.Line)
	assert.EqualValues(t, 3, r1.Start.ColumnStart)
	assert.EqualValues(t, 14, r1.Start.ColumnEnd)
	assert.EqualValues(t, 5, r1.End.Line)
	assert.EqualValues(t, 10, r1.End.ColumnStart)
	assert.EqualValues(t, 15, r1.End.ColumnEnd)

	r2 := NewRange(&Line{
		Line:        1,
		ColumnStart: 1,
		ColumnEnd:   10,
	}, &Line{
		Line:        10,
		ColumnStart: 50,
		ColumnEnd:   100,
	})

	r3 := r1.expend(&r2)
	assert.EqualValues(t, 1, r3.Start.Line)
	assert.EqualValues(t, 1, r3.Start.ColumnStart)
	assert.EqualValues(t, 10, r3.Start.ColumnEnd)
	assert.EqualValues(t, 10, r3.End.Line)
	assert.EqualValues(t, 50, r3.End.ColumnStart)
	assert.EqualValues(t, 100, r3.End.ColumnEnd)

}

func testRangeCross(t *testing.T) {
	l1, err := NewLineByYAMLNode(testNode)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, l1.Line)
	assert.EqualValues(t, 3, l1.ColumnStart)
	assert.EqualValues(t, 14, l1.ColumnEnd)

	l2 := Line{
		Line:        50,
		ColumnStart: 100,
		ColumnEnd:   150,
	}

	r1 := NewRange(l1, &l2)

	r2 := NewRange(&Line{
		Line:        5,
		ColumnStart: 1,
		ColumnEnd:   100,
	}, &Line{
		Line:        100,
		ColumnStart: 50,
		ColumnEnd:   100,
	})

	r3 := r1.expend(&r2)
	assert.EqualValues(t, 2, r3.Start.Line)
	assert.EqualValues(t, 3, r3.Start.ColumnStart)
	assert.EqualValues(t, 14, r3.Start.ColumnEnd)
	assert.EqualValues(t, 100, r3.End.Line)
	assert.EqualValues(t, 50, r3.End.ColumnStart)
	assert.EqualValues(t, 100, r3.End.ColumnEnd)
}

func testSingleLineRange1(t *testing.T) {
	l1, err := NewLineByYAMLNode(testNode)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, l1.Line)
	assert.EqualValues(t, 3, l1.ColumnStart)
	assert.EqualValues(t, 14, l1.ColumnEnd)

	l2 := Line{
		Line:        2,
		ColumnStart: 5,
		ColumnEnd:   150,
	}

	r1 := NewRange(l1, &l2)

	r2 := NewRange(&Line{
		Line:        5,
		ColumnStart: 1,
		ColumnEnd:   100,
	}, &Line{
		Line:        100,
		ColumnStart: 50,
		ColumnEnd:   100,
	})

	r3 := r1.expend(&r2)
	assert.EqualValues(t, 2, r3.Start.Line)
	assert.EqualValues(t, 3, r3.Start.ColumnStart)
	assert.EqualValues(t, 150, r3.Start.ColumnEnd)
	assert.EqualValues(t, 100, r3.End.Line)
	assert.EqualValues(t, 50, r3.End.ColumnStart)
	assert.EqualValues(t, 100, r3.End.ColumnEnd)
}

func testSingleLineRange2(t *testing.T) {
	l1, err := NewLineByYAMLNode(testNode)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, l1.Line)
	assert.EqualValues(t, 3, l1.ColumnStart)
	assert.EqualValues(t, 14, l1.ColumnEnd)

	l2 := Line{
		Line:        2,
		ColumnStart: 5,
		ColumnEnd:   150,
	}

	r1 := NewRange(l1, &l2)

	r2 := NewRange(&Line{
		Line:        2,
		ColumnStart: 1,
		ColumnEnd:   100,
	}, &Line{
		Line:        2,
		ColumnStart: 50,
		ColumnEnd:   100,
	})

	r3 := r1.expend(&r2)
	assert.EqualValues(t, 2, r3.Start.Line)
	assert.EqualValues(t, 1, r3.Start.ColumnStart)
	assert.EqualValues(t, 150, r3.Start.ColumnEnd)
	assert.EqualValues(t, 2, r3.End.Line)
	assert.EqualValues(t, 1, r3.End.ColumnStart)
	assert.EqualValues(t, 150, r3.End.ColumnEnd)
}
