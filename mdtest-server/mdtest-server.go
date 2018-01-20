package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/russross/blackfriday"
)

type Post struct {
	Title string
	Body  template.HTML
}

var (
	//копируем шаблоны, если не удалось, то выход
	post_template = template.Must(template.ParseFiles(path.Join("template", "layout.html"), path.Join("template", "post.html")))
)

func main() {
	//для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", postHandler)

	/*	s := &http.Server{
			Addr:           ":9090",
			Handler:        nil,
			ReadTimeout:    1000 * time.Second,
			WriteTimeout:   1000 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	*/
	//	http.TimeoutHandler(h, dt, msg)
	//	s.HandleFunc("/", postHandler)
	log.Println("Listening...")

	//	log.Fatal(s.ListenAndServe())
	http.ListenAndServe(":5000", nil)

	//	http.Server.Serve(l)

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fileread, _ := ioutil.ReadFile("posts/index.md")
	lines := strings.Split(string(fileread), "\n")
	title := string(lines[0])
	body := strings.Join(lines[1:len(lines)], "\n")
	body = string(blackfriday.MarkdownCommon([]byte(body)))
	post := Post{title, template.HTML(body)}
	if err := post_template.ExecuteTemplate(w, "layout", post); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
