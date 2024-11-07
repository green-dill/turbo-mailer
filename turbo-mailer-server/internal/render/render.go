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

	buf := &strings.Builder{}
	if err := tmpl.Execute(buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toTitle(s string) string {
	return cases.Title(language.English).String(s)
}

func randomPickOne(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.Intn(len(items))]
}

func now() string {
	return time.Now().Format(time.ANSIC)
}
