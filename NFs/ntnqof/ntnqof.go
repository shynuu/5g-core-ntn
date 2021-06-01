package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/version"
	"github.com/shynuu/ntn-qof/logger"
	"github.com/shynuu/ntn-qof/service"
)

var NTN = &service.NTN{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "ntn"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("NTN version: ", version.GetVersion())
	app.Usage = "-ntncfg ntn configuration file"
	app.Action = action
	app.Flags = NTN.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("NTN Run error: %v", err)
	}
}

func action(c *cli.Context) error {
	if err := NTN.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	NTN.Start()

	return nil
}
