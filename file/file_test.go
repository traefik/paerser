package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode_TOML(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "traefik-config-*.toml")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = f.Close()
	})

	_, err = f.Write([]byte(`
foo = "bar"
fii = "bir"
[yi]
`))
	require.NoError(t, err)

	element := &Yo{
		Fuu: "test",
	}

	err = Decode(f.Name(), element)
	require.NoError(t, err)

	expected := &Yo{
		Foo: "bar",
		Fii: "bir",
		Fuu: "test",
		Yi: &Yi{
			Foo: "foo",
			Fii: "fii",
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_TOML(t *testing.T) {
	content := `
foo = "bar"
fii = "bir"
[yi]
`

	element := &Yo{
		Fuu: "test",
	}

	err := DecodeContent(content, ".toml", element)
	require.NoError(t, err)

	expected := &Yo{
		Foo: "bar",
		Fii: "bir",
		Fuu: "test",
		Yi: &Yi{
			Foo: "foo",
			Fii: "fii",
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_TOML_rawSlice(t *testing.T) {
	content := `
[testData]
trustIP = [
  "10.0.0.0/8",
  "172.0.0.0/8",
  "192.0.0.0/8"
]
koo = [1, 2, 3]
soo = [1, "a", 3]
boo = [1, 2.6, 3]
buckets = [42.01, 42.02]

  [testData.Headers]
  Foo = "Bar"
`

	var element FooRaw

	err := DecodeContent(content, ".toml", &element)
	require.NoError(t, err)

	expected := FooRaw{
		TestData: map[string]interface{}{
			"Headers": map[string]interface{}{"Foo": "Bar"},
			"trustIP": []string{"10.0.0.0/8", "172.0.0.0/8", "192.0.0.0/8"},
			"koo":     []int64{1, 2, 3},
			"soo":     []string{"1", "a", "3"},
			"boo":     []float64{1, 2.6, 3},
			"buckets": []float64{42.01, 42.02},
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_TOML_rawValue(t *testing.T) {
	content := `
name = "test"
[[meta.aaa]]
	bbb = 1
`

	type Foo struct {
		Name string
		Meta map[string]interface{}
	}

	element := &Foo{}

	err := DecodeContent(content, ".toml", element)
	require.NoError(t, err)

	expected := &Foo{
		Name: "test",
		Meta: map[string]interface{}{"aaa": []interface{}{map[string]interface{}{"bbb": "1"}}},
	}
	assert.Equal(t, expected, element)
}

func TestDecode_YAML(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "traefik-config-*.yaml")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = f.Close()
	})

	_, err = f.Write([]byte(`
foo: bar
fii: bir
yi: {}
`))
	require.NoError(t, err)

	element := &Yo{
		Fuu: "test",
	}

	err = Decode(f.Name(), element)
	require.NoError(t, err)

	expected := &Yo{
		Foo: "bar",
		Fii: "bir",
		Fuu: "test",
		Yi: &Yi{
			Foo: "foo",
			Fii: "fii",
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_YAML(t *testing.T) {
	content := `
foo: bar
fii: bir
yi: {}
`

	element := &Yo{
		Fuu: "test",
	}

	err := DecodeContent(content, ".yaml", element)
	require.NoError(t, err)

	expected := &Yo{
		Foo: "bar",
		Fii: "bir",
		Fuu: "test",
		Yi: &Yi{
			Foo: "foo",
			Fii: "fii",
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_YAML_rawSlice(t *testing.T) {
	content := `
testData:
  Headers:
    Foo: Bar
  trustIP:
    - 10.0.0.0/8
    - 172.0.0.0/8
    - 192.0.0.0/8
  koo:
    - 1
    - 2
    - 3
  soo:
    - 1
    - a
    - 3
  boo:
    - 1
    - 2.6
    - 3
  buckets:
    - 42.01
    - 42.02
`

	var element FooRaw
	err := DecodeContent(content, ".yaml", &element)
	require.NoError(t, err)

	expected := FooRaw{
		TestData: map[string]interface{}{
			"Headers": map[string]interface{}{"Foo": "Bar"},
			"trustIP": []string{"10.0.0.0/8", "172.0.0.0/8", "192.0.0.0/8"},
			"koo":     []int{1, 2, 3},
			"soo":     []string{"1", "a", "3"},
			"boo":     []float64{1, 2.6, 3},
			"buckets": []float64{42.01, 42.02},
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_YAML_rawValue(t *testing.T) {
	type Foo struct {
		Name string
		Meta map[string]interface{}
	}

	testCases := []struct {
		desc     string
		content  string
		expected interface{}
	}{
		{
			desc: "simple",
			content: `
name: test
meta:
  aaa:
  - bbb: 1
`,
			expected: &Foo{
				Name: "test",
				Meta: map[string]interface{}{"aaa": []interface{}{map[string]interface{}{"bbb": "1"}}},
			},
		},
		{
			desc: "null",
			content: `
name: test
meta:
  aaa:
  - bbb: 1
  bbb: {"foo":"bar"}
  null: {"toto":"tata"}
`,
			expected: &Foo{
				Name: "test",
				Meta: map[string]interface{}{
					"aaa": []interface{}{map[string]interface{}{"bbb": "1"}},
					"bbb": map[string]interface{}{"foo": "bar"},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			element := &Foo{}

			err := DecodeContent(test.content, ".yaml", element)
			require.NoError(t, err)

			assert.Equal(t, test.expected, element)
		})
	}
}

func TestDecodeContent_JSON(t *testing.T) {
	content := `
{
  "foo": "bar",
  "fii": "bir",
  "yi": {}
}
`

	element := &Yo{
		Fuu: "test",
	}

	err := DecodeContent(content, ".json", element)
	require.NoError(t, err)

	expected := &Yo{
		Foo: "bar",
		Fii: "bir",
		Fuu: "test",
		Yi: &Yi{
			Foo: "foo",
			Fii: "fii",
		},
	}
	assert.Equal(t, expected, element)
}

func TestDecodeContent_JSON_rawValue(t *testing.T) {
	content := `
{
  "name": "test",
  "meta": {
    "aaa": [
      {
        "bbb": 1
      }
    ]
  }
}
`

	type Foo struct {
		Name string
		Meta map[string]interface{}
	}

	element := &Foo{}

	err := DecodeContent(content, ".json", element)
	require.NoError(t, err)

	expected := &Foo{
		Name: "test",
		Meta: map[string]interface{}{"aaa": []interface{}{map[string]interface{}{"bbb": "1"}}},
	}
	assert.Equal(t, expected, element)
}
