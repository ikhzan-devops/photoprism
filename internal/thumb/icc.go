package thumb

import (
	"os"
	"path"
)

/*
Possible TODO: move this into a shared pkg/ so non-thumb
consumers can also use it. However, it looks fiddly to hook that
up to `assets`, so I'm punting on that for now.
*/

func MustGetAdobeRGB1998Path() string {
	p := path.Join(IccProfilesPath, "compatibleWithAdobeRGB1998.icc")
	_, err := os.Stat(p)
	if err != nil {
		panic(err)
	}
	return p
}
