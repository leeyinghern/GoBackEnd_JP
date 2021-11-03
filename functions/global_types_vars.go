package functions

import "html/template"

type Questions struct {
	QuestionList    []Question `json:"questions"`
	CurrentQuestion int        //used to track the current question user is on
}

type Question struct {
	Question_number    int      `json:"question_number"`
	Image_link         string   `json:"image"`
	Question           string   `json:"question"`
	Question_answer    []string `json:"answers"`
	QuestionType       string   `json:"question_type"`
	TransHelperGrammar string   `json:"grammar_to_use"`
	TransHelperWords   []string `json:"helper_words"`
	UserAnswer         string   //User's answer for this question
}

var TPL *template.Template

var UserSessions = map[string]string{}          // Cookie UUID, encrypted_email
var VocabUserSessions = map[string]*Questions{} // encrypted_email, VocabUserSession=Questions{list_of_qns,currentquestion}
var TransUserSessions = map[string]*Questions{} // encrypted_email, TransUserSession=Questions{list_of_qns,currentquestion}

type GradeDataToHTML struct {
	QuestionType string
	QuestionData []Question
}
