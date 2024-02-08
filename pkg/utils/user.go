package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func UserInputWait() (bool, error) {

	var userInput string

	cCy := color.New(color.FgCyan)
	cMeg := color.New(color.FgMagenta).Add(color.Bold)
	for {
		cCy.Printf("Do you want to proced ? (Y/n): ")

		_, err := fmt.Scan(&userInput)
		if err != nil {
			fmt.Printf("\n Error reading input: %s \n", err.Error())
			return false, err
		}

		//Check user Input
		if userInput == "Y" {
			cCy.Printf("\nYou entered : %s\n", userInput)
			return true, nil
		} else if userInput == "n" {
			cCy.Printf("\nYou entered : %s\n", userInput)
			return false, nil
		} else {
			cMeg.Printf("\nInvalid input. Please enter Y or n.\n")
		}
	}
}
