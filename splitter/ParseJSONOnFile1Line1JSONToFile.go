package splitter

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

const parseJSONOnFile1Line1JSONToFileClsName = "ParseJSONOnFile1Line1JSONToFile"

//ParseJSONOnFile1Line1JSONToFile parse file json. 1 line dalam file berisi 1 JSON data. data di pecah menjadi sub json data
//outputDirectory directory tenpat hasil akan di tulis
//sourceJSONFile file json asal untuk di pecah
//removedItemPaths path yang akan di clean up dari per json data dalam file
//cleanedDestinationFileName nama file hasil di bersihkan dari sub json
//removedItemPaths path yang di hapus
func ParseJSONOnFile1Line1JSONToFile(outputDirectory string, sourceJSONFile string, cleanedDestinationFileName string, subJSONPaths []string, removedItemPaths []string, subJSONFilenameGenerator SubJSONFileNameGenerator) (err error) {
	loggerEntry := WrapLogWithClassAndMethod(nil, parseJSONOnFile1Line1JSONToFileClsName, "ParseJSONOnFile1Line1JSONToFile").WithFields(logrus.Fields{"fileName": sourceJSONFile, "destinationDirectory": outputDirectory})
	file, errOpenFile := os.Open(sourceJSONFile)
	if errOpenFile != nil {
		loggerEntry.Error("Failed to open JSON file.", errOpenFile)
		return errOpenFile
	}
	defer file.Close() // cleanup. clouse source file

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		loggerEntry.Info(fmt.Sprintf("Destination directory [%s] not exists, creating now", outputDirectory))
		errMkdir := os.Mkdir(outputDirectory, os.ModePerm)
		if errMkdir != nil {
			loggerEntry.Error("Fail to create destination directory, error ", errMkdir)
			return errMkdir
		}
	}
	fileDest, errOpenDest := os.Create(fmt.Sprintf("%s%s%s", outputDirectory, string(os.PathSeparator), cleanedDestinationFileName))
	if errOpenDest != nil {
		loggerEntry.WithField("destinationFile", cleanedDestinationFileName).Error("Error create file for destination , error: ", errOpenDest)
		return errOpenDest
	}
	defer fileDest.Close()
	destWriter := bufio.NewWriter(fileDest)
	reader := bufio.NewReader(file)
	var errReadLine = io.EOF
	var lineReaded = 0
	var flushInvoked = false

	var cleanedByte []byte
	var eofFound bool
	// var isPrefix bool
	for {
		// loggerEntry.WithField("index", lineReaded).Warn("Memproses baris:", lineReaded)
		cleanedByte, eofFound, errReadLine = parseJSONOnFile1Line1JSONToFile1LineBridge(lineReaded, reader, outputDirectory, subJSONPaths, removedItemPaths, subJSONFilenameGenerator)

		// If we're just at the EOF, break
		if errReadLine != nil {
			loggerEntry.WithField("index", lineReaded).WithError(errReadLine).Errorf("Fail to read line on index: %d", lineReaded)
			return errReadLine
		}
		writedData, errWrite := destWriter.Write(append(cleanedByte))
		flushInvoked = false
		if errWrite != nil {
			loggerEntry.WithError(errWrite).WithField("lineIndex", lineReaded).Error("Fail to write cleaned result, line: ", lineReaded)
			return errWrite
		}
		if writedData != len(cleanedByte) {
			loggerEntry.WithField("lineIndex", lineReaded).Error(fmt.Sprintf("Write to file inconsistent. return of write request [%d], data to write request: [%d] ", writedData, len(cleanedByte)+1))
		}
		if lineReaded%flushDestinationFileSize == 0 {
			flushInvoked = true
			errFlush := destWriter.Flush()
			if errFlush != nil {
				loggerEntry.WithError(errFlush).WithField("lineIndex", lineReaded).Error("Error pada saat flush destination file dalam line loop")
			}
		}
		lineReaded++
		if eofFound {
			break
		}
	}
	if !flushInvoked {
		return destWriter.Flush()
	}
	return nil
}

//parseJSONOnFile1Line1JSONToFile1LineBridge method helper. memproses 1 baris hasil dari ParseJSONOnFile1Line1JSONToFile
// pecah menjadi sub json dan kembalikan data dengan cleanup
func parseJSONOnFile1Line1JSONToFile1LineBridge(lineIndex int, reader *bufio.Reader, outputDirectory string, subJSONPaths []string, removedItemPaths []string, subJSONFilenameGenerator SubJSONFileNameGenerator) (cleanedJSONLineData []byte, isEOF bool, err error) {
	loggerEntry := WrapLogWithClassAndMethod(nil, parseJSONOnFile1Line1JSONToFileClsName, "parseJSONOnFile1Line1JSONToFile1LineBridge").WithFields(logrus.Fields{"lineIndex": lineIndex, "destinationDirectory": outputDirectory})

	/*lineReaded,*/
	l, errReadLine := reader.ReadBytes('\n')
	// l := []byte(lineReaded)
	if errReadLine != nil {
		if errReadLine == io.EOF {
			return l, true, nil
		}
		loggerEntry.WithError(errReadLine).WithField("index", lineIndex).Error("Red line index failed", errReadLine)
		return l, false, errReadLine
	}
	var extractDefinitions []ExtractJSONToFileDefinition
	for _, subPath := range subJSONPaths {
		fileDestName := fmt.Sprintf("%s%s%s", outputDirectory, string(os.PathSeparator), subJSONFilenameGenerator(subPath, lineIndex))
		extractDefinitions = append(extractDefinitions, ExtractJSONToFileDefinition{PathToExtract: subPath, DestinationFilePath: fileDestName})
	}
	markerForDeletionThisLine, errParse := ParseJSONStringToFile(l, extractDefinitions, removedItemPaths)
	if errParse != nil {
		return l, false, errParse
	}
	if markerForDeletionThisLine != nil {
		return markerForDeletionThisLine.MakeCleanedUpByte(l), false, nil
	}
	return l, false, nil

}
