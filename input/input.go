package input

import "strings"

// ParseYesNoResponse は入力文字列をYes/Noの真偽値に変換します
func ParseYesNoResponse(input string, defaultValue bool) bool {
	if input = strings.TrimSpace(strings.ToLower(input)); input == "" {
		return defaultValue
	}
	return input == "y" || input == "yes"
}

// FormatYesNoPrompt はYes/No質問のプロンプト文字列を生成します
func FormatYesNoPrompt(prompt string, defaultYes bool) string {
	yesNo := map[bool]string{true: "Y/n", false: "y/N"}[defaultYes]
	return prompt + " [" + yesNo + "]: "
}
