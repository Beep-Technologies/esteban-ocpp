// Code generated by schema-generate. DO NOT EDIT.

package ocpp16

import (
    "bytes"
    "encoding/json"
    "fmt"
    "errors"
)

// UpdateFirmwareRequest 
type UpdateFirmwareRequest struct {
  Location string `json:"location"`
  Retries int `json:"retries,omitempty"`
  RetrieveDate string `json:"retrieveDate"`
  RetryInterval int `json:"retryInterval,omitempty"`
}

func (strct *UpdateFirmwareRequest) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // "Location" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "location" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"location\": ")
	if tmp, err := json.Marshal(strct.Location); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "retries" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"retries\": ")
	if tmp, err := json.Marshal(strct.Retries); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "RetrieveDate" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "retrieveDate" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"retrieveDate\": ")
	if tmp, err := json.Marshal(strct.RetrieveDate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "retryInterval" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"retryInterval\": ")
	if tmp, err := json.Marshal(strct.RetryInterval); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *UpdateFirmwareRequest) UnmarshalJSON(b []byte) error {
    locationReceived := false
    retrieveDateReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "location":
            if err := json.Unmarshal([]byte(v), &strct.Location); err != nil {
                return err
             }
            locationReceived = true
        case "retries":
            if err := json.Unmarshal([]byte(v), &strct.Retries); err != nil {
                return err
             }
        case "retrieveDate":
            if err := json.Unmarshal([]byte(v), &strct.RetrieveDate); err != nil {
                return err
             }
            retrieveDateReceived = true
        case "retryInterval":
            if err := json.Unmarshal([]byte(v), &strct.RetryInterval); err != nil {
                return err
             }
        default:
            return fmt.Errorf("additional property not allowed: \"" + k + "\"")
        }
    }
    // check if location (a required property) was received
    if !locationReceived {
        return errors.New("\"location\" is required but was not present")
    }
    // check if retrieveDate (a required property) was received
    if !retrieveDateReceived {
        return errors.New("\"retrieveDate\" is required but was not present")
    }
    return nil
}