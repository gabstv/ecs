package rx

import "aletheia.icu/broccoli/fs"

//go:generate broccoli -var=rx -src templates -include *.tmpl

func FS() *fs.Broccoli {
	return rx
}
