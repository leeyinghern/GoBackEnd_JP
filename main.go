package main

import (
	"net/http"
	"packages/functions"
)

func main() {

	// go http.ListenAndServe(":8080", http.HandlerFunc(redirect))
	// mux := http.NewServeMux()
	http.HandleFunc("/", functions.Homepage)
	http.HandleFunc("/login", functions.Login)
	http.HandleFunc("/signup", functions.SignUpPage)
	http.HandleFunc("/menu", functions.MenuPage)
	http.HandleFunc("/vocab", functions.ServeVocabQuestionToUser)
	http.HandleFunc("/vocab/", functions.ServeVocabQuestionToUser)
	http.HandleFunc("/grade/", functions.DisplayGrade)
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// 	log.Printf("Defaulting to port %s ", port)
	// }
	// http.ListenAndServeTLS(":443", "cert.pem", "key.pem", mux)
	http.ListenAndServe(":8080", nil)
}

// func redirect(w http.ResponseWriter, req *http.Request) {
// 	http.Redirect(w, req,
// 		"https://"+req.Host+req.URL.String(),
// 		http.StatusMovedPermanently)
// }
