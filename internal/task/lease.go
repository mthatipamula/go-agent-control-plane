package task

import "time"

func (t Task) IsLeaseExpired(now time.Time) bool {
	if t.LeaseExpiresAt == nil {
		return false
	}

	return !now.Before(*t.LeaseExpiresAt)
}
