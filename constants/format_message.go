package constants

import "strings"

func FormatErrorMessage(template string, vars map[string]string) string {
    for key, value := range vars {
        placeholder := "{" + key + "}"
        template = strings.ReplaceAll(template, placeholder, value)
    }
    return template
}