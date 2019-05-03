package splitter

import (
	"github.com/sirupsen/logrus"
)

//type byteArray []byte

//ParseJSONByteToByte parse data from json. output to byte
func ParseJSONByteToByte(jsonData []byte, primaryPath string, additionalPaths []string) (primaryResult string, additionalData map[string]string, cleanedSourceData string, err error) {
	appenderMap := make(map[string]ByteBufferJSONParseResultAppender)
	appenderMap[primaryPath] = NewByteBufferJSONParseResultAppender()
	if additionalPaths != nil {
		for _, key := range additionalPaths {
			appenderMap[key] = NewByteBufferJSONParseResultAppender()
		}
	}
	appender := func(jsonPath string, bytes []byte) (errAppender error) {
		if appenderMap[jsonPath] == nil {
			return nil
		}
		theAppender := appenderMap[jsonPath]
		return theAppender.Append(bytes)
	}
	strJSONData := string(jsonData)
	rslt, err := parseJSONStringOld(strJSONData, appender)
	if err != nil {
		logrus.WithFields(logrus.Fields{"method": "ParseJSONByteToByte", "primaryPath": primaryPath}).Error("Error ParseJSONString: " + err.Error())
		return strJSONData, nil, "", err
	}
	var additionalDataRet map[string]string
	if additionalPaths != nil && len(additionalPaths) > 0 {
		additionalDataRet = make(map[string]string)
		for _, key := range additionalPaths {
			additionalDataRet[key] = string(appenderMap[key].Bytes())
		}
	}
	var primaryResultData = string(appenderMap[primaryPath].Bytes())
	// additionalDataRet = nil
	return primaryResultData, additionalDataRet, rslt, nil
}
