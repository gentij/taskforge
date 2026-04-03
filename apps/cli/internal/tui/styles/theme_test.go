package styles

import "testing"

func TestThemeRegistry_IncludesSimpleThemes(t *testing.T) {
	registry := ThemeRegistry()

	tests := []string{"simple-dark", "simple-light"}
	for _, name := range tests {
		theme, ok := registry[name]
		if !ok {
			t.Fatalf("expected theme %q in registry", name)
		}
		if theme.Name == "" {
			t.Fatalf("expected theme %q to have display name", name)
		}
	}
}
