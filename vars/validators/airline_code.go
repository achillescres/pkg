package validators

import "regexp"

var AirlineCodeMatcher = regexp.MustCompile(`^[A-Z0-9]{2}[A-Z0-9]?$`)

var IataMatcher = regexp.MustCompile(`[A-Z]{3}`)
