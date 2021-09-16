package main

import (
	"fmt"
	"github.com/darmiel/sapui5-helper/internal/cmds"
	"github.com/darmiel/sapui5-helper/pkg/s5"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	if err := (&cli.App{
		Name:        "SAPUI5 Helper",
		Usage:       "s5 [command]",
		Version:     "1.0.0-Beta",
		Description: "Generate some boilerplate",
		Commands: []*cli.Command{
			cmds.RmdNewView,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Value:   "",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.String("project") == "" {
				manifest, err := s5.ReadManifest()
				if err != nil {
					fmt.Println("WARN :: Cannot find manifest.json")
					return nil
				}
				if err := ctx.Set("project", manifest.App.ID); err != nil {
					return nil
				}
			}
			return nil
		},
		EnableBashCompletion: true,
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
