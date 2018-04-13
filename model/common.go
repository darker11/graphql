package model

type UserStatusType int64

const (
	EnableStatus UserStatusType = iota + 1
	DisableStatus
)
