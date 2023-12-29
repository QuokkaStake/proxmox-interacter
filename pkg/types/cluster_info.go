package types

import (
	"fmt"
	"strings"
	"unicode"
)

type ClusterInfos []ClusterInfo

type ClusterInfo struct {
	Name  string
	Nodes []NodeWithAssets
	Error error
}

func ParseMatchers(query string) map[string]string {
	matchers := map[string]string{}

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	items := strings.FieldsFunc(query, f)

	for len(items) > 0 {
		item := strings.Split(items[0], "=")
		if len(item) <= 1 {
			matchers["name"] = strings.Join(items, " ")
			return matchers
		}

		matchers[item[0]] = strings.Trim(item[1], "\"")

		items = items[1:]
	}

	return matchers
}

func (c ClusterInfos) FindNode(query string) (*Node, bool) {
	for _, cluster := range c {
		if cluster.Error != nil {
			continue
		}

		for _, node := range cluster.Nodes {
			if node.Node.Node == query || node.Node.ID == query {
				return &node.Node, true
			}
		}
	}

	return nil, false
}

func (c ClusterInfos) FindContainer(query string) (*Container, string, error) {
	queryParsed := ParseMatchers(query)
	containerMatcher, err := NewContainerMatcher(queryParsed)
	if err != nil {
		return nil, "", err
	}

	for _, cluster := range c {
		if cluster.Error != nil {
			continue
		}

		// Taking the first container we can find matching the filter.
		for _, node := range cluster.Nodes {
			for _, container := range node.Containers {
				if container.Matches(containerMatcher) {
					return &container, cluster.Name, nil
				}
			}
		}
	}

	return nil, "", fmt.Errorf("Container is not found!")
}
