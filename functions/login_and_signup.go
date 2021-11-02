package functions

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Type for storing user info in session and json
type User struct {
	Email          []byte
	Password       []byte
	WrongQuestions map[string][]int
	Is_Admin       bool
}

// Type for getting the list of users from JSON
type UserList struct {
	Users []User `json:"users"`
}

func init() {
	TPL = template.Must(template.New("new").ParseGlob("./templates/*.html"))

}

func Login(w http.ResponseWriter, r *http.Request) {
	// Load in values from form
	email := []byte(r.FormValue("email"))
	password := []byte(r.FormValue("password"))

	// Set header to render HTML to page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Initialize variables for checking
	user_in_database := false // bool to indicate if user record exists in db
	user_record := User{}     // Initialize User{} object to store the 'found' user in the json if record exists

	// Load user json into users_json var
	users_json := read_user_json("./users.json", &UserList{})

	for i := 0; i < len(users_json.Users); i++ {
		if err := bcrypt.CompareHashAndPassword(users_json.Users[i].Email, email); err == nil {
			user_in_database = true
			user_record = users_json.Users[i]
			break
		}
	}
	if user_in_database {
		// Check if entered password and encrypted one is correct
		err := bcrypt.CompareHashAndPassword(user_record.Password, password)

		if err != nil {
			fmt.Fprint(w, `Your password is wrong. Please try again.
			<br>
			Return to <a href=/>Login</a>
			`)
			return
		}
		cookie, err := r.Cookie("riinsan")
		if err != nil {
			cookie = CreateNewCookie(w, r)
		}
		// Create Session for user when they log in
		fmt.Println("Session successfully created")
		UserSessions[cookie.Value] = string(user_record.Email)
		TPL.ExecuteTemplate(w, "menu.html", nil)

	} else {
		fmt.Fprint(w, "Your email does not exist. Please signup <a href=/signup>here</a>")
	}

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	//Create new user and encrypt password
	email := r.FormValue("email")
	password := r.FormValue("password")
	password_confirm := r.FormValue("password2")
	if email == "" || password == "" || password_confirm == "" {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if password != password_confirm {
		fmt.Fprintf(w, `<h3> Your passwords do not match. Please sign up again</h3>`)
		return
	}
	e_email, err_email := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	e_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil || err_email != nil {
		log.Fatal(err)
		log.Fatal(err_email)
		return
	}
	// Create a new user object to save in DB

	new_user := User{Email: e_email,
		Password:       e_password,
		WrongQuestions: map[string][]int{"vocab": {}, "trans": {}},
		Is_Admin:       false}

	// Create cookie for new user
	cookie := CreateNewCookie(w, r)

	// Create session for user
	UserSessions[string(cookie.Value)] = string(e_email)

	// Update the users.json with the newly created user
	User_json_struct := read_user_json("./users.json", &UserList{})
	Array_of_user_jsons := User_json_struct.Users
	Array_of_user_jsons = append(Array_of_user_jsons, new_user)
	User_json_struct.Users = Array_of_user_jsons
	updated_json, err := json.MarshalIndent(User_json_struct, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./users.json", updated_json, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, "Successfully registered, please return to our <a href=/>Login Page</a>")
}

func read_user_json(filename string, mystruct *UserList) UserList {
	// To return a UserList struct that has the .Users property attached
	// This function just reads the json into the UserList struct
	user_json, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer user_json.Close()
	byteValue, _ := ioutil.ReadAll(user_json)
	User_json_struct := mystruct
	json.Unmarshal(byteValue, User_json_struct)
	return *User_json_struct

}
