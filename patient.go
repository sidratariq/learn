package implementation

import (
	inf "Interfaces"
	entity "Model"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type User inf.User

// ============================================================
// RegisterPatient - create a new Patient, store into chaincode state
// ============================================================
func (u *User) RegisterPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 13"+ args[0]+args[1]+args[2]+args[3]+args[4]+args[5]+args[6]+args[7]+args[8]+args[9]+args[10]+args[11]+args[12])
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
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return shim.Error("9th argument must be a non-empty string")
	}
	if len(args[9]) <= 0 {
		return shim.Error("10th argument must be a non-empty string")
	}
	if len(args[10]) <= 0 {
		return shim.Error("11th argument must be a non-empty string")
	}
	if len(args[11]) <= 0 {
		return shim.Error("12th argument must be a non-empty string")
	}
	if len(args[12]) <= 0 {
		return shim.Error("13th argument must be a non-empty string")
	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	if userrole == "Doctor" {

		patientId := strings.ToLower(args[0])
		patientSSN := strings.ToLower(args[1])
		patientUrl := strings.ToLower(args[2])
		firstname := strings.ToLower(args[3])
		lastname := strings.ToLower(args[4])
		DOB := strings.ToLower(args[5])
		ehr := strings.ToLower(args[6])
		ehrcode := strings.ToLower(args[7])
		ehrurl := strings.ToLower(args[8])
		ehrlocalkey := strings.ToLower(args[9])
		patientglobalkey := strings.ToLower(args[10])
		stime := args[11]
		etime := args[12]

		patientData, err := stub.GetState(patientId)
		if err != nil {
			return shim.Error("Patient exist no need to add again " + err.Error())
		}
		//If patient doesn't exist, make a new object and save it to the state with the EHR.
		if patientData == nil {
			fmt.Println("Patient doesn't exist so a new patient")

			var patient entity.Patient

			//Get patient object and assign values to it for storing it in state.
			patient.ObjectType = "Patient"
			patient.PatientId = patientId
			patient.PatientSSN = patientSSN
			patient.PatientUrl = patientUrl
			patient.PatientFirstname = firstname
			patient.PatientLastname = lastname
			patient.DOB = DOB
			patient.EHR = ehr
			patient.EHRCode = ehrcode
			patient.EHRUrl = ehrurl
			patient.EHRlocalKey = ehrlocalkey
			patient.PatientGlobalKey = patientglobalkey

			providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
			if err != nil {
				return shim.Error("Fails to get id " + err.Error())
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
			/////////////////////////////////
			//Get the Doctor's EHR from his object
			var ehrlist entity.EHRList
			ehrlist.Ehr = provider.ProviderEHR

			patient.EHRlist = []entity.EHRList{}
			patient.EHRlist = append(patient.EHRlist, ehrlist) //append the EHR of Doctor who is registering to the patient.

			//Add Provider consent to patient's object
			var defaultConsent entity.Consent
			defaultConsent.ObjectType = "Consent"
			defaultConsent.Provider = provider
			defaultConsent.StartTime = stime
			defaultConsent.EndTime = etime

			patient.ProviderConsent = []entity.Consent{}
			patient.ProviderConsent = append(patient.ProviderConsent, defaultConsent)

			patientJSONasBytes, err := json.Marshal(&patient)
			if err != nil {
				return shim.Error("Fails to Marshal " + err.Error())
			}

			//=== Save Patient to state ===
			err = stub.PutState(patientId, patientJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			//==== Create patientMedication object and marshal to JSON ====
			var patientdetails entity.PatientDetails
			patientdetails.Medications.ObjectType = "Medications"
			patientdetails.Medications.Patient = patient
			patientdetails.Medications.ProviderConsent = []entity.Consent{}
			patientdetails.Medications.ProviderConsent = append(patientdetails.Medications.ProviderConsent, defaultConsent)

			//==== Create patientAllergies object and marshal to JSON ====
			patientdetails.Allergies.ObjectType = "Allergies"
			patientdetails.Allergies.Patient = patient
			patientdetails.Allergies.ProviderConsent = []entity.Consent{}
			patientdetails.Allergies.ProviderConsent = append(patientdetails.Allergies.ProviderConsent, defaultConsent)

			//==== Create patientImmunizations object and marshal to JSON ====
			patientdetails.Immunization.ObjectType = "Immunizations"
			patientdetails.Immunization.Patient = patient
			patientdetails.Immunization.ProviderConsent = []entity.Consent{}
			patientdetails.Immunization.ProviderConsent = append(patientdetails.Immunization.ProviderConsent, defaultConsent)

			//==== Create patientPastMedicalHx object and marshal to JSON ====
			patientdetails.PastMedicalHx.ObjectType = "PastMedicalHx"
			patientdetails.PastMedicalHx.Patient = patient
			patientdetails.PastMedicalHx.ProviderConsent = []entity.Consent{}
			patientdetails.PastMedicalHx.ProviderConsent = append(patientdetails.PastMedicalHx.ProviderConsent, defaultConsent)

			//==== Create patientFamilyHx object and marshal to JSON ====
			patientdetails.FamilyHx.ObjectType = "FamilyHx"
			patientdetails.FamilyHx.Patient = patient
			patientdetails.FamilyHx.ProviderConsent = []entity.Consent{}
			patientdetails.FamilyHx.ProviderConsent = append(patientdetails.FamilyHx.ProviderConsent, defaultConsent)

			PatientDetailsJSONasBytes, err := json.Marshal(&patientdetails)

			if err != nil {
				return shim.Error(err.Error())
			}

			// === Save patientDetails to state ===
			err = stub.PutPrivateData("patientDetailsIn2Orgs", patientId, PatientDetailsJSONasBytes)
			//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			//  ==== Index the Patient to enable name-based range queries, e.g. return all Patients ====
			//  An 'index' is a normal key/value entry in state.
			//  The key is a composite key, with the elements that you want to range query on listed first.
			//  In our case, the composite key is based on indexName~color~name.
			//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
			//Create firstname, ID composite key
			indexName := "firstname~patientId"
			fnameIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{patient.PatientFirstname, patient.PatientId})
			if err != nil {
				return shim.Error(err.Error())
			}
			//Create lastname, ID composite key
			indexName = "lastname~patientId"
			lnameIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{patient.PatientLastname, patient.PatientId})
			if err != nil {
				return shim.Error(err.Error())
			}
			//Create DOB, ID composite key
			indexName = "dob~patientId"
			dobIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{patient.DOB, patient.PatientId})
			if err != nil {
				return shim.Error(err.Error())
			}
			//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
			//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
			value := []byte{0x00}
			stub.PutState(fnameIDIndexKey, value)
			stub.PutState(lnameIDIndexKey, value)
			stub.PutState(dobIDIndexKey, value)

			//If patient already exists, simply get the patient object and then append the new EHR into the list of EHR.
		} else if patientData != nil {

			var patient entity.Patient

			fmt.Println("This patient already exists or the ID is already assigned to some other user")
			err = json.Unmarshal(patientData, &patient)

			providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
			if err != nil {
				return shim.Error("Fails to get id " + err.Error())
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
			/*Patient already exists, now we check that if the doctor that is registering the patient and the patient's EHR are same then
			we return and error saying that the patient already exists, but if the patient's and doctor's EHR is not the same then we append
			the doctor's EHR into the patient EHR list*/

			for i := 0; i < len(patient.EHRlist); i++ {
				tempEHR := patient.EHRlist[i].Ehr
				if tempEHR == provider.ProviderEHR {
					return shim.Error("Patient Already exist in this EHR")
				}
			}
			//Will come out of loop without an error only if the patient doesn't exist in the EHR of the provider

			//Add Provider consent to patient's object
			var defaultConsent entity.Consent
			defaultConsent.ObjectType = "Consent"
			defaultConsent.Provider = provider
			defaultConsent.StartTime = stime
			defaultConsent.EndTime = etime

			patient.ProviderConsent = append(patient.ProviderConsent, defaultConsent) //General consent

			//Get patientdetails from private data
			patientDetailsBytes, err := stub.GetPrivateData("patientDetailsIn2Orgs", patient.PatientId)
			if err != nil {
				return shim.Error("Patient not found " + patient.PatientId + "role " + patient.PatientId + "patient details " + string(patientDetailsBytes))
			}

			var patientdetailsDB entity.PatientDetails
			//unmarshal them into patientDetailsDB
			err = json.Unmarshal(patientDetailsBytes, &patientdetailsDB) //unmarshal it aka JSON.parse()
			if err != nil {
				return shim.Error(err.Error())
			}

			//Append the Provider into consents within medication, allergies, immunization, pastmedical and family
			patientdetailsDB.Medications.ProviderConsent = append(patientdetailsDB.Medications.ProviderConsent, defaultConsent)
			patientdetailsDB.Allergies.ProviderConsent = append(patientdetailsDB.Allergies.ProviderConsent, defaultConsent)
			patientdetailsDB.Immunization.ProviderConsent = append(patientdetailsDB.Immunization.ProviderConsent, defaultConsent)
			patientdetailsDB.PastMedicalHx.ProviderConsent = append(patientdetailsDB.PastMedicalHx.ProviderConsent, defaultConsent)
			patientdetailsDB.FamilyHx.ProviderConsent = append(patientdetailsDB.FamilyHx.ProviderConsent, defaultConsent)

			PatientDetailsJSONasBytes, err := json.Marshal(patientdetailsDB)

			if err != nil {
				return shim.Error(err.Error())
			}

			// === Save patientDetails to state ===
			err = stub.PutPrivateData("patientDetailsIn2Orgs", patient.PatientId, PatientDetailsJSONasBytes)
			//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

			//Get the Doctor's EHR from his object
			var ehrlist entity.EHRList
			ehrlist.Ehr = provider.ProviderEHR

			patient.EHRlist = append(patient.EHRlist, ehrlist) //append the EHR of Doctor who is registering to the patient.

			patientJSONasBytes, err := json.Marshal(&patient)
			if err != nil {
				return shim.Error("Fails to Marshal " + err.Error())
			}

			//=== Save Patient to state ===
			err = stub.PutState(patientId, patientJSONasBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

		}

		// ==== Marble saved and indexed. Return success ====
		//fmt.Println("- end register patient")
		return shim.Success(nil)
	}
	return shim.Error("Unauthorized! Only Doctor can add patient")

}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func (u *User) SearchProviderRequestsByPatientInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Firstname := strings.ToLower(args[0])
	Lastname := strings.ToLower(args[1])
	DOB := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\"}}", Firstname, Lastname, DOB)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Get provider id based on firstname, lastname and dob
	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		patientData, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Fails to get patient " + err.Error())
		}

		if patientData == nil {

			return shim.Error("This patient doesn't exist")

		} else if patientData != nil {

			var patient entity.Patient

			err = json.Unmarshal(patientData, &patient)
			if err != nil {
				return shim.Error("Fails unmarhsal patient " + err.Error())
			}

			//Get ProviderRequests from list in patient and save them to temp Request
			var tempRequest []entity.Consent
			for i := 0; i < len(patient.ProviderRequest); i++ {
				tempRequest = append(tempRequest, patient.ProviderRequest[i])
			}

			//Marshal them into bytes in order to return them.
			RequestJSONasBytes, err := json.Marshal(tempRequest)
			if err != nil {
				shim.Error("Failed to marshal requests" + err.Error())
			}

			return shim.Success(RequestJSONasBytes)

		}
	}
	return shim.Error("Unauthorized! Only Patient can view Provider Requests")
}

