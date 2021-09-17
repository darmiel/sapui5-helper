package cmds

import (
	"fmt"
	"github.com/apex/log"
	"github.com/darmiel/gohtml"
	"github.com/urfave/cli/v2"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var CmdFormat = &cli.Command{
	Category: "Formatting",
	Name:     "fmt-xml",
	Aliases:  []string{"fx"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:  "ignore-pattern",
			Value: "(.*)\\.idea(.*)",
		},
		&cli.BoolFlag{
			Name:    "no-backup",
			Value:   false,
			Aliases: []string{"B"},
		},
	},
	Action: func(ctx *cli.Context) error {
		log.Info("Looking for *.xml files and formatting them ...")
		var dry = ctx.Bool("dry-run")
		var ignore = regexp.MustCompile(ctx.String("ignore-pattern"))
		var backup = !ctx.Bool("no-backup")

		if err := filepath.WalkDir(".", func(path string, info fs.DirEntry, _ error) error {
			if strings.ToLower(filepath.Ext(path)) != ".xml" {
				return nil
			}
			if ignore.MatchString(path) {
				log.Infof("Ignoring file: %s", path)
				return nil
			}
			log.Infof("Formatting file: %s", path)
			var (
				err  error
				data []byte
			)
			if data, err = os.ReadFile(path); err != nil {
				log.WithError(err).Warnf("Cannot read file '%s'", path)
				return nil
			}

			// make backup?
			if backup {
				var num = 0
				var backFile = ""
				for {
					num++
					backFile = fmt.Sprintf("%s.%d", path, num)
					if !Exists(backFile) {
						break
					}
				}
				if err = os.WriteFile(backFile, data, 0644); err != nil {
					log.WithError(err).Errorf("Cannot create backup of '%s'", path)
					return err
				}
			}

			// update data
			data = gohtml.FormatBytes(data, true)
			if len(data) == 0 {
				log.Warn("formatted data was 0 bytes.")
			}

			log.Infof("Writing to '%s':", path)
			fmt.Println(string(data))

			if !dry {
				var f *os.File
				if f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
					log.WithError(err).Warnf("Cannot open '%s'", path)
					return nil
				}

				if err = f.Truncate(0); err != nil {
					log.WithError(err).Warnf("Cannot truncate '%s'", path)
					return nil
				}

				_, err = f.Seek(0, 0)
				_, err = fmt.Fprintf(f, "%s", string(data))

				// close
				if err = f.Close(); err != nil {
					log.WithError(err).Warnf("Cannot close '%s'", path)
					return nil
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	},
}
