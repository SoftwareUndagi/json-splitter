package splitter

import (
	"bytes"
	"errors"

	"github.com/sirupsen/logrus"
)

//ParseNumberSequence parse number sequence
func ParseNumberSequence(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction) (nextIndex int, err error) {
	var container bytes.Buffer
	for i := startIndex; i < jsonStringLength; i++ {
		chr := targetData[i]
		if chr == MarkerComa || chr == MarkerCloseArray || chr == MarkerSpace || chr == MarkerCloseBrace {
			errAppend := appendToAppender(currentPath, container.Bytes())
			if errAppend != nil {
				return -1, errAppend
			}
			if chr == MarkerCloseArray || chr == MarkerCloseBrace {
				return i, nil
			}
			return i + 1, nil
		}
		container.WriteByte(chr)
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseNumberSequence",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse number value for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}

//ParseNumberSequenceWithMap parse nmber sequence dengan map of appender
func ParseNumberSequenceWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, appenderMap AppenderSinglePathMap) (nextIndex int, err error) {
	appender := appenderMap[currentPath]
	if appender == nil {
		for i := startIndex; i < jsonStringLength; i++ {
			chr := targetData[i]
			if chr == MarkerComa || chr == MarkerCloseArray || chr == MarkerSpace || chr == MarkerCloseBrace {
				if chr == MarkerCloseArray || chr == MarkerCloseBrace {
					return i, nil
				}
				return i + 1, nil
			}
		}
	} else {
		var container bytes.Buffer
		for i := startIndex; i < jsonStringLength; i++ {
			chr := targetData[i]
			if chr == MarkerComa || chr == MarkerCloseArray || chr == MarkerSpace || chr == MarkerCloseBrace {
				errAppend := appender(container.Bytes())
				if errAppend != nil {
					return -1, errAppend
				}
				if chr == MarkerCloseArray || chr == MarkerCloseBrace {
					return i, nil
				}
				return i + 1, nil
			}
			container.WriteByte(chr)
		}
	}

	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseNumberSequenceWithMap",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse number value for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}
