package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixType(t *testing.T) {
	a, b := ParsePrefixType("[]byte")
	assert.Equal(t, "[]", a)
	assert.Equal(t, "byte", b)
}
