package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	aperLogger "github.com/free5gc/aper/logger"
	"github.com/free5gc/http2_util"
	"github.com/free5gc/logger_util"
	nasLogger "github.com/free5gc/nas/logger"
	ngapLogger "github.com/free5gc/ngap/logger"
	openApiLogger "github.com/free5gc/openapi/logger"
	"github.com/free5gc/path_util"
	pathUtilLogger "github.com/free5gc/path_util/logger"
	"github.com/shynuu/ntn-qof/context"
	"github.com/shynuu/ntn-qof/factory"
	"github.com/shynuu/ntn-qof/logger"
	"github.com/shynuu/ntn-qof/producer"
	"github.com/shynuu/ntn-qof/util"
)

type NTN struct{}

type (
	// Config information.
	Config struct {
		ntncfg string
	}
)

var config Config

var ntnCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "ntncfg",
		Usage: "config file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*NTN) GetCliCmd() (flags []cli.Flag) {
	return ntnCLi
}

func (ntn *NTN) Initialize(c *cli.Context) error {
	config = Config{
		ntncfg: c.String("ntncfg"),
	}

	if config.ntncfg != "" {
		if err := factory.InitConfigFactory(config.ntncfg); err != nil {
			return err
		}
	} else {
		DefaultNtnConfigPath := path_util.Free5gcPath("free5gc/config/ntncfg.yaml")
		if err := factory.InitConfigFactory(DefaultNtnConfigPath); err != nil {
			return err
		}
	}

	ntn.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (ntn *NTN) setLogLevel() {
	if factory.QofConfig.Logger == nil {
		initLog.Warnln("NTN config without log level setting!!!")
		return
	}

	if factory.QofConfig.Logger.SMF != nil {
		if factory.QofConfig.Logger.SMF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.SMF.DebugLevel); err != nil {
				initLog.Warnf("NTN Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.SMF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("NTN Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("NTN Log level is default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.QofConfig.Logger.SMF.ReportCaller)
	}

	if factory.QofConfig.Logger.NAS != nil {
		if factory.QofConfig.Logger.NAS.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.NAS.DebugLevel); err != nil {
				nasLogger.NasLog.Warnf("NAS Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.NAS.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				nasLogger.SetLogLevel(level)
			}
		} else {
			nasLogger.NasLog.Warnln("NAS Log level not set. Default set to [info] level")
			nasLogger.SetLogLevel(logrus.InfoLevel)
		}
		nasLogger.SetReportCaller(factory.QofConfig.Logger.NAS.ReportCaller)
	}

	if factory.QofConfig.Logger.NGAP != nil {
		if factory.QofConfig.Logger.NGAP.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.NGAP.DebugLevel); err != nil {
				ngapLogger.NgapLog.Warnf("NGAP Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.NGAP.DebugLevel)
				ngapLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				ngapLogger.SetLogLevel(level)
			}
		} else {
			ngapLogger.NgapLog.Warnln("NGAP Log level not set. Default set to [info] level")
			ngapLogger.SetLogLevel(logrus.InfoLevel)
		}
		ngapLogger.SetReportCaller(factory.QofConfig.Logger.NGAP.ReportCaller)
	}

	if factory.QofConfig.Logger.Aper != nil {
		if factory.QofConfig.Logger.Aper.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.Aper.DebugLevel); err != nil {
				aperLogger.AperLog.Warnf("Aper Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.Aper.DebugLevel)
				aperLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				aperLogger.SetLogLevel(level)
			}
		} else {
			aperLogger.AperLog.Warnln("Aper Log level not set. Default set to [info] level")
			aperLogger.SetLogLevel(logrus.InfoLevel)
		}
		aperLogger.SetReportCaller(factory.QofConfig.Logger.Aper.ReportCaller)
	}

	if factory.QofConfig.Logger.PathUtil != nil {
		if factory.QofConfig.Logger.PathUtil.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.PathUtil.DebugLevel); err != nil {
				pathUtilLogger.PathLog.Warnf("PathUtil Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.PathUtil.DebugLevel)
				pathUtilLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				pathUtilLogger.SetLogLevel(level)
			}
		} else {
			pathUtilLogger.PathLog.Warnln("PathUtil Log level not set. Default set to [info] level")
			pathUtilLogger.SetLogLevel(logrus.InfoLevel)
		}
		pathUtilLogger.SetReportCaller(factory.QofConfig.Logger.PathUtil.ReportCaller)
	}

	if factory.QofConfig.Logger.OpenApi != nil {
		if factory.QofConfig.Logger.OpenApi.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.QofConfig.Logger.OpenApi.DebugLevel); err != nil {
				openApiLogger.OpenApiLog.Warnf("OpenAPI Log level [%s] is invalid, set to [info] level",
					factory.QofConfig.Logger.OpenApi.DebugLevel)
				openApiLogger.SetLogLevel(logrus.InfoLevel)
			} else {
				openApiLogger.SetLogLevel(level)
			}
		} else {
			openApiLogger.OpenApiLog.Warnln("OpenAPI Log level not set. Default set to [info] level")
			openApiLogger.SetLogLevel(logrus.InfoLevel)
		}
		openApiLogger.SetReportCaller(factory.QofConfig.Logger.OpenApi.ReportCaller)
	}

}

func (ntn *NTN) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range ntn.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (ntn *NTN) Start() {
	context.InitQofContext(&factory.QofConfig)
	// allocate id for each upf

	initLog.Infoln("Server started")
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// err := consumer.SendNFRegistration()
	// if err != nil {
	// 	retry_err := consumer.RetrySendNFRegistration(10)
	// 	if retry_err != nil {
	// 		logger.InitLog.Errorln(retry_err)
	// 		return
	// 	}
	// }

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		ntn.Terminate()
		os.Exit(0)
	}()

	producer.AddService(router)

	time.Sleep(1000 * time.Millisecond)

	HTTPAddr := fmt.Sprintf("%s:%d", context.NTN_Self().BindingIPv4, context.NTN_Self().SBIPort)
	server, err := http2_util.NewServer(HTTPAddr, util.NtnLogPath, router)

	if server == nil {
		initLog.Error("Initialize HTTP server failed:", err)
		return
	}

	if err != nil {
		initLog.Warnln("Initialize HTTP server:", err)
	}

	serverScheme := factory.QofConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(util.NtnPemPath, util.NtnKeyPath)
	}

	if err != nil {
		initLog.Fatalln("HTTP server setup failed:", err)
	}
}

func (ntn *NTN) Terminate() {
	logger.InitLog.Infof("Terminating NTN...")
	// deregister with NRF
}

func (ntn *NTN) Exec(c *cli.Context) error {
	return nil
}
