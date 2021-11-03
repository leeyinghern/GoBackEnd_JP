package functions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ServeTransQuestionToUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		// CreateUserSession(w, r)
		cookie = CreateNewCookie(w, r)
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
	}

	// Assigns the questions to the user and saves the question list in their unique session
	ServeQuestionsToUser(w, r, "trans_qns")

	// // Get the current question number and assigned trans questions
	UserSession := UserSessions[cookie.Value]
	UserQuestionNumber := TransUserSessions[UserSession].CurrentQuestion
	TransQuestions := TransUserSessions[UserSession].QuestionList

	if UserQuestionNumber <= 4 {
		CurrentUserQuestion := TransQuestions[UserQuestionNumber]
		CacheUserAnswer(w, r, cookie, &CurrentUserQuestion, &UserSession, &UserQuestionNumber, &TransQuestions, "trans_qns")
	} else {
		UserAnswer := r.PostFormValue("user_answer")
		TransUserSessions[UserSession].QuestionList[UserQuestionNumber-1].UserAnswer = UserAnswer
		http.Redirect(w, r, "/grade/trans", http.StatusSeeOther)
	}

}

func CacheTransAnswer(w *http.ResponseWriter, r *http.Request, TransQuestions []Question, QuestionType string, QuestionPassed *Question,
	UserSession string, button_val string, UserQuestionNumber int) {
	if r.Method == http.MethodPost {
		UserAnswer := r.PostFormValue("user_answer")

		TransUserSessions[UserSession].QuestionList[UserQuestionNumber-1].UserAnswer = UserAnswer
		CurrentUserQuestion := TransQuestions[UserQuestionNumber]

		NextQuestion := map[string]string{
			"question":       CurrentUserQuestion.Question,
			"next_url":       fmt.Sprintf("location.href='/trans/%s';", strconv.Itoa(UserQuestionNumber)),
			"helper_grammar": CurrentUserQuestion.TransHelperGrammar,
			"helper_words":   "[" + strings.Join(CurrentUserQuestion.TransHelperWords, "  ;   ") + "]",
			"button_value":   button_val,
		}

		TPL.ExecuteTemplate(*w, "trans.html", NextQuestion)
	} else {
		// Serve the first question to the user
		CurrentUserQuestion := TransQuestions[UserQuestionNumber]
		FirstQuestion := map[string]string{
			"question":       CurrentUserQuestion.Question,
			"next_url":       fmt.Sprintf("location.href='/trans/%s';", strconv.Itoa(UserQuestionNumber)),
			"helper_grammar": CurrentUserQuestion.TransHelperGrammar,
			"helper_words":   "[" + strings.Join(CurrentUserQuestion.TransHelperWords, "  ;  ") + "]",
			"button_value":   button_val,
		}
		// Render template here
		TPL.ExecuteTemplate(*w, "trans.html", FirstQuestion)
	}
}
