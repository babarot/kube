package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	showPrefix = flag.Bool("p", false, "show prefixes")
)

func inSlice(e string, slice []string) bool {
	for _, s := range slice {
		if s == e {
			return true
		}
	}
	return false
}

func checkCubectlCommand(command string) bool {
	return inSlice(command, []string{
		"annotate", "autoscale", "convert", "describe", "expose",
		"patch", "rollout", "top", "api-versions", "certificate",
		"cordon", "drain", "get", "port-forward", "run", "proxy",
		"apply", "cluster-info", "cp", "edit", "label", "config",
		"scale", "version", "attach", "completion", "create",
		"exec", "logs", "replace", "set", "auth", "uncordon",
		"delete", "explain", "options", "rolling-update", "taint",
	})
}

func run(args []string) int {
	var (
		candidates []string
		prefixes   = []string{"kubectl-", "kube-", "kube"}
	)

	if *showPrefix {
		fmt.Println(strings.Join(prefixes, " "))
		return 0
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		return 1
	}

	if checkCubectlCommand(args[0]) {
		err := runComamnd("kubectl", args...)
		if err != nil {
			return 1
		}
		return 0
	}

	for _, prefix := range prefixes {
		command := prefix + args[0]
		candidates = append(candidates, command)
		if _, err := exec.LookPath(command); err != nil {
			continue
		}
		err := runComamnd(command, args[1:]...)
		if err != nil {
			return 1
		}
		return 0
	}

	fmt.Fprintf(os.Stderr, "%v: not found\n", candidates)
	return 1
}

func runComamnd(command string, args ...string) error {
	if command == "" {
		return errors.New("command not found")
	}
	if runtime.GOOS == "windows" {
		return errors.New("not support platform")
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	flag.Parse()
	os.Exit(run(flag.Args()))
}
