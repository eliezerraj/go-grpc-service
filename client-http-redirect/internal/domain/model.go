package domain

import (
//	"time"
	"github.com/golang-jwt/jwt/v4"
)

type Credential struct {
	Token			string `json:"token,omitempty"`
}

type ResponseMessage struct {
	Message			string `json:"message,omitempty"`
}

type JwtData struct {
	Username	string 		`json:"username"`
	Scope		[]string 	`json:"scope"`
	jwt.RegisteredClaims
}

type InfoPod struct {
	PodName			string `json:"pod_name,omitempty"`
	ApiVersion		string `json:"version,omitempty"`
	OsPid			string `json:"os_pid,omitempty"`
	Ip				string `json:"ip,omitempty"`
	PodPath			string `json:"pod_path,omitempty"`
	Port			int 	`json:"port,omitempty"`
	JwtKey			[]byte `json:"jwt_key,omitempty"`
}

type Secret struct {
	Username		string	`json:username`
	Password		string	`json:"password"`
}

type Balance struct {
	Id				string		`json:"id,omitempty"`
	Account			string		`json:"account,omitempty"`
	Amount			int			`json:"amount,omitempty"`
	Date_Balance  	string	`json:"date_balance,omitempty"`
	Description		string		`json:"description,omitempty"`
}