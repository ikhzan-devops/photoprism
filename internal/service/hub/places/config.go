package places

import (
	"time"
)

// ApiName is the backend API name.
const ApiName = "places"

// LocationServiceUrls specifies S2 cell location query URLs.
var LocationServiceUrls = []string{
	"https://places.photoprism.app/v1/location/%s",
}

// ReverseServiceUrls specifies reverse location query URLs.
var ReverseServiceUrls = []string{
	"https://places.photoprism.app/v1/reverse",
}

// SearchServiceUrls specifies location name query URLs.
var SearchServiceUrls = []string{
	"https://places.photoprism.app/v1/search",
}

// Retries specifies the number of attempts to retry the service request.
var Retries = 2

// RetryDelay specifies the waiting time between retries.
var RetryDelay = 100 * time.Millisecond

var Key = "f60f5b25d59c397989e3cd374f81cdd7710a4fca"
var Secret = "photoprism"
var UserAgent = ""
