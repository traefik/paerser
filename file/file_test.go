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
	//nolint:dupword // false positive
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
hoo = [true, false, false, true]
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
			"trustIP": []interface{}{"10.0.0.0/8", "172.0.0.0/8", "192.0.0.0/8"},
			"koo":     []interface{}{int64(1), int64(2), int64(3)},
			"soo":     []interface{}{"1", "a", "3"},
			"boo":     []interface{}{float64(1), 2.6, float64(3)},
			"hoo":     []interface{}{true, false, false, true},
			"buckets": []interface{}{42.01, 42.02},
		},
	}
	assert.EqualValues(t, expected, element)
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

func TestDecodeTypedSliceRecursive_YAML(t *testing.T) {
	content := `
foo:
  - foo
  - bar
bar:
  foo:
  - foo
  - bar
  baz:
    foo:
    - foo
    - bar
    boz:
      foo:
      - 42
      - 42
`

	element := &map[string]interface{}{}

	err := DecodeContent(content, ".yaml", element)
	require.NoError(t, err)

	expected := &map[string]interface{}{
		"foo": []interface{}{"foo", "bar"},
		"bar": map[string]interface{}{
			"foo": []interface{}{"foo", "bar"},
			"baz": map[string]interface{}{
				"foo": []interface{}{"foo", "bar"},
				"boz": map[string]interface{}{
					"foo": []interface{}{int64(42), int64(42)},
				},
			},
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
  hoo:
    - true
    - false
    - false
    - true
  joo:
    - name: a
      value: 1
    - name: c
      value: 2
`

	var element FooRaw
	err := DecodeContent(content, ".yaml", &element)
	require.NoError(t, err)

	expected := FooRaw{
		TestData: map[string]interface{}{
			"Headers": map[string]interface{}{"Foo": "Bar"},
			"trustIP": []interface{}{"10.0.0.0/8", "172.0.0.0/8", "192.0.0.0/8"},
			"koo":     []interface{}{int64(1), int64(2), int64(3)},
			"soo":     []interface{}{"1", "a", "3"},
			"boo":     []interface{}{float64(1), 2.6, float64(3)},
			"hoo":     []interface{}{true, false, false, true},
			"buckets": []interface{}{42.01, 42.02},
			"joo": []interface{}{
				map[string]interface{}{"name": "a", "value": "1"},
				map[string]interface{}{"name": "c", "value": "2"},
			},
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
