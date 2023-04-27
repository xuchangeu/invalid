package invalid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtility(t *testing.T) {
	key := deepFieldWithDot([]string{"foo", "bar", "see"})
	assert.EqualValues(t, "foo.bar.see", key)
}
