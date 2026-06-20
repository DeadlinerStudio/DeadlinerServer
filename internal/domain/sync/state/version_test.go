package state

import "testing"

func TestCompareLogicalVersion(t *testing.T) {
	t.Run("timestamp wins first", func(t *testing.T) {
		left := LogicalVersion{TS: "2026-03-24T12:00:00Z", Ctr: 9, Dev: "B"}
		right := LogicalVersion{TS: "2026-03-24T12:00:01Z", Ctr: 0, Dev: "A"}

		got, err := CompareLogicalVersion(left, right)
		if err != nil {
			t.Fatalf("CompareLogicalVersion returned error: %v", err)
		}
		if got >= 0 {
			t.Fatalf("expected left < right, got %d", got)
		}
	})

	t.Run("counter breaks timestamp ties", func(t *testing.T) {
		left := LogicalVersion{TS: "2026-03-24T12:00:00Z", Ctr: 0, Dev: "B"}
		right := LogicalVersion{TS: "2026-03-24T12:00:00Z", Ctr: 1, Dev: "A"}

		got, err := CompareLogicalVersion(left, right)
		if err != nil {
			t.Fatalf("CompareLogicalVersion returned error: %v", err)
		}
		if got >= 0 {
			t.Fatalf("expected left < right, got %d", got)
		}
	})

	t.Run("device breaks full ties", func(t *testing.T) {
		left := LogicalVersion{TS: "2026-03-24T12:00:00Z", Ctr: 1, Dev: "A"}
		right := LogicalVersion{TS: "2026-03-24T12:00:00Z", Ctr: 1, Dev: "B"}

		got, err := CompareLogicalVersion(left, right)
		if err != nil {
			t.Fatalf("CompareLogicalVersion returned error: %v", err)
		}
		if got >= 0 {
			t.Fatalf("expected left < right, got %d", got)
		}
	})
}
