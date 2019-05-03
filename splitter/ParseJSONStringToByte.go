package splitter

import (
	"fmt"
)

//ParseJSONStringToByte parser 1 json string ke appender byte
func ParseJSONStringToByte(targetData []byte, extractDefinitions []ExtractJSONToByteDefinition, removedItemPaths []string) (markerForDeletion JSONItemRemover, err error) {
	entry := WrapLogWithClassAndMethod(nil, "ParseJSONString", "ParseJSONStringToByte")
	if extractDefinitions == nil || len(extractDefinitions) == 0 {
		entry.Error("Extract definision kosong. tidak ada item di proses")
		return nil, fmt.Errorf("No parsed field defined. no ")
	}
	appenderMap := make(map[string]AppenderSinglePathFunction)
	var theAppender []ByteBufferJSONParseResultAppender
	var idx = 0
	for _, def := range extractDefinitions {
		var bytAppender = def.Appender
		theAppender = append(theAppender, bytAppender)

		appenderMap[def.PathToExtract] = func(bytes []byte) (errLocal error) {
			return bytAppender.Append(bytes)
		}
		idx++
	}
	var markerForDeletionAct JSONItemRemover
	if len(removedItemPaths) > 0 {
		markerForDeletionAct = NewJSONItemRemover(removedItemPaths)
	}
	targetDataString := string(targetData)
	return markerForDeletionAct, parseJSONStringActualWorker(targetDataString, appenderMap, markerForDeletionAct)
}
