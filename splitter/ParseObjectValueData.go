package splitter

import (
	"errors"

	"github.com/sirupsen/logrus"
)

//ParseObjectValueData worker untuk membaca data value. string, number boolean or array
func ParseObjectValueData(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction, markerForDeletion MarkArrayForIndexOfDeletion) (nextIndex int, err error) {

	i := startIndex
	for {
		chr := targetData[i]

		if chr == ' ' {
			// kosong do nothing :D
		} else if chr == '[' {
			parseObjectValueDataCharReaded1, err := ParseArrayData(targetData, jsonStringLength, currentPath, i+1, appendToAppender, markerForDeletion)
			return parseObjectValueDataCharReaded1, err
		} else if chr == '"' {
			parseStringRslt, errStr := ParseStringSequence(targetData, jsonStringLength, currentPath, i, appendToAppender)
			return parseStringRslt, errStr
		} else if chr == '{' {
			parseObjRslt, errObj := ParseObjectData(targetData, jsonStringLength, currentPath, i, appendToAppender, markerForDeletion)
			return parseObjRslt, errObj
		} else if chr == 't' || chr == 'f' {
			parseBoolRslt, errBool := ParseTrueFalseValue(targetData, jsonStringLength, currentPath, i, appendToAppender)
			return parseBoolRslt, errBool
		} else if isNumberChar(chr) {
			parseNumberRslt, errNumber := ParseNumberSequence(targetData, jsonStringLength, currentPath, i, appendToAppender)
			return parseNumberRslt, errNumber
		}
		i++
		if i >= jsonStringLength {
			break
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseObjectValueData",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse  value object for path: " + currentPath)
	return -1, errors.New("Parse until end of string. no match. giveup")
}

//ParseObjectValueDataWithMap parse object value dengan map
func ParseObjectValueDataWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, appenderMap AppenderSinglePathMap, markerForDeletion JSONItemRemover) (nextIndex int, err error) {

	i := startIndex
	for {
		chr := targetData[i]

		if chr == ' ' {
			// kosong do nothing :D
		} else if chr == '[' {
			parseObjectValueDataCharReaded1, err := ParseArrayDataWithMap(targetData, jsonStringLength, currentPath, i+1, appenderMap, markerForDeletion)
			return parseObjectValueDataCharReaded1, err
		} else if chr == '"' {
			if i+1 >= jsonStringLength {
				println("Di sini sumber maasalah")
			}
			parseStringRslt, errStr := ParseStringSequenceWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			return parseStringRslt, errStr
		} else if chr == '{' {
			parseObjRslt, errObj := ParseObjectDataWithMap(targetData, jsonStringLength, currentPath, i, 0, appenderMap, markerForDeletion)
			return parseObjRslt, errObj
		} else if chr == 't' || chr == 'f' {
			parseBoolRslt, errBool := ParseTrueFalseValueWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			return parseBoolRslt, errBool
		} else if isNumberChar(chr) {
			parseNumberRslt, errNumber := ParseNumberSequenceWithMap(targetData, jsonStringLength, currentPath, i, appenderMap)
			return parseNumberRslt, errNumber
		}
		i++
		if i >= jsonStringLength {
			break
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseObjectValueDataWithMap",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse  value object for path: " + currentPath)
	return -1, errors.New("Parse until end of string. no match. giveup")
}
