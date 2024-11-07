package render

import (
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		content string
		params  map[string]any
		want    string
		wantErr bool
	}{
		{
			name:    "simple",
			content: "Hello, {{.name}}!",
			params:  map[string]any{"name": "John"},
			want:    "Hello, John!",
			wantErr: false,
		},
		{
			name:    "replace",
			content: "Hello, {{replace .name \"o\" \"0\"}}!",
			params:  map[string]any{"name": "John"},
			want:    "Hello, J0hn!",
			wantErr: false,
		},
		{
			name:    "trim",
			content: "Hello, {{trim .name}}!",
			params:  map[string]any{"name": "  John  "},
			want:    "Hello, John!",
			wantErr: false,
		},
		{
			name:    "upper",
			content: "Hello, {{upper .name}}!",
			params:  map[string]any{"name": "John"},
			want:    "Hello, JOHN!",
			wantErr: false,
		},
		{
			name:    "lower",
			content: "Hello, {{lower .name}}!",
			params:  map[string]any{"name": "John"},
			want:    "Hello, john!",
			wantErr: false,
		},
		{
			name:    "title",
			content: "Hello, {{title .name}}!",
			params:  map[string]any{"name": "john doe"},
			want:    "Hello, John Doe!",
			wantErr: false,
		},
		{
			name:    "contains",
			content: "{{if contains .name \"John\"}}Hello, John!{{else}}Hello, stranger!{{end}}",
			params:  map[string]any{"name": "John Doe"},
			want:    "Hello, John!",
			wantErr: false,
		},
		{
			name:    "hasPrefix",
			content: "{{if hasPrefix .name \"John\"}}Hello, John!{{else}}Hello, stranger!{{end}}",
			params:  map[string]any{"name": "John Doe"},
			want:    "Hello, John!",
			wantErr: false,
		},
		{
			name:    "hasSuffix",
			content: "{{if hasSuffix .name \"Doe\"}}Hello, John!{{else}}Hello, stranger!{{end}}",
			params:  map[string]any{"name": "John Doe"},
			want:    "Hello, John!",
			wantErr: false,
		},
		{
			name:    "randomPickOne",
			content: "Hello, {{randomPickOne .names}}!",
			params:  map[string]any{"names": []string{"John", "Jane", "Joe"}},
			wantErr: false,
		},
		{
			name:    "now",
			content: "now: {{now}}",
			wantErr: false,
		},
		{
			name:    "invalid template",
			content: "Hello, {{.name!",
			params:  map[string]any{"name": "John"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.content, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want && tt.want != "" {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
			if tt.want == "" && got != "" {
				t.Logf("Render() got = %v", got)
			}
		})
	}
}
