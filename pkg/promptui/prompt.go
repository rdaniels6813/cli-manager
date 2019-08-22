package promptui

type Prompter interface {
	PromptString(message string)
}
