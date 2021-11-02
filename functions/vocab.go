package functions

import (
	"fmt"
	"net/http"
	"strconv"
)

func ServeVocabQuestionToUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		// CreateUserSession(w, r)
		cookie = CreateNewCookie(w, r)
		http.Redirect(w, r, "/menu", http.StatusSeeOther)
	}

	// Assigns the questions to the user and saves the question list in their unique session
	ServeQuestionsToUser(w, r, "vocab_qns")

	// // Get the current question number and assigned vocab questions
	UserSession := UserSessions[cookie.Value]
	UserQuestionNumber := VocabUserSessions[UserSession].CurrentQuestion
	VocabQuestions := VocabUserSessions[UserSession].QuestionList

	if UserQuestionNumber <= 4 {
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]
		CacheUserAnswer(w, r, cookie, &CurrentUserQuestion, &UserSession, &UserQuestionNumber, &VocabQuestions, "vocab_qns")
	} else {
		// We copy this block in order to deal with the last question: question 5
		UserAnswer := r.PostFormValue("user_answer")
		CorrectAnswer := VocabQuestions[UserQuestionNumber-1].Question_answer
		var UserCorrect = false
		for _, val := range CorrectAnswer {
			if val == UserAnswer {
				UserCorrect = true
			}
		}
		if !UserCorrect {
			VocabUserSessions[UserSession].WrongAnswers[UserQuestionNumber] = UserAnswer
		}
		http.Redirect(w, r, "/grade/vocab", http.StatusSeeOther)
	}

}

func CacheVocabAnswer(w *http.ResponseWriter, r *http.Request, VocabQuestions []Question, QuestionType string, QuestionPassed *Question,
	UserSession string, button_val string, UserQuestionNumber int) {
	if r.Method == http.MethodPost {
		UserAnswer := r.PostFormValue("user_answer")
		CorrectAnswer := QuestionPassed.Question_answer
		var UserCorrect = false
		for _, val := range CorrectAnswer {
			if val == UserAnswer {
				UserCorrect = true
			}
		}
		if !UserCorrect {
			VocabUserSessions[UserSession].WrongAnswers[UserQuestionNumber] = UserAnswer
		}
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]

		NextQuestion := map[string]string{
			"question":     CurrentUserQuestion.Question,
			"next_url":     fmt.Sprintf("location.href='/vocab/%s';", strconv.Itoa(UserQuestionNumber)),
			"image":        CurrentUserQuestion.Image_link,
			"button_value": button_val,
		}

		TPL.ExecuteTemplate(*w, "vocab.html", NextQuestion)
	} else {
		// Serve the first question to the user
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]
		FirstQuestion := map[string]string{
			"question":     CurrentUserQuestion.Question,
			"next_url":     fmt.Sprintf("location.href='/vocab/%s';", strconv.Itoa(UserQuestionNumber)),
			"image":        CurrentUserQuestion.Image_link,
			"button_value": button_val,
		}
		// Render template here
		TPL.ExecuteTemplate(*w, "vocab.html", FirstQuestion)
	}
}
