package utils

import (
	"fmt"
	"html/template"
	"main/pkg/types"
	"math"
	"strings"
	"time"

	"github.com/c2h5oh/datasize"
)

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
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

func IntToBool(value int) bool {
	return value != 0
}

func FormatSize(size uint64) string {
	return datasize.ByteSize(size).HumanReadable()
}

func FormatBool(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
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
