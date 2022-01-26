package leopard

import "testing"

func TestInfo(t *testing.T) {
	Info("This should be an info log")
}

func TestDebug(t *testing.T) {
	Debug("This should be an debug log")
}

func TestWarning(t *testing.T) {
	Warning("This should be an warning log")
}

func TestError(t *testing.T) {
	Error("This should be an error log")
}
