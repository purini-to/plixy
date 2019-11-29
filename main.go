package main

import (
	"github.com/purini-to/plixy/cmd"
	"github.com/purini-to/plixy/pkg/log"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
