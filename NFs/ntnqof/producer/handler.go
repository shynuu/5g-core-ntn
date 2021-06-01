// Package producer is
package producer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shynuu/ntn-qof/context"
	"github.com/shynuu/ntn-qof/factory"
	"github.com/shynuu/ntn-qof/logger"
)

type MobileSession struct {
	RAN        string      `json:"ran" yaml:"ran" bson:"ran"`
	UPF        string      `json:"upf" yaml:"upf" bson:"upf"`
	SliceMatch *SliceMatch `json:"slice_match" yaml:"slice_match" bson:"slice_match"`
	QosMatch   *QosMatch   `json:"qos_match" yaml:"qos_match" bson:"qos_match"`
	SliceID    string      `json:"id" yaml:"id" bson:"id"`
	IPv4       string      `json:"ipv4" yaml:"ipv4" bson:"ipv4"`
}

type SliceMatch struct {
	UTEID uint32 `json:"uteid" yaml:"uteid" bson:"uteid"`
	DTEID uint32 `json:"dteid" yaml:"uteid" bson:"uteid"`
}

type QosMatch struct {
	DSCP uint8 `json:"dscp" yaml:"dscp" bson:"dscp"`
}

type PipeRule struct {
	TEID    uint32 `json:"teid" yaml:"teid" bson:"teid"`
	DSCP5   uint8  `json:"dscp_5g" yaml:"dscp_5g" bson:"dscp_5g"`
	DSCPS   uint8  `json:"dscp_satellite" yaml:"dscp_satellite" bson:"dscp_satellite"`
	SliceID uint8  `json:"slice_id" yaml:"slice_id" bson:"slice_id"`
	IPv4    string `json:"ipv4" yaml:"ipv4" bson:"ipv4"`
}

type IPipeRule struct {
	DSCP5   uint8 `json:"dscp_5g" yaml:"dscp_5g" bson:"dscp_5g"`
	DSCPS   uint8 `json:"dscp_satellite" yaml:"dscp_satellite" bson:"dscp_satellite"`
	SliceID uint8 `json:"slice_id" yaml:"slice_id" bson:"slice_id"`
}

type PDU struct {
	TEID     uint32 `json:"teid" yaml:"teid" bson:"teid"`
	DSCP5    uint8  `json:"dscp_5g" yaml:"dscp_5g" bson:"dscp_5g"`
	DSCPS    uint8  `json:"dscp_satellite" yaml:"dscp_satellite" bson:"dscp_satellite"`
	SliceID  uint8  `json:"slice_id" yaml:"slice_id" bson:"slice_id"`
	IPv4     string `json:"ipv4" yaml:"ipv4" bson:"ipv4"`
	IsRAN    bool   `json:"is_ran" yaml:"is_ran" bson:"is_ran"`
	Endpoint string `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
	Ingress  string `json:"ingress" yaml:"ingress" bson:"ingress"`
}

type ADMControl struct {
	SliceID    uint8  `json:"slice_id" yaml:"slice_id" bson:"slice_id"`
	Throughput int    `json:"throughput" yaml:"throughput" bson:"throughput"`
	Endpoint   string `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
}

type ADM struct {
	Controls []ADMControl `json:"controls" yaml:"controls" bson:"controls"`
	Aware    bool         `json:"slice_aware" yaml:"slice_aware" bson:"slice_aware"`
}

// TranslateQoS translates the 5G DSCP to Satellite DSCP
func TranslateQoS(qosMatch *QosMatch) uint8 {
	return context.NTN_Self().QoS[qosMatch.DSCP]
}

// MapSlice maps the 5G slice to the satellite slice
func MapSlice(sliceID uint8) *factory.Slice {
	for _, s := range context.NTN_Self().Slice {
		if s.SliceID == sliceID {
			return s
		}
	}
	return nil
}

// GetEgressInterface returns the classifier egress interface corresponding to a slice ID
func GetEgressInterface(classifier *factory.Classifier, ipV4 string) (string, error) {
	var iop string = ""
	ipV4 = fmt.Sprintf("%s/24", ipV4)
	_, ipnet, err := net.ParseCIDR(ipV4)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to find the address")
		return "", err
	}

	for _, i := range classifier.Egress {
		ip := net.ParseIP(i)
		if ipnet.Contains(ip) {
			iop = ip.To4().String()
			return iop, nil
		}
	}

	return "", err
}

