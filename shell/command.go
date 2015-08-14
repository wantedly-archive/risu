package shell

import (
	"os/exec"
	"os/user"
	"strings"
)

// Command executes command and returns string output and error
func Command(name string, args ...string) (string, error) {
	return CommandInDir("", name, args...)
}

// CommandInDir executes command in specified directory, and returns string output and error
func CommandInDir(dir, name string, args ...string) (string, error) {
	if dir != "" {
		usr, _ := user.Current()
		dir = strings.Replace(dir, "~", usr.HomeDir, 1)
	}
	//fmt.Println("> " + name + " " + strings.Join(args, " ") + "\t[" + dir + "]")
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
