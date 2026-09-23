package session

import "time"

type RefreshToken struct {
	hash       string
	employeeID string
	expiresAt  time.Time
	revoked    bool
	rotatedAt  time.Time
}

func NewRefreshToken(hash, employeeID string, expiresAt time.Time) RefreshToken {
	return RestoreRefreshToken(hash, employeeID, expiresAt, false, time.Time{})
}

func RestoreRefreshToken(hash, employeeID string, expiresAt time.Time, revoked bool, rotatedAt time.Time) RefreshToken {
	return RefreshToken{hash: hash, employeeID: employeeID, expiresAt: expiresAt, revoked: revoked, rotatedAt: rotatedAt}
}

func (t RefreshToken) Hash() string {
	return t.hash
}

func (t RefreshToken) EmployeeID() string {
	return t.employeeID
}

func (t RefreshToken) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t RefreshToken) IsExpired(now time.Time) bool {
	return !now.Before(t.expiresAt)
}

func (t RefreshToken) IsRevoked() bool {
	return t.revoked
}

func (t RefreshToken) RotatedAt() time.Time {
	return t.rotatedAt
}

func (t RefreshToken) RotatedWithin(now time.Time, window time.Duration) bool {
	return !t.rotatedAt.IsZero() && now.Sub(t.rotatedAt) <= window
}

func (t *RefreshToken) Revoke() {
	t.revoked = true
	t.rotatedAt = time.Time{}
}

func (t *RefreshToken) Rotate(at time.Time) {
	t.revoked = true
	t.rotatedAt = at
}
