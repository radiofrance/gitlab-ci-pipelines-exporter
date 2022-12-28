package gitlab_events

import (
	"encoding/json"
	"strings"
	"time"
)

const RFC3339Spaced = `"2006-01-02 15:04:05 MST"`

// DateTime is a type helper used to marshal/unmarshal JSON into time.Time with the proper Gitlab time format
type DateTime time.Time

func (o DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(o).Format(RFC3339Spaced)), nil
}

func (o *DateTime) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(RFC3339Spaced, string(bytes))
	if err != nil {
		return err
	}

	*o = DateTime(date)
	return nil
}

// refKind defined what kind of ref as generated an event (branch, tag, merge_request, ...)
type refKind string

const (
	BranchKind       refKind = "branch"
	MergeRequestKind refKind = "merge_request"
	TagKind          refKind = "tag"
)

func (r refKind) String() string { return string(r) }

// status defined a job or a pipeline status
type status byte

const (
	Unknown status = iota
	Created
	WaitingForResource
	Preparing
	Pending
	Running
	Success
	Failed
	Canceled
	Skipped
	Manual
	Scheduled
)

// Statuses list all possibles status for a job or a pipeline
var Statuses = [...]string{"unknown", "created", "waiting_for_resource", "preparing", "pending", "running", "success", "failed", "canceled", "skipped", "manual", "scheduled"}

// StatusFromString returns the status corresponding to the given string
func StatusFromString(str string) status {
	for i, v := range Statuses {
		if str == v {
			return status(i)
		}
	}
	return Unknown
}

func (s status) String() string { return Statuses[s] }

func (s status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *status) UnmarshalJSON(bytes []byte) error {
	str := strings.ReplaceAll(string(bytes), `"`, "")

	for i, v := range Statuses {
		if str == v {
			*s = status(i)
			break
		}
	}
	return nil
}

// Option represents an optional value; every Option is either
// `Some` and contains a value, or `None`, and does not.
// NOTE: this is only used as helper during JSON unmarshalling.
type Option[T any] []T

// Some makes an Option type containing the actual value.
func Some[T any](v T) Option[T] { return Option[T]{v} }

// None makes an Option type that doesn't have a value.
func None[T any]() Option[T] { return nil }

// IsNone returns whether the Option doesn't have a value or not.
func (o Option[T]) IsNone() bool { return len(o) == 0 }

// IsSome returns whether the Option has a value or not.
func (o Option[T]) IsSome() bool { return !o.IsNone() }

// Unwrap extracts the contained value in an Option[T] when it is the Some variant.
func (o Option[T]) Unwrap() (v T) {
	if o.IsNone() {
		return v
	}
	return o[0]
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Unwrap())
}

func (o *Option[T]) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}

	var any T
	err := json.Unmarshal(bytes, &any)
	if err != nil {
		return err
	}

	*o = Some(any)
	return nil
}
