package cmds

import (
	"fmt"
	"github.com/apex/log"
	"github.com/cheggaaa/pb/v3"
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
			Usage:   "do not actually write files",
		},
		&cli.StringFlag{
			Name:  "ignore-pattern",
			Value: "(.*)\\.idea(.*)",
			Usage: "ignore files matching regex pattern",
		},
		&cli.BoolFlag{
			Name:    "no-backup",
			Value:   false,
			Aliases: []string{"B"},
			Usage:   "do not create backup files before writing",
		},
		&cli.BoolFlag{
			Name:    "remove-empty-lines",
			Value:   false,
			Aliases: []string{"L"},
			Usage:   "remove empty lines in xml files",
		},
		&cli.BoolFlag{
			Name:    "fast-fail",
			Value:   false,
			Aliases: []string{"x"},
			Usage:   "exit if any file fails to format",
		},
	},
	Action: func(ctx *cli.Context) error {
		log.Info("Looking for *.xml files and formatting them ...")

		var (
			dry      = ctx.Bool("dry-run")
			backup   = !ctx.Bool("no-backup")
			failFast = ctx.Bool("fast-fail")
			ignore   = regexp.MustCompile(ctx.String("ignore-pattern"))
		)

		var files []string

		// collect files to process
		if err := filepath.WalkDir(".", func(path string, info fs.DirEntry, _ error) error {
			if strings.ToLower(filepath.Ext(path)) != ".xml" {
				return nil
			}
			if ignore.MatchString(path) {
				log.Infof("Ignoring file: %s", path)
				return nil
			}
			files = append(files, path)
			return nil
		}); err != nil {
			return err
		}

		// process files
		//tpl := `{{ red "Formatting:" }} {{string . "path" | green}} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}} `
		//bar := pb.ProgressBarTemplate(tpl).Start(len(files))
		bar := pb.Full.Start(len(files))
		for _, path := range files {
			bar.Increment()
			bar.Set("prefix", "Formatting: "+path+" | ")

			var (
				err  error
				data []byte
			)
			if data, err = os.ReadFile(path); err != nil {
				log.WithError(err).Warnf("Cannot read file '%s'", path)
				if failFast {
					return err
				}
				continue
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
					if failFast {
						return err
					}
					continue
				}
			}

			// update data
			data = gohtml.FormatBytes(data, !ctx.Bool("remove-empty-lines"))
			if len(data) == 0 {
				log.Warn("formatted data was 0 bytes.")
				if failFast {
					return err
				}
				continue
			}

			if !dry {
				var f *os.File
				if f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755); err != nil {
					log.WithError(err).Warnf("Cannot open '%s'", path)
					if failFast {
						return err
					}
					continue
				}

				if err = f.Truncate(0); err != nil {
					log.WithError(err).Warnf("Cannot truncate '%s'", path)
					if failFast {
						return err
					}
					continue
				}

				_, err = f.Seek(0, 0)
				_, err = fmt.Fprintf(f, "%s", string(data))

				// close
				if err = f.Close(); err != nil {
					log.WithError(err).Warnf("Cannot close '%s'", path)
					if failFast {
						return err
					}
					continue
				}
			} else {
				log.Infof("[dry-run] %s would now look like this:", path)
				fmt.Println(string(data))
			}
		}
		bar.Finish()

		return nil
	},
}
