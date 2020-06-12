package main

import (
	"log"
	"os/exec"
	"strings"
)

var capnpVersion string

func main() {
	loadCapnpVersion()
	generateImports()
	generateSchemaValue()
}

func loadCapnpVersion() {
	if capnpVersion != "" {
		return
	}

	cmd := exec.Command("git", "submodule", "status", "capnproto")
	sb := strings.Builder{}
	cmd.Stdout = &sb
	err := cmd.Run()
	dieOnErr(err)
	capnpVersion = strings.TrimSpace(sb.String())
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}