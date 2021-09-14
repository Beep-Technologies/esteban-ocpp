// Code generated by schema-generate. DO NOT EDIT.

package ocpp16

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// LocalAuthorizationListItems
type LocalAuthorizationListItems struct {
	IdTag     string     `json:"idTag"`
	IdTagInfo *IdTagInfo `json:"idTagInfo,omitempty"`
}

// SendLocalListRequest
type SendLocalListRequest struct {
	ListVersion            int                            `json:"listVersion"`
	LocalAuthorizationList []*LocalAuthorizationListItems `json:"localAuthorizationList,omitempty"`
	UpdateType             string                         `json:"updateType"`
}

func (strct *LocalAuthorizationListItems) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "IdTag" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "idTag" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"idTag\": ")
	if tmp, err := json.Marshal(strct.IdTag); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "idTagInfo" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"idTagInfo\": ")
	if tmp, err := json.Marshal(strct.IdTagInfo); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *LocalAuthorizationListItems) UnmarshalJSON(b []byte) error {
	idTagReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "idTag":
			if err := json.Unmarshal([]byte(v), &strct.IdTag); err != nil {
				return err
			}
			idTagReceived = true
		case "idTagInfo":
			if err := json.Unmarshal([]byte(v), &strct.IdTagInfo); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if idTag (a required property) was received
	if !idTagReceived {
		return errors.New("\"idTag\" is required but was not present")
	}
	return nil
}

func (strct *SendLocalListRequest) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "ListVersion" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "listVersion" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"listVersion\": ")
	if tmp, err := json.Marshal(strct.ListVersion); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "localAuthorizationList" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"localAuthorizationList\": ")
	if tmp, err := json.Marshal(strct.LocalAuthorizationList); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "UpdateType" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "updateType" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"updateType\": ")
	if tmp, err := json.Marshal(strct.UpdateType); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *SendLocalListRequest) UnmarshalJSON(b []byte) error {
	listVersionReceived := false
	updateTypeReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "listVersion":
			if err := json.Unmarshal([]byte(v), &strct.ListVersion); err != nil {
				return err
			}
			listVersionReceived = true
		case "localAuthorizationList":
			if err := json.Unmarshal([]byte(v), &strct.LocalAuthorizationList); err != nil {
				return err
			}
		case "updateType":
			if err := json.Unmarshal([]byte(v), &strct.UpdateType); err != nil {
				return err
			}
			updateTypeReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if listVersion (a required property) was received
	if !listVersionReceived {
		return errors.New("\"listVersion\" is required but was not present")
	}
	// check if updateType (a required property) was received
	if !updateTypeReceived {
		return errors.New("\"updateType\" is required but was not present")
	}
	return nil
}