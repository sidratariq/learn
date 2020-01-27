package implementation

import (
	entity "Model"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================
// RegisterPatient - create a new Provider, store into chaincode state
// ============================================================
func (u *User) RegisterProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- start register patient")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	providerId := strings.ToLower(args[0])
	providerEHR := strings.ToLower(args[1])
	providerEHRUrl := strings.ToLower(args[2])
	firstname := strings.ToLower(args[3])
	lastname := strings.ToLower(args[4])
	speciality := strings.ToLower(args[5])

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Only allow when the user is a provider
	if userrole == "Doctor" {

		// ==== Check if provider already exists ====
		providerData, err := stub.GetState(providerId)
		if err != nil {
			return shim.Error("Fails to get provider " + err.Error())
		}
		if providerData == nil {
			fmt.Println("Provider doesn't exist so a new provider")
		} else if providerData != nil {
			fmt.Println("This provider already exists or the ID is already assigned to some other user")
			return shim.Error("This provider already exists or the ID is already assigned to some other user")
		}

		//==== Create Provider object and marshal to JSON ====
		objectType := "Provider"
		provider := &entity.Provider{objectType, providerId, providerEHR, providerEHRUrl, firstname, lastname, speciality}
		//fmt.Println(Provider.firstname)

		providerJSONasBytes, err := json.Marshal(provider)
		if err != nil {
			return shim.Error(err.Error())
		}
		//Alternatively, build the marble json string manually if you don't want to use struct marshalling
		//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
		//marbleJSONasBytes := []byte(str)

		//Alternatively, build the marble json string manually if you don't want to use struct marshalling
		//patientJSONasString := `{"docType":"Patient",  "patientId": "` + patientId + `", "patientSSN": "` + patientSSN + `", "patientUrl": ` + patientUrl + `, "firstname": "` + firstname + `, "DOB": "` + DOB + `, "email": "` + email + `, "mobile": "` + mobile + `"}`
		//patientJSONasBytes := []byte(patientJSONasString)

		// === Save Provider to state ===

		//err = stub.PutPrivateData("patientDetails", providerId, providerJSONasBytes)
		err = stub.PutState(providerId, providerJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		//  ==== Index the Provider to enable name-based range queries, e.g. return all Patients ====
		//  An 'index' is a normal key/value entry in state.
		//  The key is a composite key, with the elements that you want to range query on listed first.
		//  In our case, the composite key is based on indexName~color~name.
		//  This will enable very efficient state range queries based on composite keys matching indexName~color~*

		//Create Firstname, ID composite key
		indexName := "firstname~providerId"
		fnameIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{provider.ProviderFirstname, provider.ProviderLastname, provider.ProviderId})
		if err != nil {
			return shim.Error(err.Error())
		}
		//Create Lastname, ID composite key
		indexName = "lastname~providerId"
		lnameIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{provider.ProviderLastname, provider.ProviderId})
		if err != nil {
			return shim.Error(err.Error())
		}
		//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		stub.PutState(fnameIDIndexKey, value) //Put firstname,ID composite key in state
		stub.PutState(lnameIDIndexKey, value) //Put lastname,ID composite key in state

		// ==== Marble saved and indexed. Return success ====
		//fmt.Println("- end register patient")
		return shim.Success(nil)
	}
	return shim.Error("Unauthorized! Only Provider can register!")
}

func (u *User) UpdateProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) < 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	//set the variables
	providerId := strings.ToLower(args[0])
	providerEHR := strings.ToLower(args[1])
	providerEHRUrl := strings.ToLower(args[2])
	firstname := strings.ToLower(args[3])
	lastname := strings.ToLower(args[4])
	speciality := strings.ToLower(args[5])

	//check in the state that if the provider exists or not
	providerAsByte, err := stub.GetState(providerId)
	if err != nil {
		fmt.Println("Provider doesn't exist so a new patient")
		return shim.Error("This provider doesn't exist so can't update")
	} else if providerAsByte != nil {
		fmt.Println("This provider already exists or the ID is already assigned to some other user")
	}

	//Unmarshal the provider details acquired and store them in provider object
	var provider entity.Provider
	err = json.Unmarshal(providerAsByte, &provider)

	//Set updated values into the provider object retrieved from state
	provider.ObjectType = "Provider"
	provider.ProviderId = providerId
	provider.ProviderEHR = providerEHR
	provider.ProviderEHRURL = providerEHRUrl
	provider.ProviderFirstname = firstname
	provider.ProviderLastname = lastname
	provider.Speciality = speciality

	//Marshal the object into json format
	providerJSONasBytes, err := json.Marshal(provider)

	//Save the object into state
	err = stub.PutState(providerId, providerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (u *User) GetProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\"}}", id)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (u *User) GetProviderByFirstName(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	firstname := args[0]
	indexName := "firstname~providerId"

	resultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{firstname})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		returnedProviderId := compositeKeyParts[1]

		ProviderAsBytes, _ := stub.GetState(returnedProviderId)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedProviderId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(ProviderAsBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryGetProviderByFirstName:\n%s\n", buffer.String())
	queryResult := buffer.Bytes()
	return shim.Success(queryResult)
}

