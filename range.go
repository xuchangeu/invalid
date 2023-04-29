package invalid

import (
	"errors"
	"math"
	"sort"
)

// FieldRange struct represent the specific location of the character of beginning and ending in ruleMap for both field key and field value.
// It describes a single line of range.For multiline literal value, it contains a list of FieldRange object.
type FieldRange struct {
	Line        int
	ColumnStart int
	ColumnEnd   int
}

// merge determine the position according to multiple FieldRange objects.
// func merge ranges in same line with `ColumnStart` and `ColumnEnd` in proper calculation.
//
//	eg,. [{FieldRange{Line:1,ColumnStart 10, ColumnEnd : 20}}, {FieldRange{Line:1,ColumnStart 21, ColumnEnd : 30}}] merged into
//	[{FieldRange{Line:1,ColumnStart 10, ColumnEnd : 30}}}]
//
// merge sort result by line number
func merge(origin []*FieldRange) ([]*FieldRange, error) {
	m := make(map[int]*FieldRange, 0)
	r := make([]*FieldRange, 0)
	for _, v := range origin {
		ran, find := m[v.Line]
		if find {
			result, err := appendColumn(ran, v)
			if err != nil {
				return nil, err
			} else {
				m[v.Line] = result
			}
		} else {
			m[v.Line] = v
		}
	}
	for _, v := range m {
		r = append(r, v)
	}
	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Line < r[j].Line
	})
	return r, nil
}

// be aware that a&b should in same line
func appendColumn(a, b *FieldRange) (*FieldRange, error) {
	if a.Line != b.Line {
		return nil, errors.New("line of range for merging must be in same line")
	}
	min := math.Min(float64(a.ColumnStart), float64(b.ColumnEnd))
	max := math.Max(float64(a.ColumnEnd), float64(b.ColumnEnd))
	f := &FieldRange{
		Line:        a.Line,
		ColumnStart: int(min),
		ColumnEnd:   int(max),
	}
	return f, nil
}
