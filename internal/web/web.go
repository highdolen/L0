package web

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterWebHandlers(r *mux.Router) {
	// Отдача статики
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})
}