// GetIngressInterface returns the classifier ingress interface corresponding to a slice ID
func GetIngressInterface(classifier *factory.Classifier, ipV4 string) (string, error) {
	var iop string = ""
	ipV4 = fmt.Sprintf("%s/24", ipV4)
	_, ipnet, err := net.ParseCIDR(ipV4)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to find the address")
		return "", err
	}

	for _, i := range classifier.Ingress {
		ip := net.ParseIP(i)
		if ipnet.Contains(ip) {
			iop = ip.To4().String()
			return iop, nil
		}
	}

	return "", err
}

func AdmissionControl(classifier *factory.Classifier, adm ADM, wg *sync.WaitGroup) error {

	var uri string = fmt.Sprintf("http://%s:%d/control-plane/adm", classifier.RegisterIPv4, classifier.Port)
	client := http.Client{}

	reqBody, err := json.Marshal(&adm)

	if err != nil {
		logger.PduSessLog.Errorln("Impossible to serialzie Admission Control")
		return err
	}

	resp, err := client.Post(uri,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.PduSessLog.Errorln(err)
		logger.PduSessLog.Errorln("Impossible to send rules to classifiers")
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to read the body")
		return err
	}
	logger.PduSessLog.Infof("Admission Control OK for classifier %s", classifier.RegisterIPv4)

	wg.Done()
	return nil
}

func Pipe(classifier *factory.Classifier,
	teid uint32,
	dscp5G uint8,
	dscpSatellite uint8,
	sliceID uint8,
	ipv4 string,
	endpoint string,
	ingress string,
	isRan bool,
	wg *sync.WaitGroup) error {

	var uri string = fmt.Sprintf("http://%s:%d/data-plane/pdu", classifier.RegisterIPv4, classifier.Port)
	client := http.Client{}

	pdu := &PDU{
		DSCP5:    dscp5G,
		DSCPS:    dscpSatellite,
		TEID:     teid,
		SliceID:  sliceID,
		IPv4:     ipv4,
		Endpoint: endpoint,
		Ingress:  ingress,
		IsRAN:    isRan,
	}

	reqBody, err := json.Marshal(pdu)

	if err != nil {
		logger.PduSessLog.Errorln("Impossible to serialzie PipeRule info")
		return err
	}

	resp, err := client.Post(uri,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.PduSessLog.Errorln(err)
		logger.PduSessLog.Errorln("Impossible to send rules to classifiers")
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to read the body")
		return err
	}
	logger.PduSessLog.Infoln(string(body))

	wg.Done()
	return nil
}

// func IPipe(classifier *factory.Classifier,
// 	dscp5G uint8,
// 	dscpSatellite uint8,
// 	sliceID uint8,
// 	wg *sync.WaitGroup) error {

// 	var uri string = fmt.Sprintf("http://%s:%d/data-plane/ipipe", classifier.RegisterIPv4, classifier.Port)
// 	client := http.Client{}

// 	pipe := &IPipeRule{
// 		DSCP5:   dscp5G,
// 		DSCPS:   dscpSatellite,
// 		SliceID: sliceID,
// 	}

// 	reqBody, err := json.Marshal(pipe)

// 	if err != nil {
// 		logger.PduSessLog.Errorln("Impossible to serialzie PipeRule info")
// 		return err
// 	}

// 	resp, err := client.Post(uri,
// 		"application/json", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		logger.PduSessLog.Errorln(err)
// 		logger.PduSessLog.Errorln("Impossible to send rules to classifiers")
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.PduSessLog.Errorln("Impossible to read the body")
// 		return err
// 	}
// 	logger.PduSessLog.Infoln(string(body))

// 	wg.Done()
// 	return nil

// }

// HandleAdmissionControl handle the ADM from 5G
func HandleAdmissionControl(c *gin.Context) {

	var controlPlane factory.ControlPlane

	if err := c.BindJSON(&controlPlane); err != nil {
		logger.PduSessLog.Errorln(err)
		c.JSON(500, gin.H{
			"message": "Error retrieving parameters",
		})
		return
	}

	logger.PduSessLog.Infoln("Handling Admission Control")

	var admRAN *ADM = &ADM{
		Controls: make([]ADMControl, len(context.NTN_Self().Slice)),
		Aware:    context.NTN_Self().SliceAware,
	}
	var admCN *ADM = &ADM{
		Controls: make([]ADMControl, len(context.NTN_Self().Slice)),
		Aware:    context.NTN_Self().SliceAware,
	}

	for k, sl := range context.NTN_Self().Slice {

		admCN.Controls[k] = ADMControl{
			SliceID:    sl.SliceID,
			Throughput: sl.Forward,
			Endpoint:   sl.ClassifierCNEndpoint,
		}

		admRAN.Controls[k] = ADMControl{
			SliceID:    sl.SliceID,
			Throughput: sl.Return,
			Endpoint:   sl.ClassifierRANEndpoint,
		}

	}

	classifierCN := context.NTN_Self().Classifiers.CN
	classifierRAN := context.NTN_Self().Classifiers.RAN

	var wg sync.WaitGroup

	wg.Add(2)
	go AdmissionControl(classifierCN, *admCN, &wg)
	go AdmissionControl(classifierRAN, *admRAN, &wg)

	wg.Wait()

	c.JSON(200, gin.H{
		"message": "success",
	})

}

