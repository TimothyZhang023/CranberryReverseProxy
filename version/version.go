package version

import (
	"fmt"
)

const (
	Major = "0"
	Minor = "1"
	Build = "1"
)

func Version() string {
	return fmt.Sprintf("%s.%s.%s", Major, Minor, Build)
}
