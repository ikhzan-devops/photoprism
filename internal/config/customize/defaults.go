package customize

import (
	"time"

	"github.com/photoprism/photoprism/pkg/i18n"
)

var DefaultTheme = "default"
var DefaultStartPage = "default"
var DefaultLocale = i18n.Default.Locale()
var DefaultTimezone = time.Local.String()
