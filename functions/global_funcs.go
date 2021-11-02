package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

var LoadedQuestions Questions

func Assign_Questions_To_User(w http.ResponseWriter, r *http.Request, QuestionJSON string) *http.Cookie {
	// Get cookie UUID
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		cookie = CreateNewCookie(w, r)
	}

	// Implement if else to skip the rest of the function
	if _, ok := VocabUserSessions[UserSessions[cookie.Value]]; ok {
		fmt.Println("Vocab User Session detected. Reading from previously assigned map")
		return cookie
	} else if _, ok := TransUserSessions[UserSessions[cookie.Value]]; ok {
		fmt.Println("Trans User Session detected. Reading from previously assigned map")
		return cookie
	} else {
		questions_json, err := os.Open(fmt.Sprintf("./questionbank/vocab/%s.json", QuestionJSON))
		if err != nil {
			log.Fatal(err)
		}
		defer questions_json.Close()
		// Assigns questions to user and returns pointer to cookie
		// Read opened JSON as a byte array
		byteValue, _ := ioutil.ReadAll(questions_json)

		// Unmarshal the JSON
		json.Unmarshal(byteValue, &LoadedQuestions)

		// List of questions
		QuestionList := LoadedQuestions.QuestionList
		var RandomQuestionNumbers []int
		rand.NewSource(time.Now().Unix())

		// Generate 5 random, non-overlapping question numbers
		for len(RandomQuestionNumbers) < 5 {
			r_n := rand.Intn(5)
			if CheckIfOverlappingQuestionNumber(RandomQuestionNumbers, r_n) {
				RandomQuestionNumbers = append(RandomQuestionNumbers, r_n)
			}
		}

		SelectedQuestions := []Question{}
		for _, val := range RandomQuestionNumbers {
			SelectedQuestions = append(SelectedQuestions, QuestionList[val])
		}

		// Check if session present. If not, then generate 5 qns and return
		if QuestionJSON == "vocab_qns" {
			fmt.Println("AssignQnsToUser is creating a new vocab user session")
			UserVocabQuestions := Questions{
				QuestionList: SelectedQuestions,
				QuestionType: LoadedQuestions.QuestionType,
				WrongAnswers: []int{},
				// UserAnswers:     []string{},
				CurrentQuestion: 0,
			}
			// Create new vocab user session
			VocabUserSessions[UserSessions[cookie.Value]] = &UserVocabQuestions
		}
		return cookie
	}

}

func ServeQuestionsToUser(w http.ResponseWriter, r *http.Request, QuestionType string) *Questions {

	if QuestionType == "vocab_qns" {
		cookie := Assign_Questions_To_User(w, r, "vocab_qns")
		return VocabUserSessions[UserSessions[cookie.Value]]
	} else {
		cookie := Assign_Questions_To_User(w, r, "trans_qns")
		return TransUserSessions[UserSessions[cookie.Value]]
	}
}

// Used to save user answers into their session and index the question object to store
func CacheUserAnswer(w http.ResponseWriter, r *http.Request, cookie *http.Cookie,
	QuestionPassed *Question, usersession *string, userquestionnumber *int, vocabquestions *[]Question, QuestionType string) {
	r.ParseForm()

	// Handle User Session values and button
	UserSession := *usersession
	UserQuestionNumber := *userquestionnumber
	VocabQuestions := *vocabquestions
	VocabUserSessions[UserSession].CurrentQuestion = UserQuestionNumber + 1

	var button_val string
	if UserQuestionNumber == 4 {
		button_val = "Submit and check my grade"
	} else {
		button_val = "Next Question"
	}

	//TODO: REFACTOR THIS INTO SEPERATE FUNCTION
	// If user has submitted an answer
	if QuestionType == "vocab_qns" && r.Method == http.MethodPost {
		fmt.Println("This is the user answer", r.PostFormValue("user_answer"))
		UserAnswer := r.PostFormValue("user_answer")
		CorrectAnswer := QuestionPassed.Question_answer
		var UserCorrect = false
		for _, val := range CorrectAnswer {
			if val == UserAnswer {
				UserCorrect = true
			}
		}
		if !UserCorrect {
			VocabUserSessions[UserSession].WrongAnswers = append(VocabUserSessions[UserSession].WrongAnswers, QuestionPassed.Question_number)
		}
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]

		NextQuestion := map[string]string{
			"question": CurrentUserQuestion.Question,
			"next_url": fmt.Sprintf("location.href='/vocab/%s';", strconv.Itoa(UserQuestionNumber)),
			// "question_number": fmt.Sprint(UserQuestionNumber),
			"image":        CurrentUserQuestion.Image_link,
			"button_value": button_val,
		}

		TPL.ExecuteTemplate(w, "vocab.html", NextQuestion)
	} else {
		// Serve the first question to the user
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]
		FirstQuestion := map[string]string{
			"question": CurrentUserQuestion.Question,
			"next_url": fmt.Sprintf("location.href='/vocab/%s';", strconv.Itoa(UserQuestionNumber)),
			// "question_number": fmt.Sprint(UserQuestionNumber),
			"image":        CurrentUserQuestion.Image_link,
			"button_value": button_val,
		}
		// Render template here
		TPL.ExecuteTemplate(w, "vocab.html", FirstQuestion)
	}
}

func DisplayGrade(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	UserSession := UserSessions[cookie.Value]
	if r.URL.Path == "/grade/vocab" {
		VocabUserSession := VocabUserSessions[UserSession]
		UserVocabWrongAnswerIndex := VocabUserSession.WrongAnswers
		QuestionList := VocabUserSession.QuestionList
		UserVocabWrongAnswers := map[string][]string{}
		for _, val := range UserVocabWrongAnswerIndex {
			UserVocabWrongAnswers["wrong_question"] = append(UserVocabWrongAnswers["wrong_question"], QuestionList[val].Question)
			UserVocabWrongAnswers["image_link"] = append(UserVocabWrongAnswers["image_link"], QuestionList[val].Image_link)
			UserVocabWrongAnswers["correct_answer"] = append(UserVocabWrongAnswers["correct_answer"], strings.Join(QuestionList[val].Question_answer, ", "))
		}
		TPL.ExecuteTemplate(w, "grade.html", UserVocabWrongAnswers)
	}
}
