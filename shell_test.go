package shell_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/atlas-org/shell"
)

func TestSetenv(t *testing.T) {
	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	err = sh.Setenv("FOO", "101")
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = sh.Setenv("FOO", "202")
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = sh.Setenv("FOO", "303")
	if err != nil {
		t.Fatalf(err.Error())
	}

	foo := ""

	foo = sh.Getenv("FOO")
	if foo != "303" {
		t.Fatalf("expected %q. got %q\n", "303", foo)
	}

	bar := "+++"
	bar = sh.Getenv("BAR")
	if bar != "" {
		t.Fatalf("expected %q. got %q\n", "", bar)
	}
}

func TestChdir(t *testing.T) {
	top, err := ioutil.TempDir("", "test-shell-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(top)

	err = os.Chdir(top)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	pwd, err := sh.Getwd()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if pwd != top {
		t.Fatalf("expected %q. got %q\n", top, pwd)
	}

	foo := filepath.Join(top, "foo")
	err = os.MkdirAll(foo, 0755)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = sh.Chdir(foo)
	if err != nil {
		t.Fatalf(err.Error())
	}

	pwd, err = sh.Getwd()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if pwd != foo {
		t.Fatalf("expected %q. got %q\n", foo, pwd)
	}
}

func TestSource(t *testing.T) {

	top, err := ioutil.TempDir("", "test-shell-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(top)

	err = os.Chdir(top)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	f, err := os.Create(filepath.Join(top, "test-script.sh"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(`#!/bin/sh
export FOO="101"
echo $FOO
export BAR="1234"
echo $BAR
BAZ="4321"
export BAR
export ABC='"101"'
export DEF="'101'"
export GHI=101
# EOF
`),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f.Sync()

	err = sh.Source(f.Name())
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, tt := range []struct {
		key string
		val string
	}{
		{"FOO", "101"},
		{"BAR", "1234"},
		{"BAZ", "4321"},
		{"ABC", `"101"`},
		{"DEF", `'101'`},
		{"GHI", `101`},
	} {
		val := sh.Getenv(tt.key)
		if val != tt.val {
			t.Fatalf("expected %s=%q. got %q\n", tt.key, tt.val, val)
		}
	}

}

func TestSourceExitCode(t *testing.T) {

	top, err := ioutil.TempDir("", "test-shell-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(top)

	err = os.Chdir(top)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	f, err := os.Create(filepath.Join(top, "test-script.sh"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(`#!/bin/sh
export FOO="101"
return 101
export BAR="202"
# EOF
`),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f.Sync()

	err = sh.Source(f.Name())
	if err.Error() != `error: "101"` {
		t.Fatalf("expected non-zero status code 101. got: %q\n", err.Error())
	}

	for _, tt := range []struct {
		key string
		val string
	}{
		{"FOO", "101"},
		{"BAR", ""},
	} {
		val := sh.Getenv(tt.key)
		if val != tt.val {
			t.Fatalf("expected %s=%q. got %q\n", tt.key, tt.val, val)
		}
	}
}

func TestRun(t *testing.T) {

	top, err := ioutil.TempDir("", "test-shell-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(top)

	err = os.Chdir(top)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	f, err := os.Create(filepath.Join(top, "test-script.sh"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(`#!/bin/sh
export FOO="101"
echo $FOO
export BAR="1234"
echo $BAR
BAZ="4321"
echo $BAZ
# EOF
`),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f.Sync()

	err = f.Chmod(os.FileMode(0755))
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = f.Close()
	if err != nil {
		t.Fatalf(err.Error())
	}

	out, err := sh.Run(f.Name())
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, tt := range []struct {
		key string
		val string
	}{
		{"FOO", ""},
		{"BAR", ""},
		{"BAZ", ""},
	} {
		val := sh.Getenv(tt.key)
		if val != tt.val {
			t.Fatalf("expected %s=%q. got %q\n", tt.key, tt.val, val)
		}
	}

	out_exp := []byte("101\n1234\n4321\n")
	if string(out) != string(out_exp) {
		t.Fatalf("expected run-output=%q. got: %q\n", out_exp, out)
	}
}

func TestRunExitCode(t *testing.T) {

	top, err := ioutil.TempDir("", "test-shell-")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(top)

	err = os.Chdir(top)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sh, err := shell.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer sh.Delete()

	f, err := os.Create(filepath.Join(top, "test-script.sh"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(`#!/bin/sh
export FOO="101"
echo $FOO
export BAR="1234"
echo $BAR
BAZ="4321"
echo $BAZ
exit 101
# EOF
`),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f.Sync()

	err = f.Chmod(os.FileMode(0755))
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = f.Close()
	if err != nil {
		t.Fatalf(err.Error())
	}

	out, err := sh.Run(f.Name())
	if err.Error() != `error: "101"` {
		t.Fatalf("expected non-zero status code 101. got: %q\n", err.Error())
	}

	for _, tt := range []struct {
		key string
		val string
	}{
		{"FOO", ""},
		{"BAR", ""},
		{"BAZ", ""},
	} {
		val := sh.Getenv(tt.key)
		if val != tt.val {
			t.Fatalf("expected %s=%q. got %q\n", tt.key, tt.val, val)
		}
	}

	out_exp := []byte("101\n1234\n4321\n")
	if string(out) != string(out_exp) {
		t.Fatalf("expected run-output=%q. got: %q\n", out_exp, out)
	}

}

// EOF
