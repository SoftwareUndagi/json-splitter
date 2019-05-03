package splitter

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSimpleJsonTagOpen(t *testing.T) {
	CaptureLog(t).Release()
	sampleJSON := []byte(`{"email":"gede.sutarsa@gmail.com"}`)
	appender := SimpleJSONTagOpenGenerator{}
	appender.AppendBoolean("boolVal", true)
	appender.AppendFloatingNumber("float1", 1.57)
	appender.AppendIntegerNumber("int1", 101)
	appender.AppendString("stringVal", "kuda")
	rslt := appender.AppendSimpleJSONOnStart(sampleJSON)
	logrus.Info(string(rslt))
}
