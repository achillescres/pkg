package validators

import "regexp"

var (
	AirlineCodeMatcher          = regexp.MustCompile(`^[A-Z0-9]{2}[A-Z0-9]?$`)
	UnboundedAirlineCodeMatcher = regexp.MustCompile(`[A-Z0-9]{2}[A-Z0-9]?`)
	IataMatcher                 = regexp.MustCompile(`^[A-Z]{3}$`)
	UnboundedIataMatcher        = regexp.MustCompile(`[A-Z]{3}`)
	CrontabMatcher              = regexp.MustCompile("^(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\\d+(ns|us|µs|ms|s|m|h))+)|((((\\d+,)+\\d+|(\\d+(\\/|-)\\d+)|\\d+|\\*) ?){5,7})$")
)
