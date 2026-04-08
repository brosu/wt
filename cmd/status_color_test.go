package cmd

import (
	"strings"
	"testing"
)

func TestFormatStatusLineNoColor(t *testing.T) {
	// When color is disabled, formatStatusLine should produce the same output as before
	tests := []struct {
		name  string
		entry worktreeStatus
	}{
		{
			name: "Current worktree dirty with ahead/behind",
			entry: worktreeStatus{
				Path:        "/path/to/worktree",
				Branch:      "feat/foo",
				HEAD:        "abc1234",
				Dirty:       true,
				Ahead:       2,
				Behind:      1,
				Current:     true,
				HasUpstream: true,
			},
		},
		{
			name: "Non-current clean worktree",
			entry: worktreeStatus{
				Path:        "/path/to/main",
				Branch:      "main",
				HEAD:        "def5678",
				Dirty:       false,
				Ahead:       0,
				Behind:      0,
				Current:     false,
				HasUpstream: true,
			},
		},
		{
			name: "No upstream",
			entry: worktreeStatus{
				Path:        "/path/to/wt",
				Branch:      "fix/bar",
				HEAD:        "ghi9012",
				Dirty:       false,
				Ahead:       0,
				Behind:      0,
				Current:     false,
				HasUpstream: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noColor := formatStatusLineColor(tt.entry, false)
			// Should not contain any ANSI escape sequences
			if strings.Contains(noColor, "\033[") {
				t.Errorf("formatStatusLineColor(color=false) contains ANSI codes: %q", noColor)
			}
			// Verify same content as original formatStatusLine
			original := formatStatusLine(tt.entry)
			if noColor != original {
				t.Errorf("formatStatusLineColor(color=false) = %q, want %q (same as formatStatusLine)", noColor, original)
			}
		})
	}
}

func TestFormatStatusLineWithColor(t *testing.T) {
	tests := []struct {
		name           string
		entry          worktreeStatus
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "Current dirty worktree with ahead/behind",
			entry: worktreeStatus{
				Path:        "/path/to/worktree",
				Branch:      "feat/foo",
				HEAD:        "abc1234",
				Dirty:       true,
				Ahead:       2,
				Behind:      1,
				Current:     true,
				HasUpstream: true,
			},
			wantContains: []string{
				"\033[1;36m*\033[0m",     // bold+cyan marker
				"\033[1mfeat/foo\033[0m", // bold branch
				"\033[31mdirty\033[0m",   // red dirty
				"\033[33m",               // yellow for ahead/behind
			},
		},
		{
			name: "Non-current clean worktree",
			entry: worktreeStatus{
				Path:        "/path/to/main",
				Branch:      "main",
				HEAD:        "def5678",
				Dirty:       false,
				Ahead:       0,
				Behind:      0,
				Current:     false,
				HasUpstream: true,
			},
			wantContains: []string{
				"\033[32mclean\033[0m", // green clean
			},
			wantNotContain: []string{
				"\033[1;36m*\033[0m", // should NOT have colored marker
			},
		},
		{
			name: "No upstream",
			entry: worktreeStatus{
				Path:        "/path/to/wt",
				Branch:      "fix/bar",
				HEAD:        "ghi9012",
				Dirty:       false,
				Ahead:       0,
				Behind:      0,
				Current:     false,
				HasUpstream: false,
			},
			wantContains: []string{
				"\033[2mno upstream\033[0m", // dim no upstream
			},
		},
		{
			name: "Ahead only colored yellow",
			entry: worktreeStatus{
				Path:        "/path/to/wt",
				Branch:      "feat/x",
				HEAD:        "abc123",
				Dirty:       false,
				Ahead:       5,
				Behind:      0,
				Current:     false,
				HasUpstream: true,
			},
			wantContains: []string{
				"\033[33m↑5\033[0m", // yellow ahead
			},
		},
		{
			name: "Behind only colored yellow",
			entry: worktreeStatus{
				Path:        "/path/to/wt",
				Branch:      "feat/y",
				HEAD:        "abc123",
				Dirty:       false,
				Ahead:       0,
				Behind:      3,
				Current:     false,
				HasUpstream: true,
			},
			wantContains: []string{
				"\033[33m↓3\033[0m", // yellow behind
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStatusLineColor(tt.entry, true)
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("formatStatusLineColor(color=true) = %q, want it to contain %q", got, want)
				}
			}
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("formatStatusLineColor(color=true) = %q, should NOT contain %q", got, notWant)
				}
			}
		})
	}
}
