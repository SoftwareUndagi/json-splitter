package splitter

import (
	"os"
	"testing"
)

func TestByteBufferAppenderToFile(t *testing.T) {
	appender := NewByteBufferJSONParseResultAppender()
	appender.Append([]byte("dodol"))
	appender.Append([]byte("garut"))
	appender.WriteToFile(appendFilePath(os.TempDir(), "byte-to-file.txt"))

}


