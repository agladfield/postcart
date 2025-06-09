package enum

type EmailStatusEnum int8

const (
	EmailStatusUnknown EmailStatusEnum = iota
	EmailStatusSent
	EmailStatusDelivered
	EmailStatusBounced
	EmailStatusMarkedAsSpam
)
