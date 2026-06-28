package jsonl

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
	Type    string          `json:"type"`
}

// ContentText extracts readable text from content, which may be a plain
// string or a Claude content-block array.
func (m Message) ContentText() string {
	// Try plain string
	var s string
	if json.Unmarshal(m.Content, &s) == nil {
		return s
	}
	// Try content-block array
	var blocks []contentBlock
	if json.Unmarshal(m.Content, &blocks) != nil {
		return string(m.Content) // fallback: raw JSON
	}
	var parts []string
	for _, b := range blocks {
		switch b.Type {
		case "text":
			if b.Text != "" {
				parts = append(parts, b.Text)
			}
		case "thinking":
			if b.Thinking != "" {
				parts = append(parts, "> "+b.Thinking)
			}
		case "tool_use":
			if b.Name != "" {
				parts = append(parts, "> 🔧 **"+b.Name+"**")
			}
		}
	}
	return strings.Join(parts, "\n\n")
}

type contentBlock struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Thinking string `json:"thinking"`
	Name     string `json:"name"`
}

type rawLine struct {
	Type    string          `json:"type"`
	Message json.RawMessage `json:"message"`
}

func ParseFile(path string) ([]Message, error) {
	return ParseFileRange(path, 0, 0)
}

// ParseFileRange reads lines [start, start+limit) from a jsonl file.
// If start==0 && limit==0, reads all.
func ParseFileRange(path string, start, limit int) ([]Message, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	type rawLine struct {
		Type    string          `json:"type"`
		Message json.RawMessage `json:"message"`
	}

	var out []Message
	var lineNum int
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 50*1024*1024)

	for scanner.Scan() {
		lineNum++
		// Skip lines before start
		if start > 0 && lineNum <= start {
			continue
		}
		// Stop when we have enough
		if limit > 0 && len(out) >= limit {
			break
		}
		var rl rawLine
		if err := json.Unmarshal(scanner.Bytes(), &rl); err != nil {
			continue
		}
		var m Message
		if err := json.Unmarshal(rl.Message, &m); err != nil {
			continue
		}
		m.Type = rl.Type
		out = append(out, m)
	}
	return out, scanner.Err()
}
