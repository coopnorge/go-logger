package logger

import (
	"context"
	"fmt"
	"testing"
)

func TestNoFieldsMutation(t *testing.T) {
	testCases := map[string]struct {
		apply           func(e *Entry) *Entry
		expectedEntries int
	}{
		"WithError": {
			apply: func(e *Entry) *Entry {
				return e.WithError(fmt.Errorf("some error"))
			},
			expectedEntries: 1,
		},
		"WithField": {
			apply: func(e *Entry) *Entry {
				return e.WithField("single", "value")
			},
			expectedEntries: 1,
		},
		"WithFields": {
			apply: func(e *Entry) *Entry {
				return e.WithFields(Fields{"multiple1": "value", "multiple2": "value"})
			},
			expectedEntries: 2,
		},
		"WithField+WithFields": {
			apply: func(e *Entry) *Entry {
				return e.WithField("single", "value").WithFields(Fields{"multiple1": "value", "multiple2": "value"})
			},
			expectedEntries: 3,
		},
		"WithField+WithError": {
			apply: func(e *Entry) *Entry {
				return e.WithField("single", "value").WithError(fmt.Errorf("some error"))
			},
			expectedEntries: 2,
		},
		"Override WithField": {
			apply: func(e *Entry) *Entry {
				return e.WithField("single", "value").WithField("single", "overridden-value")
			},
			expectedEntries: 1,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			entry := &Entry{}
			if len(entry.fields) != 0 {
				t.Fatalf("A blank entry should not contain any fields. Found %d", len(entry.fields))
			}

			newEntry := tc.apply(entry)
			if len(entry.fields) != 0 {
				t.Fatalf("Applying a func mutated the original entry. Found %d entries in original entry.", len(entry.fields))
			}
			if len(newEntry.fields) != tc.expectedEntries {
				t.Fatalf("Applying a func did not add fields to the new entry. Expected %d, found %d", tc.expectedEntries, len(newEntry.fields))
			}
		})
	}
}

func TestNoContextMutation(t *testing.T) {
	entry1 := &Entry{}
	if entry1.context != nil {
		t.Fatalf("A blank entry should not contain a context.")
	}
	ctx1 := context.Background()
	ctx2 := context.TODO()

	entry2 := entry1.WithContext(ctx1)
	if entry1.context != nil {
		t.Fatalf("WithContext mutated the original entry1.")
	}
	if entry2.context != ctx1 {
		t.Fatalf("WithContext did not set context in entry2.")
	}

	entry3 := entry2.WithContext(ctx2)
	if entry1.context != nil {
		t.Fatalf("WithContext mutated the original entry1.")
	}
	if entry2.context != ctx1 {
		t.Fatalf("The second WithContext mutated entry2.")
	}
	if entry3.context != ctx2 {
		t.Fatalf("The second WithContext did not set context in entry3.")
	}
}
