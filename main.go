package godnf
import (
	"fmt"
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
func (a *godnf) runner(guest func() ([]string, error), opt *Options) error {
	arguments, err := guest()
	if err != nil {
		return fmt.Errorf("runner: %v", err)
	}
	command := execCommander().Command(a.binaryPath, arguments...)
	if opt.Output != nil {
		command.Stdout = opt.Output
		command.Stderr = opt.Output
	}
	return command.Run()
}
