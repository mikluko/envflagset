package envflagset

import "os"

type EnvGetter interface {
	GetEnv(string) (string, bool)
}

type syscallEnvGetter struct{}

func (syscallEnvGetter) GetEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
