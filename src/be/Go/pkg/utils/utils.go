package utils

import (
	"strings"
)

func TitleToLink(title string) string {
	return "/wiki/" + strings.ReplaceAll(title, " ", "_")
}

func LinkToTitle(link string) string {
	return strings.ReplaceAll(strings.TrimPrefix(link, "/wiki/"), "_", " ")
}
