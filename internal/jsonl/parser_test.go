package jsonl

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile_BasicLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.jsonl")
	mustWrite(t, path, "{\"type\":\"user\",\"message\":{\"role\":\"user\",\"content\":\"hi\"}}\n"+
		"{\"type\":\"assistant\",\"message\":{\"role\":\"assistant\",\"content\":\"hello\"}}\n")

	msgs, err := ParseFile(path)
	require.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.Equal(t, "user", msgs[0].Role)
	assert.Equal(t, "hi", msgs[0].ContentText())
}

func TestParseFile_SkipsBadLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.jsonl")
	mustWrite(t, path, "{\"type\":\"user\",\"message\":{\"role\":\"user\",\"content\":\"hi\"}}\n"+
		"this is not json\n"+
		"{\"type\":\"assistant\",\"message\":{\"role\":\"assistant\",\"content\":\"hello\"}}\n")

	msgs, err := ParseFile(path)
	require.NoError(t, err) // bad lines skipped, not error
	assert.Len(t, msgs, 2)
}

func TestContentText_ToolUseBash(t *testing.T) {
	m := Message{Role: "assistant", Content: json.RawMessage(`[{"type":"tool_use","name":"Bash","input":{"command":"ls -la"}}]`)}
	assert.Equal(t, "🔧 **Bash** `ls -la`", m.ContentText())
}

func TestContentText_ToolUseRead(t *testing.T) {
	m := Message{Role: "assistant", Content: json.RawMessage(`[{"type":"tool_use","name":"Read","input":{"file_path":"/tmp/foo.go"}}]`)}
	assert.Equal(t, "🔧 **Read** `/tmp/foo.go`", m.ContentText())
}

func TestContentText_ToolUseUnknown(t *testing.T) {
	m := Message{Role: "assistant", Content: json.RawMessage(`[{"type":"tool_use","name":"WeirdTool","input":{"a":1,"b":"x"}}]`)}
	got := m.ContentText()
	assert.Contains(t, got, "🔧 **WeirdTool**")
	assert.Contains(t, got, "a")
}

func TestContentText_ToolResultString(t *testing.T) {
	m := Message{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"x","content":"hello world"}]`)}
	assert.Equal(t, "📤 hello world", m.ContentText())
}

func TestContentText_ToolResultError(t *testing.T) {
	m := Message{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"x","content":"boom","is_error":true}]`)}
	assert.Equal(t, "⚠️ boom", m.ContentText())
}

func TestContentText_ToolResultArray(t *testing.T) {
	m := Message{Role: "user", Content: json.RawMessage(`[{"type":"tool_result","tool_use_id":"x","content":[{"type":"text","text":"line1"},{"type":"text","text":"line2"}]}]`)}
	assert.Equal(t, "📤 line1\nline2", m.ContentText())
}

func TestContentText_ToolResultTruncated(t *testing.T) {
	big := strings.Repeat("x", 2000)
	m := Message{Role: "user", Content: json.RawMessage(fmt.Sprintf(`[{"type":"tool_result","content":%q}]`, big))}
	got := m.ContentText()
	assert.True(t, strings.HasPrefix(got, "📤 "))
	assert.Contains(t, got, "…(+")
}

func TestContentText_MixedBlocks(t *testing.T) {
	m := Message{Role: "assistant", Content: json.RawMessage(`[{"type":"text","text":"看看"},{"type":"tool_use","name":"Bash","input":{"command":"pwd"}}]`)}
	got := m.ContentText()
	assert.Contains(t, got, "看看")
	assert.Contains(t, got, "🔧 **Bash** `pwd`")
}
