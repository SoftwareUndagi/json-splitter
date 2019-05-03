package splitter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func generateAppender(t *testing.T) func(path string, data []byte) (err error) {
	return func(path string, data []byte) (err error) {
		strResult := fmt.Sprintf("Untuk path [%s] menerima data [%s]\n", path, data) // string(data)
		// fmt.Printf(strResult)
		t.Log("path" + strResult)
		return nil
	}
}

//TestParseStringSequence tes string parse
func TestParseStringSequence(t *testing.T) {
	appender := generateAppender(t)
	t.Log("Starting TestParseStringSequence")
	var sample1 = `sampleData"`
	var sample2 = `sample\"dodol\"garut"`
	hasil1, _ := ParseStringSequence(sample1, len(sample1), "sample1", 0, appender)
	t.Logf("1. Selesai dengan index %d", hasil1)
	hasil2, _ := ParseStringSequence(sample2, len(sample2), "sample2", 0, appender)
	t.Logf("2. Selesai dengan index %d", hasil2)
}

//TestParseTrueFalseValueOK test boolean pattern dengan nilai di temukan
func TestParseTrueFalseValueOK(t *testing.T) {
	appender := generateAppender(t)
	withTrueVal := "true,\"key\":\"dodol\""
	withFalseVal := "sesuatufalse,\"key\":\"dodol\""
	got, _ := ParseTrueFalseValue(withTrueVal, len(withTrueVal), "sample1", 0, appender)
	if got == -1 {
		t.Errorf("ParseTrueFalseValue mendapat %d seharusnya bukan -1", got)
	}
	got2, _ := ParseTrueFalseValue(withFalseVal, len(withFalseVal), "sample2", 7, appender)
	if got == -1 {
		t.Errorf("ParseTrueFalseValue mendapat %d seharusnya bukan -1", got2)
	}

}
func TestParseNumberSequence(t *testing.T) {
	appender := generateAppender(t)
	withEndSpace := "123 "
	withEndComma := "4.567,"
	withEndKurawal := "86}"
	withEndArray := "45]"
	r1, _ := ParseNumberSequence(withEndSpace, len(withEndSpace), "pth1", 0, appender)
	r2, _ := ParseNumberSequence(withEndComma, len(withEndComma), "pth2", 0, appender)
	r3, _ := ParseNumberSequence(withEndKurawal, len(withEndKurawal), "pth3", 0, appender)
	r4, _ := ParseNumberSequence(withEndArray, len(withEndArray), "pth4", 0, appender)
	v := []int{r1, r2, r3, r4}
	for i := 0; i < len(v); i++ {
		if v[i] == -1 {
			t.Errorf("Index %d tidak sesuai. Silakan cek kembali", i)
		}
	}
}
func TestParseArray(t *testing.T) {
	appender := generateAppender(t)
	numberArray := `[1,144,55,67,89], "key2": { "name": "Gede Sutarsa"}`
	stringArray := `["dodol", "garut" ], "key2": { "name": "Gede Sutarsa"}`
	booleanArray := `[true, false, true ], "key2": { "name": "Gede Sutarsa"}`
	mixedArray := `[true,"dodol" , false,1888, true, { "name": "Gede Sutarsa"} ], "key2": { "name": "Gede Sutarsa"}`
	nIndex, err1 := ParseArrayData(numberArray, len(numberArray), "1", 1, appender, nil)
	if nIndex < 0 || err1 != nil {
		t.Error("Parse array 1 gagal")
	}
	nIndex2, err2 := ParseArrayData(stringArray, len(stringArray), "2", 1, appender, nil)
	if nIndex2 < 0 || err2 != nil {
		t.Error("Parse array 2 gagal")
	}
	nIndex3, err3 := ParseArrayData(booleanArray, len(booleanArray), "3", 1, appender, nil)
	if nIndex3 < 0 || err3 != nil {
		t.Error("Parse array 3 gagal")
	}
	nIndex4, err4 := ParseArrayData(mixedArray, len(mixedArray), "4", 1, appender, nil)
	if nIndex4 < 0 || err4 != nil {
		t.Error("Parse array 4 gagal")
	}
}

func TestParseObject(t *testing.T) {
	appender := generateAppender(t)
	objectSample := `{"key1": { "name": "Gede Sutarsa"} ], "key2": { "name": "Gede Sutarsa"}}`
	nIdx, err1 := ParseObjectData(objectSample, len(objectSample), "1", 10, appender, nil)
	if nIdx < 0 || err1 != nil {
		t.Error("Parse obj1 gagal")
	}

}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func TestReadFile(t *testing.T) {
	dir, errF := filepath.Abs(filepath.Dir(os.Args[0]))
	check(errF)
	dat, err := ioutil.ReadFile(dir + "/sample-data-small.json")
	check(err)
	fmt.Print(string(dat))
}

func TestParseJsonSmallFile(t *testing.T) {
	dir, errF := filepath.Abs(filepath.Dir(os.Args[0]))
	check(errF)
	dat, err := ioutil.ReadFile(dir + "/sample-data.json")
	check(err)
	t.Log("Parse data now")

	// t.Log()

	primaryResult, additionalResult, cleanedData, errParse := ParseJSONByteToByte(dat, "hotels", nil)

	// t.Log()

	if errParse != nil {
		t.Error("Fail to parse, ", errParse)
		return
	}
	t.Log(cleanedData)
	t.Log(additionalResult)
	ioutil.WriteFile("/tmp/sample-large-sample.txt", []byte(primaryResult), 0644)

}

func TestCheckSubstring(t *testing.T) {
	sampleArray := []byte(`alpha,omega`)
	result := sampleArray[5:len(sampleArray)]
	swapStr := string(result)
	t.Log(swapStr)

}

func TestParseJSONToFile(t *testing.T) {
	CaptureLog(t).Release()
	tmpDir := os.TempDir()
	destFile := appendFilePath(tmpDir, "hotel.json")
	dir, errF := filepath.Abs(filepath.Dir(os.Args[0]))
	logrus.Warn("Output dir: " + tmpDir)
	check(errF)
	deletionMarker, orgByte, err := ParseJSONFileToFile(dir+"/sample-data.json", []ExtractJSONToFileDefinition{ExtractJSONToFileDefinition{PathToExtract: "hotels", DestinationFilePath: destFile}}, []string{"hotels"})
	if err != nil {
		t.Error(err)
	}
	if deletionMarker != nil {
		cleandedByte := deletionMarker.MakeCleanedUpByte(orgByte)
		errWriteFile := ioutil.WriteFile(appendFilePath(tmpDir, "cleaned-result.json"), cleandedByte, 0644)
		if errWriteFile != nil {
			t.Error(errWriteFile)
		}
	}

	t.Log("selesai")
}
