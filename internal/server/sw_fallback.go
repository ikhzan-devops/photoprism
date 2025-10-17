package server

import _ "embed"

// fallbackServiceWorker contains a minimal no-op service worker used when a
// Workbox-generated sw.js is not available (for example during tests or when
// frontend assets have not been built yet).
//
//go:embed sw_fallback.js
var fallbackServiceWorker []byte
