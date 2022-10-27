package log

import (
	"fmt"
	"testing"
)

func TestSetLevel(t *testing.T) {
	fmt.Println("test info level")
	SetLevel(InfoLevel)
	Info("test info")
	Infof("test infof: %s", "test")
	Error("test error")
	Errorf("test errorf: %s", "test")
	SetLevel(ErrorLevel)
	fmt.Println("test error level")
	Info("test info")
	Infof("test infof: %s", "test")
	Error("test error")
	Errorf("test errorf: %s", "test")
	fmt.Println("test disabled level")
	SetLevel(Disabled)
	Info("test info")
	Infof("test infof: %s", "test")
	Error("test error")
	Errorf("test errorf: %s", "test")
}
