package goProgressBar

import (
	"fmt"
	"math"
	"strings"
)

type ProgressBar struct {
	total          int64     //total number of elements
	current        int64     //current elements
	percentage     float64   //cyrretn percentage
	char           string    //character/Symbol for ProgressBar
	totalBarSize   int64     //total no of character to depict 100%
	currentBarSize int64     //current number of Characters to represent bars
	spinner        [4]string //characters to show a spinner
	pulse          int64     //an increasing number to be used by spinner
}

func GetNewBar(totalElements int64, currentElements int64, symbol string, totalBarSize int64) ProgressBar {

	CurrentBar := int64(math.Ceil(float64(currentElements/totalElements) * 50))
	return ProgressBar{
		total:          totalElements,
		current:        currentElements,
		char:           symbol,
		percentage:     float64(currentElements/totalElements) * 100,
		totalBarSize:   totalBarSize,
		currentBarSize: CurrentBar,
		spinner:        [4]string{"|", "/", "—", "\\"},
		pulse:          0,
	}

}

func (bar *ProgressBar) Display(current int64) {

	//update the current status
	bar.current = current

	//calculate the new current percentage
	bar.percentage = (float64(bar.current) * 100.0) / float64(bar.total)

	//calculate the number of Symbols to display
	bar.currentBarSize = int64(bar.percentage * float64(bar.totalBarSize) / 100)

	//Generate the ProgressBar string
	currentBar := strings.Repeat(bar.char, int(bar.currentBarSize))

	//Generate the Format String
	fmtString := fmt.Sprintf("\r[%%-%ds] %%3.2f%%%% %%8d/%%d {%%s}", bar.totalBarSize)
	bar.pulse = (bar.pulse + 1) % 4
	fmt.Printf(fmtString, currentBar, bar.percentage, bar.current, bar.total, bar.spinner[bar.pulse])
}

func (bar *ProgressBar) End() {
	fmt.Println()
}