func (u *User) SearchProviderRequestsByPatientSSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ssn := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

	//Returns a Key and Value
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Save those attributed in PatientUnmarshal
	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Get patient id in key based on firstname, lastname and dob
	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		patientData, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Fails to get patient " + err.Error())
		}

		if patientData == nil {

			return shim.Error("This patient doesn't exist")

		} else if patientData != nil {

			var patient entity.Patient

			err = json.Unmarshal(patientData, &patient)
			if err != nil {
				return shim.Error("Fails unmarhsal patient " + err.Error())
			}

			//Get ProviderRequests from list in patient and save them to temp Request
			var tempRequest []entity.Consent
			for i := 0; i < len(patient.ProviderRequest); i++ {
				tempRequest = append(tempRequest, patient.ProviderRequest[i])
			}

			//Marshal them into bytes in order to return them.
			RequestJSONasBytes, err := json.Marshal(tempRequest)
			if err != nil {
				shim.Error("Failed to marshal requests" + err.Error())
			}

			return shim.Success(RequestJSONasBytes)

		}
	}
	return shim.Error("Unauthorized! Only Patient can view Provider Requests")
}

func (u *User) AllowConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Get Current Patient based on token
	patientId, err := getAttribute(stub, "id")
	if err != nil {
		return shim.Error("Fails to get id " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		providerFirstname := strings.ToLower(args[0])
		providerLastname := strings.ToLower(args[1])
		providerSpeciality := strings.ToLower(args[2])

		queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"speciality\":\"%s\"}}", providerFirstname, providerLastname, providerSpeciality)

		// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

		queryResults, err := getQueryResultForQueryString(stub, queryString)
		if err != nil {
			return shim.Error(err.Error())
		}

		var tempArray []entity.ProviderUnmarshal
		err = json.Unmarshal(queryResults, &tempArray)
		if err != nil {
			return shim.Error(err.Error())
		}

		//Get provider id based on firstname, lastname and dob
		var key string
		for _, temp := range tempArray {

			key = temp.Key

		}

		providerData, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Fails to get provider " + err.Error())
		}

		if providerData == nil {

			return shim.Error("This provider doesn't exist")

		} else if providerData != nil {

			var provider entity.Provider

			err = json.Unmarshal(providerData, &provider)
			if err != nil {
				return shim.Error("Fails unmarhsal provider " + err.Error())
			}

			patientAsByte, err := stub.GetState(patientId) //Get the Patient's details from the state.
			if err != nil {
				return shim.Error("Fails to get patient: " + err.Error())
			}

			var patient entity.Patient
			err = json.Unmarshal(patientAsByte, &patient) //Store the Patient's object into patient.

			if err != nil {
				return shim.Error("Fails to unmarshal patient " + err.Error())
			}

			for i := 0; i < len(patient.ProviderConsent); i++ {
				tempID := patient.ProviderConsent[i].Provider.ProviderId
				if tempID == provider.ProviderId {
					return shim.Error("Provider already has consent.")
				}
			}
			//Will come out of loop without an error only if the Provider doesn't have consent.

			a := 0

			//Now check if the provider has a request or not and if he does, then get that request and append it into the consent.
			for i := 0; i < len(patient.ProviderRequest); i++ {
				tempID := patient.ProviderRequest[i].Provider.ProviderId
				if tempID == provider.ProviderId {

					//Add Provider consent to patient's object
					var defaultConsent entity.Consent
					defaultConsent.ObjectType = "Consent"
					defaultConsent.Provider = patient.ProviderRequest[i].Provider
					defaultConsent.StartTime = patient.ProviderRequest[i].StartTime
					defaultConsent.EndTime = patient.ProviderRequest[i].EndTime

					patient.ProviderConsent = append(patient.ProviderConsent, defaultConsent)

					//remove the request from list that is added to consent
					patient.ProviderRequest = append(patient.ProviderRequest[:i], patient.ProviderRequest[i+1:]...)
					//Marshal the object into json format
					patientJSONasBytes, err := json.Marshal(patient)

					//Save the object into state
					err = stub.PutState(patientId, patientJSONasBytes)
					if err != nil {
						return shim.Error(err.Error())
					}

					//Get patientdetails from private data
					patientDetailsBytes, err := stub.GetPrivateData("patientDetailsIn2Orgs", patient.PatientId)
					if err != nil {
						return shim.Error("Patient not found " + patient.PatientId + "role " + patient.PatientId + "patient details " + string(patientDetailsBytes))
					}

					var patientdetailsDB entity.PatientDetails
					//unmarshal them into patientDetailsDB
					err = json.Unmarshal(patientDetailsBytes, &patientdetailsDB) //unmarshal it aka JSON.parse()
					if err != nil {
						return shim.Error(err.Error())
					}

					//Append the Provider into consents within medication, allergies, immunization, pastmedical and family
					patientdetailsDB.Medications.ProviderConsent = append(patientdetailsDB.Medications.ProviderConsent, defaultConsent)
					patientdetailsDB.Allergies.ProviderConsent = append(patientdetailsDB.Allergies.ProviderConsent, defaultConsent)
					patientdetailsDB.Immunization.ProviderConsent = append(patientdetailsDB.Immunization.ProviderConsent, defaultConsent)
					patientdetailsDB.PastMedicalHx.ProviderConsent = append(patientdetailsDB.PastMedicalHx.ProviderConsent, defaultConsent)
					patientdetailsDB.FamilyHx.ProviderConsent = append(patientdetailsDB.FamilyHx.ProviderConsent, defaultConsent)

					PatientDetailsJSONasBytes, err := json.Marshal(patientdetailsDB)

					if err != nil {
						return shim.Error(err.Error())
					}

					// === Save patientDetails to state ===
					err = stub.PutPrivateData("patientDetailsIn2Orgs", patient.PatientId, PatientDetailsJSONasBytes)
					//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
					if err != nil {
						return shim.Error(err.Error())
					}

					a = 1 //Flag for checking id the provider request exists or not
				}
			}

			if a != 1 {
				return shim.Error("This Provider doesn't have any request")
			}

			return shim.Success(nil)

		}
	}
	return shim.Error("Unauthorized! Only Patient can Allow Consent")
}

