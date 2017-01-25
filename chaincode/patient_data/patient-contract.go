package main

import (
	"errors"
	"fmt"
	"time"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/base64"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	//"github.com/zekzit/hfc-testbed/chaincode/patient_data/const"
)

var myLogger = logging.MustGetLogger("patient_mgm")

type PatientChaincode struct {
}

func (t *PatientChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	// Set the admin
	adminCert, errGetData := stub.GetCallerMetadata()
	if errGetData != nil {
		myLogger.Debug("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	if len(adminCert) == 0 {
		myLogger.Debug("Invalid admin certificate. Empty.")
		return nil, errors.New("Invalid admin certificate. Empty.")
	}

	currentDateTime := time.Now().Format(time.UnixDate)

	var adminUserInfo []UserInfo
	adminUserInfo = append(adminUserInfo, UserInfo{"Name", "Seksit Disaro"})

	adminUser := User{adminCert, ROLE_ADMIN, adminUserInfo ,currentDateTime }
	compoundKey, _ := t.createCompoundKey("User", []string{ROLE_ADMIN, base64.URLEncoding.EncodeToString(adminUser.Key)})

	myLogger.Debug("The administrator is [%x]", compoundKey)

	adminJSONBytes, _ := json.Marshal(adminUser)
	stub.PutState(compoundKey, adminJSONBytes)

	// Initialize user list array
	var userListArray []string
	userListArray = append(userListArray, compoundKey)
	fmt.Println("UserListArray =",userListArray)
	fmt.Println("Number of user =",len(userListArray),"(should be 1)")
	jsonAsBytes, _ := json.Marshal(userListArray)
	errPutState := stub.PutState("UserList",jsonAsBytes)
	fmt.Println("errPutState =",errPutState)

	return nil, nil
}

func (t *PatientChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, params := stub.GetFunctionAndParameters()

	if function == "init" {
		fmt.Println("Init again!")
		return t.Init(stub)
 	}else if function == "create_admin" {
		fmt.Println("Create Admin invoke!")
		return t.createAdmin(stub, params)
	} else if function == "create_patient" {
		fmt.Println("Create Patient invoke!")
		return t.createPatient(stub, params)
	} else if function == "create_hcp" {
		fmt.Println("Create HCP invoke!")
		return t.createHealthcareProvider(stub, params)
	} else if function == "list_users" {
		fmt.Println("List users invoke!")
		return t.listPatients(stub)
	} else if function == "append_medical_data" {
		fmt.Println("Append medical data invoke!")
	} else if function == "request_permission" {
		fmt.Println("Request acess permission invoke!")
	} else if function == "read_medical_data" {
		fmt.Println("Read medical data invoke!")
	} else if function == "grant_permission" {
		fmt.Println("Grant access permission invoke!")
	} else if function=="get_admin" {
		fmt.Println("GET ADMIN CALLED!! (Is that that serious?)")
		result, err := stub.GetState("ADMIN5MEUCIQCV3LNIwSEVtkk5pVuztd3UlmDzSYkkY0rqEaqOeNf1EAIgICetPpM88ocF11jDbqiIriNzwsIWB955xoeH5gJSq8A=")
		fmt.Println("Result = ",result)
		fmt.Println("Error = ", err)
		adminUser := User{}
		json.Unmarshal(result, adminUser)
		fmt.Println(adminUser);
		return result, err
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *PatientChaincode) Query(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, _ := stub.GetFunctionAndParameters()

	if function == "list_users" {
		fmt.Println("List users invoke!")
		return t.listPatients(stub)
	}

	return nil, errors.New("Received unknown query function")
}

func (t *PatientChaincode) listPatientsOld(stub shim.ChaincodeStubInterface) ([]byte,error) {
	objectType := "User"
	partialKeysForQuery := []string{ROLE_PATIENT}

	keysIter, _ := t.partialCompoundKeyQuery(stub, objectType, partialKeysForQuery)
	defer keysIter.Close()

	var users []User
	for keysIter.HasNext() {
		_, userJSONBytes, _ := keysIter.Next()
		user := User{}
		json.Unmarshal(userJSONBytes, &user)
		users = append(users, user)
	}

	return json.Marshal(users)
}

func (t *PatientChaincode) listPatients(stub shim.ChaincodeStubInterface) ([]byte,error) {
	var userListArray []string
	rawArray, err := stub.GetState("UserList")

	fmt.Println("UserList state =", rawArray)
	fmt.Println("err =", err)

	errUnmarshal := json.Unmarshal(rawArray,&userListArray)
	fmt.Println("Unmarshal error =",errUnmarshal)
	fmt.Println("Unmarshal userListArray =",userListArray)
	fmt.Println("Number of users =", len(userListArray))

	for i := 0; i<len(userListArray); i++ {
		fmt.Println(userListArray[i])
	}

	return rawArray, err
}

// Create user for this chain based on input role parameter. Return Unique ID for created user.
func (t *PatientChaincode) createUserGeneric(stub shim.ChaincodeStubInterface, role string, params []string) ([]byte,error) {
	cert, _ := stub.GetCallerMetadata()
	currentDateTime := time.Now().Format(time.UnixDate)
	inputHashData :=  string(cert[:]) + time.Now().Format(time.UnixDate)

	// SHA1 Hash function for key
	h := sha1.New()
	h.Write([]byte(inputHashData))
	bs := h.Sum(nil)

	compoundKey, _ := t.createCompoundKey("User", []string{role, base64.URLEncoding.EncodeToString(bs)})

	// Create user object
	var userInfo []UserInfo
	userInfo = append(userInfo, UserInfo{"Name", "Hello World"})
	userObj := User{bs, role, userInfo ,currentDateTime}
	userObjJsonBytes, _ := json.Marshal(userObj)

	stub.PutState(compoundKey, userObjJsonBytes)
	fmt.Println("User created, Role = "+role+" Compounded Key = "+compoundKey+"\n")

	// Append user list array
	var userListArray []string
	rawArray, err := stub.GetState("UserList")
	json.Unmarshal(rawArray,&userListArray)
	userListArray = append(userListArray, compoundKey)
	jsonAsBytes, _ := json.Marshal(userListArray)
	err = stub.PutState("UserList", jsonAsBytes)

	return userObjJsonBytes, err
}

func (t *PatientChaincode) createAdmin(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	cert, _ := stub.GetCallerMetadata()
	if t.checkRole(string(cert)) == ROLE_ADMIN {
		return t.createUserGeneric(stub, ROLE_ADMIN, params)
	} else {
		fmt.Println("")
		return nil, errors.New("Cannot read caller's data or caller has no sufficient permission.")
	}
}

func (t *PatientChaincode) createPatient(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	return t.createUserGeneric(stub, ROLE_PATIENT, params)
}

func (t *PatientChaincode) createHealthcareProvider(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	return t.createUserGeneric(stub, ROLE_HCP, params)
}

// Check key for role then return ROLE constant
func (t *PatientChaincode) checkRole(key string) string {
	return ROLE_ADMIN
}

// Check patient key, HCP key and then check for authorization return AuthConst
func (t *PatientChaincode) checkAuthorization(patientKey string, hcpKey string, checkAuthConst int) int {
	return AUTHORIZE_NOTHING
}

func (t *PatientChaincode) readPatientMedicalData(patientKey string, readerKey string, medicalDataType string, limit int) []MedicalData {
	return nil
}

func (t *PatientChaincode) appendPatientMedicalData(patientKey string, appenderKey string,
	medicalDataType string, medicalDataId string, medicalData []byte) bool {
	return true
}

func (t *PatientChaincode) appendAuditLog(subjectKey string, actionConst string, objectKey string,
	result string, reasonOrRemark string) bool {
	return true
}

func (t *PatientChaincode) requestAccessPermission(patientKey string, healthcareProviderKey string,
	requestAuthConst int) bool {
	return true
}

func (t *PatientChaincode) grantAccessPermission(patientKey string, healthcareProviderKey string,
	grantAuthConst int) bool {
	return true
}

// ============================================================================================================================
// Utility functions (may become chaincode APIs)
// ============================================================================================================================

func (t *PatientChaincode) createCompoundKey(objectType string, keys []string) (string, error) {
	var keyBuffer bytes.Buffer
	keyBuffer.WriteString(objectType)
	for _, key := range keys {
		keyBuffer.WriteString(strconv.Itoa(len(key)))
		keyBuffer.WriteString(key)
	}
	return keyBuffer.String(), nil
}

func (t *PatientChaincode) partialCompoundKeyQuery(stub shim.ChaincodeStubInterface, objectType string, keys []string) (shim.StateRangeQueryIteratorInterface, error) {
	// TODO - call RangeQueryState() based on the partial keys and pass back the iterator

	keyString, _ := t.createCompoundKey(objectType, keys)
	fmt.Println("Query range from "+keyString+"1 - "+keyString+":")
	keysIter, err := stub.RangeQueryState(keyString+"1", keyString+":")
	if err != nil {
		return nil, fmt.Errorf("Error fetching rows: %s", err)
	}
	defer keysIter.Close()

	return keysIter, err
}

func main() {
	err := shim.Start(new(PatientChaincode))
	if err != nil {
		fmt.Println("Error starting Patient chaincode: %s", err)
	}
}
