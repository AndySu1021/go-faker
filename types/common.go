package types

type Status int8

const (
	StatusEnabled Status = iota + 1
	StatusDisabled
)
