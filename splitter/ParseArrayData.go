package splitter

import (
	"errors"

	"github.com/sirupsen/logrus"
)

//ParseArrayData parse array data
func ParseArrayData(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction, markerForDeletion MarkArrayForIndexOfDeletion) (nextIndex int, err error) {
	i := startIndex

	for {
		chr := targetData[i]
		var charReaded int //
		charReaded = i + 1
		if chr == ']' {
			return i + 1, nil
		} else if chr == '[' {
			charReaded1, errArray := ParseArrayData(targetData, jsonStringLength, currentPath, i+1, appendToAppender, markerForDeletion)
			if errArray != nil {
				return -1, errArray
			}
			charReaded = charReaded1
		} else if chr == '"' {
			charReaded2, errSeq := ParseStringSequence(targetData, jsonStringLength, currentPath, i, appendToAppender)
			if errSeq != nil {
				return -1, errSeq
			}
			charReaded = charReaded2
		} else if chr == '{' {
			charReaded3, err := ParseObjectData(targetData, jsonStringLength, currentPath, i, appendToAppender, markerForDeletion)
			if err != nil {
				return -1, err
			}
			charReaded = charReaded3
		} else if chr == 't' || chr == 'f' {
			charReaded4, errTrueFalse := ParseTrueFalseValue(targetData, jsonStringLength, currentPath, i, appendToAppender)
			if errTrueFalse != nil {
				return -1, err
			}
			charReaded = charReaded4
		} else if isNumberChar(chr) {
			charReaded5, errNumber := ParseNumberSequence(targetData, jsonStringLength, currentPath, i, appendToAppender)
			if errNumber != nil {
				return -1, err
			}
			charReaded = charReaded5
		}
		if charReaded == -1 {
			return -1, errors.New("Sub proccess return -1")
		}
		i = charReaded
		if i >= jsonStringLength {
			break
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseArrayData",
		"stringToParse": string(targetData[startIndex:(jsonStringLength)])}).Error("Fail to parse array value for path: " + currentPath + ", pass until end of string no no close array found")
	return -1, errors.New("Parse until end of string. no match. giveup")
}

//ParseArrayDataWithMap proses array dengan
func ParseArrayDataWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, appenderMap AppenderSinglePathMap, markerForDeletion JSONItemRemover) (nextIndex int, err error) {
	i := startIndex
	arrayIndex := 0
	for {
		chr := targetData[i]
		var charReaded int //
		charReaded = i + 1
		if chr == ',' {
			arrayIndex++
		} else if chr == ']' {
			return i + 1, nil
		} else if chr == '[' {
			charReaded1, errArray := ParseArrayDataWithMap(targetData, jsonStringLength, currentPath, i+1, appenderMap, markerForDeletion)
			if errArray != nil {
				return -1, errArray
			}
			charReaded = charReaded1
		} else if chr == '"' {
			charReaded2, errSeq := ParseStringSequenceWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			if errSeq != nil {
				return -1, errSeq
			}
			charReaded = charReaded2
		} else if chr == '{' {
			charReaded3, err := ParseObjectDataWithMap(targetData, jsonStringLength, currentPath, i, arrayIndex, appenderMap, markerForDeletion)
			if err != nil {
				return -1, err
			}
			charReaded = charReaded3
		} else if chr == 't' || chr == 'f' {
			charReaded4, errTrueFalse := ParseTrueFalseValueWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			if errTrueFalse != nil {
				return -1, err
			}
			charReaded = charReaded4
		} else if isNumberChar(chr) {
			charReaded5, errNumber := ParseNumberSequenceWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			if errNumber != nil {
				return -1, err
			}
			charReaded = charReaded5
		}
		if charReaded == -1 {
			return -1, errors.New("Sub proccess return -1")
		}
		i = charReaded
		if i >= jsonStringLength {
			break
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseArrayDataWithMap",
		"stringToParse": string(targetData[startIndex:(jsonStringLength)])}).Error("Fail to parse array value for path: " + currentPath + ", pass until end of string no no close array found")
	return -1, errors.New("Parse until end of string. no match. giveup")
}
