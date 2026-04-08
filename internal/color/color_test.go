package color

import (
	"testing"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name string
		text string
		code string
		want string
	}{
		{
			name: "Bold text",
			text: "hello",
			code: Bold,
			want: "\033[1mhello\033[0m",
		},
		{
			name: "Red text",
			text: "dirty",
			code: Red,
			want: "\033[31mdirty\033[0m",
		},
		{
			name: "Green text",
			text: "clean",
			code: Green,
			want: "\033[32mclean\033[0m",
		},
		{
			name: "Yellow text",
			text: "3",
			code: Yellow,
			want: "\033[33m3\033[0m",
		},
		{
			name: "Cyan text",
			text: "*",
			code: Cyan,
			want: "\033[36m*\033[0m",
		},
		{
			name: "Dim text",
			text: "no upstream",
			code: Dim,
			want: "\033[2mno upstream\033[0m",
		},
		{
			name: "Empty string",
			text: "",
			code: Bold,
			want: "\033[1m\033[0m",
		},
		{
			name: "Combined bold+cyan",
			text: "*",
			code: Bold + ";" + CyanRaw,
			want: "\033[1;36m*\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Colorize(tt.text, tt.code)
			if got != tt.want {
				t.Errorf("Colorize(%q, %q) = %q, want %q", tt.text, tt.code, got, tt.want)
			}
		})
	}
}

func TestIsEnabledRespectsNO_COLOR(t *testing.T) {
	// When NO_COLOR is set, color should be disabled regardless
	t.Setenv("NO_COLOR", "1")
	if IsEnabled() {
		t.Error("IsEnabled() = true, want false when NO_COLOR is set")
	}
}

func TestIsEnabledRespectsEmptyNO_COLOR(t *testing.T) {
	// When NO_COLOR is set to empty string, it still counts per no-color.org spec
	t.Setenv("NO_COLOR", "")
	if IsEnabled() {
		t.Error("IsEnabled() = true, want false when NO_COLOR is set (even empty)")
	}
}
