package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var defaultWriter io.Writer = os.Stdout

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func SetOutput(w io.Writer) {
	if w != nil {
		defaultWriter = w
	}
}

func isColorDisabled() bool {
	return os.Getenv("NO_COLOR") != ""
}

func colorize(code, text string) string {
	if isColorDisabled() {
		return text
	}
	return code + text + Reset
}

func PrintBanner() {
	title := colorize(Cyan+Bold, "🤖 DevStack CLI v1.0.0")
	subtitle := colorize(Bold, "AI-Powered Dev-Environment & Winget Automation")
	width := 64

	fmt.Fprintln(defaultWriter, colorize(Cyan, "┌"+strings.Repeat("─", width)+"┐"))
	fmt.Fprintf(defaultWriter, "%s  %-78s %s\n", colorize(Cyan, "│"), title, colorize(Cyan, "│"))
	fmt.Fprintf(defaultWriter, "%s  %-62s %s\n", colorize(Cyan, "│"), subtitle, colorize(Cyan, "│"))
	fmt.Fprintln(defaultWriter, colorize(Cyan, "└"+strings.Repeat("─", width)+"┘"))
}

func PrintHeader(title string) {
	fmt.Fprintf(defaultWriter, "\n%s=== %s ===%s\n\n", colorize(Cyan+Bold, ""), title, "")
}

func PrintSuccess(msg string) {
	fmt.Fprintf(defaultWriter, "%s %s\n", colorize(Green+Bold, "✔"), msg)
}

func PrintWarning(msg string) {
	fmt.Fprintf(defaultWriter, "%s %s\n", colorize(Yellow+Bold, "⚠"), msg)
}

func PrintError(msg string) {
	fmt.Fprintf(defaultWriter, "%s %s\n", colorize(Red+Bold, "✖"), msg)
}

func PrintInfo(msg string) {
	fmt.Fprintf(defaultWriter, "%s %s\n", colorize(Blue+Bold, "ℹ"), msg)
}

func PrintStatusBadge(label, status string, isInstalled bool) {
	badge := colorize(Green, "[INSTALADO]")
	if !isInstalled {
		badge = colorize(Yellow, "[FALTANDO]")
	}
	fmt.Fprintf(defaultWriter, "  • %-25s %s (%s)\n", label, badge, status)
}

func PrintBox(title string, lines []string) {
	width := 60
	fmt.Fprintln(defaultWriter, colorize(Cyan, "┌"+strings.Repeat("─", width)+"┐"))
	fmt.Fprintf(defaultWriter, "%s %-58s %s\n", colorize(Cyan, "│"), colorize(Bold, title), colorize(Cyan, "│"))
	fmt.Fprintln(defaultWriter, colorize(Cyan, "├"+strings.Repeat("─", width)+"┤"))
	for _, l := range lines {
		fmt.Fprintf(defaultWriter, "%s %-58s %s\n", colorize(Cyan, "│"), l, colorize(Cyan, "│"))
	}
	fmt.Fprintln(defaultWriter, colorize(Cyan, "└"+strings.Repeat("─", width)+"┘"))
}
