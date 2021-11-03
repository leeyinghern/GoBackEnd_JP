package functions

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func Assign_Questions_To_User(w http.ResponseWriter, r *http.Request, QuestionJSON string) *http.Cookie {
	// Get cookie UUID
	cookie, err := r.Cookie("riinsan")
	if err == http.ErrNoCookie {
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

		var filepath string
		if QuestionJSON == "vocab_qns" {
			filepath = fmt.Sprintf("./questionbank/vocab/%s.json", QuestionJSON)
		} else if QuestionJSON == "trans_qns" {
			filepath = fmt.Sprintf("./questionbank/sentence/%s.json", QuestionJSON)
		}
		var LoadedQuestions Questions
		ReadJson(filepath, &LoadedQuestions)
		// List of questions
		QuestionList := LoadedQuestions.QuestionList
		SelectedQuestions := GenerateNewQuestions(len(QuestionList), QuestionList)

		// Check if session present. If not, then generate 5 qns and return
		if QuestionJSON == "vocab_qns" {
			fmt.Println("AssignQnsToUser is creating a new vocab user session")
			UserVocabQuestions := Questions{
				QuestionList: SelectedQuestions,
				// QuestionType:    LoadedQuestions.QuestionType,
				WrongAnswers:    map[int]string{},
				CurrentQuestion: 0,
			}
			// Create new vocab user session
			VocabUserSessions[UserSessions[cookie.Value]] = &UserVocabQuestions
		} else if QuestionJSON == "trans_qns" {
			fmt.Println("AssignQnsToUser is creating a new trans user session")
			UserTransQuestions := Questions{
				QuestionList: SelectedQuestions,
				// QuestionType:    LoadedQuestions.QuestionType,
				WrongAnswers:    map[int]string{},
				CurrentQuestion: 0,
			}
			// Create new vocab user session
			TransUserSessions[UserSessions[cookie.Value]] = &UserTransQuestions
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
	QuestionPassed *Question, usersession *string, userquestionnumber *int, questionslist *[]Question, QuestionType string) {
	r.ParseForm()

	// Handle User Session values and button
	UserSession := *usersession
	UserQuestionNumber := *userquestionnumber

	var button_val string
	if UserQuestionNumber == 4 {
		button_val = "Submit and check my grade"
	} else {
		button_val = "Next Question"
	}

	// If user has submitted an answer
	if QuestionType == "vocab_qns" {
		VocabQuestions := *questionslist
		VocabUserSessions[UserSession].CurrentQuestion = UserQuestionNumber + 1
		CacheVocabAnswer(&w, r, VocabQuestions, QuestionType, QuestionPassed,
			UserSession, button_val, UserQuestionNumber)
	} else if QuestionType == "trans_qns" {
		TransQuestions := *questionslist
		TransUserSessions[UserSession].CurrentQuestion = UserQuestionNumber + 1
		CacheTransAnswer(&w, r, TransQuestions, QuestionType, QuestionPassed,
			UserSession, button_val, UserQuestionNumber)
	}
}

func DisplayGrade(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("riinsan")
	if err != nil || cookie == nil {
		cookie = CreateNewCookie(w, r)
	}
	UserSession := UserSessions[cookie.Value]
	if r.URL.Path == "/grade/vocab" {
		VocabUserSession := VocabUserSessions[UserSession]
		UserVocabWrongAnswerIndex := []int{}
		for question_number := range VocabUserSession.WrongAnswers {
			UserVocabWrongAnswerIndex = append(UserVocabWrongAnswerIndex, question_number-1)
		}
		QuestionList := VocabUserSession.QuestionList
		UserVocabWrongAnswers := map[string][]string{}
		fmt.Println("This is it", UserVocabWrongAnswerIndex)
		for _, val := range UserVocabWrongAnswerIndex {
			UserVocabWrongAnswers["wrong_question"] = append(UserVocabWrongAnswers["wrong_question"], QuestionList[val].Question)
			UserVocabWrongAnswers["image_link"] = append(UserVocabWrongAnswers["image_link"], QuestionList[val].Image_link)
			UserVocabWrongAnswers["correct_answer"] = append(UserVocabWrongAnswers["correct_answer"], strings.Join(QuestionList[val].Question_answer, ", "))
		}
		TPL.ExecuteTemplate(w, "grade.html", UserVocabWrongAnswers)
	} else if r.URL.Path == "/grade/trans" {
		TransUserSession := TransUserSessions[UserSession]
		UserTransAnswerIndex := []int{}
		for question_number := range TransUserSession.WrongAnswers {
			UserTransAnswerIndex = append(UserTransAnswerIndex, question_number-1)
		}
		QuestionList := TransUserSession.QuestionList
		UserTransAnswers := map[string][]string{}
		fmt.Println("Trans Questions asked", UserTransAnswerIndex)
		for _, val := range UserTransAnswerIndex {
			UserTransAnswers["user_trans_answer"] = append(UserTransAnswers["user_trans_answer"], QuestionList[val].Question)
			UserTransAnswers["grammar_to_use"] = append(UserTransAnswers["grammar_to_use"], QuestionList[val].TransHelperGrammar)
			UserTransAnswers["correct_answer"] = append(UserTransAnswers["correct_answer"], strings.Join(QuestionList[val].Question_answer, "/\n "))
		}
		TPL.ExecuteTemplate(w, "grade.html", UserTransAnswers)
	}
}

func TakeANewTest(w *http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("riinsan")
	if err != nil {
		cookie = CreateNewCookie(*w, r)
	}
	UserSession := UserSessions[cookie.Value]
	_, UserVocabSessionCheck := VocabUserSessions[UserSession]
	_, UserTransSessionCheck := TransUserSessions[UserSession]
	if UserVocabSessionCheck {
		delete(VocabUserSessions, UserSession)
	}
	if UserTransSessionCheck {
		delete(TransUserSessions, UserSession)
	}
}

func GenerateNewQuestions(QuestionListLength int, QuestionList []Question) []Question {
	var RandomQuestionNumbers []int

	//TODO: CHECK IF THIS WORKS
	rand.NewSource(time.Now().Unix())

	// Generate 5 random, non-overlapping question numbers
	for len(RandomQuestionNumbers) < 5 {
		r_n := rand.Intn(QuestionListLength)
		if CheckIfOverlappingQuestionNumber(RandomQuestionNumbers, r_n) {
			RandomQuestionNumbers = append(RandomQuestionNumbers, r_n)
		}
	}
	fmt.Println("this is the randomly generated question index", RandomQuestionNumbers)

	SelectedQuestions := []Question{}
	for _, val := range RandomQuestionNumbers {
		SelectedQuestions = append(SelectedQuestions, QuestionList[val])
	}
	return SelectedQuestions
}
