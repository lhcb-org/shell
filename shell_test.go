package shell_test

import (
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

// EOF
