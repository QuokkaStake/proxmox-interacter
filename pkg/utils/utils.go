package utils

import (
	"fmt"
	"html/template"
	"main/pkg/types"
	"math"
	"regexp"
	"strings"
	"time"
)

func NormalizeString(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.ToLower(reg.ReplaceAllString(input, ""))
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func Find[T any](slice []T, f func(T) bool) (*T, bool) {
	for _, e := range slice {
		if f(e) {
			return &e, true
		}
	}
	return nil, false
}

func Map[T, V any](slice []T, f func(T) V) []V {
	n := make([]V, len(slice))
	for index, e := range slice {
		n[index] = f(e)
	}
	return n
}

func SerializeQueryString(qs map[string]string) string {
	tmp := make([]string, len(qs))
	counter := 0

	for key, value := range qs {
		tmp[counter] = key + "=" + value
		counter++
	}

	return strings.Join(tmp, "&")
}

func MergeMaps(first, second map[string]string) map[string]string {
	for key, value := range second {
		first[key] = value
	}

	return first
}

func SerializeLink(link types.Link) template.HTML {
	if link.Href == "" {
		return template.HTML(link.Name)
	}

	return template.HTML(fmt.Sprintf(
		"<a href=\"%s\">%s</a>",
		link.Href,
		link.Name,
	))
}

func FormatSize(size int64) string {
	sizeFloat := float64(size)

	if size > 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f TB", sizeFloat/1024/1024/1024/1024)
	}

	if size > 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", sizeFloat/1024/1024/1024)
	}

	if size > 1024*1024 {
		return fmt.Sprintf("%.2f MB", sizeFloat/1024/1024)
	}

	if size > 1024 {
		return fmt.Sprintf("%.2f KB", sizeFloat/1024)
	}

	return fmt.Sprintf("%.2f B", sizeFloat)
}

func FormatDuration(durationInt int64) string {
	duration := time.Duration(durationInt * 1_000_000_000)

	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
