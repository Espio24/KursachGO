package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var mass map[string]*Model

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

type Model struct {
	Data      string
	Aidi      string
	Name      string
	Zagolovok string
	Image     string
	Text      string
}

func NewModel(Data, Aidi, Name, Zagolovok, Image, Text string) *Model {
	return &Model{Data, Aidi, Name, Zagolovok, Image, Text}
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

	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	t, err := template.ParseFiles("pages/edit.html")
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
	t, err := template.ParseFiles("pages/view.html")
	t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/post.html")
	if err != nil {
		fmt.Print(err)
	}

	massiv := []*Model{}

	id := r.FormValue("id")
	post, found := mass[id]
	massiv = append(massiv, post)

	if !found {
		http.NotFound(w, r)
	}

	t.Execute(w, massiv)
}

func generalHandler(w http.ResponseWriter, r *http.Request) {

	category_name, err := ioutil.ReadFile("Text/category_name.txt")
	if err != nil {
		fmt.Print(err)
	}

	category := strings.Split(string(category_name), "\n")
	var lengt int = len(category)
	var a = make([][]string, lengt)
	var categoryReqest bool = false
	for i := 0; i < len(category); i++ {
		a[i] = strings.Split(string(category[i]), "/")
		namePage := strings.Split(string(category[i]), "/")
		if ("/"+namePage[1]+"/") == r.RequestURI || ("/"+namePage[1]) == r.RequestURI {
			categoryReqest = true
		}
	}
	if categoryReqest {
		categoryHandler(w, r)
	} else {
		t, err := template.ParseFiles("pages/general.html")
		if err != nil {
			fmt.Print(err)
		}
		t.Execute(w, a)
	}

}

func categoryHandler(w http.ResponseWriter, r *http.Request) {

	category_name, err := ioutil.ReadFile("Text/" + strings.Replace(r.RequestURI, "/", "", -1) + ".txt")
	if err != nil {
		fmt.Print(err)
	}

	header_name := strings.Split(string(category_name), "\n")

	var Data string
	var Aidi string
	var Name string
	var Zagolovok string
	var Image string
	var Text string

	var n int = 0

	massiv := []*Model{}

	for i := 0; i < len(header_name); i++ {
		switch n {
		case 0:
			TexnicText := strings.Split(string(header_name[i]), " ")
			Data = TexnicText[0]
			Aidi = TexnicText[1]
			Name = TexnicText[2]
		case 1:
			Zagolovok = header_name[i]
		case 2:
			Image = header_name[i]
		case 3:
			Text = header_name[i]
			post := NewModel(Data, Aidi, Name, Zagolovok, Image, Text)
			mass[post.Aidi] = post
			massiv = append(massiv, post)
		case 4:
			n = -1
		}
		n++
	}

	t, err := template.ParseFiles("pages/category.html")
	if err != nil {
		fmt.Print(err)
	}

	t.Execute(w, massiv)

}

func main() {
	mass = make(map[string]*Model, 0)

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/viewPost/", viewPostHandler)
	http.HandleFunc("/", generalHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
