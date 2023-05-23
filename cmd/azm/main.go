package main

import (
	"log"
	"os"

	v2 "github.com/aserto-dev/azm/v2"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func main() {
	var manifest string
	flag.StringVarP(&manifest, "manifest", "m", "manifest", "manifest file path")
	flag.Parse()

	r, err := os.Open(manifest)
	if err != nil {
		log.Fatal(err)
	}

	m, err := v2.Model().Read(r)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.NewEncoder(os.Stdout).Encode(m); err != nil {
		log.Fatal(err)
	}
}
