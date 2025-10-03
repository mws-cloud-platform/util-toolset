package golden

import (
	"flag"
)

var (
	updateFlag = flag.Bool("update", false, "update golden test files")
)

const (
	updateMessage = "Actual golden data differs from Expected one. Run with -update to see the diff"
)

func IsUpdate() bool {
	return *updateFlag
}
