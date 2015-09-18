package helpers

import (
	"fmt"
	"os"
	"strings"
)

var validAnswers = map[string]bool{
	"y":   true,
	"yes": true,
	"n":   false,
	"no":  false,
}

// PHIPrompt asks the user if they are eligible to download PHI. If accepted,
// the previous operation will continue. Otherwise the program will exit.
func PHIPrompt() {
	var answer string
	for {
		fmt.Println("This operation might result in PHI data being downloaded and decrypted to your local machine. By entering \"y\" at the prompt below, you warrant that you have the necessary privileges to view the data, have taken all necessary precautions to secure this data, and absolve Catalyze of any issues that might arise from its loss.")
		fmt.Print("Do you wish to proceed? (y/n) ")
		fmt.Scanln(&answer)
		fmt.Println("")
		if _, contains := validAnswers[strings.ToLower(answer)]; !contains {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if !validAnswers[strings.ToLower(answer)] {
		fmt.Println("Exiting")
		os.Exit(1)
	}
}

// YesNoPrompt outputs a given message and waits for a user to answer `y/n`.
// If yes, flow continues as normal. If no, the program is quit. The given
// message SHOULD contain the string "(y/n)" or some other form of y/n
// indicating that the user needs to type in y or n. This method does not do
// that for you. The message will not have a new line appended to it. If you
// require a newline, add this to the given message.
func YesNoPrompt(message string) {
	var answer string
	for {
		fmt.Printf(message)
		fmt.Scanln(&answer)
		fmt.Println("")
		if _, contains := validAnswers[strings.ToLower(answer)]; !contains {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if !validAnswers[strings.ToLower(answer)] {
		fmt.Println("Exiting")
		os.Exit(1)
	}
}
