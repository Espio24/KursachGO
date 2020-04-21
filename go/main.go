package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path[5:])
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	/*fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	"<form action=\"/save/%s\" method=\"POST\">"+
	"<textarea name=\"body\">%s</textarea><br>"+
	"<input type=\"submit\" value=\"Save\">"+
	"</form>",
	p.Title, p.Title, p.Body)*/
	t, err := template.ParseFiles("edit.html")
	if err != nil {
		log.Println(err)
	}

	var buf *bytes.Buffer = &bytes.Buffer{}

	err = t.Execute(buf, p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("Server error!"))
		return
	}
	//w.Write(buf.Bytes())
	buf.WriteTo(w)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	t, err := template.ParseFiles("view.html")
	t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

func generalHandler(w http.ResponseWriter, r *http.Request) {
	/*title := r.URL.Path[len("/"):]
	p, _ := loadPage(title)
	t, err := template.ParseFiles("general.html")
	t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}*/

	const tpl = `<html><head><title>Проверка статуса</title></head><body> <h1>Статус</h1><div>{{.}}</div></body></html>`
	b, err := ioutil.ReadFile("category_name.txt")
	if err != nil {
		fmt.Print(err)
	}
	s := string(b)
	t, err := template.ParseFiles("general.html")
	if err != nil {
		fmt.Print(err)
	}
	t.Execute(w, s)

}

func main() {

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/", generalHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
