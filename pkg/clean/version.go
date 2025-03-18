package clean

import (
	"fmt"
	"strings"

	"github.com/photoprism/photoprism/pkg/txt"
)

// Version parses and returns a semantic version string.
func Version(s string) string {
	if s == "" {
		return ""
	}

	if strings.Contains(s, ":") {
		split := strings.Split(s, ":")

		if len(split) > 1 {
			s = split[1]
		}
	}

	if v := strings.Split(s, "."); len(v) < 3 {
		return ""
	} else {
		patch, _, _ := strings.Cut(v[2], "+")
		return fmt.Sprintf("v%d.%d.%d", txt.UInt(Numeric(v[0])), txt.UInt(Numeric(v[1])), txt.UInt(patch))
	}
}
