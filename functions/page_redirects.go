package functions

import (
	"net/http"
)

func Home_page(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./templates/"))
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	TPL.ExecuteTemplate(w, "index.html", nil)
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	TPL.ExecuteTemplate(w, "signup.html", nil)
	SignUp(w, r)
}

func MenuPage(w http.ResponseWriter, r *http.Request) {
	TPL.ExecuteTemplate(w, "menu.html", nil)
}
