package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/Nudm_SubscriberDataManagement"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/shynuu/qof/factory"
	"github.com/shynuu/qof/logger"
)

func init() {
	qofContext.NfInstanceID = uuid.New().String()
}

var qofContext QOFContext

type QOFContext struct {
	Name         string
	NfInstanceID string

	QoS   map[int32]uint16
	Slice []*factory.Slice

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
	NtnUri                         string
	NFManagementClient             *Nnrf_NFManagement.APIClient
	NFDiscoveryClient              *Nnrf_NFDiscovery.APIClient
	SubscriberDataManagementClient *Nudm_SubscriberDataManagement.APIClient
}

func InitQofContext(config *factory.Config) {
	if config == nil {
		logger.CtxLog.Error("Config is nil")
		return
	}

	logger.CtxLog.Infof("qofconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	if configuration.QofName != "" {
		qofContext.Name = configuration.QofName
	}

	qofContext.QoS = configuration.QoS
	qofContext.Slice = configuration.Slice
	qofContext.NtnUri = configuration.NtnUri

	sbi := configuration.Sbi
	if sbi == nil {
		logger.CtxLog.Errorln("Configuration needs \"sbi\" value")
		return
	} else {
		qofContext.URIScheme = models.UriScheme(sbi.Scheme)
		qofContext.RegisterIPv4 = factory.QOF_DEFAULT_IPV4 // default localhost
		qofContext.SBIPort = factory.QOF_DEFAULT_PORT_INT  // default port
		if sbi.RegisterIPv4 != "" {
			qofContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			qofContext.SBIPort = sbi.Port
		}

		if tls := sbi.TLS; tls != nil {
			qofContext.Key = tls.Key
			qofContext.PEM = tls.PEM
		}

		qofContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if qofContext.BindingIPv4 != "" {
			logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			qofContext.BindingIPv4 = sbi.BindingIPv4
			if qofContext.BindingIPv4 == "" {
				logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				qofContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}

	if configuration.NrfUri != "" {
		qofContext.NrfUri = configuration.NrfUri
	} else {
		logger.CtxLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		qofContext.NrfUri = fmt.Sprintf("%s://%s:%d", qofContext.URIScheme, "127.0.0.1", 29510)
	}

	// Set client and set url
	ManagementConfig := Nnrf_NFManagement.NewConfiguration()
	ManagementConfig.SetBasePath(QOF_Self().NrfUri)
	qofContext.NFManagementClient = Nnrf_NFManagement.NewAPIClient(ManagementConfig)

	NFDiscovryConfig := Nnrf_NFDiscovery.NewConfiguration()
	NFDiscovryConfig.SetBasePath(QOF_Self().NrfUri)
	qofContext.NFDiscoveryClient = Nnrf_NFDiscovery.NewAPIClient(NFDiscovryConfig)

}

func QOF_Self() *QOFContext {
	return &qofContext
}

func InitDefaultSlice() error {

	logger.PduSessLog.Infoln("Handling Default Slice")

	var url string = fmt.Sprintf("%s/ntn-session/admission-control", QOF_Self().NtnUri)

	var controlPlaneInfo *factory.ControlPlane

	for _, sl := range QOF_Self().Slice {
		if sl.Default {
			u64, _ := strconv.ParseUint(sl.ID, 10, 8)
			u8 := uint8(u64)
			controlPlaneInfo = &factory.ControlPlane{
				RAN: sl.RAN,
				CN:  sl.AMF,
				ID:  u8,
			}
		}
	}

	client := http.Client{}

	reqBody, err := json.Marshal(controlPlaneInfo)

	if err != nil {
		logger.PduSessLog.Errorln("Impossible to serialzie Default Slice Info")
		return err
	}

	resp, err := client.Post(url,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.PduSessLog.Errorln(err)
		logger.PduSessLog.Errorln("Impossible to post session Info to NTN QOF")
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to read the body")
		return err
	}
	logger.PduSessLog.Infoln(string(body))
	return nil
}
