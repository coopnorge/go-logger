package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	//mage:import
	_ "github.com/coopnorge/mage/targets/golib"
	"github.com/magefile/mage/sh"
)

// moduleSet is the multimod module set in versions.yaml that is released in
// lockstep.
const moduleSet = "go-logger"

// Consistency syncs intra-repository replace directives, tidies every module,
// verifies versions.yaml and fails if that changed any go.mod, go.sum or
// versions.yaml.
func Consistency(_ context.Context) error {
	err := sh.RunV("go", "tool", "crosslink", "--root", ".",
		"--exclude", "github.com/coopnorge/go-logger-e2e-tests")
	if err != nil {
		return err
	}
	if err := sh.RunV("go", "tool", "multimod", "verify"); err != nil {
		return err
	}
	dirs, err := moduleDirs()
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err := sh.RunWithV(map[string]string{"GOWORK": "off"}, "go", "-C", dir, "mod", "tidy")
		if err != nil {
			return err
		}
	}
	return sh.RunV("git", "diff", "--exit-code", "--", "*go.mod", "*go.sum", "versions.yaml")
}

// Prerelease updates the module versions in all go.mod files to the version in
// versions.yaml and commits them to a new branch. Run it on a clean main.
func Prerelease(_ context.Context) error {
	return sh.RunV("go", "tool", "multimod", "prerelease", "--module-set-names", moduleSet)
}

// Tag creates the git tags for all modules at the given commit, for example
// `mage tag $(git rev-parse HEAD)`, and prints them. Push them afterwards.
func Tag(_ context.Context, commit string) error {
	if commit == "" {
		return fmt.Errorf("commit hash is required")
	}
	return sh.RunV("go", "tool", "multimod", "tag",
		"--module-set-name", moduleSet, "--commit-hash", commit, "--print-tags")
}

// moduleDirs returns the directory of every go.mod in the repository.
func moduleDirs() ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == "go.mod" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return dirs, err
}