//Using parametrized Rich Query
func (u *User) GetPatientById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

func (u *User) UpdatePatientById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	if len(args) < 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	//set the variables
	patientId := strings.ToLower(args[0])
	patientSSN := strings.ToLower(args[1])
	patientUrl := strings.ToLower(args[2])
	firstname := strings.ToLower(args[3])
	lastname := strings.ToLower(args[4])
	DOB := strings.ToLower(args[5])

	//check in the state that if the patient exists or not
	patientAsByte, err := stub.GetState(patientId)
	if err != nil {
		fmt.Println("Patient doesn't exist so a new patient")
		return shim.Error("This patient doesn't exist so can't update")
	} else if patientAsByte != nil {
		fmt.Println("This patient already exists or the ID is already assigned to some other user")
	}

	//Unmarshal the patient details acquired and store them in patient object
	var patient entity.Patient
	err = json.Unmarshal(patientAsByte, &patient)

	//Set updated values into the patient object retrieved from state
	patient.ObjectType = "Patient"
	patient.PatientId = patientId
	patient.PatientSSN = patientSSN
	patient.PatientUrl = patientUrl
	patient.PatientFirstname = firstname
	patient.PatientLastname = lastname
	patient.DOB = DOB

	//Marshal the object into json fromet
	patientJSONasBytes, err := json.Marshal(patient)

	//Save the object into state
	err = stub.PutState(patientId, patientJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//Using Composite Key
func (u *User) GetPatientByFirstName(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	firstname := args[0]
	indexName := "firstname~patientId"

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

		returnedPatientId := compositeKeyParts[1]

		PatientAsBytes, _ := stub.GetState(returnedPatientId)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedPatientId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(PatientAsBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryGetPatientByFirstName:\n%s\n", buffer.String())
	queryResult := buffer.Bytes()
	return shim.Success(queryResult)
}

//Using composite key
func (u *User) GetPatientByLastName(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	lastname := args[0]
	indexName := "lastname~patientId"

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

		returnedPatientId := compositeKeyParts[1]

		PatientAsBytes, _ := stub.GetState(returnedPatientId)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedPatientId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(PatientAsBytes))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryGetPatientByLastName:\n%s\n", buffer.String())
	queryResult := buffer.Bytes()
	return shim.Success(queryResult)
}

//By parametrized Rich Query
func (u *User) GetPatientByInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fname := strings.ToLower(args[0])
	lname := strings.ToLower(args[1])
	dob := strings.ToLower(args[2])
	curr := args[3] //current time as argument

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

		//Get Provider
		providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
		if err != nil {
			return shim.Error("Fails to get id " + err.Error())
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

		var patientinfo entity.ShowPatient
		//Search in patient's providerconsent list that if the current provider is allowed to access the patient's data.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {

				layout := "2006-01-02T15:04:05" //Set layout to parse the string time saved in state.

				current, err := time.Parse(layout, curr) //Get current time in a certain layout
				if err != nil {
					return shim.Error("Error in parsing current-time!")
				}

				starttime, err := time.Parse(layout, patient.ProviderConsent[i].StartTime) //Parse Start time
				if err != nil {
					return shim.Error("Error in parsing start-time!")
				}

				endtime, err := time.Parse(layout, patient.ProviderConsent[i].EndTime) //Parse End time
				if err != nil {
					return shim.Error("Error in parsing end-time!")
				}

				if inTimeSpan(starttime, endtime, current) { //If current time is in between start and end time then
					patientinfo.ObjectType = "PatientInformation"
					patientinfo.PatientId = patient.PatientId
					patientinfo.PatientSSN = patient.PatientSSN
					patientinfo.PatientUrl = patient.PatientUrl
					patientinfo.PatientFirstname = patient.PatientFirstname
					patientinfo.PatientLastname = patient.PatientLastname
					patientinfo.DOB = patient.DOB

					PatientInfoJSONasBytes, err := json.Marshal(patientinfo)
					if err != nil {
						shim.Error("Failed to marshal patient information" + err.Error())
					}

					return shim.Success(PatientInfoJSONasBytes)

				} else if !inTimeSpan(starttime, endtime, current) { //If not then
					return shim.Error("Provider consent time is over!")
				}
			}
		}

		return shim.Error("Provider doesn't have consent to view the patient")
	}
	return shim.Error("Unauthorized! Only Doctor can search patient")
}

