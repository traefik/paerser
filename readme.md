# Paerser

[![Package documentation](https://img.shields.io/badge/go.dev-docs-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/traefik/paerser)
[![Build Status](https://github.com/traefik/paerser/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/paerser/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/traefik/paerser)](https://goreportcard.com/report/github.com/traefik/paerser)

## Features

Loads configuration from many sources:

- CLI flags
- Configuration files in YAML, TOML, JSON format
- Environment variables

It also provides a simple CLI commands handling system.

## Examples

### Configuration

#### CLI Flags

```go
package flag_test

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/traefik/paerser/flag"
)

type ConfigExample struct {
	Foo   string
	Bar   Bar
	Great bool
}

type Bar struct {
	Sub  *Sub
	List []string
}

type Sub struct {
	Name  string
	Value int
}

func ExampleDecode() {
	args := []string{
		"--foo=aaa",
		"--great=true",
		"--bar.list=AAA,BBB",
		"--bar.sub.name=bbb",
		"--bar.sub.value=6",
	}

	config := ConfigExample{}

	err := flag.Decode(args, &config)
	if err != nil {
		log.Fatal(err)
	}

	spew.Config = spew.ConfigState{
		Indent:                  "\t",
		DisablePointerAddresses: true,
	}
	spew.Dump(config)

	// Output:
	// (flag_test.ConfigExample) {
	// 	Foo: (string) (len=3) "aaa",
	// 	Bar: (flag_test.Bar) {
	// 		Sub: (*flag_test.Sub)({
	// 			Name: (string) (len=3) "bbb",
	// 			Value: (int) 6
	// 		}),
	// 		List: ([]string) (len=2 cap=2) {
	// 			(string) (len=3) "AAA",
	// 			(string) (len=3) "BBB"
	// 		}
	// 	},
	// 	Great: (bool) true
	// }
}
```

#### File

```go
package file_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/traefik/paerser/file"
)

type ConfigExample struct {
	Foo   string
	Bar   Bar
	Great bool
}

type Bar struct {
	Sub  *Sub
	List []string
}

type Sub struct {
	Name  string
	Value int
}

func ExampleDecode() {
	tempFile, err := ioutil.TempFile("", "paeser-*.yml")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tempFile.Name())

	data := `
foo: aaa
bar:
  sub:
    name: bbb
    value: 6
  list:
  - AAA
  - BBB
great: true
`

	_, err = fmt.Fprint(tempFile, data)
	if err != nil {
		log.Fatal(err)
	}

	// Read configuration file

	filePath := tempFile.Name()

	config := ConfigExample{}

	err = file.Decode(filePath, &config)
	if err != nil {
		log.Fatal(err)
	}

	spew.Config = spew.ConfigState{
		Indent:                  "\t",
		DisablePointerAddresses: true,
	}
	spew.Dump(config)

	// Output:
	// (file_test.ConfigExample) {
	// 	Foo: (string) (len=3) "aaa",
	// 	Bar: (file_test.Bar) {
	// 		Sub: (*file_test.Sub)({
	// 			Name: (string) (len=3) "bbb",
	// 			Value: (int) 6
	// 		}),
	// 		List: ([]string) (len=2 cap=2) {
	// 			(string) (len=3) "AAA",
	// 			(string) (len=3) "BBB"
	// 		}
	// 	},
	// 	Great: (bool) true
	// }
}

```

#### Environment Variables

```go
package env_test

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/traefik/paerser/env"
)

type ConfigExample struct {
	Foo   string
	Bar   Bar
	Great bool
}

type Bar struct {
	Sub  *Sub
	List []string
}

type Sub struct {
	Name  string
	Value int
}

func ExampleDecode() {
	_ = os.Setenv("MYAPP_FOO", "aaa")
	_ = os.Setenv("MYAPP_GREAT", "true")
	_ = os.Setenv("MYAPP_BAR_LIST", "AAA,BBB")
	_ = os.Setenv("MYAPP_BAR_SUB_NAME", "bbb")
	_ = os.Setenv("MYAPP_BAR_SUB_VALUE", "6")

	config := ConfigExample{}

	err := env.Decode(os.Environ(), "MYAPP_", &config)
	if err != nil {
		log.Fatal(err)
	}

	spew.Config = spew.ConfigState{
		Indent:                  "\t",
		DisablePointerAddresses: true,
	}
	spew.Dump(config)

	// Output:
	// (env_test.ConfigExample) {
	// 	Foo: (string) (len=3) "aaa",
	// 	Bar: (env_test.Bar) {
	// 		Sub: (*env_test.Sub)({
	// 			Name: (string) (len=3) "bbb",
	// 			Value: (int) 6
	// 		}),
	// 		List: ([]string) (len=2 cap=2) {
	// 			(string) (len=3) "AAA",
	// 			(string) (len=3) "BBB"
	// 		}
	// 	},
	// 	Great: (bool) true
	// }
}
```

### CLI Commands

```go
TODO
```
