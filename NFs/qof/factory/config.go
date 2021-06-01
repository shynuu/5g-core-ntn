/*
 * QOF Configuration Factory
 */

package factory

import (
	"github.com/free5gc/logger_util"
	"github.com/free5gc/openapi/models"
)

const (
	QOF_EXPECTED_CONFIG_VERSION = "1.0.0"
)

type Config struct {
	Info          *Info               `yaml:"info"`
	Configuration *Configuration      `yaml:"configuration"`
	Logger        *logger_util.Logger `yaml:"logger"`
}

type Configuration struct {
	QofName         string           `yaml:"QofName,omitempty"`
	Sbi             *Sbi             `yaml:"sbi,omitempty"`
	NrfUri          string           `yaml:"nrfUri,omitempty"`
	NtnUri          string           `yaml:"ntnUri,omitempty"`
	ServiceNameList []string         `yaml:"serviceNameList,omitempty"`
	ULCL            bool             `yaml:"ulcl,omitempty"`
	QoS             map[int32]uint16 `yaml:"qos,omitempty"`
	Slice           []*Slice         `yaml:"slice,omitempty"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

const (
	QOF_DEFAULT_IPV4     = "127.0.0.1"
	QOF_DEFAULT_PORT     = "8000"
	QOF_DEFAULT_PORT_INT = 8000
)

type ControlPlane struct {
	RAN string `yaml:"ran" json:"ran"`
	CN  string `yaml:"cn" json:"cn"`
	ID  uint8  `yaml:"id" json:"id"`
}

type Slice struct {
	SNssai  *models.Snssai `yaml:"sNssai"`
	RAN     string         `yaml:"ran"`
	CN      string         `yaml:"cn"`
	ID      string         `yaml:"id"`
	AMF     string         `yaml:"amf,omitempty"`
	Default bool           `yaml:"default"`
}

type QoS struct {
	Var5QI int32  `yaml:"5qi,omitempty"`
	DSCP   uint16 `yaml:"dscp,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	TLS          *TLS   `yaml:"tls"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
}

type TLS struct {
	PEM string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

type QOFSessionInfo struct {
	SessionID int32 `json:"sessionid" yaml:"sessionid" bson:"sessionid"`
	Snssai    *models.Snssai
	Supi      string `json:"supi" yaml:"supi" bson:"supi"`
	UTEID     uint32 `json:"uteid" yaml:"supi" bson:"supi"`
	DTEID     uint32 `json:"dteid" yaml:"supi" bson:"supi"`
	IPv4      string `json:"ipv4" yaml:"ipv4" bson:"ipv4"`
	Var5QI    int32  `json:"var5qi" yaml:"var5qi" bson:"var5qi"`
}

type NTNSession struct {
	RAN        string      `json:"ran" yaml:"ran" bson:"ran"`
	UPF        string      `json:"upf" yaml:"upf" bson:"upf"`
	SliceMatch *SliceMatch `json:"slice_match" yaml:"slice_match" bson:"slice_match"`
	QosMatch   *QosMatch   `json:"qos_match" yaml:"qos_match" bson:"qos_match"`
	SliceID    string      `json:"id" yaml:"id" bson:"id"`
	IPv4       string      `json:"ipv4" yaml:"ipv4" bson:"ipv4"`
}

type SliceMatch struct {
	UTEID uint32 `json:"uteid" yaml:"supi" bson:"supi"`
	DTEID uint32 `json:"dteid" yaml:"supi" bson:"supi"`
}

type QosMatch struct {
	DSCP uint16 `json:"dscp" yaml:"dscp" bson:"dscp"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
