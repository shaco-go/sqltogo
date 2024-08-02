package sqlparse

import (
	"github.com/samber/lo"
	"github.com/shaco-go/sqltogo/internal/sqlparse/utils"
)

type NamingStyle int

const (
	PascalCase NamingStyle = iota + 1
	CamelCase
	KebabCase
	SnakeCase
)

type Conf struct {
	NamingStyle NamingStyle `json:"naming_style"`
	Tags        []Tag       `json:"tags"`
	Mapping     []Mapping   `json:"mapping"`
	Comment     bool        `json:"comment"`
	MappingRela map[string]string
}

type Tag struct {
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type Mapping struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Null     bool   `json:"null"`
	Unsigned bool   `json:"unsigned"`
}

func NewDefTags() []Tag {
	return []Tag{
		{
			Name:   "form",
			Enable: true,
		},
		{
			Name:   "json",
			Enable: true,
		},
		{
			Name:   "gorm",
			Enable: true,
		},
	}
}

// NewDefMapping 默认的映射关系
func NewDefMapping() []Mapping {
	var m = make([]Mapping, 0, 30)

	// 整数类型
	m = append(m, Mapping{Name: "tinyint", Type: "tinyint", Value: "int8"})
	m = append(m, Mapping{Name: "smallint", Type: "smallint", Value: "int16"})
	m = append(m, Mapping{Name: "mediumint", Type: "mediumint", Value: "int32"})
	m = append(m, Mapping{Name: "int", Type: "int", Value: "int32"})
	m = append(m, Mapping{Name: "bigint", Type: "bigint", Value: "int64"})

	// 引用类型的整数
	m = append(m, Mapping{Name: "tinyint:null", Type: "tinyint", Value: "*int8"})
	m = append(m, Mapping{Name: "smallint:null", Type: "smallint", Value: "*int16"})
	m = append(m, Mapping{Name: "mediumint:null", Type: "mediumint", Value: "*int32"})
	m = append(m, Mapping{Name: "int:null", Type: "int", Value: "*int32"})
	m = append(m, Mapping{Name: "bigint:null", Type: "bigint", Value: "*int64"})

	// 无符号整数类型
	m = append(m, Mapping{Name: "tinyint,unsigned", Type: "tinyint", Value: "uint8"})
	m = append(m, Mapping{Name: "smallint,unsigned", Type: "smallint", Value: "uint16"})
	m = append(m, Mapping{Name: "mediumint,unsigned", Type: "mediumint", Value: "uint32"})
	m = append(m, Mapping{Name: "int,unsigned", Type: "int", Value: "uint32"})
	m = append(m, Mapping{Name: "bigint,unsigned", Type: "bigint", Value: "uint64"})

	// 引用类型的无符号整数
	m = append(m, Mapping{Name: "tinyint:null,unsigned", Type: "tinyint", Value: "*uint8"})
	m = append(m, Mapping{Name: "smallint:null,unsigned", Type: "smallint", Value: "*uint16"})
	m = append(m, Mapping{Name: "mediumint:null,unsigned", Type: "mediumint", Value: "*uint32"})
	m = append(m, Mapping{Name: "int:null,unsigned", Type: "int", Value: "*uint32"})
	m = append(m, Mapping{Name: "bigint:null,unsigned", Type: "bigint", Value: "*uint64"})

	// 浮点类型
	m = append(m, Mapping{Name: "float", Type: "float", Value: "float32"})
	m = append(m, Mapping{Name: "double", Type: "double", Value: "float64"})

	// 引用类型的浮点数
	m = append(m, Mapping{Name: "float:null", Type: "float", Value: "*float32"})
	m = append(m, Mapping{Name: "double:null", Type: "double", Value: "*float64"})

	// 定点类型
	m = append(m, Mapping{Name: "decimal", Type: "decimal", Value: "float64"})

	// 引用类型的定点数
	m = append(m, Mapping{Name: "decimal:null", Type: "decimal", Value: "*float64"})

	// 日期和时间类型
	m = append(m, Mapping{Name: "date", Type: "string", Value: "string"})
	m = append(m, Mapping{Name: "datetime", Type: "string", Value: "string"})
	m = append(m, Mapping{Name: "timestamp", Type: "string", Value: "string"})
	m = append(m, Mapping{Name: "time", Type: "string", Value: "string"})
	m = append(m, Mapping{Name: "year", Type: "string", Value: "string"})

	// 字符串类型
	m = append(m, Mapping{Name: "char", Type: "char", Value: "string"})
	m = append(m, Mapping{Name: "varchar", Type: "varchar", Value: "string"})

	// 文本类型
	m = append(m, Mapping{Name: "tinytext", Type: "tinytext", Value: "string"})
	m = append(m, Mapping{Name: "text", Type: "text", Value: "string"})
	m = append(m, Mapping{Name: "mediumtext", Type: "mediumtext", Value: "string"})
	m = append(m, Mapping{Name: "longtext", Type: "longtext", Value: "string"})

	// 二进制类型
	m = append(m, Mapping{Name: "binary", Type: "binary", Value: "[]byte"})
	m = append(m, Mapping{Name: "varbinary", Type: "varbinary", Value: "[]byte"})

	// 文本型二进制数据
	m = append(m, Mapping{Name: "tinyblob", Type: "tinyblob", Value: "[]byte"})
	m = append(m, Mapping{Name: "blob", Type: "blob", Value: "[]byte"})
	m = append(m, Mapping{Name: "mediumblob", Type: "mediumblob", Value: "[]byte"})
	m = append(m, Mapping{Name: "longblob", Type: "longblob", Value: "[]byte"})

	// 枚举和集合类型
	m = append(m, Mapping{Name: "enum", Type: "enum", Value: "string"})
	m = append(m, Mapping{Name: "set", Type: "set", Value: "string"})

	// JSON 类型
	m = append(m, Mapping{Name: "json", Type: "json", Value: "string"})

	return m
}

// NewConf 默认的配置
func NewConf(c *Conf) *Conf {
	conf := c
	if conf == nil {
		conf = &Conf{
			NamingStyle: SnakeCase,
			Tags:        NewDefTags(),
			Mapping:     NewDefMapping(),
			Comment:     true,
		}
	}
	conf.MappingRela = conf.convMap()
	return conf
}

func (c *Conf) convMap() map[string]string {
	var m = make(map[string]string, len(c.Mapping))
	for _, item := range c.Mapping {
		m[utils.MappingToKey(item)] = item.Value
	}
	return m
}

// ToFieldStyle 字段转风格
func (c *Conf) ToFieldStyle(field string) string {
	switch c.NamingStyle {
	case PascalCase:
		return lo.PascalCase(field)
	case CamelCase:
		return lo.CamelCase(field)
	case KebabCase:
		return lo.KebabCase(field)
	case SnakeCase:
		return lo.SnakeCase(field)
	default:
		return field
	}
}
