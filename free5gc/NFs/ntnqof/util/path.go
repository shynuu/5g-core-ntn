//+build !debug

package util

import (
	"github.com/free5gc/path_util"
)

var (
	NtnLogPath           = path_util.Free5gcPath("free5gc/ntnsslkey.log")
	NtnPemPath           = path_util.Free5gcPath("free5gc/support/TLS/ntn.pem")
	NtnKeyPath           = path_util.Free5gcPath("free5gc/support/TLS/ntn.key")
	DefaultNtnConfigPath = path_util.Free5gcPath("free5gc/config/ntncfg.yaml")
)
