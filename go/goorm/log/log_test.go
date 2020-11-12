package log

import (
	"os"
	"testing"
)

func TestSetLevel(t *testing.T) {
	SetLogLevel(ErrorLevel)

	if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("fail to set log level", ErrorLevel)
	}

	SetLogLevel(Disabled)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() == os.Stdout {
		t.Fatal("fail to set log level", Disabled)
	}
}
