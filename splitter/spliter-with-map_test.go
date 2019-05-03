package splitter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SoftwareUndagi/go-common-libs/common"
)

func generateSingleAppender(t *testing.T, path string) func(bytes []byte) (err error) {
	return func(bytes []byte) (err error) {
		vData := string(bytes)
		t.Log("Path: " + path + ", append: " + vData)
		return nil
	}
}

func generateSingleAppenderMap(t *testing.T, paths []string) AppenderSinglePathMap {
	rtvl := make(map[string]AppenderSinglePathFunction)
	for _, key := range paths {
		rtvl[key] = generateSingleAppender(t, key)
	}
	return rtvl
}

func TestBooleanScannerMap(t *testing.T) {

	appenderMap := generateSingleAppenderMap(t, []string{"sample1", "sample2"})
	withTrueVal := "true,\"key\":\"dodol\""
	withFalseVal := "sesuatufalse,\"key\":\"dodol\""
	got, _ := ParseTrueFalseValueWithMap(withTrueVal, len(withTrueVal), "sample1", 0, appenderMap)
	if got == -1 {
		t.Errorf("ParseTrueFalseValue mendapat %d seharusnya bukan -1", got)
	}
	got2, _ := ParseTrueFalseValueWithMap(withFalseVal, len(withFalseVal), "sample2", 7, appenderMap)
	if got == -1 {
		t.Errorf("ParseTrueFalseValue mendapat %d seharusnya bukan -1", got2)
	}
}

func TestParseStringSequenceWithMap(t *testing.T) {
	appender := generateSingleAppenderMap(t, []string{"sample1", "sample2"})
	t.Log("Starting TestParseStringSequence")
	var sample1 = `sampleData"`
	var sample2 = `sample\"dodol\"garut"`
	hasil1, _ := ParseStringSequenceWithMap(sample1, len(sample1), "sample1", 0, appender)
	t.Logf("1. Selesai dengan index %d", hasil1)
	hasil2, _ := ParseStringSequenceWithMap(sample2, len(sample2), "sample2", 0, appender)
	t.Logf("2. Selesai dengan index %d", hasil2)
}

func TestParseNumberSequenceWithMap(t *testing.T) {
	appender := generateSingleAppenderMap(t, []string{"pth1", "pth2", "pth3", "pth4"})
	withEndSpace := "123 "
	withEndComma := "4.567,"
	withEndKurawal := "86}"
	withEndArray := "45]"
	r1, _ := ParseNumberSequenceWithMap(withEndSpace, len(withEndSpace), "pth1", 0, appender)
	r2, _ := ParseNumberSequenceWithMap(withEndComma, len(withEndComma), "pth2", 0, appender)
	r3, _ := ParseNumberSequenceWithMap(withEndKurawal, len(withEndKurawal), "pth3", 0, appender)
	r4, _ := ParseNumberSequenceWithMap(withEndArray, len(withEndArray), "pth4", 0, appender)
	v := []int{r1, r2, r3, r4}
	for i := 0; i < len(v); i++ {
		if v[i] == -1 {
			t.Errorf("Index %d tidak sesuai. Silakan cek kembali", i)
		}
	}
}

func TestParseArrayWithMap(t *testing.T) {
	appender := generateSingleAppenderMap(t, []string{"1", "2", "3", "4"})
	numberArray := `[1,144,55,67,89], "key2": { "name": "Gede Sutarsa"}`
	stringArray := `["dodol", "garut" ], "key2": { "name": "Gede Sutarsa"}`
	booleanArray := `[true, false, true ], "key2": { "name": "Gede Sutarsa"}`
	mixedArray := `[true,"dodol" , false,1888, true, { "name": "Gede Sutarsa"} ], "key2": { "name": "Gede Sutarsa"}`
	nIndex, err1 := ParseArrayDataWithMap(numberArray, len(numberArray), "1", 1, appender, nil)
	if nIndex < 0 || err1 != nil {
		t.Error("Parse array 1 gagal")
	}
	nIndex2, err2 := ParseArrayDataWithMap(stringArray, len(stringArray), "2", 1, appender, nil)
	if nIndex2 < 0 || err2 != nil {
		t.Error("Parse array 2 gagal")
	}
	nIndex3, err3 := ParseArrayDataWithMap(booleanArray, len(booleanArray), "3", 1, appender, nil)
	if nIndex3 < 0 || err3 != nil {
		t.Error("Parse array 3 gagal")
	}
	nIndex4, err4 := ParseArrayDataWithMap(mixedArray, len(mixedArray), "4", 1, appender, nil)
	if nIndex4 < 0 || err4 != nil {
		t.Error("Parse array 4 gagal")
	}
}

