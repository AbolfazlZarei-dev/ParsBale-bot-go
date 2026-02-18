package ParsBale

import "fmt"

func Bold(text string) string {
	return fmt.Sprintf(" *%s* ", text)
}

func Italic(text string) string {
	return fmt.Sprintf(" _%s_ ", text)
}

func Link(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func Code(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func PreText(text string) string {
	return fmt.Sprintf("```\n%s\n```", text)
}

func Spoiler(text, description string) string {
	return fmt.Sprintf("```‌[%s]‌%s‌```", text, description)
}
