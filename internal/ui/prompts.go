package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Confirm asks the user for a yes/no confirmation.
// Returns true if the user enters 'y' or 'yes' (case-insensitive).
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", color.YellowString(prompt))

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "n" || response == "no" {
			return false
		}

		if response == "y" || response == "yes" {
			return true
		}

		fmt.Println("Please enter 'y' or 'n'.")
	}
}

// Ask asks the user for a string input with a default value.
func Ask(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	displayPrompt := prompt
	if defaultValue != "" {
		displayPrompt = fmt.Sprintf("%s [%s]", prompt, defaultValue)
	}
	fmt.Printf("%s: ", color.CyanString(displayPrompt))

	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	response = strings.TrimSpace(response)
	if response == "" {
		return defaultValue
	}

	return response
}
