package functions

import "html/template"

type Questions struct {
	QuestionList    []Question     `json:"questions"`
	QuestionType    string         `json:"question_type"`
	WrongAnswers    map[int]string // Wrong answer takes in assigned qn number + user_ans
	CurrentQuestion int            //used to track the current question user is on
}

type Question struct {
	Question_number int      `json:"question_number"`
	Image_link      string   `json:"image"`
	Question        string   `json:"question"`
	Question_answer []string `json:"answers"`
}

var TPL *template.Template

var UserSessions = map[string]string{}          // Cookie UUID, encrypted_email
var VocabUserSessions = map[string]*Questions{} // encrypted_email, VocabUserSession
var TransUserSessions = map[string]*Questions{} // encrypted_email, TransUserSession
