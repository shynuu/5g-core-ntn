package consumer

// func NTN5GSessionCreate(ntnSession *factory.NTNSession) error {

// logger.PduSessLog.Infoln("Handling NTN 5G Session Create")

// var url string = context.NTN_Self().QofUri

// client := http.Client{}

// reqBody, err := json.Marshal(ntnSession)

// if err != nil {
// 	logger.PduSessLog.Errorln("Impossible to serialzie NTN session Info")
// 	return err
// }

// resp, err := client.Post(url,
// 	"application/json", bytes.NewBuffer(reqBody))
// if err != nil {
// 	logger.PduSessLog.Errorln(err)
// 	logger.PduSessLog.Errorln("Impossible to post session Info to NTN QOF")
// 	return err
// }
// defer resp.Body.Close()
// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	logger.PduSessLog.Errorln("Impossible to read the body")
// 	return err
// }
// logger.PduSessLog.Infoln(string(body))
// return nil
// }
