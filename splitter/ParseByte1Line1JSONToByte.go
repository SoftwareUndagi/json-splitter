package splitter

//ParseByte1Line1JSONToByte parse dari dalam byte( )
import (
	"bufio"
	"io"

	"github.com/sirupsen/logrus"
)

const parseByte1Line1JSONToByteClsName = "ParseByte1Line1JSONToByte"

//ParseByte1Line1JSONToByte parse json line by line dari bufio.Reader
func ParseByte1Line1JSONToByte(jsonSourceReader *bufio.Reader, subJSONMetadata []ParseByte1Line1JSONToByteBridgeSubJSONData, extractedDataPaths []string, cleanedResultAppender AppenderSinglePathFunction, customValueGenerator AdditionalJSONDataGenerator, logEntry *logrus.Entry) (err error) {
	loggerEntry := WrapLogWithClassAndMethod(nil, parseByte1Line1JSONToByteClsName, "ParseByte1Line1JSONToByte")
	//cleanedJSONAppender := NewByteBufferJSONParseResultAppender()
	var rowIndex int
	var plainRowIndex int
	eofFound := false
	for {
		rowData, errRead := jsonSourceReader.ReadBytes('\n')
		if errRead != nil {
			if errRead == io.EOF {
				eofFound = true
			} else {
				loggerEntry.WithError(err).Errorf("Error membaca buffer, row filled index [%d]. error: %s", plainRowIndex, err.Error())
				return
			}
		}
		plainRowIndex = plainRowIndex + 1
		if len(rowData) == 0 {
			//string kosong skip to next
			if eofFound {
				break
			}
			continue
		}
		cleanedUpResult, err := parseByte1Line1JSONToByteBridge(rowIndex, rowData, subJSONMetadata, extractedDataPaths, customValueGenerator, logEntry)
		if err != nil {
			return err
		}
		cleanedResultAppender(cleanedUpResult)
		rowIndex = rowIndex + 1
		if eofFound {
			break
		}
	}

	return
}

//ParseByte1Line1JSONToByteBridgeSubJSONData wrapper sub json data. path extract + modificator data( misal untuk menambahkan parent key ke dalam json data). ini untuk memudahkan parse sub data
type ParseByte1Line1JSONToByteBridgeSubJSONData struct {
	//SubJSONPath path untuk di extract
	SubJSONPath string
	//Modificator method untuk modify data.jika misal tidak ada handler maka pergunakan NoModificatorJSONDataHandler
	Modificator ModificatorJSONData
	//FinalAppender setelah selesai di clean up hasil di masukan ke appender mana
	FinalAppender AppenderSinglePathFunction
	//DoNotRemovePath kalau di flag true, data tidak akan di cleanup dari json
	DoNotRemovePath bool
}

//parseByte1Line1JSONToByteBridge worker parse 1 line. data di pecah menjadi subnya, dengan kembalian ada yang sudah di bersihkan. item yang hapus dari data json
// sesuai dengan parameter removedItemPaths.
// urutan data dalam subJsonData sesuai dengan urutan pada data subJSONPaths
// cleanedDataModifier(cek fungsi : GenerateCustomDataForCleanedData ) di pergunakan untuk mengeneratekan setiap data primary key misalnya. untuk kemudahan parsing data
//
func parseByte1Line1JSONToByteBridge(lineIndex int, jsonData []byte, subJSONMetadata []ParseByte1Line1JSONToByteBridgeSubJSONData, extractedDataPaths []string, customValueGenerator AdditionalJSONDataGenerator, logEntry *logrus.Entry) (cleanedUpResult []byte, err error) {
	var subPathExtractor []ByteBufferJSONParseResultAppender
	extractPathValueContainer := make(map[string][]byte)
	appenderMap := make(map[string]AppenderSinglePathFunction)
	var removedItemPaths []string
	for _, subPath := range subJSONMetadata {
		theAppend := NewByteBufferJSONParseResultAppender()
		subPathExtractor = append(subPathExtractor, theAppend)
		appenderMap[subPath.SubJSONPath] = theAppend.Append
		if !subPath.DoNotRemovePath {
			removedItemPaths = append(removedItemPaths, subPath.SubJSONPath)
		}
	}

	if extractedDataPaths != nil && len(extractedDataPaths) > 0 {
		for _, extractPath := range extractedDataPaths {
			extrFunc := func(bytes []byte) (err error) {
				extractPathValueContainer[extractPath] = bytes
				return nil
			}
			appenderMap[extractPath] = extrFunc
		}
	}
	var markerForDeletionAct JSONItemRemover
	if len(removedItemPaths) > 0 {
		markerForDeletionAct = NewJSONItemRemover(removedItemPaths)
	}
	errParse := parseJSONStringActualWorker(string(jsonData), appenderMap, markerForDeletionAct)

	if errParse != nil {
		err = errParse
		return
	}
	cleanedUpResult = markerForDeletionAct.MakeCleanedUpByte(jsonData)
	var generatedDataParameter map[string]interface{}
	if customValueGenerator != nil {
		var errGenCustomParam error
		cleanedUpResult, generatedDataParameter, errGenCustomParam = customValueGenerator(cleanedUpResult, lineIndex, extractPathValueContainer)
		if errGenCustomParam != nil {
			logEntry.WithError(errGenCustomParam).WithField("index", lineIndex).Errorf("Gagal generate custom value untuk index : %d, error: %s", lineIndex, errGenCustomParam.Error())
			err = errGenCustomParam
			return
		}

	}
	for i := 0; i < len(subJSONMetadata); i++ {
		tmpAppender := subPathExtractor[i]
		subJSONMetadata := subJSONMetadata[i]
		modifyJSONByteSendToActualAppender(lineIndex, subJSONMetadata, tmpAppender, extractPathValueContainer, generatedDataParameter, logEntry)
	}
	return
}

//modifyJSONByteSendToActualAppender modifikasi json. dan kirim ke actual appender
func modifyJSONByteSendToActualAppender(lineIndex int, subJSONMetadata ParseByte1Line1JSONToByteBridgeSubJSONData, temporaryAppender ByteBufferJSONParseResultAppender, extractPathValueContainer map[string][]byte, generatedDataParameter map[string]interface{}, logEntry *logrus.Entry) (err error) {
	if subJSONMetadata.Modificator == nil {
		return subJSONMetadata.FinalAppender(temporaryAppender.Bytes())
	}
	for {
		byteDt, err := temporaryAppender.ReadBytes()
		if err != nil {
			if err != io.EOF {
				logEntry.WithError(err).WithField("index", lineIndex).Errorf("Gagal read data dari buffer. error: %s", err.Error())
				return err
			}
			break
		}
		if len(byteDt) > 0 {
			cleanedResult, errClean := subJSONMetadata.Modificator(byteDt, lineIndex, extractPathValueContainer, generatedDataParameter)
			if errClean != nil {
				return errClean
			}
			subJSONMetadata.FinalAppender(cleanedResult)
		}
	}
	return nil
}
