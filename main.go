package main

import (
	"fmt"
	"os"
	"os/exec"
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
	if checkCubectlCommand(args[0]) {
		out, err := exec.Command("kubectl", args...).Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		if len(out) > 0 {
			fmt.Print(string(out))
		}
		return 0
	}

	var (
		candidates []string
		prefixes   = []string{"kubectl-", "kube-", "kube"}
	)
	for _, prefix := range prefixes {
		command := prefix + args[0]
		candidates = append(candidates, command)
		if _, err := exec.LookPath(command); err != nil {
			continue
		}
		out, err := exec.Command(command, args[1:]...).Output()
		if err != nil {
			panic(err)
		}
		if len(out) > 0 {
			fmt.Print(string(out))
		}
		return 1
	}

	fmt.Fprintf(os.Stderr, "%v: not found\n", candidates)
	return 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		os.Exit(1)
	}
	os.Exit(run(os.Args[1:]))
}
