package template

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

type Fields map[string]interface{}

type tmpl struct {
	t           *template.Template
	defaultData Fields
	files       []string
}

var templates = struct {
	sync.RWMutex
	m map[string]tmpl
}{m: make(map[string]tmpl)}

var (
	templatesDir string
	devmode      bool
)

func init() {
	templatesDir = ""
	devmode = false
}

func Develop(flag bool) {
	devmode = flag
}

func SetDirectory(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Fatal("path not found: %s", err)
	}
	templatesDir = path
}

func Load(name string, defaultdata Fields, filenames ...string) {
	// Zero out templatesDir if hot-reloading in development mode.
	tdir := ""
	if !devmode {
		tdir = templatesDir
	}

	var files []string

	for _, file := range filenames {
		files = append(files, tdir+file)
	}

	templates.m[name] = tmpl{
		template.Must(template.ParseFiles(files...)),
		defaultdata,
		files,
	}
}

func PrintLoaded() {
	for name, t := range templates.m {
		fmt.Printf("%s: %v\n", name, t)
	}
}

// Render loads file template t, parses and replaces the fields f in t and writes a response to w.
func Render(w http.ResponseWriter, name string, f Fields) error {
	if devmode {
		templates.RLock()
	}
	_, ok := templates.m[name]
	if devmode {
		templates.RUnlock()
	}
	if !ok {
		http.Error(w, fmt.Sprintf("%s does not exist.", name), http.StatusInternalServerError)
		return fmt.Errorf("%s is not a valid template", name)
	}

	// Hot-reload template files if in development mode.
	if devmode {
		templates.Lock()
		Load(name, templates.m[name].defaultData, templates.m[name].files...)
		templates.Unlock()
	}
	if devmode {
		templates.RLock()
	}
	tmpl := templates.m[name]
	if devmode {
		templates.RUnlock()
	}

	// Merge default data fields
	if f == nil {
		f = tmpl.defaultData
	} else {
		for key, value := range tmpl.defaultData {
			if _, ok := f[key]; !ok {
				f[key] = value
			}
		}
	}

	err := tmpl.t.Execute(w, f)
	if err != nil {
		return fmt.Errorf("error executing template %s: %v", name, err)
	}

	return nil
}
