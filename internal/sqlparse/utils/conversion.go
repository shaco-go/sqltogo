package utils

import (
	"github.com/pkg/errors"
	"github.com/shaco-go/sqltogo/internal/sqlparse"
	"strings"
)

func MappingToKey(m sqlparse.Mapping) string {
	var key = m.Type
	var suffix []string
	if m.Null {
		suffix = append(suffix, "null")
	}
	if m.Unsigned {
		suffix = append(suffix, "unsigned")
	}
	if len(suffix) > 0 {
		key = key + ":" + strings.Join(suffix, ",")
	}
	return key
}

func ColumnToKey(c sqlparse.Column) string {
	return MappingToKey(sqlparse.Mapping{
		Name:     c.Type,
		Type:     c.Type,
		Value:    "",
		Null:     c.Null,
		Unsigned: c.Unsigned,
	})
}

func KeyToMapping(key string, val string) (*sqlparse.Mapping, error) {
	split := strings.Split(key, ":")
	var m = &sqlparse.Mapping{Name: key, Value: val}
	if len(split) == 0 {
		return nil, errors.New("key is empty")
	}
	m.Type = split[0]
	if len(split) == 1 {
		return m, nil
	}
	if len(split) == 2 {
		sp := strings.Split(split[1], ",")
		for _, val := range sp {
			switch strings.ToLower(val) {
			case "null":
				m.Null = true
			case "unsigned":
				m.Unsigned = true
			}
		}
	}
	return nil, errors.New("key is invalid")
}
