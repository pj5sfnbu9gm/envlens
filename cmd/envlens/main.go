package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envlens/internal/loader"
	"github.com/yourorg/envlens/internal/resolver"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "envlens: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("envlens", flag.ContinueOnError)
	var targetFlags multiFlag
	fs.Var(&targetFlags, "target", "target in name=path format (repeatable)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine working directory: %w", err)
	}

	targets, err := resolver.ResolveTargets([]string(targetFlags), cwd)
	if err != nil {
		return err
	}

	if err := resolver.ValidatePaths(targets); err != nil {
		return err
	}

	for _, t := range targets {
		env, err := loader.LoadFile(t.Path)
		if err != nil {
			return fmt.Errorf("loading target %q: %w", t.Name, err)
		}
		fmt.Printf("[%s] loaded %d keys from %s\n", t.Name, len(env), t.Path)
		for k, v := range env {
			_ = v
			_ = k
		}
	}

	fmt.Println("targets resolved successfully — diff/audit coming soon")
	return nil
}

// multiFlag is a flag.Value that accumulates repeated string flags.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ", ") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}
