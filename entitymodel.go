	package entity

	// SimpleChaincode example simple Chaincode implementation
	type SimpleChaincode struct {
	}


	type marble struct {
		ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
		Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
		Color      string `json:"color"`
		Size       int    `json:"size"`
		Owner      string `json:"owner"`
	}


