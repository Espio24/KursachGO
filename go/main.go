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

func generalHandler(w http.ResponseWriter, r *http.Request) {

	category_name, err := ioutil.ReadFile("Text/category_name.txt")
	if err != nil {
		fmt.Print(err)
	}

	category := strings.Split(string(category_name), "\n")

	t, err := template.ParseFiles("pages/general.html")
	if err != nil {
		fmt.Print(err)
	}
	t.Execute(w, category)

}

func categoryHandler(w http.ResponseWriter, r *http.Request) {

	category_name, err := ioutil.ReadFile("Text/programmirovanie.txt")
	if err != nil {
		fmt.Print(err)
	}

	header_name := strings.Split(string(category_name), "\n")

	/*for i := 0; i < len(header_name); i++ {

		if header_name[i] != "" {

			configLine := strings.Split(string(header_name[i]), "\n")

			newHeader := Header{HeaderName: configLine[1]}
			configs = append(configs, newHeader)
		}
	}*/

	t, err := template.ParseFiles("pages/category.html")
	if err != nil {
		fmt.Print(err)
	}

	t.Execute(w, header_name)

}

func main() {

	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/", generalHandler)
	http.HandleFunc("/category/", categoryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