func (u *User) GetProviderByLastName(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lastname := args[0]
	indexName := "lastname~providerId"

	resultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{lastname})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		returnedProviderId := compositeKeyParts[1]

		ProviderAsBytes, _ := stub.GetState(returnedProviderId)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedProviderId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(ProviderAsBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryGetProviderByLastName:\n%s\n", buffer.String())
	queryResult := buffer.Bytes()
	return shim.Success(queryResult)
}

// func (u *User) UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {

// 	if len(args) < 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	logger := logging.NewLogger("log")

// 	logger.Debugf("log is debugging: %s", args[0])

// 	var patientDetails entity.PatientDetailsUnmarshal
// 	var val []byte = []byte("`" + args[0] + "`")

// 	s, err1 := strconv.Unquote(string(val))

// 	if err1 != nil {
// 		return shim.Error("Error in unquote -->\n" + s + "-->\n" + args[0] + err1.Error())
// 	}

// 	err := json.Unmarshal([]byte(s), &patientDetails)

// 	if err != nil {
// 		return shim.Error("Error in unmarshal input json -->\n" + s + "-->\n" + args[0] + err.Error())
// 	}

// 	patientId, err := getAttribute(stub, "id")
// 	if err != nil {
// 		return shim.Error("Fail to get Attribute from private DB " + err.Error())
// 	}

// 	patientDetailsAsBytes, err := stub.GetPrivateData("patientDetails", "123")

// 	if err != nil {
// 		return shim.Error("Fail to get patint from private DB " + err.Error())
// 	}

// 	var patientDetailsDB entity.PatientDetails

// 	err = json.Unmarshal(patientDetailsAsBytes, &patientDetailsDB) //unmarshal it aka JSON.parse()
// 	if err != nil {
// 		return shim.Error("ID" + err.Error())
// 	}

// 	patientDetailsDB.Allergies.ProviderConsent = append(patientDetailsDB.Allergies.ProviderConsent, patientDetails.Allergies.ProviderConsent[0])
// 	patientDetailsDB.Immunization.ProviderConsent = append(patientDetailsDB.Immunization.ProviderConsent, patientDetails.Immunization.ProviderConsent[0])
// 	patientDetailsDB.Medications.ProviderConsent = append(patientDetailsDB.Medications.ProviderConsent, patientDetails.Medications.ProviderConsent[0])
// 	patientDetailsDB.PastMedicalHx.ProviderConsent = append(patientDetailsDB.PastMedicalHx.ProviderConsent, patientDetails.PastMedicalHx.ProviderConsent[0])

// 	PatientDetailsJSONasBytes, err := json.Marshal(&patientDetailsDB)

// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	// === Save patientDetails to state ===
// 	// err = stub.PutPrivateData("patientDetails", patientId, PatientDetailsJSONasBytes)
// 	// if err != nil {
// 	// 	return shim.Error( "Error in put private data in one org" +err.Error())
// 	// }

// 	err = stub.PutPrivateData("patientDetailsIn2Orgs", patientId, PatientDetailsJSONasBytes)
// 	if err != nil {
// 		return shim.Error("Error in put private data in two org " + err.Error())
// 	}

// 	return shim.Success([]byte("Success"))
// }

func (u *User) RegisterProviderRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	fname := strings.ToLower(args[0])
	lname := strings.ToLower(args[1])
	dob := strings.ToLower(args[2])
	ssn := strings.ToLower(args[3])
	stime := args[4]
	etime := args[5]

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\",\"patientssn\":\"%s\"}}", fname, lname, dob, ssn)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	if userrole == "Doctor" {

		var tempArray []entity.PatientUnmarshal
		err = json.Unmarshal(queryResults, &tempArray)
		if err != nil {
			return shim.Error(err.Error())
		}

		var key string
		for _, temp := range tempArray {

			key = temp.Key

		}

		//Get Current Provider
		providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
		if err != nil {
			return shim.Error("Fails to get provider id " + err.Error())
		}

		providerAsByte, err := stub.GetState(providerId) //Get the Doctor's details from the state.
		if err != nil {
			return shim.Error("Fails to get provider: " + err.Error())
		}

		var provider entity.Provider
		err = json.Unmarshal(providerAsByte, &provider) //Store the Doctor's object into provider.

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		//Get Patient
		var patient entity.Patient
		patientAsByte, err := stub.GetState(key)
		err = json.Unmarshal(patientAsByte, &patient)

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		//Search in patient's providerconsent list that if the current provider already has consent or not.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {
				return shim.Error("Provider already has consent.")
			}
		}

		//Search in patient's providerrequest list that if the current provider is already registered or not.
		for i := 0; i < len(patient.ProviderRequest); i++ {
			tempRequest := patient.ProviderRequest[i].Provider.ProviderId
			if tempRequest == provider.ProviderId {
				return shim.Error("The Provider's Access request has already been registered!")
			}
		}

		//if not, then add his request to the list.
		var defaultRequest entity.Consent
		defaultRequest.ObjectType = "ProviderRequest"
		defaultRequest.Provider = provider
		defaultRequest.StartTime = stime
		defaultRequest.EndTime = etime

		patient.ProviderRequest = append(patient.ProviderRequest, defaultRequest)

		//Marshal the object into json format
		patientJSONasBytes, err := json.Marshal(patient)

		//Save the object into state
		err = stub.PutState(key, patientJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}
	return shim.Error("Unauthorized! Only Doctor can Register Provider Request")

}

