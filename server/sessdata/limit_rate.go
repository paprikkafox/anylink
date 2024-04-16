package sessdata

import (
	"context"

	"golang.org/x/time/rate"
)

type LimitRater struct {
	limit *rate.Limiter
}

// lim: token generation rate
// burst: maximum allowed burst rate
func NewLimitRater(lim, burst int) *LimitRater {
	limit := rate.NewLimiter(rate.Limit(lim), burst)
	return &LimitRater{limit: limit}
}

// bt cannot exceed the burst size
func (l *LimitRater) Wait(bt int) error {
	return l.limit.WaitN(context.Background(), bt)
}
