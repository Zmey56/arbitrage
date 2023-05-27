package withoutcardamong

import (
	"strings"
)

func CommonPair(pairOne, pairTwo []string) []string {
	var commonPair []string

	for _, pairO := range pairOne {
		for _, pairT := range pairTwo {
			if strings.ToLower(pairO) == strings.ToLower(pairT) {
				commonPair = append(commonPair, pairT)
			}
		}
	}
	return commonPair
}
