package customize

import (
	"time"

	"github.com/photoprism/photoprism/pkg/i18n"
)

var (
	DefaultTheme     = "default"
	DefaultStartPage = "default"
	DefaultMapsStyle = "default"
	DefaultLocale    = i18n.Default.Locale()
	DefaultTimezone  = time.Local.String()
)
