package scaffold

import (
	"testing"
)

func TestCheckCommandExists(t *testing.T) {
	// "go" should always exist when running go tests
	if !CheckCommand("go") {
		t.Error("expected 'go' command to be found")
	}
}

func TestCheckCommandMissing(t *testing.T) {
	if CheckCommand("nonexistent_command_xyz_12345") {
		t.Error("expected nonexistent command to not be found")
	}
}

func TestCheckExistingState(t *testing.T) {
	log, _, _ := testLogger()
	dir := t.TempDir()

	// Empty directory should return no existing items
	existing := CheckExistingState(dir, log)
	if len(existing) != 0 {
		t.Errorf("expected empty existing state, got %v", existing)
	}
}
