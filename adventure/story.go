package adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <title>Choose Your Own Adventure</title>
		<style>
			body {
				font-family: Helvetica, Arial, sans-serif;
			}

			h1 {
				position: relative;
				text-align: center;
			}

			.page {
				background-color: #fffcf6;
				border: 1px solid #eee;
				box-shadow: 0 10px 6px -6px #777;
				margin: 40px auto;
				max-width: 500px;
				padding: 80px;
				width: 80%;
			}

			ul {
				border-top: 1px dotted #ccc;
				padding: 10px 0 0 0;
			}

			li {
				padding-top: 10px;
			}

			a,
			a:visited {
				color: #6295b5;
				text-decoration: none;
			}

			a:active,
			a:hover {
				color: #7792a2;
			}

			p {
				text-indent: 1em;
			}
		</style>
    </head>
    <body>
        <section class="page">
			<h1>{{.Title}}</h1>
        	{{range .Paragraphs}}
            	<p>{{.}}</p>
        	{{end}}
        	<ul>
            	{{range .Options}}
                	<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
            	{{end}}
        	</ul>
		</section>
    </body>
</html>
`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)

		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}

		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)

	var story Story

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
