package main

import (
	"errors"
	"fmt"
	"time"
	"bytes"
	"crypto/sha1"
	"encoding/json"
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

	myLogger.Debug("The administrator is [%x]", adminCert)

	currentDateTime := time.Now().Format(time.UnixDate)

	var adminUserInfo []UserInfo
	adminUserInfo = append(adminUserInfo, UserInfo{"Name", "Seksit Disaro"})

	adminUser := &User{adminCert, ROLE_ADMIN, adminUserInfo ,currentDateTime }

	compoundKey, _ := t.createCompoundKey("User", []string{ROLE_ADMIN, string(adminUser.Key)})
	adminJSONBytes, _ := json.Marshal(adminUser)

	stub.PutState(compoundKey, adminJSONBytes)

	return nil, nil
}

func (t *PatientChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, params := stub.GetFunctionAndParameters()

	if function == "create_admin" {
		fmt.Printf("Create Admin invoke!")
		t.createAdmin(stub, params)
	} else if function == "create_patient" {
		fmt.Printf("Create Patient invoke!")
		t.createPatient(stub, params)
	} else if function == "create_hcp" {
		fmt.Printf("Create HCP invoke!")
		t.createHealthcareProvider(stub, params)
	} else if function == "list_users" {
		fmt.Printf("List users invoke!")
		t.listPatients(stub)
	} else if function == "append_medical_data" {
		fmt.Printf("Append medical data invoke!")
	} else if function == "request_permission" {
		fmt.Printf("Request acess permission invoke!")
	} else if function == "read_medical_data" {
		fmt.Printf("Read medical data invoke!")
	} else if function == "grant_permission" {
		fmt.Printf("Grant access permission invoke!")
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *PatientChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

func (t *PatientChaincode) listPatients(stub shim.ChaincodeStubInterface) ([]byte,error) {
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

// Create user for this chain based on input role parameter. Return Unique ID for created user.
func (t *PatientChaincode) createUserGeneric(stub shim.ChaincodeStubInterface, role string, params []string) ([]byte,error) {
	cert, _ := stub.GetCallerMetadata()
	currentDateTime := time.Now().Format(time.UnixDate)
	inputHashData :=  string(cert[:]) + time.Now().Format(time.UnixDate)

	// SHA1 Hash function for key
	h := sha1.New()
	h.Write([]byte(inputHashData))
	bs := h.Sum(nil)

	compoundKey, _ := t.createCompoundKey("User", []string{role, string(bs)})

	// Create user object
	var userInfo []UserInfo
	userInfo = append(userInfo, UserInfo{"Name", "Hello World"})
	userObj := &User{bs, role, userInfo ,currentDateTime}
	userObjJsonBytes, _ := json.Marshal(userObj)

	stub.PutState(compoundKey, userObjJsonBytes)
	fmt.Printf("User created, Role = "+role+" Compounded Key = "+compoundKey)

	return userObjJsonBytes, nil
}

func (t *PatientChaincode) createAdmin(stub shim.ChaincodeStubInterface, params []string) ([]byte, error) {
	cert, _ := stub.GetCallerMetadata()
	if t.checkRole(string(cert)) == ROLE_ADMIN {
		return t.createUserGeneric(stub, ROLE_ADMIN, params)
	} else {
		fmt.Printf("")
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
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}