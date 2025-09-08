package tracer

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/guregu/null"
)

func IsNullableType(t reflect.Type) bool {
	nullable := []any{
		null.String{},
		null.Bool{},
		null.Int{},
		null.Float{},
		null.Time{},
		sql.NullBool{},
		sql.NullString{},
		sql.NullByte{},
		sql.NullFloat64{},
		sql.NullTime{},
		sql.NullInt16{},
		sql.NullInt32{},
		sql.NullInt64{},
	}

	for _, nt := range nullable {
		if t.AssignableTo(reflect.TypeOf(nt)) {
			return true
		}
	}

	return false
}

func ParseNullableType(value any) string {
	switch reflect.TypeOf(value).String() {
	case "null.String":
		if v, ok := value.(null.String); ok {
			return v.ValueOrZero()
		}
	case "null.Bool":
		if v, ok := value.(null.Bool); ok {
			return fmt.Sprint(v.ValueOrZero())
		}
	case "null.Int":
		if v, ok := value.(null.Int); ok {
			return fmt.Sprint(v.ValueOrZero())
		}
	case "null.Float":
		if v, ok := value.(null.Float); ok {
			return fmt.Sprint(v.ValueOrZero())
		}
	case "null.Time":
		if v, ok := value.(null.Time); ok {
			return fmt.Sprint(v.ValueOrZero())
		}
	case "sql.NullBool":
		if v, ok := value.(sql.NullBool); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullString":
		if v, ok := value.(sql.NullString); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullByte":
		if v, ok := value.(sql.NullByte); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullFloat64":
		if v, ok := value.(sql.NullFloat64); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullTime":
		if v, ok := value.(sql.NullTime); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullInt16":
		if v, ok := value.(sql.NullInt16); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullInt32":
		if v, ok := value.(sql.NullInt32); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	case "sql.NullInt64":
		if v, ok := value.(sql.NullInt64); ok {
			realValue, err := v.Value()
			if err == nil {
				return fmt.Sprint(realValue)
			}
		}
	}
	return ""
}
