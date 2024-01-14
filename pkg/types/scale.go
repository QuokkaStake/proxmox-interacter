package types

import (
	"fmt"
	"strconv"

	"github.com/c2h5oh/datasize"
)

type ScaleMatcher struct {
	Name string
	Node string
	ID   string

	Memory int64
	Swap   int64
	CPU    int64
}

func NewScaleMatcher(matchers map[string]string) (ScaleMatcher, error) {
	matcher := ScaleMatcher{}

	for matcherKey, matcherValue := range matchers {
		if matcherKey == "node" {
			matcher.Node = matcherValue
		} else if matcherKey == "name" {
			matcher.Name = matcherValue
		} else if matcherKey == "id" {
			matcher.ID = matcherValue
		} else if matcherKey == "memory" {
			bytes, err := datasize.Parse([]byte(matcherValue))
			if err != nil {
				return matcher, fmt.Errorf("error parsing memory size: %s", err)
			}
			matcher.Memory = int64(bytes.Bytes())
		} else if matcherKey == "swap" {
			bytes, err := datasize.Parse([]byte(matcherValue))
			if err != nil {
				return matcher, fmt.Errorf("error parsing swap size: %s", err)
			}
			matcher.Swap = int64(bytes.Bytes())
		} else if matcherKey == "cpu" {
			cores, err := strconv.ParseInt(matcherValue, 10, 64)
			if err != nil {
				return matcher, fmt.Errorf("error parsing cores: %s", err)
			}
			matcher.CPU = cores
		} else {
			return matcher, fmt.Errorf(
				"expected one of the keys 'node', 'name', 'id', 'memory', 'swap', 'cpu', but got '%s'",
				matcherKey,
			)
		}
	}

	return matcher, nil
}

func (s ScaleMatcher) CPUChanged(c Container) bool {
	if s.CPU == 0 {
		return false
	}

	return s.CPU != c.MaxCPU
}

func (s ScaleMatcher) MemoryChanged(c Container) bool {
	if s.Memory == 0 {
		return false
	}

	return s.Memory != c.MaxMemory
}

func (s ScaleMatcher) SwapChanged(config *ContainerConfig) bool {
	if s.Swap == 0 || !config.SwapPresent {
		return false
	}

	return s.Swap != config.Swap
}

func (s ScaleMatcher) AnythingChanged(c Container) bool {
	return s.CPUChanged(c) || s.MemoryChanged(c)
}
