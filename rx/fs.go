package rx

import "github.com/montanaflynn/broccoli/fs"

//go:generate broccoli -var=rx -src templates -include *.tmpl

func FS() *fs.Broccoli {
	return rx
}
