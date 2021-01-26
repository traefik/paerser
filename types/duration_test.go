package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_Set(t *testing.T) {
	testCases := []struct {
		desc     string
		value    string
		assert   require.ErrorAssertionFunc
		expected Duration
	}{
		{
			desc:   "empty",
			value:  "",
			assert: require.Error,
		},
		{
			desc:     "duration",
			value:    "2m",
			assert:   require.NoError,
			expected: Duration(2 * time.Minute),
		},
		{
			desc:     "integer",
			value:    "2",
			assert:   require.NoError,
			expected: Duration(2 * time.Second),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var d Duration

			err := d.Set(test.value)
			test.assert(t, err)

			assert.Equal(t, test.expected, d)
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		desc     string
		text     []byte
		assert   require.ErrorAssertionFunc
		expected Duration
	}{
		{
			desc:   "empty",
			text:   []byte(""),
			assert: require.Error,
		},
		{
			desc:     "duration",
			text:     []byte(`"2m"`),
			assert:   require.NoError,
			expected: Duration(2 * time.Minute),
		},
		{
			desc:     "integer",
			text:     []byte(`2`),
			assert:   require.NoError,
			expected: Duration(2 * time.Second),
		},
		{
			desc:   "bad format",
			text:   []byte(`"2"`),
			assert: require.Error,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var d Duration

			err := d.UnmarshalJSON(test.text)
			test.assert(t, err)

			assert.Equal(t, test.expected, d)
		})
	}
}
