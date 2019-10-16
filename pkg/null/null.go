// github.com/guregu/null

package null

import "database/sql/driver"

type NUll interface {
    Scan( interface{}) error
    Value() (driver.Value, error)
    MarshalJSON() ([]byte, error)
    UnmarshalJSON( []byte) error
    MarshalText() ([]byte, error)
    UnmarshalText(text []byte) error
    IsZero() bool
}