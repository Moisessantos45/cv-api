package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONLinkSocial []SocialLink
type JSONStringArray []string
type JSONLinkArray []Link

func (m JSONLinkSocial) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (a JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a JSONLinkArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (m *JSONLinkSocial) Scan(value any) error {
	if value == nil {
		*m = JSONLinkSocial{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, m)
}

func (a *JSONStringArray) Scan(value any) error {
	if value == nil {
		*a = JSONStringArray{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}

func (a *JSONLinkArray) Scan(value any) error {
	if value == nil {
		*a = JSONLinkArray{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, a)
}
