// Package hk_artemis_sdk hikvision artemis sdk
package hk_artemis_sdk

import (
	"github.com/emirpasic/gods/maps/treemap"
	"strings"
	"unicode"
)

func isBlankString(str string) bool {
	return removeStringSpace(str) == ""
}

func removeStringSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func isHeaderToSign(headerName string, signHeaderPrefixList []string) bool {
	if isBlankString(headerName) {
		return false
	} else {
		if strings.HasPrefix(headerName, "x-ca-") {
			return true
		} else {
			if len(signHeaderPrefixList) != 0 {
				for _, v := range signHeaderPrefixList {
					if strings.EqualFold(headerName, v) {
						return true
					}
				}
			}
			return false
		}
	}
}

func mapToTreeMap(data map[string]string) *treemap.Map {
	treeMap := treemap.NewWithStringComparator()
	for k := range data {
		treeMap.Put(k, data[k])
	}
	return treeMap
}


