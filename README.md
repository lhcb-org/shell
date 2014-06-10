shell
=====

``shell`` allows to programmatically talk to local shell.

## Installation

```sh
$ go get github.com/lhcb-org/shell
```

## Documentation

Documentation is available on [godoc](http://godoc.org/github.com/lhcb-org/shell)

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/lhcb-org/shell"
)

func assert(exp, val string) {
	if val != exp {
		panic(fmt.Errorf("expected %q. got %q", exp, val))
	}
}

func main() {
	sh, err := shell.New()
	if err != nil {
		panic(err)
	}
	defer sh.Delete()

	fmt.Printf(">>> starting shell...\n")

	err = sh.Setenv("FOO", "666")
	if err != nil {
		panic(err)
	}

    foo := sh.Getenv("FOO")
	fmt.Printf("FOO=%q\n", foo)
	assert("666", foo)

	err = sh.Setenv("FOO", "42")
	if err != nil {
		panic(err)
	}
    
    foo = sh.Getenv("FOO")
	fmt.Printf("FOO=%q\n", foo)
	assert("42", foo)

	not := sh.Getenv("__NOT_THERE__")
	fmt.Printf("__NOT_THERE__=%q\n", not)
	assert("", not)

    err = sh.Source("./test-script.sh")
	if err != nil {
		panic(err)
	}

    out, err := sh.Run("/bin/ls", ".")
	if err != nil {
		panic(err)
	}
	fmt.Printf("/bin/ls: %q\n", out)

	_, err = sh.Run("/bin/ls", "/dev/fail")
	if err == nil {
		panic(fmt.Errorf("expected an error listing /dev/fail"))
	}
}
```

