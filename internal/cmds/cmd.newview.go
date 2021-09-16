package cmds

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/darmiel/sapui5-helper/pkg/gen"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strings"
)

var RmdNewView = &cli.Command{
	Category: "Boilerplate",
	Name:     "newview",
	Aliases:  []string{"nv"},
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "controller-dir",
			Value:   "controller",
			Aliases: []string{"c"},
		},
		&cli.PathFlag{
			Name:    "view-dir",
			Value:   "view",
			Aliases: []string{"v"},
		},
	},
	Action: func(ctx *cli.Context) error {
		resp := struct {
			Typ       string `survey:"typ"`
			NameSpace string `survey:"namespace"`
			View      string `survey:"view"`
		}{}

		if err := survey.Ask([]*survey.Question{
			{
				Name:     "namespace",
				Validate: survey.Required,
				Prompt: &survey.Input{
					Message: "Namespace",
					Default: ctx.String("project"),
				},
			},
			{
				Name:     "view",
				Validate: survey.Required,
				Prompt: &survey.Input{
					Message: "View Name",
					Default: "MyView",
				},
			},
			{
				Name:     "typ",
				Validate: survey.Required,
				Prompt: &survey.Select{
					Message: "Select Type",
					Options: []string{
						"XML",
						"JS",
					},
					Default: "XML",
				},
			},
		}, &resp); err != nil {
			return err
		}

		// replace spaces in view name with ""
		resp.View = strings.ReplaceAll(resp.View, " ", "")

		var (
			viewName    string
			viewContent string
			ctlName     string
			ctlContent  string
		)

		switch resp.Typ {
		case "XML":
			viewName = fmt.Sprintf("%s.view.xml", resp.View)
			viewContent = gen.GenerateXMLView(resp.NameSpace, resp.View)
		case "JS":
			viewName = fmt.Sprintf("%s.view.js", resp.View)
			viewContent = gen.GenerateJSView(resp.NameSpace, resp.View)
		}

		ctlName = fmt.Sprintf("%s.controller.js", resp.View)
		ctlContent = gen.GenerateController(resp.NameSpace, resp.View)

		var (
			viewDir = ctx.Path("view-dir")
			ctlDir  = ctx.Path("controller-dir")
		)

		// inside webapp folder?
		if Exists("webapp") {
			viewDir = path.Join("webapp", viewDir)
			ctlDir = path.Join("webapp", ctlDir)
		}

		// check if view dir exists
		if !Exists(viewDir) {
			// create view dir
			if err := os.Mkdir(viewDir, 0755); err != nil {
				panic(err)
			}
		}
		if !Exists(ctlDir) {
			if err := os.Mkdir(ctlDir, 0755); err != nil {
				panic(err)
			}
		}

		// write to files
		if err := os.WriteFile(path.Join(viewDir, viewName), []byte(viewContent), 0755); err != nil {
			fmt.Println("WARN :: Cannot write view:", err)
			return err
		}
		if err := os.WriteFile(path.Join(ctlDir, ctlName), []byte(ctlContent), 0755); err != nil {
			fmt.Println("WARN :: Cannot write controller:", err)
			return err
		}

		fmt.Println("OK :: Written view and controller.")
		return nil
	},
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
