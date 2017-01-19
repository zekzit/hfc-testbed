package patient_data

import (
	"errors"
	"fmt"
	"time"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	//"github.com/zekzit/hfc-testbed/chaincode/patient_data/const"
)

var myLogger = logging.MustGetLogger("patient_mgm")

type PatientChaincode struct {
}

func (t *PatientChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	// Create User table
	errUser := stub.CreateTable("User", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Key", Type: shim.ColumnDefinition_BYTES, Key: true},
		&shim.ColumnDefinition{Name: "Role", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Info", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CreatedDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if errUser != nil {
		return nil, errors.New("Failed creating User table.")
	}

	// Create AuditLog table
	errAudit := stub.CreateTable("AuditLog", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Action", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Subject", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Object", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Result", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Reason", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if errAudit != nil {
		return nil, errors.New("Failed creating AuditLog table.")
	}

	// Create Authorization table
	errAuth := stub.CreateTable("Authorization", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Patient", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "HCP", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "RequestedPermission", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "GrantedPermission", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "RequestedDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "GrantedDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if errAuth != nil {
		return nil, errors.New("Failed creating Authorization table.")
	}

	// Create MedicalData table
	errMedData := stub.CreateTable("MedicalData", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Patient", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "HCP", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "MedicalDataType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "MedicalDataId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "MedicalData", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CreatedDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ModifiedDateTime", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if errMedData != nil {
		return nil, errors.New("Failed creating MedicalData table.")
	}

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

	var currentDateTime = time.Now().Format(time.UnixDate)

	ok, errInsertRow := stub.InsertRow("User",shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bytes{Bytes: adminCert}},
			&shim.Column{Value: &shim.Column_String_{String_: ROLE_ADMIN}},
			&shim.Column{Value: &shim.Column_String_{String_: ""}},
			&shim.Column{Value: &shim.Column_String_{String_: currentDateTime}}},
	})

	if ok && errInsertRow == nil {
		return nil, errors.New("Asset was already assigned.")
	}

	return nil, nil
}

func (t *PatientChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	function, _ := stub.GetFunctionAndParameters()

	if function == "assign" {
		fmt.Printf("Assign invoke!")
	} else if function == "transfer" {
		fmt.Printf("Transfer invoke!")
	}
	return nil, errors.New("Received unknown function invocation")
}

func (t *PatientChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

func main() {
	err := shim.Start(new(PatientChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}