func TestParseObjectWithMap(t *testing.T) {
	appender := generateSingleAppenderMap(t, []string{"key1", "key2"})
	objectSample := `{"key1": { "name": "Gede Sutarsa"} ], "key2": { "name": "Gede Sutarsa"}}`
	nIdx, err1 := ParseObjectDataWithMap(objectSample, len(objectSample), "key1", 10, 0, appender, nil)
	if nIdx < 0 || err1 != nil {
		t.Error("Parse obj1 gagal")
	}

}

func TestGenerateOsPath(t *testing.T) {
	x2 := fmt.Sprintf("dodol%sgarut", string(os.PathSeparator))
	println(x2)
}

func TestParseJSON1line1JSON(t *testing.T) {

	dir, errF := filepath.Abs(filepath.Dir(os.Args[0]))
	tmpDir := os.TempDir()
	check(errF)
	sprt := string(os.PathSeparator)
	skr := time.Now().UnixNano() / 1000000
	// tempDir := os.TempDir()
	outDir := common.AppendFilePath(tmpDir, fmt.Sprintf("parse_out%d", skr))
	sourceFile := fmt.Sprintf("%s%ssample-1line-1json.txt", dir, sprt)
	f := func(JSONPath string, JSONRowIndex int) string {
		return fmt.Sprintf("%s-%d.json", JSONPath, JSONRowIndex)
	}
	errParse := ParseJSONOnFile1Line1JSONToFile(outDir, sourceFile, "cleaned-1line-1json.txt",
		[]string{"rooms", "facilities", "images", "interestPoints", "phones", "terminals", "wildcards"},
		[]string{"rooms", "facilities", "images", "interestPoints", "phones", "terminals", "wildcards"}, f)
	if errParse != nil {
		t.Error(errParse)
	}

}

func TestParseJSONToSingleJson(t *testing.T) {
	CaptureLog(t).Release()
	tmpDir := os.TempDir()

	//dir, errF := filepath.Abs(filepath.Dir(os.Args[0]))
	//check(errF)
	//deletionMarker, orgByte, err := ParseJSONFileToFile(dir+"/sample-singlejson.json",
	deletionMarker, orgByte, err := ParseJSONFileToFile("C:\\Users\\gedesutarsa\\Documents\\go-projects-root\\src\\bham-server\\splitter\\sample-singlejson.json",
		[]ExtractJSONToFileDefinition{
			ExtractJSONToFileDefinition{PathToExtract: "rooms", DestinationFilePath: common.AppendFilePath(tmpDir, "rooms.json")},
			ExtractJSONToFileDefinition{PathToExtract: "facilities", DestinationFilePath: common.AppendFilePath(tmpDir, "facilities.json")},
			ExtractJSONToFileDefinition{PathToExtract: "images", DestinationFilePath: common.AppendFilePath(tmpDir, "images.json")},
			ExtractJSONToFileDefinition{PathToExtract: "interestPoints", DestinationFilePath: common.AppendFilePath(tmpDir, "interestPoints.json")},
			ExtractJSONToFileDefinition{PathToExtract: "phones", DestinationFilePath: common.AppendFilePath(tmpDir, "phones.json")},
			ExtractJSONToFileDefinition{PathToExtract: "terminals", DestinationFilePath: common.AppendFilePath(tmpDir, "terminals.json")},
			ExtractJSONToFileDefinition{PathToExtract: "wildcards", DestinationFilePath: common.AppendFilePath(tmpDir, "wildcards.json")}},
		[]string{"rooms", "facilities", "images", "interestPoints", "phones", "terminals", "wildcards"})
	if err != nil {
		t.Error(err)
	}
	if deletionMarker != nil {
		cleandedByte := deletionMarker.MakeCleanedUpByte(orgByte)
		errWriteFile := ioutil.WriteFile(common.AppendFilePath(tmpDir, "cleaned-result.json"), cleandedByte, 0644)
		if errWriteFile != nil {
			t.Error(errWriteFile)
		}
	}

	t.Log("selesai")
}
