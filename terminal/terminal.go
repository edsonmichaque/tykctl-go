package terminal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

// Terminal represents terminal capabilities
type Terminal struct {
	isTTY    bool
	Width    int
	Height   int
	Color    bool
	NoColor  bool
	ForceTTY bool
}

// New creates a new terminal instance
func New() *Terminal {
	return &Terminal{
		isTTY:    isTTY(),
		Width:    getWidth(),
		Height:   getHeight(),
		Color:    getColor(),
		NoColor:  getNoColor(),
		ForceTTY: getForceTTY(),
	}
}

// IsTTY checks if the output is a TTY
func isTTY() bool {
	// Check if stdout is a TTY
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// GetWidth returns the terminal width
func getWidth() int {
	if width := os.Getenv("COLUMNS"); width != "" {
		if w, err := strconv.Atoi(width); err == nil {
			return w
		}
	}
	return 80 // default width
}

// GetHeight returns the terminal height
func getHeight() int {
	if height := os.Getenv("LINES"); height != "" {
		if h, err := strconv.Atoi(height); err == nil {
			return h
		}
	}
	return 24 // default height
}

// GetColor returns whether color output is enabled
func getColor() bool {
	// Check for NO_COLOR environment variable
	if noColor := os.Getenv("NO_COLOR"); noColor != "" {
		return false
	}
	
	// Check for CLICOLOR environment variable
	if clicolor := os.Getenv("CLICOLOR"); clicolor != "" {
		if clicolor == "0" {
			return false
		}
		return true
	}
	
	// Check for CLICOLOR_FORCE environment variable
	if clicolorForce := os.Getenv("CLICOLOR_FORCE"); clicolorForce != "" {
		if clicolorForce == "1" {
			return true
		}
	}
	
	// Default to true if TTY
	return isTTY()
}

// GetNoColor returns whether color is explicitly disabled
func getNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

// GetForceTTY returns whether TTY is forced
func getForceTTY() bool {
	return os.Getenv("TYKCTL_FORCE_TTY") == "1"
}

// IsTTY returns whether the output is a TTY
func (t *Terminal) IsTTY() bool {
	if t.ForceTTY {
		return true
	}
	return t.isTTY
}

// GetWidth returns the terminal width
func (t *Terminal) GetWidth() int {
	return t.Width
}

// GetHeight returns the terminal height
func (t *Terminal) GetHeight() int {
	return t.Height
}

// SupportsColor returns whether color is supported
func (t *Terminal) SupportsColor() bool {
	if t.NoColor {
		return false
	}
	return t.Color && t.IsTTY()
}

// GetSize returns the terminal size
func (t *Terminal) GetSize() (width, height int) {
	return t.Width, t.Height
}

// IsInteractive returns whether the terminal is interactive
func (t *Terminal) IsInteractive() bool {
	return t.IsTTY() && t.SupportsColor()
}

// Write writes to the terminal with color support
func (t *Terminal) Write(data []byte) (int, error) {
	if t.SupportsColor() {
		return os.Stdout.Write(data)
	}
	// Strip ANSI color codes if color is not supported
	stripped := stripANSI(string(data))
	return os.Stdout.Write([]byte(stripped))
}

// WriteString writes a string to the terminal
func (t *Terminal) WriteString(s string) (int, error) {
	return t.Write([]byte(s))
}

// Printf prints formatted output to the terminal
func (t *Terminal) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stdout, format, args...)
}

// Println prints a line to the terminal
func (t *Terminal) Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stdout, args...)
}

// stripANSI removes ANSI color codes from a string
func stripANSI(s string) string {
	// Simple ANSI escape sequence removal
	// This is a basic implementation - for production use, consider a more robust solution
	for {
		start := strings.Index(s, "\x1b[")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "m")
		if end == -1 {
			break
		}
		s = s[:start] + s[start+end+1:]
	}
	return s
}

// Color constants
const (
	ColorReset  = "\x1b[0m"
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorBlue   = "\x1b[34m"
	ColorPurple = "\x1b[35m"
	ColorCyan   = "\x1b[36m"
	ColorWhite  = "\x1b[37m"
	ColorGray   = "\x1b[90m"
)

// Colorize applies color to text if color is supported
func (t *Terminal) Colorize(text, color string) string {
	if !t.SupportsColor() {
		return text
	}
	return color + text + ColorReset
}

// Red returns red colored text
func (t *Terminal) Red(text string) string {
	return t.Colorize(text, ColorRed)
}

// Green returns green colored text
func (t *Terminal) Green(text string) string {
	return t.Colorize(text, ColorGreen)
}

// Yellow returns yellow colored text
func (t *Terminal) Yellow(text string) string {
	return t.Colorize(text, ColorYellow)
}

// Blue returns blue colored text
func (t *Terminal) Blue(text string) string {
	return t.Colorize(text, ColorBlue)
}

// Purple returns purple colored text
func (t *Terminal) Purple(text string) string {
	return t.Colorize(text, ColorPurple)
}

// Cyan returns cyan colored text
func (t *Terminal) Cyan(text string) string {
	return t.Colorize(text, ColorCyan)
}

// Gray returns gray colored text
func (t *Terminal) Gray(text string) string {
	return t.Colorize(text, ColorGray)
}

// IsTerminal checks if the file descriptor is a terminal
func IsTerminal(fd uintptr) bool {
	return term.IsTerminal(fd)
}
