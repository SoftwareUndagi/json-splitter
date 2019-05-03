package splitter

import (
	"bytes"
	"fmt"
	"strings"
)

//SimpleJSONTagOpenGenerator builder tag json. pada open di tambahkan tag
type SimpleJSONTagOpenGenerator struct {
	count  int
	buffer bytes.Buffer
}

//AppendSimpleNumber append number yang sudah dalam bentuk string
func (p *SimpleJSONTagOpenGenerator) AppendSimpleNumber(key string, value string) {
	if p.count > 0 {
		p.buffer.WriteByte(',')
	}
	p.buffer.WriteString(fmt.Sprintf(`"%s":%s`, key, value))
	p.count = p.count + 1
}

//AppendIntegerNumber add int number. int, int16,int32 di buat dengan interface untuk kemudahan
func (p *SimpleJSONTagOpenGenerator) AppendIntegerNumber(key string, value interface{}) {
	if p.count > 0 {
		p.buffer.WriteByte(',')
	}
	p.buffer.WriteString(fmt.Sprintf(`"%s":%d`, key, value))
	p.count = p.count + 1
}

//AppendFloatingNumber add floating number ke dalam json data
func (p *SimpleJSONTagOpenGenerator) AppendFloatingNumber(key string, value interface{}) {
	if p.count > 0 {
		p.buffer.WriteByte(',')
	}
	p.buffer.WriteString(fmt.Sprintf(`"%s":%f`, key, value))
	p.count = p.count + 1
}

//AppendString add string value
func (p *SimpleJSONTagOpenGenerator) AppendString(key string, value string) {
	if p.count > 0 {
		p.buffer.WriteByte(',')
	}
	if strings.Contains(value, `"`) {
		value = strings.ReplaceAll(value, `"`, "\\\"")
	}
	p.buffer.WriteString(fmt.Sprintf(`"%s":"%s"`, key, value))
	p.count = p.count + 1
}

//AppendBoolean add boolean value ke simple json tag
func (p *SimpleJSONTagOpenGenerator) AppendBoolean(key string, value bool) {
	if p.count > 0 {
		p.buffer.WriteByte(',')
	}

	p.buffer.WriteString(fmt.Sprintf(`"%s":"%t"`, key, value))
	p.count = p.count + 1
}

//AppendSimpleJSONOnStart add simple json di awal data.data di tambah dari apa yang sudah di append dalam struct
func (p *SimpleJSONTagOpenGenerator) AppendSimpleJSONOnStart(bytes []byte) []byte {
	rtvl := []byte("{")
	rtvl = append(rtvl, p.buffer.Bytes()...)
	rtvl = append(rtvl, ',')
	lenOfBytes := len(bytes)
	for i := 0; i < lenOfBytes; i++ {
		if bytes[i] == '{' {
			rtvl = append(rtvl, bytes[i+1:lenOfBytes]...)
			break
		}
	}
	return rtvl
}
