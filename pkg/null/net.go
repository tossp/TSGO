package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type IP struct {
	IP    net.IP
	Valid bool
}

func (n *IP) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case []byte:
		if n.IP = net.ParseIP(string(x)); n.IP == nil {
			err = fmt.Errorf("null: cannot scan type %T into null.IP: %v %v", value, value)
		}
	case string:
		if n.IP = net.ParseIP(x); n.IP == nil {
			err = fmt.Errorf("null: cannot scan type %T into null.IP: %v %v", value, value)
		}
	case net.IP:
		n.IP = x
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.IP: %v", value, value)
	}
	n.Valid = err == nil
	return err
}
func (n IP) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.IP.String(), nil
}

func (n IP) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.IP.String())
}
func (n *IP) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	defer func() {
		n.Valid = err == nil
	}()
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case []byte:
		if n.IP = net.ParseIP(string(x)); n.IP == nil {
			err = fmt.Errorf("null: cannot scan type %T into null.IP: %v %v", v, v)
		}
	case string:
		if n.IP = net.ParseIP(x); n.IP == nil {
			err = fmt.Errorf("null: cannot scan type %T into null.IP: %v %v", v, v)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, n)
	case nil:
		n.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.IP", reflect.TypeOf(v).Name())
	}

	return err
}

func (n IP) MarshalText() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return n.IP.MarshalText()
}
func (n *IP) UnmarshalText(data []byte) error {
	err := n.IP.UnmarshalText(data)
	n.Valid = err == nil
	return err
}

func (n IP) IsZero() bool {
	return !n.Valid
}