func (u *User) GetPatientBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ssn := strings.ToLower(args[0])
	curr := args[1] //current time as argument

	queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

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

		//Get Provider
		providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
		if err != nil {
			return shim.Error("Fails to get id " + err.Error())
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

		var patientinfo entity.ShowPatient
		//Search in patient's providerconsent list that if the current provider is allowed to access the patient's data.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {

				layout := "2006-01-02T15:04:05" //Set layout to parse the string time saved in state.

				current, err := time.Parse(layout, curr) //Get current time in a certain layout
				if err != nil {
					return shim.Error("Error in parsing current-time!")
				}

				starttime, err := time.Parse(layout, patient.ProviderConsent[i].StartTime) //Parse Start time
				if err != nil {
					return shim.Error("Error in parsing start-time!")
				}

				endtime, err := time.Parse(layout, patient.ProviderConsent[i].EndTime) //Parse End time
				if err != nil {
					return shim.Error("Error in parsing end-time!")
				}

				if inTimeSpan(starttime, endtime, current) { //If current time is in between start and end time then
					patientinfo.ObjectType = "PatientInformation"
					patientinfo.PatientId = patient.PatientId
					patientinfo.PatientSSN = patient.PatientSSN
					patientinfo.PatientUrl = patient.PatientUrl
					patientinfo.PatientFirstname = patient.PatientFirstname
					patientinfo.PatientLastname = patient.PatientLastname
					patientinfo.DOB = patient.DOB

					PatientInfoJSONasBytes, err := json.Marshal(patientinfo)
					if err != nil {
						shim.Error("Failed to marshal patient information" + err.Error())
					}

					return shim.Success(PatientInfoJSONasBytes)

				} else if !inTimeSpan(starttime, endtime, current) { //If not then
					return shim.Error("Provider consent time is over!")
				}
			}
		}

		return shim.Error("Provider doesn't have consent to view the patient")
	}
	return shim.Error("Unauthorized! Only Doctor can search patient")
}

