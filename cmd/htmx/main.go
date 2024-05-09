package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			{Name: "first", Email: "first@"},
			{Name: "second", Email: "second@"},
			{Name: "third", Email: "third@"},
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	templates := template.Must(template.ParseGlob("views/*.html"))

	page := newPage()

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index", page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	router.HandleFunc("POST /contacts", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "email already exists"

			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := templates.ExecuteTemplate(w, "form", formData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		contact := Contact{name, email}
		page.Data.Contacts = append(page.Data.Contacts, contact)

		if err := templates.ExecuteTemplate(w, "form", newFormData()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := templates.ExecuteTemplate(w, "oob-contact", contact); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Starting server on port :8080")
	log.Fatal(server.ListenAndServe())
}