func (u *User) RemoveProviderRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fname := strings.ToLower(args[0])
	lname := strings.ToLower(args[1])
	dob := strings.ToLower(args[2])
	ssn := strings.ToLower(args[3])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\",\"patientssn\":\"%s\"}}", fname, lname, dob, ssn)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	if userrole == "Doctor" {

		var tempArray []entity.PatientUnmarshal
		err = json.Unmarshal(queryResults, &tempArray)
		if err != nil {
			return shim.Error(err.Error())
		}

		var key string
		for _, temp := range tempArray {

			key = temp.Key

		}

		//Get Current Provider
		providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
		if err != nil {
			return shim.Error("Fails to get provider id " + err.Error())
		}

		providerAsByte, err := stub.GetState(providerId) //Get the Doctor's details from the state.
		if err != nil {
			return shim.Error("Fails to get provider: " + err.Error())
		}

		var provider entity.Provider
		err = json.Unmarshal(providerAsByte, &provider) //Store the Doctor's object into provider.

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		//Get Patient
		var patient entity.Patient
		patientAsByte, err := stub.GetState(key)
		err = json.Unmarshal(patientAsByte, &patient)

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		//Search in patient's providerconsent list that if the current provider already has consent or not.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {
				return shim.Error("Provider already has consent.")
			}
		}

		a := 0 //flag
		//Search in patient's providerrequest list that if the current provider is already registered or not.
		for i := 0; i < len(patient.ProviderRequest); i++ {
			tempRequest := patient.ProviderRequest[i].Provider.ProviderId
			if tempRequest == provider.ProviderId {
				//remove the request from list
				patient.ProviderRequest = append(patient.ProviderRequest[:i], patient.ProviderRequest[i+1:]...)
				a = 1
			}
		}

		if a != 1 {
			return shim.Error("Provider doesn't have any requests!")
		}
		//Marshal the object into bytes
		patientJSONasBytes, err := json.Marshal(patient)

		//Save the object into state
		err = stub.PutState(key, patientJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}
	return shim.Error("Unauthorized! Only Doctor can Remove Request")

}

func (u *User) GetProviderAccessStatusForProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// "bob"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fname := strings.ToLower(args[0])
	lname := strings.ToLower(args[1])
	dob := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\"}}", fname, lname, dob)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	//Get Current Provider based on token
	providerId, err := getAttribute(stub, "id") //Get the current provider ID from identity.
	if err != nil {
		return shim.Error("Fails to get id " + err.Error())
	}

	//Only allow when the user is a provider
	if userrole == "Doctor" {

		providerAsByte, err := stub.GetState(providerId) //Get the provider's details from the state.
		if err != nil {
			return shim.Error("Fails to get provider: " + err.Error())
		}

		var provider entity.Provider
		err = json.Unmarshal(providerAsByte, &provider) //Store the Provider's object into provider.

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		//Get Patient
		var patient entity.Patient
		patientAsByte, err := stub.GetState(key)
		err = json.Unmarshal(patientAsByte, &patient)

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		a := 0 //Flag
		//Search in patient's providerconsent list that if the current provider is in the list of consent, if in list then append to another list of consent and return that list.
		var tempConsentlist []entity.Consent
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {
				tempConsentlist = append(tempConsentlist, patient.ProviderConsent[i])
				a = 1
			}
		}

		if a != 1 {
			return shim.Error("Provider doesn't exist in consent list")
		}

		//Marshal them into bytes in order to return them.
		ConsentJSONasBytes, err := json.Marshal(tempConsentlist)
		if err != nil {
			shim.Error("Failed to marshal consent" + err.Error())
		}

		return shim.Success(ConsentJSONasBytes)

	}
	return shim.Error("Unauthorized! Only Provider can get Provider Access Status!")
}
