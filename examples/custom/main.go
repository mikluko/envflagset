package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikluko/envflagset"
)

func main() {
	fs := flag.NewFlagSet("example", flag.ExitOnError)

	addr := fs.String("service-addr", "localhost", "address")
	port := fs.Int("service-port", 8080, "port")
	debug := fs.Bool("debug", false, "debug")
	debugShort := fs.Bool("d", false, "shorthand for debug")

	efs := envflagset.EnvFlagSet{
		FlagSet:   fs,
		Prefix:    "ENVFS_",
		MinLength: 3,
	}
	err := efs.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	fmt.Println("service-addr:", *addr)
	fmt.Println("service-port:", *port)
	fmt.Println("debug:", *debug || *debugShort)
}
