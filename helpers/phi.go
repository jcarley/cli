package helpers

import (
	"fmt"
	"os"
)

// PHIPrompt asks the user if they are eligible to download PHI. If accepted,
// the previous operation will continue. Otherwise the program will exit.
func PHIPrompt() {
	var answer string
	for {
		fmt.Println("This operation might result in PHI data being downloaded and decrypted to your local machine. By entering \"y\" at the prompt below, you warrant that you have the necessary privileges to view the data, have taken all necessary precautions to secure this data, and absolve Catalyze of any issues that might arise from its loss.")
		fmt.Print("Do you wish to proceed? (y/n) ")
		fmt.Scanln(&answer)
		fmt.Println("")
		if answer != "y" && answer != "n" {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if answer == "n" {
		fmt.Println("Exiting")
		os.Exit(1)
	}
}
