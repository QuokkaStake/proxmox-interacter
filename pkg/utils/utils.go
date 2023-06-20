package utils

import (
	"fmt"
	"html/template"
	"main/pkg/logger"
	"main/pkg/types"
	"regexp"
	"strconv"
	"strings"
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

func StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Str("value", s).Msg("Could not parse float")
	}

	return f
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
