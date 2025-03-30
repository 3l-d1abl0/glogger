package progressbar

import (
	"fmt"
	"math"
	"os"
	"strings"

	"golang.org/x/term"
)

type count struct {
	current int64 //current count
	total   int64 //total count
}

type ProgressBar struct {
	status1        count    //status1
	status2        count    //status2
	percentage     float64  //progress % based on status1
	char           string   //character/Symbol for ProgressBar
	totalBarSize   int64    //total no of character to depict 100%
	currentBarSize int64    //current number of Characters to represent bars
	spinner        []string //characters to show a spinner
	pulse          int64    //an increasing number to be used by spinner
	minWidthForBar int      //minimum terminal width to display the progress bar
	lastTermWidth  int      //last known terminal width
}

func GetNewBar(totalElements int64, currentElements int64, symbol string, totalBarSize int64, currentCount int64, totalCount int64) ProgressBar {
	// Prevent division by zero
	var currentBar int64 = 0
	var percentage float64 = 0

	if totalElements > 0 {
		currentBar = int64(math.Ceil(float64(currentElements) / float64(totalElements) * 50))
		percentage = float64(currentElements) / float64(totalElements) * 100
	}

	return ProgressBar{
		status1:        count{currentElements, totalElements},
		status2:        count{currentCount, totalCount},
		char:           symbol,
		percentage:     percentage,
		totalBarSize:   totalBarSize,
		currentBarSize: currentBar,
		//spinner:        [4]string{"|", "/", "—", "\\"},
		spinner: []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
		//spinner:        []string{"☀️", "🌤", "⛅️", "🌥", "☁️", "🌦", "🌧", "⛈", "🌩", "🌨", "❄️"},
		//spinner:        []string{"🌎", "🌍", "🌏"},
		//spinner:        []string{"🡹", "🡽", "🡺", "🡾", "🡻", "🡿", "🡸", "🡼"},
		pulse:          0,
		minWidthForBar: 80,  // Minimum width to display the progress bar
		lastTermWidth:  120, // Default terminal width if we can't detect it
	}
}

// getTerminalWidth returns the current terminal width or the last known width if it can't be detected
func (bar *ProgressBar) getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return bar.lastTermWidth // fallback to last known width
	}
	bar.lastTermWidth = width
	return width
}

func (bar *ProgressBar) Display(currentStatus1 int64, currentStatus2 int64, totalBytes int) {
	//update the current status
	bar.status1.current = currentStatus1

	// Prevent division by zero
	if bar.status1.total > 0 {
		//calculate the new current percentage
		bar.percentage = (float64(bar.status1.current) * 100.0) / float64(bar.status1.total)
	} else {
		bar.percentage = 0
	}

	//calculate the number of Symbols to display
	bar.currentBarSize = int64(bar.percentage * float64(bar.totalBarSize) / 100)

	// Get current terminal width
	termWidth := bar.getTerminalWidth()

	// Calculate values for display
	var currentVal float32 = float32(currentStatus2) / (1024 * 1024)
	var totalVal float32 = float32(bar.status2.total) / (1024 * 1024)
	var totalMBytesSec float32 = float32(totalBytes) / (1024 * 1024)

	// Update pulse for spinner
	bar.pulse = (bar.pulse + 1) % int64(len(bar.spinner))

	// Check if terminal is wide enough to display the progress bar
	if termWidth < bar.minWidthForBar {
		// Compact mode - show only progress details without the bar
		fmt.Printf("\r{%s} %3.2f%% %d/%d @%.2f MB/s [%.2f/%.2f MB]",
			bar.spinner[bar.pulse], bar.percentage, bar.status1.current,
			bar.status1.total, totalMBytesSec, currentVal, totalVal)
	} else {
		// Calculate dynamic bar size based on available terminal width
		// Reserve space for other elements (approx 60 chars for text, numbers, etc.)
		availableSpace := termWidth - 60
		dynamicBarSize := bar.totalBarSize

		if availableSpace > 10 && availableSpace < int(bar.totalBarSize) {
			dynamicBarSize = int64(availableSpace)
		}

		// Generate the progress bar string with dynamic size
		currentBar := strings.Repeat(bar.char, int(bar.currentBarSize))

		// Generate the format string with dynamic bar size
		fmtString := fmt.Sprintf("\r%%s [%%-%ds] %%3.2f%%%%  %%4d/%%d @%%2.2f MB/s [%%2.2f/%%2.2f MB]", dynamicBarSize)

		// Display the progress bar
		fmt.Printf(fmtString, bar.spinner[bar.pulse], currentBar, bar.percentage,
			bar.status1.current, bar.status1.total, totalMBytesSec, currentVal, totalVal)
	}
}

func (bar *ProgressBar) End() {
	fmt.Println()
}
