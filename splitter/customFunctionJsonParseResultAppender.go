package splitter

//customFunctionJSONParseResultAppender appender custom. append di handle dengan method custom
type customFunctionJSONParseResultAppender struct {
	//actualAppender actual appender yang bertugas append json data
	actualAppender *AppenderFunction
}
