package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	if err := (&cli.App{
		Name:                   "SAPUI5 Helper",
		Usage:                  "s5 [command]",
		Version:                "1.0.0-Beta",
		Description:            "Generate some boilerplate",
		Commands:               nil,
		Flags:                  nil,
		EnableBashCompletion:   true,
		Action:                 nil,
		Authors: []*cli.Author{
			{
				Name:  "darmiel",
				Email: "hi@d2a.io",
			},
		},
		UseShortOptionHandling: true,
	}).Run(os.Args); err != nil {
		panic(err)
	}
}