// HandleSessionCreateQof handles the PDU Session creation on the satellite side
func HandleSessionCreateQof(c *gin.Context) {

	logger.PduSessLog.Infoln("Handling PDU Session Creation")

	var mobileSession MobileSession

	if err := c.BindJSON(&mobileSession); err != nil {
		logger.PduSessLog.Errorln(err)
		c.JSON(500, gin.H{
			"message": "Error retrieving parameters",
		})
		return
	}

	logger.PduSessLog.Infof("New 5G session created for slice %s", mobileSession.SliceID)

	// Translate the 5G DSCP to Satellite DSCP
	dscp5G := mobileSession.QosMatch.DSCP
	dscpSatellite := TranslateQoS(mobileSession.QosMatch)
	logger.PduSessLog.Infof("Slice ID: %s, DSCP 5G: %d, DSCP SAT: %d", mobileSession.SliceID, dscp5G, dscpSatellite)
	logger.PduSessLog.Infof("RAN EP: %s, CN EP: %s", mobileSession.RAN, mobileSession.UPF)

	u64, _ := strconv.ParseUint(mobileSession.SliceID, 10, 64)

	// Get the ST endpoint and the GW endpoint
	sliceSatellite := MapSlice(uint8(u64))
	// logger.PduSessLog.Infof("Getting ST endpoint %s and GW endpoint %s for Slice ID %d", sliceSatellite.StEndpoint, sliceSatellite.GwEndpoint, sliceSatellite.SliceId)
	// logger.PduSessLog.Infof("5G RAN endpoint %s and 5G CN endpoint %s for Slice ID %d", mobileSession.RAN, mobileSession.UPF, uint8(u64))

	classifierCN := context.NTN_Self().Classifiers.CN
	classifierRAN := context.NTN_Self().Classifiers.RAN

	// Get the Egress interfaces for the IPipe operation
	// classifierRANEgress, _ := GetEgressInterface(classifierRAN, sliceSatellite.StEndpoint)
	// classifierCNEgress, _ := GetEgressInterface(classifierCN, sliceSatellite.GwEndpoint)
	// logger.PduSessLog.Infof("Egress interface of RAN %s and CN %s", classifierRANEgress, classifierCNEgress)

	// Get the Ingress interfaces for the Pipe operation
	classifierRANIngress, _ := GetIngressInterface(classifierRAN, mobileSession.RAN)
	classifierCNIngress, _ := GetIngressInterface(classifierCN, mobileSession.UPF)
	logger.PduSessLog.Infof("Classifier CN: %s, Classifier RAN: %s", classifierCNIngress, classifierRANIngress)
	// logger.PduSessLog.Infof("Ingress interface of RAN %s and CN %s", classifierRANIngress, classifierCNIngress)

	// logger.PduSessLog.Infof("Got IP of UE %s", mobileSession.IPv4)

	var wg sync.WaitGroup
	wg.Add(2)

	// Programm the forward link
	go Pipe(classifierCN, mobileSession.SliceMatch.DTEID, dscp5G, dscpSatellite, sliceSatellite.SliceID, mobileSession.UPF, sliceSatellite.ClassifierCNEndpoint, classifierCNIngress, false, &wg)
	// go IPipe(classifierRAN, dscp5G, dscpSatellite, classifierRANEgress, mobileSession.RAN, sliceSatellite.SliceID, mobileSession.UPF, mobileSession.IPv4, true, &wg)

	// Programm the return link
	go Pipe(classifierRAN, mobileSession.SliceMatch.UTEID, dscp5G, dscpSatellite, sliceSatellite.SliceID, mobileSession.UPF, sliceSatellite.ClassifierRANEndpoint, classifierRANIngress, true, &wg)
	// go IPipe(classifierCN, dscp5G, dscpSatellite, classifierCNEgress, mobileSession.UPF, sliceSatellite.SliceID, mobileSession.RAN, mobileSession.IPv4, false, &wg)

	wg.Wait()

	c.JSON(200, gin.H{
		"message": "success",
	})
}
