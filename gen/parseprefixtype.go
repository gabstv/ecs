package gen

import "strings"

func ParsePrefixType(v string) (prefix, mtype string) {
	bprefix := make([]rune, 0)
	btype := make([]rune, 0)
	r := []rune(v)
	inprefix := false
	for i := len(r) - 1; i >= 0; i-- {
		if !(r[i] >= 'a' && r[i] <= 'z') && !(r[i] >= 'A' && r[i] <= 'Z') && !(r[i] >= '0' && r[i] <= '9') {
			if !inprefix {
				inprefix = true
			}
		}
		if !inprefix {
			btype = append([]rune{r[i]}, btype...)
		} else {
			bprefix = append([]rune{r[i]}, bprefix...)
		}
	}
	prefix = string(bprefix)
	mtype = string(btype)
	return
}

func LinePrefixMatch(line string, prefixes ...string) (rest string, match bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, prefix) {
			rest = line[len(prefix):]
			match = true
			return
		}
	}
	rest = line
	match = false
	return
}
