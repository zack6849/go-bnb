package errs

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"GoBNB/internal/convert"

	"github.com/stretchr/testify/assert"
)

func TestCollectsEveryFailureNotJustTheFirst(t *testing.T) {
	record := map[string]string{
		"good": "42",
		"bad":  "not-a-number",
		//"absent" is deliberately not here
	}

	ec := CreateEchoChamber()
	good := ec.ConvertValue(record, "good", strconv.Atoi)
	bad := ec.ConvertValue(record, "bad", strconv.Atoi)
	absent := ec.ConvertValue(record, "absent", strconv.Atoi)

	assert.Equal(t, 42, good)
	assert.Equal(t, 0, bad)
	//should have 2 errors, bad & absent
	assert.Equal(t, 2, len(ec.errs))

	if good != 42 {
		t.Errorf("good = %d, want 42", good)
	}
	if bad != 0 || absent != 0 {
		t.Errorf("failed conversions should yield zero, got %d and %d", bad, absent)
	}
	if n := len(ec.Errors()); n != 2 {
		t.Fatalf("collected %d errors, want 2", n)
	}

	summary := ec.Summary().Error()
	// The point of collecting: one report naming both problems.
	if !strings.Contains(summary, `"bad"`) || !strings.Contains(summary, `"absent"`) {
		t.Errorf("summary should name both fields, got:\n%s", summary)
	}
	// A missing column must not look like an empty one -- that distinction is
	// what makes a silently absent CSV header diagnosable.
	if !strings.Contains(summary, "not present") {
		t.Errorf("summary should flag the absent field distinctly, got:\n%s", summary)
	}
}

func TestSummaryIsNilWhenClean(t *testing.T) {
	ec := CreateEchoChamber()
	if v := ec.ConvertValue(map[string]string{"n": "7"}, "n", strconv.Atoi); v != 7 {
		t.Errorf("v = %d, want 7", v)
	}
	if !ec.Empty() {
		t.Error("HasErrors on a clean collector")
	}
	if ec.Summary() != nil {
		t.Errorf("Summary = %v, want nil", ec.Summary())
	}
}

// Summary must wrap rather than flatten, so callers can still match on
// sentinel errors -- this is what fmt.Errorf on a concatenated string loses.
func TestSummaryPreservesUnwrapping(t *testing.T) {
	ec := CreateEchoChamber()
	ec.ConvertValue(map[string]string{"n": "xyz"}, "n", strconv.Atoi)
	if !errors.Is(ec.Summary(), strconv.ErrSyntax) {
		t.Error("Summary lost the underlying strconv.ErrSyntax")
	}
}

// A message containing a % verb must survive verbatim.
func TestSummaryDoesNotReinterpretFormatVerbs(t *testing.T) {
	ec := CreateEchoChamber()
	ec.PushError(errors.New("occupancy hit 100%d of capacity"))
	if got := ec.Summary().Error(); !strings.Contains(got, "100%d") {
		t.Errorf("Summary mangled a %% in the message: %s", got)
	}
}

func TestOptionalFieldToleratesMissingAndEmpty(t *testing.T) {
	ec := CreateEchoChamber()
	record := map[string]string{"empty": "", "bad": "1.5.9"}

	if v := ec.ConvertFieldOrZeroValue(record, "gone", convert.Float(64)); v != 0 {
		t.Errorf("missing = %v, want 0", v)
	}
	if v := ec.ConvertFieldOrZeroValue(record, "empty", convert.Float(64)); v != 0 {
		t.Errorf("empty = %v, want 0", v)
	}
	if !ec.Empty() {
		t.Fatalf("missing/empty optionals should not error, got: %v", ec.Summary())
	}
	// Present but malformed is still a real error.
	ec.ConvertFieldOrZeroValue(record, "bad", convert.Float(64))
	if !ec.HasErrors() {
		t.Error("a present-but-malformed optional should be reported")
	}
}

func TestCollectHandlesArbitraryShapes(t *testing.T) {
	ec := CreateEchoChamber()
	ok := ec.Collect(func() (string, error) { return "value", nil })
	bad := ec.Collect(func() (string, error) { return "ignored", errors.New("boom") })

	assert.Equal(t, "value", ok)
	assert.Equal(t, "", bad)
	assert.Equal(t, 1, len(ec.errs))
	assert.ErrorContains(t, ec.Summary(), "boom")
}
