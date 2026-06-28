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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []Message
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 50*1024*1024) // up to 50MB lines

	for scanner.Scan() {
		var rl rawLine
		if err := json.Unmarshal(scanner.Bytes(), &rl); err != nil {
			continue // skip bad lines
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
