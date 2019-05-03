package splitter

import (
	"time"

	"github.com/sirupsen/logrus"
)

//WrapLogWithClassAndMethod add class + method to class
func WrapLogWithClassAndMethod(baseEntry *logrus.Entry, className string, methodName string) *logrus.Entry {
	markerFields := logrus.Fields{"className": className, "methodName": methodName}
	if baseEntry == nil {
		return logrus.WithFields(markerFields)
	}
	return baseEntry.WithFields(markerFields)
}

//WrapLogWithUsername menambahkan username dalam log
func WrapLogWithUsername(baseEntry *logrus.Entry, username string) *logrus.Entry {
	markerFields := logrus.Fields{"username": username}
	if baseEntry == nil {
		return logrus.WithFields(markerFields)
	}
	return baseEntry.WithFields(markerFields)
}

//EchoFunctionDuration print duration. ini untuk di masukan dalam deffer
func EchoFunctionDuration(logEntry *logrus.Entry, startTime time.Time) {
	t2 := time.Now()
	diff := t2.Sub(startTime)
	logEntry.Infof("Duration : %v", diff)
}
