package splitter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestModifyJSONSampler(t *testing.T) {
	sampleData := []byte(`{"code":233146,"name":"Tanaya Bed & Breakfast"_+}`)
	var modificator ModificatorJSONData
	modificator = func(bytes []byte, lineIndex int, parameterData map[string][]byte, generatedData map[string]interface{}) (result []byte, err error) {
		appendedJSON := fmt.Sprintf(",\"hotelCode\":\"%s\",\"hotelName\": \"%s\" }", parameterData["hotelCode"], parameterData["hotelName"])
		dpn := bytes[0 : len(bytes)-1]
		blk := []byte(appendedJSON)
		result = append(dpn, blk...)
		return
	}
	CaptureLog(t).Release()
	mapParam := make(map[string][]byte)
	mapParam["hotelCode"] = []byte("12345")
	mapParam["hotelName"] = []byte("Sample Dodol")

	mdf, errMdf := modificator(sampleData, 0, mapParam, nil)
	if errMdf != nil {
		logrus.WithError(errMdf).Error(errMdf)
	}
	logrus.Infof("Result :: %s", string(mdf))
}

func TestModifySliceByIndex(t *testing.T) {
	CaptureLog(t).Release()
	dataProcess := []byte("Dodol garut")
	for i := len(dataProcess); i < len(dataProcess)+10; i++ {
		dataProcess[i] = 'a'
	}
	logrus.Infof("Result :: [%s]", string(dataProcess))
}

func TestParseLineByLineWithModify(t *testing.T) {
	logEntry := logrus.WithField("method", "TestParseLineByLineWithModify")
	thePath := "C:\\Users\\gedesutarsa\\Documents\\go-projects-root\\src\\bham-server\\splitter\\sample-1line-1json.txt"
	file, err := os.Open(thePath)
	if err != nil {
		t.Error(err)
		return
	}
	//  func(bytes []byte, lineIndex int, parameterData map[string][]byte, generateDataParameter map[string]interface{}) (result []byte, err error)
	facModifier := func(bytes []byte, lineIndex int, parameterData map[string][]byte, generatedData map[string]interface{}) (result []byte, err error) {
		sD := SimpleJSONTagOpenGenerator{}
		sD.AppendString("hotelCode", string(parameterData["code"]))
		result = sD.AppendSimpleJSONOnStart(bytes)
		return
	}
	imgModifier := func(bytes []byte, lineIndex int, parameterData map[string][]byte, generatedData map[string]interface{}) (result []byte, err error) {
		sD := SimpleJSONTagOpenGenerator{}
		sD.AppendString("hotelCode", string(parameterData["code"]))
		result = sD.AppendSimpleJSONOnStart(bytes)
		return
	}
	facAppender := NewFileBackedJSONParseResultAppender(appendFilePath(os.TempDir(), "facilities.json"))
	imagesAppender := NewFileBackedJSONParseResultAppender(appendFilePath(os.TempDir(), "images.json"))
	finalAppender := NewFileBackedJSONParseResultAppender(appendFilePath(os.TempDir(), "cleaned-data.json"))

	facDef := ParseByte1Line1JSONToByteBridgeSubJSONData{SubJSONPath: "facilities", FinalAppender: facAppender.Append, Modificator: facModifier}
	imgDef := ParseByte1Line1JSONToByteBridgeSubJSONData{SubJSONPath: "images", FinalAppender: imagesAppender.Append, Modificator: imgModifier}
	reader := bufio.NewReader(file)
	cleanupFunc := GenerateCleanupAppander(logEntry, facAppender, finalAppender, imagesAppender)
	defer cleanupFunc()
	// facAppender.OpenAppender()
	// finalAppender.OpenAppender()
	errParse := ParseByte1Line1JSONToByte(reader, []ParseByte1Line1JSONToByteBridgeSubJSONData{facDef, imgDef}, []string{"code"}, finalAppender.Append, nil, logEntry)
	if errParse != nil {
		logrus.Error(errParse)
	}

}

func TestReplaceString(t *testing.T) {
	CaptureLog(t).Release()
	sample := `Sample dengan petik`

	x := strings.ReplaceAll(sample, `"`, "\\\"")
	logrus.Info(x)

}
