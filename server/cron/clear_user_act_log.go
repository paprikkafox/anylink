package cron

import (
	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
)

// Clear user activity log
func ClearUserActLog() {
	lifeDay, timesUp := isClearTime()
	if !timesUp {
		return
	}
	// When the audit log is saved permanently, exit
	if lifeDay <= 0 {
		return
	}
	affected, err := dbdata.UserActLogIns.ClearUserActLog(getTimeAgo(lifeDay))
	base.Info("Cron ClearUserActLog: ", affected, err)
}
