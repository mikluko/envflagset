package envflagset

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// EnvFlagSet represents a envflag object that contains several settings.
type EnvFlagSet struct {
	// FlagSet defines flag.FlagSet to operate on. If not provided, flag.CommandLine will be used.
	FlagSet *flag.FlagSet

	// Prefix defines environment variable prefix to use
	Prefix string

	// MinLength defines minimal flag name length to use in mapping.
	MinLength int

	// Map defines custom flag name to environment variable mappings. Prefix or MinLength are not taken into account.
	Map map[string]string

	// UpdateUsage switches environment variable names in usage message.
	UpdateUsage bool

	// Env defines environment lookup method. If not defined, syscall environment lookup will be used which is probably
	// the best option for the majority of use cases.
	Env EnvGetter

	once sync.Once
}

var std = &EnvFlagSet{
	UpdateUsage: true,
	MinLength:   3,
}

func (ef *EnvFlagSet) init() {
	if ef.FlagSet == nil {
		ef.FlagSet = flag.CommandLine
	}
	if ef.Map == nil {
		ef.Map = make(map[string]string)
	}
	if ef.Env == nil {
		ef.Env = syscallEnvGetter{}
	}
}

// flagEnvName builds environment variable names for a flag taking Prefix and Map into account
func (ef *EnvFlagSet) flagEnvName(f *flag.Flag) string {
	e, ok := ef.Map[f.Name]
	if ok {
		return e
	}
	e = strings.ToUpper(f.Name)
	e = strings.ReplaceAll(e, "-", "_")
	e = strings.ReplaceAll(e, ".", "_")
	e = fmt.Sprintf("%s%s", ef.Prefix, e)
	return e
}

// flagEnvValue looks up environment variable value for a flag
func (ef *EnvFlagSet) flagEnvValue(f *flag.Flag) (string, bool) {
	return ef.Env.GetEnv(ef.flagEnvName(f))
}

func (ef *EnvFlagSet) updateFlagUsage(f *flag.Flag) {
	if len(f.Name) < ef.MinLength {
		return
	}
	f.Usage = fmt.Sprintf("[%s] %s", ef.flagEnvName(f), f.Usage)
}

func (ef *EnvFlagSet) updateFlagValue(f *flag.Flag) error {
	if len(f.Name) < ef.MinLength {
		return nil
	}
	v, ok := ef.flagEnvValue(f)
	if !ok {
		return nil // environment variable not set for the flag
	}
	return ef.FlagSet.Set(f.Name, v)
}

// Process updates FlagSet with values from the environment.
// NOTICE: flag.Parse() will not be called by this function.
func (ef *EnvFlagSet) Process() error {
	ef.once.Do(ef.init)
	if ef.FlagSet.Parsed() {
		return errors.New("flag set has already been parsed")
	}
	if ef.UpdateUsage {
		ef.FlagSet.VisitAll(ef.updateFlagUsage)
	}
	var me *multierror.Error
	ef.FlagSet.VisitAll(func(f *flag.Flag) {
		err := ef.updateFlagValue(f)
		if err != nil {
			me = multierror.Append(me, err)
		}
	})
	if me != nil {
		return me
	}
	return nil
}

// Parse parses flag definitions from env and the argument list.
func (ef *EnvFlagSet) Parse(arguments []string) error {
	if err := ef.Process(); err != nil {
		return err
	}
	return ef.FlagSet.Parse(arguments)
}

// Parse parses the command-line flags from env and os.Args[1:].
func Parse() error {
	return std.Parse(os.Args[1:])
}

// SetPrefix sets prefix on default EnvFlagSet instance
func SetPrefix(p string) {
	std.Prefix = p
}
