package cron

import (
	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
)

// Clear access audit logs
func ClearAudit() {
	lifeDay, timesUp := isClearTime()
	if !timesUp {
		return
	}
	// When the audit log is saved permanently, exit
	if lifeDay <= 0 {
		return
	}
	affected, err := dbdata.ClearAccessAudit(getTimeAgo(lifeDay))
	base.Info("Cron ClearAudit: ", affected, err)
}
