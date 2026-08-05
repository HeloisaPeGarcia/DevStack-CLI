package ui_test

import (
	"bytes"
	"devstack/pkg/ui"
	"os"
	"strings"
	"testing"
)

func TestPrintStatusBadge_Output(t *testing.T) {
	buf := new(bytes.Buffer)
	ui.SetOutput(buf)
	defer ui.SetOutput(os.Stdout)

	ui.PrintStatusBadge("Go SDK", "go version go1.22", true)
	out := buf.String()

	if !strings.Contains(out, "Go SDK") {
		t.Errorf("esperava conter 'Go SDK', obteve: %s", out)
	}
	if !strings.Contains(out, "[INSTALADO]") {
		t.Errorf("esperava conter '[INSTALADO]', obteve: %s", out)
	}
}

func TestPrintStatusBadge_NoColorEnv(t *testing.T) {
	buf := new(bytes.Buffer)
	ui.SetOutput(buf)
	defer ui.SetOutput(os.Stdout)

	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	ui.PrintStatusBadge("Docker", "Faltando", false)
	out := buf.String()

	if strings.Contains(out, "\033[") {
		t.Errorf("não esperava códigos ANSI quando NO_COLOR está ativo, obteve: %s", out)
	}
}
