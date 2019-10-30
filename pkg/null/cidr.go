package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type CIDR struct {
	CIDR  *net.IPNet
	Valid bool
}

func (n *CIDR) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case []byte:
		if _, n.CIDR, err = net.ParseCIDR(string(x)); err != nil {
			err = fmt.Errorf("null: cannot scan type %T into null.CIDR: %v %v", value, value)
		}
	case string:
		if _, n.CIDR, err = net.ParseCIDR(x); err != nil {
			err = fmt.Errorf("null: cannot scan type %T into null.CIDR: %v %v", value, value)
		}
	case *net.IPNet:
		n.CIDR = x
	case net.IPNet:
		n.CIDR = &x
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.CIDR: %v", value, value)
	}
	n.Valid = err == nil
	return err
}
func (n CIDR) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.CIDR.String(), nil
}

func (n CIDR) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.CIDR.String())
}
func (n *CIDR) UnmarshalJSON(data []byte) error {
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
		if _, n.CIDR, err = net.ParseCIDR(string(x)); err != nil {
			err = fmt.Errorf("null: cannot scan type %T into null.CIDR: %v %v", v, v)
		}
	case string:
		if _, n.CIDR, err = net.ParseCIDR(x); err != nil {
			err = fmt.Errorf("null: cannot scan type %T into null.CIDR: %v %v", v, v)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, n)
	case nil:
		n.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.CIDR", reflect.TypeOf(v).Name())
	}

	return err
}

func (n CIDR) MarshalText() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", n.CIDR.String())), nil
}
func (n *CIDR) UnmarshalText(data []byte) (err error) {
	_, n.CIDR, err = net.ParseCIDR(string(data))
	n.Valid = err == nil
	return err
}

func (n CIDR) IsZero() bool {
	return !n.Valid
}

func ParseCIDR(i string) (v *CIDR, err error) {
	v = new(CIDR)
	err = v.Scan(i)
	return
}
