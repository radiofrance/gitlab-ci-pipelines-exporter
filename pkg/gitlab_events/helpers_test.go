package gitlab_events

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefKind_String(t *testing.T) {
	assert.Equal(t, "branch", BranchKind.String())
	assert.Equal(t, "tag", TagKind.String())
	assert.Equal(t, "merge_request", MergeRequestKind.String())
}

func TestStatusFromString(t *testing.T) {
	tcases := map[string]status{
		"unknown":              Unknown,
		"created":              Created,
		"waiting_for_resource": WaitingForResource,
		"preparing":            Preparing,
		"pending":              Pending,
		"running":              Running,
		"success":              Success,
		"failed":               Failed,
		"canceled":             Canceled,
		"skipped":              Skipped,
		"manual":               Manual,
		"scheduled":            Scheduled,
	}

	for str, status := range tcases {
		assert.Equal(t, status, StatusFromString(str))
	}
}

func TestStatus_String(t *testing.T) {
	assert.Equal(t, "unknown", Unknown.String())
	assert.Equal(t, "created", Created.String())
	assert.Equal(t, "waiting_for_resource", WaitingForResource.String())
	assert.Equal(t, "preparing", Preparing.String())
	assert.Equal(t, "pending", Pending.String())
	assert.Equal(t, "running", Running.String())
	assert.Equal(t, "success", Success.String())
	assert.Equal(t, "failed", Failed.String())
	assert.Equal(t, "canceled", Canceled.String())
	assert.Equal(t, "skipped", Skipped.String())
	assert.Equal(t, "manual", Manual.String())
	assert.Equal(t, "scheduled", Scheduled.String())
}

func TestStatus_MarshalJSON(t *testing.T) {
	type T struct {
		Valid   status `json:"valid,omitempty"`
		Unknown status `json:"unknown,omitempty"`
		Null    status `json:"null"`
		Empty   status `json:"empty,omitempty"`
	}

	expected := `{"valid":"created","null":"unknown"}`
	actual := T{Valid: Created, Unknown: Unknown, Null: Unknown, Empty: Unknown}

	bytes, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.Equal(t, expected, string(bytes))
}

func TestStatus_UnmarshalJSON(t *testing.T) {
	type T struct {
		Valid   status `json:"valid,omitempty"`
		Unknown status `json:"unknown,omitempty"`
		Null    status `json:"null,omitempty"`
		Empty   status `json:"empty,omitempty"`
	}

	expected := T{Valid: Created, Unknown: Unknown, Null: Unknown, Empty: Unknown}

	var actual T
	err := json.Unmarshal([]byte(`{"valid":"created","unknown":"...","null":null}`), &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestOption_Some(t *testing.T) {
	var tcases = []any{
		1, 0, 2.98, 3 + 2i,
		"string", "",
		true, false,
	}

	for _, tcase := range tcases {
		t.Run(reflect.TypeOf(tcase).String(), func(t *testing.T) {
			x := Some(tcase)

			assert.True(t, x.IsSome())
			assert.Equal(t, tcase, x.Unwrap())
		})
	}
}

func TestOption_None(t *testing.T) {
	x := None[string]()
	assert.True(t, x.IsNone())
	assert.Equal(t, "", x.Unwrap())
}

func TestOption_Unwrap(t *testing.T) {
	var x Option[string]

	require.True(t, x.IsNone())
	assert.Equal(t, "", x.Unwrap())

	x = Some("")
	require.False(t, x.IsNone())
	require.True(t, x.IsSome())
	assert.Equal(t, "", x.Unwrap())

	x = Some("str")
	assert.Equal(t, "str", x.Unwrap())
}

func TestOption_MarshalJSON(t *testing.T) {
	type T struct {
		OptSome      Option[string] `json:"opt_int"`
		OptNone      Option[bool]   `json:"opt_none"`
		OptEmpty     Option[string] `json:"opt_empty"`
		OptOmitEmpty Option[string] `json:"opt_omit_empty,omitempty"`
	}

	expected := `{"opt_int":"str","opt_none":false,"opt_empty":""}`
	actual := T{OptSome: Some("str"), OptEmpty: Some("")}

	bytes, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.Equal(t, expected, string(bytes))
}

func TestOption_UnmarshalJSON(t *testing.T) {
	type T struct {
		OptSome    Option[string] `json:"opt_int"`
		OptNone    Option[string] `json:"opt_none"`
		OptEmpty   Option[string] `json:"opt_empty"`
		OptNull    Option[string] `json:"opt_null"`
		OptInvalid Option[string] `json:"opt_invalid"`
	}

	expect := T{OptSome: Some("str"), OptEmpty: Some("")}

	var actual T
	err := json.Unmarshal([]byte(`{"opt_int": "str","opt_empty": "", "opt_null": null}`), &actual)
	require.NoError(t, err)

	assert.Equal(t, expect, actual)

	err = json.Unmarshal([]byte(`{"opt_invalid": false}`), &actual)
	require.EqualError(t, err, "json: cannot unmarshal bool into Go struct field T.opt_invalid of type string")
}
