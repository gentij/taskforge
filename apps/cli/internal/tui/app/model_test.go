package app

import (
	"strings"
	"testing"
	"time"

	"github.com/gentij/taskforge/apps/cli/internal/config"
	"github.com/gentij/taskforge/apps/cli/internal/tui/data"
)

func TestRefreshView_PreservesSelectionByRowID(t *testing.T) {
	now := time.Now()
	m := NewModel(nil, "", false, config.Config{}, "")
	m.view = ViewWorkflows
	m.store = data.Store{
		Workflows: []data.Workflow{
			{ID: "wf_a", Name: "A", Active: true, LatestVersion: 1, UpdatedAt: now.Add(-3 * time.Hour)},
			{ID: "wf_b", Name: "B", Active: true, LatestVersion: 1, UpdatedAt: now.Add(-2 * time.Hour)},
			{ID: "wf_c", Name: "C", Active: true, LatestVersion: 1, UpdatedAt: now.Add(-1 * time.Hour)},
		},
	}

	m.refreshView()
	if len(m.filteredRowIDs) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(m.filteredRowIDs))
	}

	m.table.SetCursor(1)
	selectedID := m.selectedRowID()
	if selectedID == "" {
		t.Fatal("expected selected row id")
	}

	for i := range m.store.Workflows {
		if m.store.Workflows[i].ID == selectedID {
			m.store.Workflows[i].UpdatedAt = now.Add(4 * time.Hour)
		}
	}

	m.refreshView()
	if got := m.selectedRowID(); got != selectedID {
		t.Fatalf("selection not preserved: got %q, want %q", got, selectedID)
	}
}

func TestParseJSONObject_RequiresObject(t *testing.T) {
	if _, err := parseJSONObject(`{"ok":true}`); err != nil {
		t.Fatalf("expected object to parse, got error: %v", err)
	}
	if _, err := parseJSONObject(`[1,2,3]`); err == nil {
		t.Fatal("expected array payload to be rejected")
	}
}

func TestDeleteConfirmModal_RequiresExactPhrase(t *testing.T) {
	m := NewModel(nil, "", false, config.Config{}, "")
	m.openDeleteConfirmModal("Delete Workflow", "Delete test workflow", "DELETE wf_test", "wf_test", "")

	if got := m.actionModalValidationError(); got == "" {
		t.Fatal("expected validation error when phrase is empty")
	}

	m.action.Confirm.SetValue("DELETE something-else")
	if got := m.actionModalValidationError(); got == "" {
		t.Fatal("expected validation error for mismatched phrase")
	}

	m.action.Confirm.SetValue("DELETE wf_test")
	if got := m.actionModalValidationError(); got != "" {
		t.Fatalf("expected valid confirmation phrase, got error: %q", got)
	}
}

func TestSubmitDeleteConfirmDispatchesDeleteCmd(t *testing.T) {
	m := NewModel(nil, "", false, config.Config{}, "")
	m.openDeleteConfirmModal("Delete Trigger", "Delete test trigger", "DELETE trg_test", "wf_test", "trg_test")
	m.action.Confirm.SetValue("DELETE trg_test")

	cmd := m.submitActionModal()
	if cmd == nil {
		t.Fatal("expected delete command to be returned")
	}

	msg := cmd()
	result, ok := msg.(mutationResultMsg)
	if !ok {
		t.Fatalf("expected mutationResultMsg, got %T", msg)
	}
	if result.err == nil {
		t.Fatal("expected command to fail without API client")
	}
	if !strings.Contains(result.err.Error(), "api client unavailable") {
		t.Fatalf("unexpected error: %v", result.err)
	}
}
