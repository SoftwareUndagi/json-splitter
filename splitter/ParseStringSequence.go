package splitter

import (
	"bytes"
	"errors"

	"github.com/sirupsen/logrus"
)

//ParseStringSequence parse string sequence
//jsonStringLength berisi length dari targetData. untuk optimasi karnea data cukup di hitung sekali saja
func ParseStringSequence(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction) (nextIndex int, err error) {
	var container bytes.Buffer
	var idxActualStart = startIndex
	if targetData[idxActualStart] == '"' {
		container.WriteByte('"')
		idxActualStart = startIndex + 1
	}
	for i := idxActualStart; i < jsonStringLength; i++ {
		chr := targetData[i]
		if chr == MarkerCharDoubleQuote {
			container.WriteByte(chr)
			errAPpend := appendToAppender(currentPath, container.Bytes())
			return i + 1, errAPpend
		} else if MarkerCharBackSlash == chr {
			container.WriteByte(chr)
			container.WriteByte(targetData[i+1])
			i++

		} else {
			container.WriteByte(chr)
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseStringSequence",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to run parse string for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}

//ParseStringSequenceWithMap parse string sequence dengan appender beruapa map
func ParseStringSequenceWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, appenderMap AppenderSinglePathMap) (nextIndex int, err error) {
	appender := appenderMap[currentPath]
	if appender == nil { // demi efisiensi, kalau tidak di scan, tidak perlu membuat container
		var idxActualStart = startIndex
		if targetData[idxActualStart] == '"' {
			idxActualStart = startIndex + 1
		}
		for i := idxActualStart; i < jsonStringLength; i++ {
			chr := targetData[i]
			if chr == MarkerCharDoubleQuote {
				return i + 1, nil
			} else if MarkerCharBackSlash == chr {
				i++
			}
		}
	} else {
		var container bytes.Buffer
		var idxActualStart = startIndex
		if targetData[idxActualStart] == '"' {
			container.WriteByte('"')
			idxActualStart = startIndex + 1
		}
		for i := idxActualStart; i < jsonStringLength; i++ {
			chr := targetData[i]
			if chr == MarkerCharDoubleQuote {
				container.WriteByte(chr)
				errAPpend := appender(container.Bytes())
				return i + 1, errAPpend
			} else if MarkerCharBackSlash == chr {
				container.WriteByte(chr)
				container.WriteByte(targetData[i+1])
				i++

			} else {
				container.WriteByte(chr)
			}
		}
	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseStringSequenceWithMap",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to run parse string for path: " + currentPath)
	return -1, errors.New("Parsed until end of file, giveup")
}
