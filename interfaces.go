package Interfaces

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type User struct {
}

//Repository repository interface
type InterfacePatient interface {
	RegisterPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientByInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response
	UpdatePatientById(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientByFirstName(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientById(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientByLastName(stub shim.ChaincodeStubInterface, args []string) pb.Response
	AllowConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response
	UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response
	SearchProviderRequestsByPatientInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response
	SearchProviderRequestsByPatientSSN(stub shim.ChaincodeStubInterface, args []string) pb.Response
	RevokeProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderAccessStatusForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientAllInfoBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetAllProviderAccessByPatientSSN(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetAllProviderAccessByPatientInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response
}

type InterfaceProvider interface {
	RegisterProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response
	UpdateProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderByFirstName(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderByLastName(stub shim.ChaincodeStubInterface, args []string) pb.Response
	RegisterProviderRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response
	RemoveProviderRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderAccessStatusForProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response
}

