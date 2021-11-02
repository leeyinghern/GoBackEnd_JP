package functions

import (
	"log"
	"net/http"
	"strconv"
)

func Handle_user_cookies(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("site-cookie")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "site-cookie",
			Value: "0",
		}
	}
	count, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	count++
	cookie.Value = strconv.Itoa(count)
	http.SetCookie(w, cookie)
}
