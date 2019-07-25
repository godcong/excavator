package main

import (
	"os"
	"sort"

	"github.com/godcong/excavator"
	"github.com/urfave/cli"
)

func main() {
	app := cli.App{
		Version: "v0.0.1",
		Name:    "excavator",
		Usage:   "excavator a dictionary",
		Action: func(c *cli.Context) error {
			url := ""
			if c.NArg() > 0 {
				url = c.Args().Get(0)
			}
			excavator.New(url, "")
			return nil
		},
		Flags: mainFlags(),
	}
	app.Commands = []cli.Command{}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		return
	}
}

func mainFlags() (flags []cli.Flag) {
	flags = []cli.Flag{
		cli.StringFlag{
			Name:  "workspace",
			Usage: "set workspace to storage temp file",
		},
	}
	return flags
}
