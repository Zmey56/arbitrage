package Interexchange

import (
	"strings"
)

func CheckMatchesPair(asset string, pair_assets, coin []string) []string {
	tmp := []string{}
	for _, pair := range pair_assets {
		if strings.HasPrefix(strings.ToLower(pair), strings.ToLower(asset)) {
			//log.Println("Start", strings.ToLower(pair[len(asset):]))
			if findslice(pair[len(asset):], coin) {
				tmp = append(tmp, pair)
			}
		} else if strings.HasSuffix(strings.ToLower(pair), strings.ToLower(asset)) {
			//log.Println("END", strings.ToLower(pair[:(len(pair)-len(asset))]), findslice(pair[:(len(pair)-len(asset))], coin))
			if findslice(pair[:(len(pair)-len(asset))], coin) {
				tmp = append(tmp, pair)
			}
		}
	}
	return tmp
}

func findslice(l string, sl []string) bool {
	for _, j := range sl {
		if strings.ToLower(l) == strings.ToLower(j) {
			return true
		}
	}
	return false
}
