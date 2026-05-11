package render

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var converter = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Table,
		extension.Strikethrough,
		extension.TaskList,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

// ToHTML converts CommonMark/GFM markdown to sanitized HTML.
// An empty source returns an empty string without error.
func ToHTML(source string) (string, error) {
	if source == "" {
		return "", nil
	}
	var buf bytes.Buffer
	if err := converter.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
