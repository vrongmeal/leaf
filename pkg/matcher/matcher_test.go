package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatcher(t *testing.T) {
	patterns := []string{
		"a",
		"b/",
		"!b/c",
		"/c/**/*",
		"d/**/*.e",
	}

	matcher, err := NewMatcher(patterns)
	require.NoError(t, err)

	assert.True(t, matcher.Match("a", false))
	assert.True(t, matcher.Match("b", true))
	assert.True(t, matcher.Match("c/d/e", false))
	assert.True(t, matcher.Match("d/g/h/i.e", false))

	assert.False(t, matcher.Match("e", false))
	assert.False(t, matcher.Match("b/", false))
	assert.False(t, matcher.Match("b/c", true))
	assert.False(t, matcher.Match("d/g/h/i.j", false))
}
