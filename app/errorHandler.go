package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

type Err struct {
	Mssg string
	Status  int
}

func ErrorHandler(w http.ResponseWriter, errMessage string, status int) {
	var Err Err
	tmpl, err := template.ParseFiles("./static/templates/error.html")
	if err != nil {
		fmt.Println("error parsing the error.html file\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	Err.Mssg = errMessage
	Err.Status = status
	if err := tmpl.Execute(w, Err); err != nil {
		fmt.Println("error excuting the error.html file\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
