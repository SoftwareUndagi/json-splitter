package splitter

import (
	"fmt"
	"io/ioutil"
)

const parseJSONStringClassName = "ParseJSONString"

//flushDestinationFileSize sise autimatic flush tiap berapa baris
const flushDestinationFileSize = 50

//ExtractJSONToFileDefinition definsi extract json
type ExtractJSONToFileDefinition struct {
	//PathToExtract path di extract dari json string
	PathToExtract string
	//DestinationFilePath file tempat menulis hasil
	DestinationFilePath string
}

//ExtractJSONToByteDefinition definisi extract data to byte
type ExtractJSONToByteDefinition struct {
	//PathToExtract path di extract dari json string
	PathToExtract string
	//Appender appender hasil parse
	Appender ByteBufferJSONParseResultAppender
}

//ParseJSONFileToFile baca file json, dan parse json. hasil di output ke json
func ParseJSONFileToFile(jsonFilePath string, extractDefinitions []ExtractJSONToFileDefinition, removedItemPaths []string) (markerForDeletion JSONItemRemover, originalJSONData []byte, err error) {
	dat, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, dat, err
	}
	markerForDeletion, errParse := ParseJSONStringToFile(dat, extractDefinitions, removedItemPaths)
	return markerForDeletion, dat, errParse
}

//SubJSONFileNameGenerator generator sub data . untuk menulis data sub json yang di pecah dari induk
type SubJSONFileNameGenerator func(JSONPath string, JSONRowIndex int) string

//closeAllAppenders close all appender
func closeAllAppenders(appenders []JSONParseSingleResultAppender) {
	if appenders == nil {
		return
	}
	entry := WrapLogWithClassAndMethod(nil, "ParseJSONString", "closeAllAppender")
	var i = 0
	for _, appd := range appenders {
		errClose := appd.Close()
		if errClose != nil {
			entry.WithField("appenderIndex", i).Error(errClose)
		}
		i++
	}
}

//parseJSONStringActualWorker backbone. facade parser
func parseJSONStringActualWorker(targetData string, appenderMap AppenderSinglePathMap, markerForDeletion JSONItemRemover) (err error) {
	i := 0
	jsonStringLength := len(targetData)
	for {
		if i >= jsonStringLength {
			break
		}
		chr := targetData[i]
		if chr == ' ' || chr == '}' || chr == ',' || chr == '\n' || chr == '\r' {
			i++
			continue
		} else if chr == '[' {
			i, err = ParseArrayDataWithMap(targetData, jsonStringLength, "", i+1, appenderMap, markerForDeletion)
			if err != nil {
				return err
			}
		} else if chr == '{' {
			i, err = ParseObjectDataWithMap(targetData, jsonStringLength, "", i, 0, appenderMap, markerForDeletion)
			if err != nil {
				return err
			}
		} else {
			WrapLogWithClassAndMethod(nil, commonSpliterFileName, "ParseJSONString").Error("Pada root, menemukan char di luar spec")
			return fmt.Errorf("Error was encounter while parsing main json data")
		}
	}
	return nil
}

//parseJSONStringOld parse and split json data. deprecateditem
func parseJSONStringOld(targetData string, appendToAppender AppenderFunction) (cleanedResult string, err error) {
	i := 0
	jsonStringLength := len(targetData)
	markerDummy := func(jsonPath string, startIndex int, endIndex int) {
		// println("Markerdummy:", jsonPath, ".startindex:", startIndex, ".end index:", endIndex)
	}
	for {
		if i >= jsonStringLength {
			break
		}
		chr := targetData[i]
		if chr == ' ' || chr == '}' || chr == ',' {
			continue
		} else if chr == '[' {
			i, err = ParseArrayData(targetData, jsonStringLength, "", i+1, appendToAppender, markerDummy)
			if err != nil {
				return targetData, err
			}
		} else if chr == '{' {
			i, err = ParseObjectData(targetData, jsonStringLength, "", i, appendToAppender, markerDummy)
			if err != nil {
				return targetData, err
			}
		} else {
			WrapLogWithClassAndMethod(nil, commonSpliterFileName, "ParseJSONString").Error("Pada root, menemukan char di luar spec")
			return targetData, fmt.Errorf("Error was encounter while parsing main json data")
		}
	}
	return targetData, nil
}
