package render

import (
	"math/rand"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	templateFuncs = template.FuncMap{
		"replace":       strings.ReplaceAll,
		"trim":          strings.TrimSpace,
		"upper":         strings.ToUpper,
		"lower":         strings.ToLower,
		"title":         toTitle,
		"contains":      strings.Contains,
		"hasPrefix":     strings.HasPrefix,
		"hasSuffix":     strings.HasSuffix,
		"rand":          rand.Intn,
		"randFloat":     rand.Float64,
		"now":           now,
		"dict":          dict,
		"slice":         slice,
		"append":        _append,
		"randomPickOne": randomPickOne,
	}
)

func Render(content string, params any) (string, error) {
	tmpl, err := template.New("tpl").
		Funcs(templateFuncs).
		Parse(content)
	if err != nil {
		return "", err
	}

	if params == nil {
		params = map[string]any{}
	}

	buf := &strings.Builder{}
	if err := tmpl.Execute(buf, params); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func toTitle(s string) string {
	return cases.Title(language.English).String(s)
}

func randomPickOne(items []any) any {
	if len(items) == 0 {
		return nil
	}
	return items[rand.Intn(len(items))]
}

func now() string {
	return time.Now().Format(time.ANSIC)
}

func dict(values ...any) map[string]any {
	if len(values)%2 != 0 {
		return nil
	}
	d := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		d[values[i].(string)] = values[i+1]
	}
	return d
}

func slice(items ...any) []any {
	return items
}

func _append(slice []any, items ...any) []any {
	return append(slice, items...)
}