func (u *User) GetPatientAllInfoBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ssn := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

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

		//Get Provider
		providerId, err := getAttribute(stub, "id") //Get the current doctor ID from identity.
		if err != nil {
			return shim.Error("Fails to get id " + err.Error())
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

		//Search in patient's providerconsent list that if the current provider is allowed to access the patient's data.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempEHR := patient.ProviderConsent[i].Provider.ProviderId
			if tempEHR == provider.ProviderId {
				return shim.Success(queryResults)
			}
		}

		return shim.Error("Provider doesn't have consent to view the patient")
	}
	return shim.Error("Unauthorized! Only Doctor can search patient")
}

func (u *User) UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// for searching provider object
	Firstname := strings.ToLower(args[0])
	Lastname := strings.ToLower(args[1])
	Speciality := strings.ToLower(args[2])
	stime := args[3]
	etime := args[4]

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"speciality\":\"%s\"}}", Firstname, Lastname, Speciality)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.ProviderUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Get Current Patient based on token
	patientId, err := getAttribute(stub, "id") //Get the current patient ID from identity.
	if err != nil {
		return shim.Error("Fails to get id " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		patientAsByte, err := stub.GetState(patientId) //Get the patient's details from the state.
		if err != nil {
			return shim.Error("Fails to get patient: " + err.Error())
		}

		var patient entity.Patient
		err = json.Unmarshal(patientAsByte, &patient) //Store the Patient's object into patient.

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		//Get Provider
		var provider entity.Provider
		providerAsByte, err := stub.GetState(key)
		err = json.Unmarshal(providerAsByte, &provider)

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		a := 0
		//Search in patient's providerconsent list that if the current provider is allowed to access the patient's data.
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {
				patient.ProviderConsent[i].StartTime = stime
				patient.ProviderConsent[i].EndTime = etime

				//Marshal the object into json format
				patientJSONasBytes, err := json.Marshal(patient)
				if err != nil {
					return shim.Error(err.Error())
				}

				//Save the object into state replacing the old one
				err = stub.PutState(patientId, patientJSONasBytes)
				if err != nil {
					return shim.Error(err.Error())
				}

				//Get patientdetails from private data
				patientDetailsBytes, err := stub.GetPrivateData("patientDetailsIn2Orgs", patient.PatientId)
				if err != nil {
					return shim.Error("Patient not found " + patient.PatientId + "role " + patient.PatientId + "patient details " + string(patientDetailsBytes))
				}

				var patientdetailsDB entity.PatientDetails
				//unmarshal them into patientDetailsDB
				err = json.Unmarshal(patientDetailsBytes, &patientdetailsDB) //unmarshal it aka JSON.parse()
				if err != nil {
					return shim.Error(err.Error())
				}

				//Update the Provider Starttime and Endtime into consents within medication, allergies, immunization, pastmedical and family
				patientdetailsDB.Medications.ProviderConsent[i].StartTime = stime
				patientdetailsDB.Medications.ProviderConsent[i].EndTime = etime

				patientdetailsDB.Allergies.ProviderConsent[i].StartTime = stime
				patientdetailsDB.Allergies.ProviderConsent[i].EndTime = etime

				patientdetailsDB.Immunization.ProviderConsent[i].StartTime = stime
				patientdetailsDB.Immunization.ProviderConsent[i].EndTime = etime

				patientdetailsDB.PastMedicalHx.ProviderConsent[i].StartTime = stime
				patientdetailsDB.PastMedicalHx.ProviderConsent[i].EndTime = etime

				patientdetailsDB.FamilyHx.ProviderConsent[i].StartTime = stime
				patientdetailsDB.FamilyHx.ProviderConsent[i].EndTime = etime

				PatientDetailsJSONasBytes, err := json.Marshal(patientdetailsDB)

				if err != nil {
					return shim.Error(err.Error())
				}

				// === Save patientDetails to state ===
				err = stub.PutPrivateData("patientDetailsIn2Orgs", patient.PatientId, PatientDetailsJSONasBytes)
				//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
				if err != nil {
					return shim.Error(err.Error())
				}

				a = 1
			}
		}

		if a != 1 {
			return shim.Error("Provider's Consent is not registered")
		}

		return shim.Success(nil)

	}
	return shim.Error("Unauthorized! Only Patient can Update Provider Access!")

}

