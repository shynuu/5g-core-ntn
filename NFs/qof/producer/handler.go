package producer

import (
	"errors"

	"github.com/free5gc/openapi/models"
	"github.com/gin-gonic/gin"
	"github.com/shynuu/qof/consumer"
	"github.com/shynuu/qof/context"
	"github.com/shynuu/qof/factory"
	"github.com/shynuu/qof/logger"
)

// TranslateSnssai returns the UPF and RAN IP
func TranslateSnssai(Snssai *models.Snssai) (upf string, ran string, id string, err error) {
	for _, v := range context.QOF_Self().Slice {
		if v.SNssai.Sst == Snssai.Sst && v.SNssai.Sd == Snssai.Sd {
			upf = v.CN
			ran = v.RAN
			id = v.ID
			return upf, ran, id, nil
		}
	}
	err = errors.New("impossibe to find a correct translation for S-NSSAI")
	return "", "", "", err
}

// Translate5QI translates the 5QI in DSCP
func Translate5QI(var5qi int32) (dscp uint16) {
	return context.QOF_Self().QoS[var5qi]
}

// HandleSessionCreateQof processes
func HandleSessionCreateQof(c *gin.Context) {

	logger.PduSessLog.Infoln("Handling Session Create from 5G QOF")

	var sessionInfo factory.QOFSessionInfo

	if err := c.BindJSON(&sessionInfo); err != nil {
		logger.PduSessLog.Errorln(err)
		c.JSON(500, gin.H{
			"message": "Error retrieving parameters",
		})
	}

	upf, ran, id, err := TranslateSnssai(sessionInfo.Snssai)
	if err != nil {
		logger.PduSessLog.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	dscp := Translate5QI(sessionInfo.Var5QI)
	var ntnSession *factory.NTNSession = &factory.NTNSession{
		UPF:      upf,
		RAN:      ran,
		QosMatch: &factory.QosMatch{DSCP: dscp},
		SliceMatch: &factory.SliceMatch{
			UTEID: sessionInfo.UTEID,
			DTEID: sessionInfo.DTEID,
		},
		SliceID: id,
		IPv4:    sessionInfo.IPv4,
	}

	consumer.NTN5GSessionCreate(ntnSession)

	c.JSON(200, gin.H{
		"message": "success",
	})
}
