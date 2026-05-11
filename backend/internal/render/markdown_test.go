package render_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/todoist/backend/internal/render"
)

func TestToHTML(t *testing.T) {
	t.Run("empty string returns empty string", func(t *testing.T) {
		got, err := render.ToHTML("")
		require.NoError(t, err)
		assert.Equal(t, "", got)
	})

	t.Run("bold text", func(t *testing.T) {
		got, err := render.ToHTML("**bold**")
		require.NoError(t, err)
		assert.Contains(t, got, "<strong>bold</strong>")
	})

	t.Run("links are rendered", func(t *testing.T) {
		got, err := render.ToHTML("[example](https://example.com)")
		require.NoError(t, err)
		assert.Contains(t, got, `href="https://example.com"`)
	})

	t.Run("strikethrough extension (GFM)", func(t *testing.T) {
		got, err := render.ToHTML("~~deleted~~")
		require.NoError(t, err)
		assert.Contains(t, got, "<del>deleted</del>")
	})

	t.Run("task list extension", func(t *testing.T) {
		got, err := render.ToHTML("- [x] done\n- [ ] todo")
		require.NoError(t, err)
		assert.Contains(t, got, `type="checkbox"`)
	})

	t.Run("tables extension", func(t *testing.T) {
		md := "| A | B |\n|---|---|\n| 1 | 2 |"
		got, err := render.ToHTML(md)
		require.NoError(t, err)
		assert.True(t, strings.Contains(got, "<table>") || strings.Contains(got, "<table "))
	})

	t.Run("multiple paragraphs", func(t *testing.T) {
		got, err := render.ToHTML("Para 1\n\nPara 2")
		require.NoError(t, err)
		assert.Contains(t, got, "<p>")
	})
}
