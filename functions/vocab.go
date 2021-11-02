package functions

import (
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
	if UserQuestionNumber < 4 {
		CurrentUserQuestion := VocabQuestions[UserQuestionNumber]
		CacheUserAnswer(w, r, cookie, &CurrentUserQuestion, "vocab_qns")
	} else {
		http.Redirect(w, r, "/grade/vocab", http.StatusSeeOther)
	}

	// if UserQuestionNumber < 5 {
	// 	CurrentUserQuestion := VocabQuestions[UserQuestionNumber]
	// 	VocabUserSessions[UserSession].CurrentQuestion = UserQuestionNumber + 1
	// 	// var button_val string
	// 	// if UserQuestionNumber == 4 {
	// 	// 	button_val = "Submit and check my grade"
	// 	// } else {
	// 	// 	button_val = "Next Question"
	// 	// }
	// 	// QuestionToPass := map[string]string{
	// 	// 	"question":        CurrentUserQuestion.Question,
	// 	// 	"question_number": fmt.Sprint(UserQuestionNumber),
	// 	// 	"image":           CurrentUserQuestion.Image_link,
	// 	// 	// "answers":         strings.Join(CurrentUserQuestion.Question_answer, ","),
	// 	// 	"button_value": button_val,
	// 	// }

	// 	// TPL.ExecuteTemplate(w, "vocab.html", QuestionToPass)
	// 	CacheUserAnswer(w, r, cookie, &CurrentUserQuestion, "vocab_qns")
	// } else {
	// 	// TPL.ExecuteTemplate(w, "grade.html", nil)
	// 	http.Redirect(w, r, "/grade", http.StatusSeeOther)
	// }

}
