package functions

import (
	"fmt"
	"net/http"
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

	fmt.Println("Current user question number:", UserQuestionNumber)

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
			VocabUserSessions[UserSession].WrongAnswers = append(VocabUserSessions[UserSession].WrongAnswers, VocabQuestions[UserQuestionNumber-1].Question_number)
		}
		http.Redirect(w, r, "/grade/vocab", http.StatusSeeOther)
	}

}
