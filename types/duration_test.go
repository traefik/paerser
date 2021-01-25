package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration_Set(t *testing.T) {
	testCases := []struct {
		desc        string
		s           string
		expDuration Duration
		expErr      bool
	}{
		{
			desc:   "empty",
			s:      "",
			expErr: true,
		},
		{
			desc:        "duration",
			s:           "2m",
			expDuration: Duration(2 * time.Minute),
		},
		{
			desc:        "integer",
			s:           "2",
			expDuration: Duration(2 * time.Second),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var d Duration

			err := d.Set(test.s)
			if test.expErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expDuration, d)
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		desc        string
		text        []byte
		expDuration Duration
		expErr      bool
	}{
		{
			desc:   "empty",
			text:   []byte(""),
			expErr: true,
		},
		{
			desc:        "duration",
			text:        []byte(`"2m"`),
			expDuration: Duration(2 * time.Minute),
		},
		{
			desc:        "integer",
			text:        []byte(`2`),
			expDuration: Duration(2 * time.Second),
		},

		{
			desc:   "bad format",
			text:   []byte(`"2"`),
			expErr: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var d Duration

			err := d.UnmarshalJSON(test.text)
			if test.expErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expDuration, d)
		})
	}
}
