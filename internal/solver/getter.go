package recipetree

import (
	"sort"
)

func GetAllElements() []string {
	var elements []string
	for name := range elementsMapGlobal {
		elements = append(elements, name)
	}
	sort.Strings(elements)
	return elements
}