func (u *User) RevokeProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// for searching provider object
	Firstname := strings.ToLower(args[0])
	Lastname := strings.ToLower(args[1])
	Speciality := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"speciality\":\"%s\"}}", Firstname, Lastname, Speciality)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.ProviderUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Get Current Patient based on token
	patientId, err := getAttribute(stub, "id") //Get the current patient ID from identity.
	if err != nil {
		return shim.Error("Fails to get id " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		patientAsByte, err := stub.GetState(patientId) //Get the patient's details from the state.
		if err != nil {
			return shim.Error("Fails to get patient: " + err.Error())
		}

		var patient entity.Patient
		err = json.Unmarshal(patientAsByte, &patient) //Store the Patient's object into patient.

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		//Get Provider
		var provider entity.Provider
		providerAsByte, err := stub.GetState(key)
		err = json.Unmarshal(providerAsByte, &provider)

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		a := 0 //Flag
		//Search in patient's providerconsent list that if the current provider is in the list or not and if he is then remove him
		for i := 0; i < len(patient.ProviderConsent); i++ {
			tempConsent := patient.ProviderConsent[i].Provider.ProviderId
			if tempConsent == provider.ProviderId {
				patient.ProviderConsent = append(patient.ProviderConsent[:i], patient.ProviderConsent[i+1:]...)

				//Marshal the object into json format
				patientJSONasBytes, err := json.Marshal(patient)
				if err != nil {
					return shim.Error(err.Error())
				}

				//Save the object into state replacing the old one
				err = stub.PutState(patientId, patientJSONasBytes)
				if err != nil {
					return shim.Error(err.Error())
				}

				//Get patientdetails from private data
				patientDetailsBytes, err := stub.GetPrivateData("patientDetailsIn2Orgs", patient.PatientId)
				if err != nil {
					return shim.Error("Patient not found " + patient.PatientId + "role " + patient.PatientId + "patient details " + string(patientDetailsBytes))
				}

				var patientdetailsDB entity.PatientDetails
				//unmarshal them into patientDetailsDB
				err = json.Unmarshal(patientDetailsBytes, &patientdetailsDB) //unmarshal it aka JSON.parse()
				if err != nil {
					return shim.Error(err.Error())
				}

				//Remove Consent of the current provider from the consent list within medication, allergies, immunization, pastmedical and family
				patientdetailsDB.Medications.ProviderConsent = append(patientdetailsDB.Medications.ProviderConsent[:i], patientdetailsDB.Medications.ProviderConsent[i+1:]...)

				patientdetailsDB.Allergies.ProviderConsent = append(patientdetailsDB.Allergies.ProviderConsent[:i], patientdetailsDB.Allergies.ProviderConsent[i+1:]...)

				patientdetailsDB.Immunization.ProviderConsent = append(patientdetailsDB.Immunization.ProviderConsent[:i], patientdetailsDB.Immunization.ProviderConsent[i+1:]...)

				patientdetailsDB.PastMedicalHx.ProviderConsent = append(patientdetailsDB.PastMedicalHx.ProviderConsent[:i], patientdetailsDB.PastMedicalHx.ProviderConsent[i+1:]...)

				patientdetailsDB.FamilyHx.ProviderConsent = append(patientdetailsDB.FamilyHx.ProviderConsent[:i], patientdetailsDB.FamilyHx.ProviderConsent[i+1:]...)

				PatientDetailsJSONasBytes, err := json.Marshal(patientdetailsDB)

				if err != nil {
					return shim.Error(err.Error())
				}

				// === Save patientDetails to state ===
				err = stub.PutPrivateData("patientDetailsIn2Orgs", patient.PatientId, PatientDetailsJSONasBytes)
				//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
				if err != nil {
					return shim.Error(err.Error())
				}

				a = 1
			}
		}

		if a != 1 {
			return shim.Error("Provider's Consent is not registered")
		}

		return shim.Success(nil)

	}
	return shim.Error("Unauthorized! Only Patient can Revoke Provider Access!")

}

func (u *User) GetProviderAccessStatusForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// for searching provider object
	Firstname := strings.ToLower(args[0])
	Lastname := strings.ToLower(args[1])
	Speciality := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"speciality\":\"%s\"}}", Firstname, Lastname, Speciality)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.ProviderUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	//Get Current Patient based on token
	patientId, err := getAttribute(stub, "id") //Get the current patient ID from identity.
	if err != nil {
		return shim.Error("Fails to get id " + err.Error())
	}

	//Only allow when the user is a patient.
	if userrole == "Patient" {

		patientAsByte, err := stub.GetState(patientId) //Get the patient's details from the state.
		if err != nil {
			return shim.Error("Fails to get patient: " + err.Error())
		}

		var patient entity.Patient
		err = json.Unmarshal(patientAsByte, &patient) //Store the Patient's object into patient.

		if err != nil {
			return shim.Error("Fails to unmarshal patient " + err.Error())
		}

		//Get Provider
		var provider entity.Provider
		providerAsByte, err := stub.GetState(key)
		err = json.Unmarshal(providerAsByte, &provider)

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
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
	return shim.Error("Unauthorized! Only Patient can get Provider Access Status!")
}

func (u *User) GetAllProviderAccessByPatientInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Firstname := strings.ToLower(args[0])
	Lastname := strings.ToLower(args[1])
	DOB := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\"}}", Firstname, Lastname, DOB)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Get provider id based on firstname, lastname and dob
	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	if userrole == "Patient" {

		patientData, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Fails to get patient " + err.Error())
		}

		if patientData == nil {

			return shim.Error("This patient doesn't exist")

		} else if patientData != nil {

			var patient entity.Patient

			err = json.Unmarshal(patientData, &patient)
			if err != nil {
				return shim.Error("Fails unmarhsal patient " + err.Error())
			}

			//Get ProviderRequests from list in patient and save them to temp Request
			var tempConsent []entity.Consent
			for i := 0; i < len(patient.ProviderConsent); i++ {
				tempConsent = append(tempConsent, patient.ProviderConsent[i])
			}

			//Marshal them into bytes in order to return them.
			ConsentJSONasBytes, err := json.Marshal(tempConsent)
			if err != nil {
				shim.Error("Failed to marshal consents" + err.Error())
			}

			return shim.Success(ConsentJSONasBytes)

		}
	}
	return shim.Error("Unauthorized! Only Patient can get All Provider Access Status!")
}

