package goodbc_python_go

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/alexbrainman/odbc"
)

var (
	compareTime time.Time
)

func init() {
	_ = odbc.Driver{}
}

type multiField struct {
	Name          string
	Type          string
	IsNull        bool
	BoolValue     bool
	IntValue      int64
	FloatValue    float64
	StringValue   string
	DateTimeValue string
}

func buildRow(columns []string, values []interface{}) ([]multiField, error) {
	row := make([]multiField, len(columns))

	for i, value := range values {
		field := &row[i]
		field.Name = columns[i]

		typeOf := reflect.TypeOf(value)
		if typeOf == nil {
			field.Type = "null"
			field.IsNull = true
			continue
		}

		kind := typeOf.Kind()
		switch kind {
		case reflect.Bool:
			field.Type = "bool"
			field.BoolValue = value.(bool)
		case reflect.Uint8:
			field.Type = "int"
			field.IntValue = int64(value.(uint8))
		case reflect.Uint16:
			field.Type = "int"
			field.IntValue = int64(value.(uint16))
		case reflect.Uint32:
			field.Type = "int"
			field.IntValue = int64(value.(uint32))
		case reflect.Uint64:
			field.Type = "int"
			field.IntValue = int64(value.(uint64))
		case reflect.Int8:
			field.Type = "int"
			field.IntValue = int64(value.(int8))
		case reflect.Int16:
			field.Type = "int"
			field.IntValue = int64(value.(int16))
		case reflect.Int32:
			field.Type = "int"
			field.IntValue = int64(value.(int32))
		case reflect.Int64:
			field.Type = "int"
			field.IntValue = value.(int64)
		case reflect.Float32:
			field.Type = "float"
			field.FloatValue = float64(value.(float32))
		case reflect.Float64:
			field.Type = "float"
			field.FloatValue = value.(float64)
		case reflect.Slice:
			field.Type = "string"
			elemKind := typeOf.Elem().Kind()
			switch elemKind {
			case reflect.Uint8:
				field.StringValue = string([]byte(value.([]uint8)))
			default:
				return row, fmt.Errorf("unsupported type: %+v of %+v", kind, elemKind)
			}
		case reflect.TypeOf(compareTime).Kind():
			field.Type = "datetime"
			field.DateTimeValue = value.(time.Time).Format(time.RFC3339)
		default:
			return row, fmt.Errorf("unsupported type: %+v", kind)
		}
	}

	return row, nil
}

type session struct {
	db             *sql.DB
	dataSourceName string
	connected      bool
}

func newSession(dataSourceName string) session {
	return session{
		dataSourceName: dataSourceName,
	}
}

func (s *session) connect() error {
	if s.connected {
		return nil
	}

	db, err := sql.Open(
		"odbc",
		s.dataSourceName,
	)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	s.db = db

	s.connected = true

	return nil
}

func (s *session) query(query string) (*sql.Rows, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, err
}

func (s *session) fetchAll(rows *sql.Rows) ([][]multiField, error) {
	records := make([][]multiField, 0)

	columns, err := rows.Columns()
	if err != nil {
		return records, err
	}

	defer rows.Close()

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return records, err
			}
			break
		}

		values := make([]interface{}, len(columns))

		scanArgs := make([]interface{}, len(values))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return records, err
		}

		row, err := buildRow(columns, values)
		if err != nil {
			return records, err
		}

		records = append(records, row)
	}

	return records, nil
}

func (s *session) execute(query string) (*sql.Result, error) {
	result, err := s.db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *session) close() error {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}

	s.connected = false

	return nil
}
