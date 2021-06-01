package context

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/Nudm_SubscriberDataManagement"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/shynuu/ntn-qof/factory"
	"github.com/shynuu/ntn-qof/logger"
)

func init() {
	ntnContext.NfInstanceID = uuid.New().String()
}

var ntnContext NTNContext

type NTNContext struct {
	Name         string
	NfInstanceID string

	QoS         map[uint8]uint8
	Slice       []*factory.Slice
	Classifiers *factory.Classifiers
	SliceAware  bool

	URIScheme    models.UriScheme
	BindingIPv4  string
	RegisterIPv4 string
	SBIPort      int
	CPNodeID     pfcpType.NodeID

	UDMProfile models.NfProfile

	UPNodeIDs []pfcpType.NodeID
	Key       string
	PEM       string
	KeyLog    string

	NrfUri                         string
	QofUri                         string
	NFManagementClient             *Nnrf_NFManagement.APIClient
	NFDiscoveryClient              *Nnrf_NFDiscovery.APIClient
	SubscriberDataManagementClient *Nudm_SubscriberDataManagement.APIClient
}

func InitQofContext(config *factory.Config) {
	if config == nil {
		logger.CtxLog.Error("Config is nil")
		return
	}

	logger.CtxLog.Infof("ntnconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	if configuration.NtnName != "" {
		ntnContext.Name = configuration.NtnName
	}

	// ntnContext.QofUri = configuration.QofUri
	ntnContext.Slice = configuration.Slice
	ntnContext.QoS = configuration.QoS
	ntnContext.Classifiers = configuration.Classifiers
	ntnContext.SliceAware = configuration.SliceAware

	sbi := configuration.Sbi
	if sbi == nil {
		logger.CtxLog.Errorln("Configuration needs \"sbi\" value")
		return
	} else {
		ntnContext.URIScheme = models.UriScheme(sbi.Scheme)
		ntnContext.RegisterIPv4 = factory.QOF_DEFAULT_IPV4 // default localhost
		ntnContext.SBIPort = factory.QOF_DEFAULT_PORT_INT  // default port
		if sbi.RegisterIPv4 != "" {
			ntnContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			ntnContext.SBIPort = sbi.Port
		}

		if tls := sbi.TLS; tls != nil {
			ntnContext.Key = tls.Key
			ntnContext.PEM = tls.PEM
		}

		ntnContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if ntnContext.BindingIPv4 != "" {
			logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			ntnContext.BindingIPv4 = sbi.BindingIPv4
			if ntnContext.BindingIPv4 == "" {
				logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				ntnContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}

	if configuration.NrfUri != "" {
		ntnContext.NrfUri = configuration.NrfUri
	} else {
		logger.CtxLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		ntnContext.NrfUri = fmt.Sprintf("%s://%s:%d", ntnContext.URIScheme, "127.0.0.1", 29510)
	}

	// Set client and set url
	ManagementConfig := Nnrf_NFManagement.NewConfiguration()
	ManagementConfig.SetBasePath(NTN_Self().NrfUri)
	ntnContext.NFManagementClient = Nnrf_NFManagement.NewAPIClient(ManagementConfig)

	NFDiscovryConfig := Nnrf_NFDiscovery.NewConfiguration()
	NFDiscovryConfig.SetBasePath(NTN_Self().NrfUri)
	ntnContext.NFDiscoveryClient = Nnrf_NFDiscovery.NewAPIClient(NFDiscovryConfig)

}

func NTN_Self() *NTNContext {
	return &ntnContext
}
