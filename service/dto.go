package service

import "fmt"

type Generic struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Source string `json:"-"`
}

type SubmitForm struct {
	Name         string `form:"name"`
	Email        string `form:"email"`
	MobileNo     string `form:"mobileNo"`
	Company      string `form:"company"`
	Message      string `form:"message"`
	LinkedInUrl  string `form:"linkedInUrl"`
	PortfolioUrl string `form:"portfolioUrl"`
	
}

func (g Generic) Error() string {
	return fmt.Sprintf("code: %d, msg: %s\n", g.Code, g.Msg)
}
