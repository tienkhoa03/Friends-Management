package utils

import "regexp"

func ExtractEmails(text string) []string {
	emailRegex := `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`
	re := regexp.MustCompile(emailRegex)
	return re.FindAllString(text, -1)
}
