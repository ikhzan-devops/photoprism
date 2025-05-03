package batch

type Action = string

const (
	ActionRemove Action = "remove"
	ActionKeep   Action = "keep"
	ActionChange Action = "change"
)
