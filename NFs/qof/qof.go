package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/free5gc/version"
	"github.com/shynuu/qof/logger"
	"github.com/shynuu/qof/service"
)

var QOF = &service.QOF{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "qof"
	fmt.Print(app.Name, "\n")
	appLog.Infoln("QOF version: ", version.GetVersion())
	app.Usage = "-free5gccfg common configuration file -qofcfg qof configuration file"
	app.Action = action
	app.Flags = QOF.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("QOF Run error: %v", err)
	}
}

func action(c *cli.Context) error {
	if err := QOF.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	QOF.Start()

	return nil
}
