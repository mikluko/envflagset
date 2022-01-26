package main

import (
	"flag"
	"fmt"

	"github.com/mikluko/envflagset"
)

func main() {
	addr := flag.String("service-addr", "localhost", "address")
	port := flag.Int("service-port", 8080, "port")
	debug := flag.Bool("debug", false, "debug")
	debugShort := flag.Bool("d", false, "shorthand for debug")

	envflagset.SetPrefix("ENVFS_")
	err := envflagset.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println("service-addr:", *addr)
	fmt.Println("service-port:", *port)
	fmt.Println("debug:", *debug || *debugShort)
}
