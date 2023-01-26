package gitlab

// Kind defines the kind of ref that generated an event (branch, tag, merge_request, ...)
type Kind string

const (
	KindBranch       Kind = "branch"
	KindMergeRequest Kind = "merge_request"
	KindTag          Kind = "tag"
)

func (k Kind) String() string { return string(k) }

// Status defines a job or a pipeline status.
type Status byte

const (
	StatusUnknown Status = iota
	StatusCreated
	StatusWaitingForResource
	StatusPreparing
	StatusPending
	StatusRunning
	StatusSuccess
	StatusFailed
	StatusCanceled
	StatusSkipped
	StatusManual
	StatusScheduled
)

// Statuses list all possibles status for a job or a pipeline.
var Statuses = [...]string{
	"unknown",
	"created",
	"waiting_for_resource",
	"preparing",
	"pending",
	"running",
	"success",
	"failed",
	"canceled",
	"skipped",
	"manual",
	"scheduled",
}

// StatusFromString returns the status corresponding to the given string.
func StatusFromString(str string) Status {
	for i, v := range Statuses {
		if str == v {
			return Status(i)
		}
	}
	return StatusUnknown
}

func (s Status) String() string { return Statuses[s] }
