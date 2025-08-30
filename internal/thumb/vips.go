package thumb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// Vips generates a new thumbnail with the requested size and returns the file name and a buffer with the image bytes,
// or an error if thumbnail generation failed. For more information on libvips, see https://github.com/libvips/libvips.
func Vips(imageName string, imageBuffer []byte, hash, thumbPath string, width, height int, opts ...ResampleOption) (thumbName string, thumbBuffer []byte, err error) {
	if len(hash) < 4 {
		return "", nil, fmt.Errorf("thumb: invalid file hash %s", clean.Log(hash))
	}

	if len(imageName) < 4 {
		return "", nil, fmt.Errorf("thumb: invalid file name %s", clean.Log(imageName))
	}

	if InvalidSize(width) {
		return "", nil, fmt.Errorf("thumb: width has an invalid value (%d)", width)
	}

	if InvalidSize(height) {
		return "", nil, fmt.Errorf("thumb: height has an invalid value (%d)", height)
	}

	// Get thumb cache filename.
	thumbName, err = FileName(hash, thumbPath, width, height, opts...)

	if err != nil {
		log.Debugf("thumb: %s in %s (filename)", err, clean.Log(filepath.Base(imageName)))
		return "", nil, err
	}

	// Initialize libvips before using it.
	VipsInit()

	// Load image from file or buffer.
	var img *vips.ImageRef

	if len(imageBuffer) == 0 {
		if img, err = vips.LoadImageFromFile(imageName, VipsImportParams()); err != nil {
			log.Debugf("vips: %s in %s (load image from file)", err, clean.Log(filepath.Base(imageName)))
			return "", nil, err
		}
	} else if img, err = vips.LoadImageFromBuffer(imageBuffer, VipsImportParams()); err != nil {
		log.Debugf("vips: %s in %s (load image from buffer)", err, clean.Log(filepath.Base(imageName)))
		return "", nil, err
	}

	// Set resample options.
	var method ResampleOption
	var size vips.Size

	method, _, _ = ResampleOptions(opts...)

	// Choose thumbnail crop.
	var crop vips.Interesting
	if method == ResampleFillTopLeft {
		crop = vips.InterestingLow
		size = vips.SizeBoth
	} else if method == ResampleFillBottomRight {
		crop = vips.InterestingHigh
		size = vips.SizeBoth
	} else if method == ResampleFit {
		crop = vips.InterestingNone
		size = vips.SizeDown
	} else if method == ResampleFillCenter || method == ResampleResize {
		crop = vips.InterestingCentre
		size = vips.SizeBoth
	}

	if err = vipsSetIccProfileForInteropIndex(img, clean.Log(filepath.Base(imageName))); err != nil {
		log.Debugf("vips: %s in %s (set icc profile for interop index tag)", err, clean.Log(filepath.Base(imageName)))
	}

	// Create thumbnail image.
	if err = img.ThumbnailWithSize(width, height, crop, size); err != nil {
		log.Debugf("vips: %s in %s (create thumbnail)", err, clean.Log(filepath.Base(imageName)))
		return "", nil, err
	}

	// Remove metadata from thumbnail.
	if err = img.RemoveMetadata(); err != nil {
		log.Debugf("vips: %s in %s (remove metadata)", err, clean.Log(filepath.Base(imageName)))
		return "", nil, err
	}

	// Export to standard image format.
	switch fs.FileType(thumbName) {
	case fs.ImagePng:
		thumbBuffer, _, err = img.ExportPng(VipsPngExportParams(width, height))
	default:
		thumbBuffer, _, err = img.ExportJpeg(VipsJpegExportParams(width, height))
	}

	// Check if export failed.
	if err != nil {
		log.Debugf("vips: %s in %s (export thumbnail)", err, clean.Log(filepath.Base(imageName)))
		return "", thumbBuffer, err
	}

	// Write thumbnail to file.
	if err = os.WriteFile(thumbName, thumbBuffer, fs.ModeFile); err != nil {
		log.Debugf("vips: %s in %s (write thumbnail to file)", err, clean.Log(filepath.Base(imageName)))
		return "", thumbBuffer, err
	}

	return thumbName, thumbBuffer, nil
}

// VipsImportParams provides parameters for opening files with libvips.
func VipsImportParams() *vips.ImportParams {
	params := &vips.ImportParams{}
	params.AutoRotate.Set(true)
	params.FailOnError.Set(false)
	return params
}

