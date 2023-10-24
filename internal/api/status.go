package api

// ProcessingStatus is a type for the current processing status
type ProcessingStatus string

const (
	ProcessingStatusNotStarted ProcessingStatus = "NOT STARTED"
	ProcessingStatusProcessing ProcessingStatus = "PROCESSING"
	ProcessingStatusPaused     ProcessingStatus = "PAUSED"
	ProcessingStatusFinished   ProcessingStatus = "FINISHED"
)
