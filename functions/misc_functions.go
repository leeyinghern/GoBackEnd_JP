package functions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func CheckIfOverlappingQuestionNumber(a []int, n int) bool {
	// Returns false if there is overlap
	// Returns true otherwise
	for _, val := range a {
		if val == n {
			return false
		}
	}
	return true
}

func ReadJson(filepath string, loadedqns *Questions) *Questions {
	questions_json, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer questions_json.Close()
	// Assigns questions to user and returns pointer to cookie
	// Read opened JSON as a byte array
	byteValue, _ := ioutil.ReadAll(questions_json)

	// Unmarshal the JSON
	json.Unmarshal(byteValue, loadedqns)

	return loadedqns
}

func CreateNewCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	// Create cookie
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		cookie = &http.Cookie{
			Name:   "riinsan",
			Value:  uuid.NewString(),
			MaxAge: 60,
		}
	}
	http.SetCookie(w, cookie)
	return cookie
}

func DeleteUserCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "riinsan",
		Value:  uuid.NewString(),
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
