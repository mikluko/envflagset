package envflagset_test

import (
	"bytes"
	"flag"
	"testing"

	"github.com/mikluko/envflagset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	fs := flag.NewFlagSet("testing", flag.PanicOnError)

	addr := fs.String("service-addr", "localhost", "address")
	port := fs.Int("service-port", 8080, "port")
	debug := fs.Bool("debug", false, "debug")
	verbose := fs.Bool("v", false, "alias for debug")

	efs := envflagset.EnvFlagSet{
		FlagSet:   fs,
		MinLength: 3,
		Prefix:    "TEST_",
		Env: environMock{
			"TEST_SERVICE_ADDR": "example.com",
			"TEST_SERVICE_PORT": "9090",
			"TEST_DEBUG":        "true",
			"TEST_V":            "true", // should be ignored because of MinLength
		},
	}

	err := efs.Parse(nil)
	require.NoError(t, err)

	assert.Equal(t, "example.com", *addr)
	assert.Equal(t, 9090, *port)
	assert.True(t, *debug)
	assert.False(t, *verbose)
}

func TestUsage(t *testing.T) {
	var buf bytes.Buffer

	fs := flag.NewFlagSet("testing", flag.PanicOnError)
	fs.SetOutput(&buf)

	_ = fs.String("service-addr", "localhost", "address")
	_ = fs.Int("service-port", 8080, "port")
	_ = fs.Bool("debug", false, "debug")
	_ = fs.Bool("v", false, "alias for debug")

	efs := envflagset.EnvFlagSet{
		FlagSet:     fs,
		UpdateUsage: true,
		Prefix:      "TEST_",
		MinLength:   3,
	}

	err := efs.Process()
	require.NoError(t, err)

	fs.Usage()

	require.Contains(t, buf.String(), "[TEST_SERVICE_ADDR]")
	require.Contains(t, buf.String(), "[TEST_SERVICE_PORT]")
	require.Contains(t, buf.String(), "[TEST_DEBUG]")
	require.NotContains(t, buf.String(), "[TEST_V]")
}

type environMock map[string]string

func (e environMock) GetEnv(name string) (string, bool) {
	k, v := e[name]
	return k, v
}
