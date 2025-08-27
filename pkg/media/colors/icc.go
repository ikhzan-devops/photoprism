package colors

import (
	_ "embed"
	"fmt"
	"os"
	"path"
)

// this was pulled from the Debian package
// icc-profiles-free, where it is published with
// an MIT-style licence:
// https://salsa.debian.org/debian/icc-profiles-free/-/blob/master/icc-profiles-openicc/default_profiles/base/compatibleWithAdobeRGB1998.icc?ref_type=heads
//
//go:embed icc/compatibleWithAdobeRGB1998.icc
var compatibleWithAdobeRGB1998 []byte

var compatibleWithAdobeRGB1998Path = ""

var temporaryDirectory string = ""

func getTemporaryDirectory() (string, error) {
	if temporaryDirectory != "" {
		return temporaryDirectory, nil
	}
	var err error
	temporaryDirectory, err = os.MkdirTemp("", "photoprism-icc-")
	if err != nil {
		return "", fmt.Errorf("%w creating temp dir for ICC profiles", err)
	}
	return temporaryDirectory, nil
}

func GetAdobeRGB1998Path() (p string, err error) {
	if compatibleWithAdobeRGB1998Path != "" {
		return compatibleWithAdobeRGB1998Path, nil
	}

	tempDir, err := getTemporaryDirectory()
	if err != nil {
		return "", err
	}
	p = path.Join(tempDir, "compatibleWithAdobeRGB1998.icc")

	err = os.WriteFile(p, compatibleWithAdobeRGB1998, 0644)
	if err != nil {
		return "", fmt.Errorf("%w writing icc profile to temp dir", err)
	}
	compatibleWithAdobeRGB1998Path = p
	return p, nil
}
