package splitter

import (
	"fmt"
)

//ParseJSONStringToFile parse json ke file
func ParseJSONStringToFile(targetData []byte, extractDefinitions []ExtractJSONToFileDefinition, removedItemPaths []string) (markerForDeletion JSONItemRemover, err error) {
	entry := WrapLogWithClassAndMethod(nil, "ParseJSONString", "ParseJSONStringToFile")
	if extractDefinitions == nil || len(extractDefinitions) == 0 {
		logEntry := WrapLogWithClassAndMethod(nil, "ParseJSONString", "ParseJSONStringToFile")
		logEntry.Error("Extract definision kosong. tidak ada item di proses")
		return nil, fmt.Errorf("No parsed field defined. no ")
	}
	appenderMap := make(map[string]AppenderSinglePathFunction)
	var theAppender []JSONParseSingleResultAppender
	var idx = 0
	for _, def := range extractDefinitions {
		fileAppender := NewFileBackedJSONParseResultAppender(def.DestinationFilePath)

		errOpen := fileAppender.OpenAppender()
		if errOpen != nil {
			entry.WithField("appenderIndex", idx).Error(errOpen)
			return nil, errOpen
		}
		appenderMap[def.PathToExtract] = func(bytes []byte) (err error) {
			return fileAppender.Append(bytes)
		}

		theAppender = append(theAppender, fileAppender)
		idx++
	}
	defer closeAllAppenders(theAppender)
	var markerForDeletionAct JSONItemRemover
	if len(removedItemPaths) > 0 {
		markerForDeletionAct = NewJSONItemRemover(removedItemPaths)
	}
	targetDataString := string(targetData)
	return markerForDeletionAct, parseJSONStringActualWorker(targetDataString, appenderMap, markerForDeletionAct)
}
