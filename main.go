package godnf

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/CREDOProject/sharedutils/shell"
)

var execCommander = shell.New

// godnf represents the DNF client.
type godnf struct {
	binaryPath string
}

// Options represents the configuration options for the running command.
type Options struct {
	Verbose      bool
	DryRun       bool
	Output       io.Writer
	NotAssumeYes bool
}

// Returns a new godnf value, which represents an initialized DNF client.
func New(binaryPath string) *godnf {
	return &godnf{binaryPath}
}

var (
	errPackageNameNotSpecified = errors.New("packageName was not specified.")
)

// Install a dnf package from its packageName.
func (a *godnf) Install(packageName string, opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Install: %v", errPackageNameNotSpecified)
		}
		return []string{"install", packageName}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// Update a packages from is packageName. If packageName is empty, updates all
// the packages in the system.
func (a *godnf) Update(packageName string, opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return []string{"update"}, nil
		}
		return []string{"update", packageName}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// Obtains a list of dependencies from a packageName.
func (a *godnf) Depends(packageName string, opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Depends: %v", errPackageNameNotSpecified)
		}
		return []string{"repoquery", "--deplist", packageName}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// Remove a package from its packageName.
func (a *godnf) Remove(packageName string, opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Remove: %v", errPackageNameNotSpecified)
		}
		return []string{"remove", packageName}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// Search a package from its packageName.
func (a *godnf) Search(packageName string, opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Remove: %v", errPackageNameNotSpecified)
		}
		return []string{}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// List all installed packages.
func (a *godnf) List(opt *Options) error {
	_, err := a.runner(func() ([]string, error) {
		return []string{"list", "installed"}, nil
	}, func(string) ([]Package, error) {
		return nil, nil
	}, opt)
	return err
}

// runner runs a guest command with opt *Options.
func (a *godnf) runner(
	guest func() ([]string, error),
	parser func(string) ([]Package, error),
	opt *Options,
) ([]Package, error) {
	arguments, err := guest()
	if err != nil {
		return nil, fmt.Errorf("runner: %v", err)
	}
	arguments = append(arguments, processOptions(opt)...)
	command := execCommander().Command(a.binaryPath, arguments...)
	if opt.Output != nil {
		command.Stdout = opt.Output
		command.Stderr = opt.Output
	}
	return nil, command.Run()
}

// processOptions returns a slice of command-line flags to be passed to dnf.
func processOptions(opt *Options) []string {
	args := []string{}
	if opt.DryRun {
		args = append(args, "--setopt", "tsflags=test")
	}
	if opt.Verbose {
		args = append(args, "--verbose")
	}
	if !opt.NotAssumeYes {
		args = append(args, "--assumeyes")
	}
	return args
}

// Package represents a DNF package.
type Package struct {
	Name    string
	Version string
	Path    string
}
