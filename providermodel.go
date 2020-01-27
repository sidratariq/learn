package entity

// Information of Provider
type Provider struct {
	ObjectType        string `json:docType"`
	ProviderId        string `json:"providerId"`
	ProviderEHR       string `json:"providerehr"`
	ProviderEHRURL    string `json:"providerehrurl"`
	ProviderFirstname string `json:"firstname"` //docType is used to distinguish the various types of objects in state database
	ProviderLastname  string `json:"lastname"`  //the fieldtags are needed to keep case from bouncing around
	Speciality        string `json:"speciality"`
}

type Consent struct {
	ObjectType string   `json:docType"`
	Provider   Provider `json:"provider"`
	StartTime  string   `json:"starttime"`
	EndTime    string   `json:"endtime"`
}

type ProviderUnmarshal struct {
	Key    string
	Record Patient `json:"Provider"`
}
