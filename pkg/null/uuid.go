package null

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
    "reflect"

    uuid "github.com/satori/go.uuid"
)

var (
    NamespaceTS = uuid.Must(uuid.FromString("19860912-9dad-11d1-80b4-00c04fd430c8"))
    NamespaceZH = uuid.Must(uuid.FromString("00000001-1986-0912-0511-013419290962"))
)

// Time is a nullable uuid.UUID. It supports SQL and JSON serialization.
// It will marshal to null if null.
type UUID struct {
    uuid.NullUUID
}

// Scan implements the Scanner interface.
func (u *UUID) Scan(value interface{}) error {
    var err error
    switch x := value.(type) {
    case nil, []byte, string:
        err = u.NullUUID.Scan(value)
    case uuid.NullUUID:
        u.NullUUID = x
        return nil
    case uuid.UUID:
        u.UUID = x
    default:
        err = fmt.Errorf("null: cannot scan type %T into null.UUID: %v", value, value)
    }
    u.Valid = err == nil
    return err
}

// Value implements the driver Valuer interface.
func (u UUID) Value() (driver.Value, error) {
    return u.NullUUID.Value()
}

// NewTime creates a new UUID.
func NewUUID(u uuid.UUID, valid bool) UUID {
    return UUID{
        NullUUID: uuid.NullUUID{UUID: u, Valid: valid,},
    }
}

// TimeFrom creates a new Time that will always be valid.
func UUIDFrom(t uuid.UUID) UUID {
    return NewUUID(t, true)
}

// TimeFromPtr creates a new Time that will be null if t is nil.
func UUIDFromPtr(t *uuid.UUID) UUID {
    if t == nil {
        return NewUUID(uuid.UUID{}, false)
    }
    return NewUUID(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (u UUID) ValueOrZero() uuid.UUID {
    if !u.Valid {
        return uuid.UUID{}
    }
    return u.UUID
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (u UUID) MarshalJSON() (data []byte, err error) {
    if !u.Valid {
        return []byte("null"), nil
    }
    data = []byte( fmt.Sprintf(`"%s"`, u.UUID.String()))
    return
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (u *UUID) UnmarshalJSON(data []byte) error {
    var err error
    var v interface{}
    if err = json.Unmarshal(data, &v); err != nil {
        return err
    }
    switch x := v.(type) {
    case []byte:
        err = u.UUID.UnmarshalText(x)
    case string:
        err = u.UUID.UnmarshalText([]byte(x))
    case map[string]interface{}:
        ui, uiOK := x["UUID"].(string)
        valid, validOK := x["Valid"].(bool)
        if !uiOK || !validOK {
            return fmt.Errorf(`json: unmarshalling object into Go value of type null.UUID requires key "UUID" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["UUID"], x["Valid"])
        }
        err = u.UUID.UnmarshalText([]byte(ui))
        u.Valid = valid
        return err
    case nil:
        u.Valid = false
        return nil
    default:
        err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.UUID", reflect.TypeOf(v).Name())
    }
    u.Valid = err == nil
    return err
}

func (u UUID) MarshalText() ([]byte, error) {
    if !u.Valid {
        return []byte("null"), nil
    }
    return u.UUID.MarshalText()
}

func (u *UUID) UnmarshalText(text []byte) error {
    str := string(text)
    if str == "" || str == "null" {
        u.Valid = false
        return nil
    }
    if err := u.UUID.UnmarshalText(text); err != nil {
        return err
    }
    u.Valid = true
    return nil
}

// SetValid changes this Time's value and sets it to be non-null.
func (u *UUID) SetValid(v uuid.UUID) {
    u.UUID = v
    u.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (u UUID) Ptr() *uuid.UUID {
    if !u.Valid {
        return nil
    }
    return &u.UUID
}

// IsZero returns true for invalid Times, hopefully for future omitempty support.
// A non-null Time with a zero value will not be considered zero.
func (u UUID) IsZero() bool {
    return !u.Valid
}

func (u *UUID) NewV1() {
    u.UUID = uuid.NewV1()
    u.Valid = true
}
func (u *UUID) NewV2(domain byte) {
    u.UUID = uuid.NewV2(domain)
    u.Valid = true
}
func (u *UUID) NewV3(ns uuid.UUID, name string) {
    u.UUID = uuid.NewV3(ns, name)
    u.Valid = true
}
func (u *UUID) NewV4() {
    u.UUID = uuid.NewV4()
    u.Valid = true
}
func (u *UUID) NewV5(ns uuid.UUID, name string) {
    u.UUID = uuid.NewV5(ns, name)
    u.Valid = true
}
func (u UUID) NamespaceDNS() uuid.UUID {
    return uuid.NamespaceDNS
}
func (u UUID) NamespaceURL() uuid.UUID {
    return uuid.NamespaceURL
}
func (u UUID) NamespaceOID() uuid.UUID {
    return uuid.NamespaceOID
}

func (u UUID) NamespaceX500() uuid.UUID {
    return uuid.NamespaceX500
}
func (u UUID) NamespaceTS() uuid.UUID {
    return NamespaceTS
}
func (u UUID) NamespaceZH() uuid.UUID {
    return NamespaceZH
}

func NewUuidV4() uuid.UUID {
    return uuid.NewV4()
}
