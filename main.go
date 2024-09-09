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

func (a *godnf) Install(packageName string, opt *Options) error {
	return a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Install: %v", errPackageNameNotSpecified)
		}
		return []string{"install", packageName}, nil
	}, opt)
}

func (a *godnf) Update(packageName string, opt *Options) error {
	return a.runner(func() ([]string, error) {
		return []string{}, nil
	}, opt)
}

func (a *godnf) Depends(packageName string, opt *Options) error {
	return a.runner(func() ([]string, error) {
		return []string{}, nil
	}, opt)
}

func (a *godnf) Remove(packageName string, opt *Options) error {
	return a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Remove: %v", errPackageNameNotSpecified)
		}
		return []string{""}, nil
	}, opt)
}

func (a *godnf) Search(packageName string, opt *Options) error {
	return a.runner(func() ([]string, error) {
		if strings.TrimSpace(packageName) == "" {
			return nil, fmt.Errorf("Remove: %v", errPackageNameNotSpecified)
		}
		return []string{}, nil
	}, opt)
}

func (a *godnf) List(opt *Options) error {
	return a.runner(func() ([]string, error) {
		return []string{}, nil
	}, opt)
}

func (a *godnf) runner(guest func() ([]string, error), opt *Options) error {
	arguments, err := guest()
	if err != nil {
		return fmt.Errorf("runner: %v", err)
	}
	arguments = append(arguments, processOptions(opt)...)
	command := execCommander().Command(a.binaryPath, arguments...)
	if opt.Output != nil {
		command.Stdout = opt.Output
		command.Stderr = opt.Output
	}
	return command.Run()
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
