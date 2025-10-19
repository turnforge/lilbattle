package templates

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"strings"

	gotl "github.com/panyam/goutils/template"
	"github.com/panyam/templar"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

//go:embed *.html panels/*.html
var embeddedTemplates embed.FS
var Templates *templar.TemplateGroup

// SetupTemplates initializes the Templar template group
func init() {
	// Create a new template group
	group := templar.NewTemplateGroup()
	Templates = group

	// Set up the file appitem loader with multiple paths
	group.Loader = templar.NewEmbedFSLoader(embeddedTemplates)
	group.AddFuncs(gotl.DefaultFuncMap())
	group.AddFuncs(template.FuncMap{
		"AsHtmlAttribs": func(m map[string]string) template.HTML {
			return `a = 'b' c = 'd'`
		},
		"Indented": func(nspaces int, code string) (formatted string) {
			lines := (strings.Split(strings.TrimSpace(code), "\n"))
			return strings.Join(lines, "<br/>")
		},
		"dset": func(d map[string]any, key string, value any) map[string]any {
			d[key] = value
			return d
		},
		"lset": func(a []any, index int, value any) []any {
			a[index] = value
			return a
		},
		"safeHTMLAttr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"ToJson": func(v any) template.JS {
			if v == nil {
				return template.JS("null")
			}
			// Use protojson.Marshal for protobuf types, regular json.Marshal for others
			// Check if it's a protobuf message using proto.Message interface
			if msg, ok := v.(proto.Message); ok {
				jsonBytes, err := protojson.Marshal(msg)
				if err == nil {
					return template.JS(jsonBytes)
				}
				log.Printf("Error marshaling protobuf to JSON: %v", err)
			}
			// Fall back to regular JSON marshaling for non-protobuf types
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("Error marshaling to JSON: %v", err)
				return template.JS("null")
			}
			return template.JS(jsonBytes)
		},
	})

	// Preload common templates to ensure they're available
	commonTemplates := []string{}

	for _, tmpl := range commonTemplates {
		// Use defer to catch panics from MustLoad
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Template not found (will create): %s", tmpl)
				}
			}()
			group.MustLoad(tmpl, "")
		}()
	}
}
