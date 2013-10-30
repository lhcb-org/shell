package main

import (
	"fmt"
	//"io"
	"os"
	//"time"

	"github.com/atlas-org/shell"
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

	err = sh.Setenv("TOTO", "101")
	if err != nil {
		panic(err)
	}
	err = sh.Setenv("TOTO", "1011")
	if err != nil {
		panic(err)
	}
	err = sh.Setenv("TOTO", "1012")
	if err != nil {
		panic(err)
	}
	{
		val := sh.Getenv("TOTO")
		fmt.Fprintf(os.Stdout, "TOTO=%q\n", val)
		assert("1012", val)
	}
	err = sh.Setenv("TOTO", "1011")
	if err != nil {
		panic(err)
	}
	{
		val := sh.Getenv("TATA")
		fmt.Fprintf(os.Stdout, "TATA=%q\n", val)
		assert("", val)
	}
	{
		val := sh.Getenv("TOTO")
		fmt.Fprintf(os.Stdout, "TOTO=%q\n", val)
		assert("1011", val)
	}
	err = sh.Source("./test-script.sh")
	if err != nil {
		panic(err)
	}
	{
		val := sh.Getenv("TITI")
		fmt.Fprintf(os.Stdout, "TITI=%q\n", val)
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
