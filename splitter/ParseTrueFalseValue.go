package splitter

import (
	"errors"

	"github.com/sirupsen/logrus"
)

//ParseTrueFalseValue parse boolean value
func ParseTrueFalseValue(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction) (nextIndex int, err error) {

	// var container [5]byte
	// var indexMarker = 0
	if jsonStringLength >= startIndex+4 && targetData[startIndex] == 't' && targetData[startIndex+1] == 'r' && targetData[startIndex+2] == 'u' && targetData[startIndex+3] == 'e' {
		errAppendTrue := appendToAppender(currentPath, sliceTrue)
		return startIndex + 4, errAppendTrue
	} else if jsonStringLength > startIndex+5 && targetData[startIndex] == 'f' && targetData[startIndex+1] == 'a' && targetData[startIndex+2] == 'l' && targetData[startIndex+3] == 's' && targetData[startIndex+4] == 'e' {
		errAppendFalse := appendToAppender(currentPath, sliceFalse)
		return startIndex + 5, errAppendFalse
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseTrueFalseValue",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to run true false value for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}

//ParseTrueFalseValueWithMap parse boolean value, appender di kirim dengan map( tunggal)
func ParseTrueFalseValueWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, appenderMap AppenderSinglePathMap) (nextIndex int, err error) {
	if jsonStringLength >= startIndex+4 && targetData[startIndex] == 't' && targetData[startIndex+1] == 'r' && targetData[startIndex+2] == 'u' && targetData[startIndex+3] == 'e' {
		appender := appenderMap[currentPath]
		if appender != nil {
			return startIndex + 4, appender(sliceTrue)
		}
		return startIndex + 4, nil
	} else if jsonStringLength > startIndex+5 && targetData[startIndex] == 'f' && targetData[startIndex+1] == 'a' && targetData[startIndex+2] == 'l' && targetData[startIndex+3] == 's' && targetData[startIndex+4] == 'e' {
		appender := appenderMap[currentPath]
		if appender != nil {
			return startIndex + 4, appender(sliceFalse)
		}
		return startIndex + 5, nil
	}
	WrapLogWithClassAndMethod(nil, "ParseTrueFalseValue", "ParseTrueFalseValueWithMap").WithFields(logrus.Fields{"startIndex": startIndex,
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to run true false value for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}
