package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	smf_context "github.com/free5gc/smf/context"
	"github.com/free5gc/smf/logger"
)

/*
SendSessionQOF Read the profile of a given NF Instance
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param nfInstanceID Unique ID of the NF Instance
@return models.NfProfile
*/
func SendSessionQOF(sessionInfo *smf_context.QOFSessionInfo) error {

	var url string = fmt.Sprintf("%s/qof-session/new-session", smf_context.SMF_Self().QofUri)

	client := http.Client{}

	reqBody, err := json.Marshal(sessionInfo)

	if err != nil {
		logger.PduSessLog.Errorln("Impossible to serialzie session Info")
		return err
	}
	logger.PduSessLog.Infoln(string(reqBody))
	resp, err := client.Post(url,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logger.PduSessLog.Errorln(err)
		logger.PduSessLog.Errorln("Impossible to post session Info to 5G QOF")
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.PduSessLog.Errorln("Impossible to read the fle")
		return err
	}
	logger.PduSessLog.Infoln(string(body))
	return nil
}
