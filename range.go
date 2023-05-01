package invalid

import (
	"errors"
	"gopkg.in/yaml.v3"
	"math"
	"strings"
)

type Line struct {
	Line        uint
	ColumnStart uint
	ColumnEnd   uint
}

func NewLineByYAMLNode(node *yaml.Node) (*Line, error) {
	if strings.Contains(node.Value, "\n") {
		return nil, errors.New("string parameter must not contains break line")
	}

	c := uint(len(node.Value))
	if node.Style == yaml.DoubleQuotedStyle || node.Style == yaml.SingleQuotedStyle {
		c += 2
	}

	return &Line{
		Line:        uint(node.Line),
		ColumnStart: uint(node.Column),
		ColumnEnd:   uint(node.Column) + c,
	}, nil
}

func NewRange(l1, l2 *Line) Range {
	if l1.Line < l2.Line {
		return Range{
			Start: l1,
			End:   l2,
		}
	} else if l1.Line == l2.Line {

		start := math.Min(float64(l1.ColumnStart), float64(l2.ColumnStart))
		end := math.Max(float64(l1.ColumnEnd), float64(l2.ColumnEnd))

		return Range{
			Start: &Line{
				Line:        l1.Line,
				ColumnStart: uint(start),
				ColumnEnd:   uint(end),
			},
			End: &Line{
				Line:        l1.Line,
				ColumnStart: uint(start),
				ColumnEnd:   uint(end),
			},
		}

	} else {
		return Range{
			Start: l2,
			End:   l1,
		}
	}
}

// Range
type Range struct {
	Start *Line
	End   *Line
}

func (r *Range) expend(r2 *Range) *Range {

	var start, end *Line
	if r.Start.Line < r2.Start.Line {
		start = r.Start
	} else if r.Start.Line == r2.Start.Line {
		minStart := math.Min(float64(r.Start.ColumnStart), float64(r2.Start.ColumnStart))
		maxEnd := math.Max(float64(r.Start.ColumnEnd), float64(r2.Start.ColumnEnd))
		start = &Line{
			Line:        r.Start.Line,
			ColumnStart: uint(minStart),
			ColumnEnd:   uint(maxEnd),
		}
	} else {
		start = r2.Start
	}

	if r.End.Line < r2.End.Line {
		end = r2.End
	} else if r.End.Line == r2.End.Line {
		minStart := math.Min(float64(r.End.ColumnStart), float64(r2.End.ColumnStart))
		maxEnd := math.Max(float64(r.End.ColumnEnd), float64(r2.End.ColumnEnd))
		end = &Line{
			Line:        r.End.Line,
			ColumnStart: uint(minStart),
			ColumnEnd:   uint(maxEnd),
		}
	} else {
		end = r.End
	}

	return &Range{
		Start: start,
		End:   end,
	}
}
