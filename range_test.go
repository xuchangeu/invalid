package invalid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange(t *testing.T) {
	r1 := &FieldRange{
		Line:        1,
		ColumnStart: 2,
		ColumnEnd:   3,
	}

	r2 := &FieldRange{
		Line:        1,
		ColumnStart: 4,
		ColumnEnd:   6,
	}

	result, err := merge([]*FieldRange{r1, r2})
	assert.Nil(t, err)
	assert.EqualValuesf(t, 1, len(result), "length of result is ：%d", len(result))
	assert.EqualValuesf(t, 1, result[0].Line, "line of result is : %d", result[0].Line)
	assert.EqualValuesf(t, 2, result[0].ColumnStart, "start column of result is : %d", result[0].ColumnStart)
	assert.EqualValuesf(t, 6, result[0].ColumnEnd, "end column of result is : %d", result[0].ColumnEnd)

	r1 = &FieldRange{
		Line:        1,
		ColumnStart: 3,
		ColumnEnd:   6,
	}

	r2 = &FieldRange{
		Line:        2,
		ColumnStart: 5,
		ColumnEnd:   8,
	}

	result, err = merge([]*FieldRange{r1, r2})
	assert.Nil(t, err)
	assert.EqualValuesf(t, 2, len(result), "length of result is ：%d", len(result))
	assert.EqualValuesf(t, 1, result[0].Line, "line of result is : %d", result[0].Line)
	assert.EqualValuesf(t, 3, result[0].ColumnStart, "start column of result is : %d", result[0].ColumnStart)
	assert.EqualValuesf(t, 6, result[0].ColumnEnd, "end column of result is : %d", result[0].ColumnEnd)

	assert.EqualValuesf(t, 2, result[1].Line, "line of result is : %d", result[1].Line)
	assert.EqualValuesf(t, 5, result[1].ColumnStart, "start column of result is : %d", result[1].ColumnStart)
	assert.EqualValuesf(t, 8, result[1].ColumnEnd, "end column of result is : %d", result[1].ColumnEnd)

	r1 = &FieldRange{
		Line:        1,
		ColumnStart: 3,
		ColumnEnd:   6,
	}

	r2 = &FieldRange{
		Line:        2,
		ColumnStart: 5,
		ColumnEnd:   8,
	}

	r3 := &FieldRange{
		Line:        1,
		ColumnStart: 8,
		ColumnEnd:   10,
	}

	r4 := &FieldRange{
		Line:        2,
		ColumnStart: 10,
		ColumnEnd:   20,
	}

	result, err = merge([]*FieldRange{r1, r2, r3, r4})
	assert.Nil(t, err)
	assert.EqualValuesf(t, 2, len(result), "length of result is ：%d", len(result))
	assert.EqualValuesf(t, 1, result[0].Line, "line of result is : %d", result[0].Line)
	assert.EqualValuesf(t, 3, result[0].ColumnStart, "start column of result is : %d", result[0].ColumnStart)
	assert.EqualValuesf(t, 10, result[0].ColumnEnd, "end column of result is : %d", result[0].ColumnEnd)

	assert.EqualValuesf(t, 2, result[1].Line, "line of result is : %d", result[1].Line)
	assert.EqualValuesf(t, 5, result[1].ColumnStart, "start column of result is : %d", result[1].ColumnStart)
	assert.EqualValuesf(t, 20, result[1].ColumnEnd, "end column of result is : %d", result[1].ColumnEnd)
}
