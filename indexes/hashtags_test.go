package indexes

import (
	"testing"
)

func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "single hashtag",
			text:     "hello #world",
			expected: []string{"world"},
		},
		{
			name:     "multiple hashtags",
			text:     "#golang #scuttlebutt #p2p",
			expected: []string{"golang", "scuttlebutt", "p2p"},
		},
		{
			name:     "no hashtags",
			text:     "hello world",
			expected: []string{},
		},
		{
			name:     "hashtag at end",
			text:     "this is cool #test",
			expected: []string{"test"},
		},
		{
			name:     "hashtag at start",
			text:     "#start middle end",
			expected: []string{"start"},
		},
		{
			name:     "duplicate hashtags",
			text:     "#test #test #other",
			expected: []string{"test", "test", "other"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractHashtags(tt.text)
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractHashtags(%q) = %v, want %v", tt.text, result, tt.expected)
			}
			for i, tag := range result {
				if i >= len(tt.expected) || tag != tt.expected[i] {
					t.Errorf("ExtractHashtags(%q)[%d] = %q, want %q", tt.text, i, tag, tt.expected[i])
				}
			}
		})
	}
}

func TestExtractMentions(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "single mention",
			text:     "hello @alice",
			expected: []string{"alice"},
		},
		{
			name:     "multiple mentions",
			text:     "@alice @bob @charlie",
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "no mentions",
			text:     "hello world",
			expected: []string{},
		},
		{
			name:     "mention with path",
			text:     "check @alice/bob",
			expected: []string{"alice/bob"},
		},
		{
			name:     "mention with dots",
			text:     "hey @alice.key.ed25519",
			expected: []string{"alice.key.ed25519"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractMentions(tt.text)
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractMentions(%q) = %v, want %v", tt.text, result, tt.expected)
			}
			for i, mention := range result {
				if i >= len(tt.expected) || mention != tt.expected[i] {
					t.Errorf("ExtractMentions(%q)[%d] = %q, want %q", tt.text, i, mention, tt.expected[i])
				}
			}
		})
	}
}
