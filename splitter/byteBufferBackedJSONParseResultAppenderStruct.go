package splitter

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type byteBufferJSONParseResultAppenderStruct struct {
	//resultContainer penampung hasil parsing
	resultContainer *bytes.Buffer
}

func (p *byteBufferJSONParseResultAppenderStruct) Bytes() []byte {
	return p.resultContainer.Bytes()
}

//Close close appender
func (p *byteBufferJSONParseResultAppenderStruct) Close() (err error) {

	return nil
}
func (p *byteBufferJSONParseResultAppenderStruct) OpenAppender() (err error) {

	return nil
}

//ReadBytes read byte dengan delimiter(\n)
func (p *byteBufferJSONParseResultAppenderStruct) ReadBytes() (line []byte, err error) {
	return p.resultContainer.ReadBytes('\n')
}

//WriteToFile write bytes to file. salin semua content ke file
func (p *byteBufferJSONParseResultAppenderStruct) WriteToFile(DestinationFilePath string) (err error) {
	return ioutil.WriteFile(DestinationFilePath, p.resultContainer.Bytes(), 0644)
}
func (p *byteBufferJSONParseResultAppenderStruct) Append(data []byte) (err error) {
	var actData = data
	if data != nil && data[len(data)-1] != '\n' {
		actData = append(data, '\n')
	}
	n, err := p.resultContainer.Write(actData)
	if err != nil {
		return err
	}
	if n != len(actData) {
		return fmt.Errorf("Appender data length:%d , write result %d", len(data), n)
	}
	return nil
}

//NewByteBufferJSONParseResultAppender generate new byte array driven appender
func NewByteBufferJSONParseResultAppender() ByteBufferJSONParseResultAppender {
	bytCntr := byteBufferJSONParseResultAppenderStruct{}
	bytCntr.resultContainer = &bytes.Buffer{}
	return &bytCntr //resultContainer: bytes.NewBuffer()}
}
