package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrontabMatcher(t *testing.T) {
	assert := assert.New(t)

	assert.True(CrontabMatcher("* * * * *", false))
	assert.True(CrontabMatcher("1 1 * * *", false))
	assert.True(CrontabMatcher("1 1 * * *", true))
	assert.True(CrontabMatcher("1 1 * * 2", false))
	assert.True(CrontabMatcher("1 1 */3 * *", false))
	assert.True(CrontabMatcher("1 1 1 1 *", false))
	assert.False(CrontabMatcher("1 1 1 *", false))
	assert.False(CrontabMatcher("100 * * * *", false))
	assert.True(CrontabMatcher("TZ=Europe/Moscow 0 1 * * *", true))
	assert.True(CrontabMatcher("CRON_TZ=Europe/Moscow 0 1 * * *", true))
}
