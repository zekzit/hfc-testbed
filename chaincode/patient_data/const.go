package main

const ROLE_ADMIN = "ADMIN"
const ROLE_PATIENT = "PATIENT"
const ROLE_HCP = "HCP"

const AUTHORIZE_NOTHING = 0
const AUTHORIZE_READABLE = 1
const AUTHORIZE_APPENDABLE = 2

const MEDICALDATATYPE_LAB = "LAB"
const MEDICALDATATYPE_DRUG = "DRUG"
const MEDICALDATATYPE_FINANCIAL = "FINANCIAL"

const ACTION_REQUEST_PERMISSION = 0
const ACTION_GRANT_PERMISSION = 1
const ACTION_READ_DATA = 2
const ACTION_APPEND_DATA = 4

type User struct {
	Key			[]byte		`json:"key"`
	Role			string		`json:"role"`
	Info			[]UserInfo	`json:"info"`
	CreatedDateTime		string 		`json:"create_datetime"`
}

type UserInfo struct {
	Key	string	`json:"key"`
	Value	string	`json:"value"`
}

type AuditLog struct {
	Action	string	`json:"action"`
	Subject	string	`json:"subject"`
	Object	string	`json:"object"`
	Result	string	`json:"result"`
	Reason	string	`json:"reason"`
}

type Authorization struct {
	Patient	string	`json:"patient"`
	HCP	string	`json:"hcp"`
	RequestedPermission	string	`json:"requested_permission"`
	GrantedPermission	string	`json:"granted_permission"`
	RequestedDateTime	string	`json:"requested_datetime"`
	GrantedDateTime	string	`json:"granted_datetime"`
}

type MedicalData struct {
	Patient	string	`json:"patient"`
	HCP	string	`json:"hcp"`
	MedicalDataType	string	`json:"medical_data_type"`
	MedicalDataId	string	`json:"medical_data_id"`
	MedicalData	string	`json:"medical_data"`
	CreatedDateTime	string	`json:"created_datetime"`
	ModifiedDateTime	string	`json:"modified_datetime"`
}