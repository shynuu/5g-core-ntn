//+build !debug

package util

import (
	"github.com/free5gc/path_util"
)

var (
	QofLogPath           = path_util.Free5gcPath("free5gc/qofsslkey.log")
	QofPemPath           = path_util.Free5gcPath("free5gc/support/TLS/qof.pem")
	QofKeyPath           = path_util.Free5gcPath("free5gc/support/TLS/qof.key")
	DefaultQofConfigPath = path_util.Free5gcPath("free5gc/config/qofcfg.yaml")
)
