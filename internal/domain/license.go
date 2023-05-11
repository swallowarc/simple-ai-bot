package domain

import (
	"time"
)

type LicenseState string

const (
	LicensePendingTTL = time.Hour * 24 // 1 day

	LicenseStateNone     LicenseState = "none"
	LicenseStatePending  LicenseState = "pending"
	LicenseStateRejected LicenseState = "rejected"
	LicenseStateApproved LicenseState = "approved"
)

func (ls LicenseState) IsApproved() bool {
	return ls == LicenseStateApproved
}

func (ls LicenseState) String() string {
	return string(ls)
}

type License struct {
	EventSource EventSource  `json:"event_source"`
	State       LicenseState `json:"state"`
}

func NewLicense(es EventSource) License {
	return License{
		EventSource: es,
		State:       LicenseStatePending,
	}
}
