package main

import (
	"helper/command"
	"helper/log"
	"os"
	"path/filepath"

	"github.com/teris-io/cli"
)

func main() {
	app := cli.New("Helper").
		WithOption(cli.NewOption("verbose", "Verbose execution").WithChar('v').WithType(cli.TypeBool)).
		WithCommand(newCmdTrans())

	os.Exit(app.Run(os.Args, os.Stdout))
}

func newCmdTrans() cli.Command {
	trans := cli.NewCommand("trans", "trans xcframework to flat framework").
		WithArg(cli.NewArg("xcframework", "xcframework path")).
		WithOption(cli.NewOption("output", "output path").WithType(cli.TypeString)).
		WithAction(func(args []string, options map[string]string) int {
			output := options["output"]
			xcframework := args[0]
			if !filepath.IsAbs(xcframework) {
				absPath, err := filepath.Abs(xcframework)
				if nil == err {
					xcframework = absPath
				}
			}
			err := command.TransXCFramework(xcframework, output)
			if nil != err {
				log.Error(err)
				return 1
			}
			return 0
		})

	return trans
}
