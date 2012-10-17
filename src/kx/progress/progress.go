package progress

import (
	"os"
    "runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
    "sync"
)

const BAR = "="

type Progress struct {
    terminalWidth, total int
    bar string
    lock *sync.Mutex
}

// Progress constructor
func New(total int) *Progress {
    this := new(Progress)
    this.bar = BAR
    this.total = total
	this.terminalWidth = terminalWidth()
    this.lock = new(sync.Mutex)
    return this
}

// Set the progress bar to display
func (this Progress) SetBar(bar string) {
    this.bar = bar
}

// Show current progress bar
func (this Progress) ShowProgress(current int) {
    bar := progress(current, this.total, this.terminalWidth)
    os.Stdout.Write([]byte(bar + "\r"))
    os.Stdout.Sync()
}

func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func terminalWidth() int {
	winsize := getWinsize()
	return int(winsize.Col)
}

func progress(current, total, cols int) string {
	prefix := strconv.Itoa(current) + " / " + strconv.Itoa(total)
	bar_start := " ["
	bar_end := "] "

	bar_size := cols - len(prefix+bar_start+bar_end)
	amount := int(float32(current) / (float32(total) / float32(bar_size)))
	remain := bar_size - amount

	bar := strings.Repeat(BAR, amount) + strings.Repeat(" ", remain)
	return bold(prefix) + bar_start + bar + bar_end
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinsize() *winsize {
	ws := new(winsize)

    var TIOCGWINSZ int
    if runtime.GOOS == "darwin" {
        TIOCGWINSZ = 1074295912
    } else if runtime.GOOS == "linux" {
        TIOCGWINSZ = 0x5413
    } else {
        panic("Not supported platform")
    }

	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	return ws
}
