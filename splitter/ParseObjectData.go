package splitter

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

func findDoublQuote(targetData string, jsonStringLength int, iVal int) (retResult int, retdoubleDotFound bool, retkeyComplete bool, retLocalKey string) {
	iLocal := iVal
	j := iVal
	localKey := strings.Builder{}
	doubleDotFound := false
	var keyComplete bool //= false
	for {
		chr2 := targetData[j]
		if !keyComplete {
			if chr2 == MarkerCharDoubleQuote {
				keyComplete = true
			} else {
				localKey.WriteByte(chr2)
			}
		}
		if chr2 == MarkerCharDoubleDot {
			doubleDotFound = true
			break
		}
		j++
		iLocal++
		if j >= jsonStringLength {
			break
		}
	}
	return iLocal, doubleDotFound, keyComplete, localKey.String()
}

//ParseObjectData parse JSON object
func ParseObjectData(targetData string, jsonStringLength int, currentPath string, startIndex int, appendToAppender AppenderFunction, markerForDeletion MarkArrayForIndexOfDeletion) (nextIndex int, err error) {

	var i int
	if targetData[startIndex] == '{' {
		i = startIndex + 1
	} else {
		i = startIndex
	}
	for {
		chr := targetData[i]
		if chr == ' ' || chr == ',' || chr == '\n' || chr == '\r' {
			i++
			continue
		} else if chr == '}' { // chr == }(kurung tutup)
			subArr := targetData[startIndex : i+1]
			err := appendToAppender(currentPath, []byte(subArr))
			return i + 1, err
		} else if chr == '"' {
			//localStartIndex := i
			iRs, doubleDotFound, _, localKey := findDoublQuote(targetData, jsonStringLength, i+1)
			i = iRs
			if !doubleDotFound {
				logrus.WithFields(logrus.Fields{
					"json":        string(targetData),
					"startIndex":  startIndex,
					"currentPath": currentPath}).Error("Gagal parse object. : tidak di temukan. return -1")
				return -1, errors.New("Parse sampai akhir data. : tidak di temukan")
			}
			nextPath := localKey
			if len(currentPath) > 0 {
				nextPath = currentPath + "." + localKey
			}
			iSwap, errParseValue := ParseObjectValueData(targetData, jsonStringLength, nextPath, i+1, appendToAppender, markerForDeletion)
			if errParseValue != nil {
				return iSwap, errParseValue
			}

			if markerForDeletion != nil {
				markerForDeletion(nextPath, FindCommaToPrev(targetData, startIndex), i+1)
			}
			i = iSwap

		} else {
			i++
		}

		//all done keluar
		if i >= jsonStringLength {
			break
		}

	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseObjectData",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse  object for path: " + currentPath)
	return -1, errors.New("Parse sampai akhir data. parse object tidak sesuai")
}

//ParseObjectDataWithMap parse object
//indexOnArray kalau di di panggil dalam array. ini akan di isi. ini untuk mencegah menghapus keluar dari array
func ParseObjectDataWithMap(targetData string, jsonStringLength int, currentPath string, startIndex int, indexOnArray int, appenderMap AppenderSinglePathMap, markerForDeletion JSONItemRemover) (nextIndex int, err error) {

	var i int
	if targetData[startIndex] == '{' {
		i = startIndex + 1
	} else {
		i = startIndex
	}
	for {
		chr := targetData[i]
		if chr == ' ' || chr == ',' || chr == '\n' || chr == '\r' {
			i++
			continue
		} else if chr == '}' { // chr == }(kurung tutup)
			if markerForDeletion != nil { // tandai ini item untuk di hapus jg
				if markerForDeletion.IsRemovedPath(currentPath) {
					var startRemoveIndex = startIndex
					if indexOnArray > 0 {
						startRemoveIndex = FindCommaToPrev(targetData, startIndex)
					}
					markerForDeletion.AddRangeToRemove(startRemoveIndex, i)
				}
			}
			subArr := targetData[startIndex : i+1]
			theAppender := appenderMap[currentPath]
			if theAppender != nil {
				errApnd := theAppender([]byte(subArr))
				return i + 1, errApnd
			}
			return i + 1, nil
		} else if chr == '"' {
			//localStartIndex := i
			iRs, doubleDotFound, _, localKey := findDoublQuote(targetData, jsonStringLength, i+1)
			i = iRs
			if !doubleDotFound {
				logrus.WithFields(logrus.Fields{
					"json":        string(targetData),
					"startIndex":  startIndex,
					"currentPath": currentPath}).Error("Gagal parse object. : tidak di temukan. return -1")
				return -1, errors.New("Parse sampai akhir data. : tidak di temukan")
			}
			nextPath := localKey
			if len(currentPath) > 0 {
				nextPath = currentPath + "." + localKey
			}
			//logrus.WithFields(logrus.Fields{"nextPath": nextPath, "index": i + 1}).Info("ParseObjectValueDataWithMap starting")
			iSwap, errParseValue := ParseObjectValueDataWithMap(targetData, jsonStringLength, nextPath, i+1, appenderMap, markerForDeletion)
			if errParseValue != nil {
				return iSwap, errParseValue
			}

			i = iSwap

		} else {
			i++
		}

		//all done keluar
		if i >= jsonStringLength {
			break
		}

	}
	logrus.WithFields(logrus.Fields{"startIndex": startIndex, "method": "ParseObjectData",
		"stringToParse": string(targetData[startIndex:(jsonStringLength - 1)])}).Error("Fail to parse  object for path: " + currentPath)
	return -1, errors.New("Parse sampai akhir data. parse object tidak sesuai")
}