// VipsPngExportParams returns PNG image encoding parameters for libvips.
func VipsPngExportParams(width, height int) *vips.PngExportParams {
	params := vips.NewPngExportParams()
	params.Filter = vips.PngFilterNone
	params.Interlace = false
	params.Palette = false

	// Set compression depending on image size.
	if width > 20 || height > 20 {
		params.Compression = 6
	} else {
		params.Compression = 0
	}

	return params
}

// VipsJpegExportParams returns JPEG image encoding parameters for libvips.
func VipsJpegExportParams(width, height int) *vips.JpegExportParams {
	params := vips.NewJpegExportParams()
	params.Quality = JpegQuality(width, height).Int()
	params.Interlace = true
	params.SubsampleMode = vips.VipsForeignSubsampleAuto

	// Enable quality enhancements depending on image size.
	if width > 150 || height > 150 {
		params.OptimizeCoding = true
		// The following options can only be set if libvips
		// has been compiled with the "--with-mozjpeg" flag:
		// params.QuantTable = 3
		// params.TrellisQuant = true
		// params.OvershootDeringing = true
		// params.OptimizeScans = true
	}

	return params
}

// VipsRotate rotates a vips image based on the Exif orientation.
func VipsRotate(img *vips.ImageRef, orientation int) error {
	var err error

	switch orientation {
	case OrientationUnspecified:
		// Do nothing.
	case OrientationNormal:
		// Do nothing.
	case OrientationFlipH:
		err = img.Flip(vips.DirectionHorizontal)
	case OrientationFlipV:
		err = img.Flip(vips.DirectionVertical)
	case OrientationRotate90:
		// Rotate the image 90 degrees counter-clockwise.
		err = img.Rotate(vips.Angle270)
	case OrientationRotate180:
		err = img.Rotate(vips.Angle180)
	case OrientationRotate270:
		// Rotate the image 270 degrees counter-clockwise.
		err = img.Rotate(vips.Angle90)
	case OrientationTranspose:
		err = img.Flip(vips.DirectionHorizontal)
		if err == nil {
			// Rotate the image 90 degrees counter-clockwise.
			err = img.Rotate(vips.Angle270)
		}
	case OrientationTransverse:
		err = img.Flip(vips.DirectionVertical)
		if err == nil {
			// Rotate the image 90 degrees counter-clockwise.
			err = img.Rotate(vips.Angle270)
		}
	default:
		log.Debugf("vips: invalid orientation %d (rotate image)", orientation)
	}

	return err
}

func vipsSetIccProfileForInteropIndex(img *vips.ImageRef, logName string) error {

	// Many cameras will define a JPEG's colour space by setting the InteroperabilityIndex
	// tag instead of embedding an inline ICC profile.
	// We detect this and embed explicit icc profiles for thumbs of such images, for the benefit of Vips
	// and web browsers, none of which pay any attention to the InteropIndex tag.
	iiFull := img.GetString("exif-ifd4-InteroperabilityIndex")
	if iiFull == "" {
		return nil
	}

	// according to my reading, I think [:4] should be e.g. "R98\x00".
	// However, vips always returns [:4] = "R98 ", e.g. space instead of null.
	// I'm pulling [:3] instead to paper over this - the exif spec says "4 bytes
	// incl null terminator" so I think this is safe.
	ii := iiFull[:3]
	log.Tracef("interopindex: %s read exif and got interopindex %s, %s", logName, ii, iiFull)

	if img.HasICCProfile() {
		log.Debugf("interopindex: %s has both an interop index tag and an embedded ICC profile. ignoring.", logName)
		return nil
	}

	fallbackProfile := ""
	switch ii {
	case "R03":
		// adobe rgb
		fallbackProfile = MustGetAdobeRGB1998Path()
	case "R98":
		// srgb
		// we could logically embed an srgb profile in the image here, but
		// there's no value in doing so; everything assumes srgb anyway.
	case "THM":
		// a thumbnail file. I can't find a ref on what colour space
		// this is, so I'm assuming without evidence that they are also srgb.
	default:
		log.Debugf("interopindex: %s has unknown interop index %s", logName, ii)
	}
	if fallbackProfile == "" {
		return nil
	}
	return img.TransformICCProfileWithFallback(fallbackProfile, fallbackProfile) // icc profile gets embedded here
}
