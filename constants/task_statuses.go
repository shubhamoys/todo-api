package constants

// TaskStatuses defines all available statuses for task in the system
var TaskStatuses = struct {
	Todo       string
	InProgress string
	Complete   string
}{
	Todo:       "Todo",
	InProgress: "In Progress",
	Complete:   "Complete",
}

// IsValidTaskStatus checks if the provided status is valid
func IsValidTaskStatus(status string) bool {
	return status == TaskStatuses.Todo ||
		status == TaskStatuses.InProgress ||
		status == TaskStatuses.Complete
}
