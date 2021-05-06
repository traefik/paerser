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

func TestDuration_MarshalJSON(t *testing.T) {
	testCases := []struct {
		desc     string
		dur      Duration
		expected []byte
	}{
		{
			desc:     "1 second",
			dur:      Duration(time.Second),
			expected: []byte(`"1s"`),
		},
		{
			desc:     "1 millisecond",
			dur:      Duration(time.Millisecond),
			expected: []byte(`"1ms"`),
		},
		{
			desc:     "1 nanosecond",
			dur:      Duration(time.Nanosecond),
			expected: []byte(`"1ns"`),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			b, err := test.dur.MarshalJSON()
			require.NoError(t, err)

			assert.Equal(t, test.expected, b)
		})
	}
}

func TestDuration_JSON_bijection(t *testing.T) {
	testCases := []struct {
		desc string
		dur  Duration
	}{
		{
			desc: "1 second",
			dur:  Duration(time.Second),
		},
		{
			desc: "1 millisecond",
			dur:  Duration(time.Millisecond),
		},
		{
			desc: "1 nanosecond",
			dur:  Duration(time.Nanosecond),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			b, err := test.dur.MarshalJSON()
			require.NoError(t, err)

			var ud Duration
			err = ud.UnmarshalJSON(b)
			require.NoError(t, err)

			assert.Equal(t, test.dur, ud)
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
