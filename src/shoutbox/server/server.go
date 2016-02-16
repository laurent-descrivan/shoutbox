package server

import (
	"io/ioutil"
	"net/http"
)

const LINES_FILE_PATH = `data/lines.txt`

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html;charset=utf-8")

		bytes, err := ioutil.ReadFile(LINES_FILE_PATH)
		if err != nil {
			bytes = []byte{}
		}

		data := map[string]interface{}{
			"txt": string(bytes),
		}

		INDEX_TEMPLATE.Execute(w, data)
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		text := r.FormValue("text")
		if len(text) > 1<<16 {
			text = text[:1<<16]
		}

		ioutil.WriteFile(LINES_FILE_PATH, []byte(text), 0755)
		ReloadText()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	return mux
}
