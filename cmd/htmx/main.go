package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Count struct {
	Count int
}

func main() {
	templates := template.Must(template.ParseGlob("views/*.html"))

	count := &Count{Count: 0}

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index", count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	router.HandleFunc("POST /count", func(w http.ResponseWriter, r *http.Request) {
		count.Count++

		if err := templates.ExecuteTemplate(w, "count", count); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("Starting server on port :8080")
	server.ListenAndServe()
}