func (u *User) GetAllProviderAccessByPatientSSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ssn := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Get provider id based on firstname, lastname and dob
	var key string
	for _, temp := range tempArray {

		key = temp.Key

	}

	userrole, err := getAttribute(stub, "userrole")
	if err != nil {
		return shim.Error("Fails to get userrole " + err.Error())
	}

	if userrole == "Patient" {

		patientData, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Fails to get patient " + err.Error())
		}

		if patientData == nil {

			return shim.Error("This patient doesn't exist")

		} else if patientData != nil {

			var patient entity.Patient

			err = json.Unmarshal(patientData, &patient)
			if err != nil {
				return shim.Error("Fails unmarhsal patient " + err.Error())
			}

			//Get ProviderRequests from list in patient and save them to temp Request
			var tempConsent []entity.Consent
			for i := 0; i < len(patient.ProviderConsent); i++ {
				tempConsent = append(tempConsent, patient.ProviderConsent[i])
			}

			//Marshal them into bytes in order to return them.
			ConsentJSONasBytes, err := json.Marshal(tempConsent)
			if err != nil {
				shim.Error("Failed to marshal consents" + err.Error())
			}

			return shim.Success(ConsentJSONasBytes)

		}
	}
	return shim.Error("Unauthorized! Only Patient can get All Provider Access Status!")
}
