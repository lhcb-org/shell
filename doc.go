// shell allows to programmatically manage a local shell (/bin/sh).
//
// A typical usage might look like so:
//
//  sh, err := shell.New()
//  err = sh.Setenv("FOO", "42")
//  val := sh.Getenv("FOO")
//  err = sh.Run("/some/script.sh", "-f", "blah.txt")
package shell
