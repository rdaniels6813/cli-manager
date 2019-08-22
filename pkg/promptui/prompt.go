package promptui

import (
	"github.com/manifoldco/promptui"
)

// Prompter allows a user to prompt for information from the user
type Prompter interface {
	PromptString(message string) (string, error)
	PromptPassword(message string) (string, error)
}

// CLIPrompter prompts the user for input from the CLI
type CLIPrompter struct{}

// PromptString prompts the user for a string using the message provided
func (c *CLIPrompter) PromptString(message string) (string, error) {
	prompt := promptui.Prompt{
		Label: message,
	}
	return prompt.Run()
}

// PromptPassword prompts the user for a password that will have it's input masked
func (c *CLIPrompter) PromptPassword(message string) (string, error) {
	prompt := promptui.Prompt{
		Label: message,
		Mask:  '*',
	}
	return prompt.Run()
}
