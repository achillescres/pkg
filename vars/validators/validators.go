package validators

import (
	"github.com/adhocore/gronx"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
)

var (
	AirlineCodeMatcher          = regexp.MustCompile(`^[A-Z0-9]{2}[A-Z0-9]?$`)
	UnboundedAirlineCodeMatcher = regexp.MustCompile(`[A-Z0-9]{2}[A-Z0-9]?`)
	IataMatcher                 = regexp.MustCompile(`^[A-Z]{3}$`)
	UnboundedIataMatcher        = regexp.MustCompile(`[A-Z]{3}`)
)

// Crontab validates crontab expression like "* 1 * * *"
// tzAllowed permits CRON_TZ=UTC and TZ=Europe/Moscow and validates it
func Crontab(crontab string, tzAllowed bool) bool {
	if tzAllowed {
		if strings.HasPrefix(crontab, "TZ=") || strings.HasPrefix(crontab, "CRON_TZ=") {
			crontab = strings.TrimPrefix(strings.TrimPrefix(crontab, "TZ="), "CRON_TZ=")
			two := strings.SplitN(crontab, " ", 2)
			_, err := time.LoadLocation(two[0])
			if err != nil {
				return false
			}
			crontab = two[1]
		}
	}
	return gronx.IsValid(crontab)
}
