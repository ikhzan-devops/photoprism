package customize

import (
	"github.com/photoprism/photoprism/pkg/i18n"
	"github.com/photoprism/photoprism/pkg/time/tz"
)

var (
	DefaultTheme     = "default"
	DefaultStartPage = "default"
	DefaultMapsStyle = "default"
	DefaultLanguage  = i18n.Default.Locale()
	DefaultTimeZone  = tz.Local
)
