/*
 * QOF Configuration Factory
 */

package factory

import (
	"github.com/free5gc/logger_util"
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
	NtnName         string          `yaml:"NtnName,omitempty"`
	Sbi             *Sbi            `yaml:"sbi,omitempty"`
	NrfUri          string          `yaml:"nrfUri,omitempty"`
	ServiceNameList []string        `yaml:"serviceNameList,omitempty"`
	ULCL            bool            `yaml:"ulcl,omitempty"`
	QoS             map[uint8]uint8 `yaml:"qos,omitempty"`
	SliceAware      bool            `yaml:"slice_aware,omitempty"`
	Slice           []*Slice        `yaml:"slice,omitempty"`
	Classifiers     *Classifiers    `yaml:"classifiers,omitempty"`
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
	RAN     string `yaml:"ran" json:"ran"`
	CN      string `yaml:"cn" json:"cn"`
	SliceID uint8  `yaml:"id" json:"id"`
}

type Slice struct {
	SliceID               uint8  `yaml:"id"`
	ClassifierRANEndpoint string `yaml:"classifier-ran-endpoint"`
	ClassifierCNEndpoint  string `yaml:"classifier-cn-endpoint"`
	Forward               int    `yaml:"forward"`
	Return                int    `yaml:"return"`
}

type Classifiers struct {
	RAN *Classifier `yaml:"ran"`
	CN  *Classifier `yaml:"cn"`
}

type Classifier struct {
	RegisterIPv4 string   `yaml:"registerIPv4,omitempty"`
	Port         int      `yaml:"port,omitempty"`
	Ingress      []string `yaml:"ingress"`
	Egress       []string `yaml:"egress"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	TLS          *TLS   `yaml:"tls"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"` // IP that is registered at NRF.
	BindingIPv4  string `yaml:"bindingIPv4,omitempty"`  // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty"`
}

type TLS struct {
	PEM string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
