package splitter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

//commonSpliterFileName nama file common-splitter-tool
const commonSpliterFileName = "common-spliter-tool"

//MarkerCharDoubleDot char/ byte :( titik dua)
const MarkerCharDoubleDot = ':'

//MarkerCharDoubleQuote char "( petik dua)
const MarkerCharDoubleQuote = '"'

//MarkerCharBackSlash string \(garis miringbalik)
const MarkerCharBackSlash = '\\'

//MarkerComa mrker koma
const MarkerComa = ','

//MarkerCloseBrace marker }
const MarkerCloseBrace = '}'

//MarkerCloseArray penanda tutup
const MarkerCloseArray = ']'

//MarkerSpace space
const MarkerSpace = ' '

//MarkerSlashR penanda akhir baris (\r)
const MarkerSlashR = '\r'

//MarkerSlashN penanda selesai \n
const MarkerSlashN = '\n'

// untuk compare dengan byte
var comparatorTrue = [5]byte{'t', 'r', 'u', 'e'}
var sliceTrue = []byte{'t', 'r', 'u', 'e'}

// checker false
var comparatorFalse = [5]byte{'f', 'a', 'l', 's', 'e'}
var sliceFalse = []byte{'f', 'a', 'l', 's', 'e'}

//AppenderFunction definition of appender function
type AppenderFunction func(path string, bytes []byte) (err error)

//GenerateCustomDataForCleanedData cleaned data mungkin bisa di buatkan generated data. misal id untuk data bisa di generate dengan uuid. data ini bisa di inject ke dalam json di bersihkan pada bagian awal
// sehingga child data sudah bisa tahu id dari parent data. akan memudahkan untuk populate data, tanpa perlu melakukan select ke dalam database
//lineIndex index dari data di proses
type GenerateCustomDataForCleanedData func(lineIndex int, rawJSONData []byte) (generatedValues map[string][]byte, err error)

//ModificatorJSONData method untuk transform byte json menjadi json dengan tambahan tertentu. misal tambahan key dari parent di masukan ke dalam detail
// parameter parameterData di dapat dari parsig data json. sesuai dengan definisi parsing . ini akan di kirimkan dalam setiap parsing
// parameter data generated value di generate dengan code. misal id dari data
type ModificatorJSONData func(bytes []byte, lineIndex int, parameterData map[string][]byte, generateDataParameter map[string]interface{}) (result []byte, err error)

//AdditionalJSONDataGenerator generator data tambahan utnuk json. misal data sudah di cleanup, di generatekan ID dari data. ini untuk memudahkan transfer data dari dalam model hierarchy(master detail)
//bytes = data dta json mentah, in case perlu parse json
//lineIndex = line index dari data di proses
//parameterData = data hasil parse dari detail
//JSONDataWithGeneratedValue = json data setelah di tambahi byte yang di perlukan
type AdditionalJSONDataGenerator func(bytes []byte, lineIndex int, parameterData map[string][]byte) (JSONDataWithGeneratedValue []byte, generateDataParameter map[string]interface{}, err error)

//NoModificatorJSONDataHandler default handler untuk ModificatorJSONData. tanpa ada modifikasi sama sekali
func NoModificatorJSONDataHandler(bytes []byte, parameterData map[string][]byte) (result []byte, err error) {
	result = bytes
	return
}

//AppenderSinglePathFunction appender single path. ini untuk di masukan ke dalam map. untuk 1 appender saja
type AppenderSinglePathFunction func(bytes []byte) (err error)

//AppenderSinglePathMap alias untuk map[string]AppenderSinglePathFunction
type AppenderSinglePathMap map[string]AppenderSinglePathFunction

//JSONResultAppender inteface untuk append saja
type JSONResultAppender interface {

	//Append add data to appender
	Append(path string, bytes []byte) (err error)
}

//JSONParseSingleResultAppender result. bisa byte atau
type JSONParseSingleResultAppender interface {
	//Append add data to appender
	Append(bytes []byte) (err error)
	//Close close appender kalau ada
	Close() (err error)
	//OpenAppender open appender. and start writing
	OpenAppender() (err error)
}

//ByteBufferJSONParseResultAppender appender dengan backend file
type ByteBufferJSONParseResultAppender interface {
	JSONParseSingleResultAppender
	//Bytes get bytes writed
	Bytes() []byte
	//WriteToFile write bytes to file. salin semua content ke file
	WriteToFile(DestinationFilePath string) (err error)
	//ReadBytes read byte dengan delimiter(\n)
	ReadBytes() (line []byte, err error)
}

//FileJSONParseResultAppender appender dengan backend file
type FileJSONParseResultAppender interface {
	JSONParseSingleResultAppender
	DestinationFilePath() string
	//SetFlushSize size flush
	SetFlushSize(flushSize int)
	//GetFlushSize ukuran flush
	GetFlushSize() int
}

//MarkArrayForIndexOfDeletion interface untuk menandai data untuk di remove. Use case nya
// sample json:
// { level1: { children1: ["child1" , "child2"  , "child3"]}}
// ide nya dalam json string, children1 mau di takeout ke file yang berbeda, implementasi method ini : jsonPath=children1 , startIndex akan di isi dengan 11(start of children1 index 11) , endIndex = 56( closing array)
type MarkArrayForIndexOfDeletion func(jsonPath string, startIndex int, endIndex int)

//CleanupMethod mehtod untuk cleanup
type CleanupMethod func()

//isNumberChar check byte number atau bukan
func isNumberChar(chr byte) bool {
	return chr == '0' || chr == '1' || chr == '2' || chr == '3' || chr == '4' || chr == '5' || chr == '6' || chr == '7' || chr == '8' || chr == '9'
}

//GenerateCleanupAppander cleanup appender json
func GenerateCleanupAppander(logEntry *logrus.Entry, appenders ...JSONParseSingleResultAppender) (handler CleanupMethod) {
	handler = func() {
		i := 0
		for _, appender := range appenders {
			errCls := appender.Close()
			if errCls != nil {
				logEntry.WithError(errCls).WithField("index", i).Errorf("Error close index %d, error %s", i, errCls.Error())
			}
			i++
		}
	}
	return
}

//FindCommaToPrev mencari koma dari index di tentukan ke belakang
func FindCommaToPrev(targetData string, startScanIndex int) int {
	for i := startScanIndex - 1; i >= 0; i-- {
		chr := targetData[i]
		if chr == MarkerComa {
			return i
		}
		if chr == MarkerSpace {
			return i
		}
	}
	return startScanIndex
}

func anomalyEcho(functionName string, targetData []byte, jsonStringLength int, currentPath string, startIndex int) {
	if len(targetData) > 18554 && targetData[18554] == 13 {
		var msg = fmt.Sprintf("Anomali [%s]. Start index : %d", functionName, startIndex)
		var theJSON = string(targetData[startIndex:jsonStringLength])
		println(msg)
		println(theJSON)
	}
}
