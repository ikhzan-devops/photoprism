package customize

import (
	"time"

	"github.com/photoprism/photoprism/pkg/i18n"
)

var (
	DefaultTheme     = "default"
	DefaultStartPage = "default"
	DefaultMapsStyle = "default"
	DefaultLanguage  = i18n.Default.Locale()
	DefaultTimeZone  = time.Local.String()
)
