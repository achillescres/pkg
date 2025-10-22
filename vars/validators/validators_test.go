package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrontabMatcher(t *testing.T) {
	assert := assert.New(t)

	assert.True(Crontab("* * * * *", false))
	assert.True(Crontab("1 1 * * *", false))
	assert.True(Crontab("1 1 * * *", true))
	assert.True(Crontab("1 1 * * 2", false))
	assert.True(Crontab("1 1 */3 * *", false))
	assert.True(Crontab("1 1 1 1 *", false))

	assert.False(Crontab("1 1 1 *", false))
	assert.False(Crontab("100 * * * *", false))
	assert.False(Crontab("CRON_TZ=Europe/Moscow 0 1 * * *", false))
	assert.False(Crontab("*", false))
	assert.False(Crontab("", false))
	assert.False(Crontab("C", false))

	assert.True(Crontab("TZ=Europe/Moscow 0 1 * * *", true))
	assert.True(Crontab("CRON_TZ=Europe/Moscow 0 1 * * *", true))
	assert.False(Crontab("CRON_TZ=Europe/Moscow 0 1 C * *", true))
}
