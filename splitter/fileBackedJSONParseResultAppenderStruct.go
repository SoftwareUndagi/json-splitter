package splitter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

//FileBackedJSONParseResultAppenderFlushSize ukuran flush ke file. default akan di untuk file baced appender
const FileBackedJSONParseResultAppenderFlushSize = 100

//fileBackedJSONParseResultAppenderStruct struct appender to file
type fileBackedJSONParseResultAppenderStruct struct {
	// file tujuan file di tulis
	destinationFile string
	//writer buffered writer to file
	writer *bufio.Writer
	//fileRef file sebagai reference untuk close
	fileRef *os.File
	//notFlushLineCount line yang blm di flush
	notFlushLineCount int
	//flushSize per berapa row data akan di flush ke phisiical file. default akan di isi dengan
	flushSize int
}

//SetFlushSize size flush
func (p *fileBackedJSONParseResultAppenderStruct) SetFlushSize(flushSize int) {
	p.flushSize = flushSize
}

//GetFlushSize ukuran flush
func (p *fileBackedJSONParseResultAppenderStruct) GetFlushSize() int {
	return p.flushSize
}

//NewFileBackedJSONParseResultAppender generate new file appender
func NewFileBackedJSONParseResultAppender(destinationPath string) FileJSONParseResultAppender {
	return &fileBackedJSONParseResultAppenderStruct{destinationFile: destinationPath, writer: nil, fileRef: nil, notFlushLineCount: 0, flushSize: FileBackedJSONParseResultAppenderFlushSize}
}

func (p *fileBackedJSONParseResultAppenderStruct) Append(data []byte) (err error) {
	p.notFlushLineCount++
	if p.fileRef == nil {
		tryOpenErr := p.OpenAppender()
		if tryOpenErr != nil {
			return tryOpenErr
		}
	}
	var actData = data
	if data != nil && data[len(data)-1] != '\n' {
		actData = append(data, '\n')
	}
	n, err := p.writer.Write(actData)
	if err != nil {
		return err
	}
	if n != len(actData) {
		return fmt.Errorf("Appended data length:%d , write result %d", len(data), n)
	}
	return nil
}

//runFlush run flush
func (p *fileBackedJSONParseResultAppenderStruct) runFlush() {
	if p.notFlushLineCount > 0 {
		p.writer.Flush()
	}
	p.notFlushLineCount = 0

}

//OpenAppender open file for appender
func (p *fileBackedJSONParseResultAppenderStruct) OpenAppender() (err error) {
	file, errOpen := os.Create(p.destinationFile)
	if errOpen != nil {
		logrus.WithField("filePath", p.destinationFile).WithField("struct", "fileBackedJSONParseResultAppenderStruct").WithField("method", "Close").Error(errOpen)
		return errOpen
	}
	p.fileRef = file
	p.writer = bufio.NewWriter(file)
	return nil
}

func (p *fileBackedJSONParseResultAppenderStruct) Close() (err error) {
	p.runFlush()
	errCloseFile := p.fileRef.Close()
	return errCloseFile
}

//DestinationFilePath getter file destination path dari json appender
func (p *fileBackedJSONParseResultAppenderStruct) DestinationFilePath() string {
	return p.destinationFile
}
