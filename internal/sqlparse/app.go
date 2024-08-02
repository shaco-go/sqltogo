package sqlparse

import (
	_ "embed"
	"fmt"
	"github.com/samber/lo"
	"github.com/shaco-go/sqltogo/internal/sqlparse/utils"
	"go/format"
	"html/template"
	"strings"
)

//go:embed template/model.tpl
var modelTpl string

type App struct {
	table *Table
	conf  *Conf
}

func NewApp(ddl string, conf *Conf) (*App, error) {
	var t Table
	var a = &App{
		conf: NewConf(conf),
	}
	table, err := t.parse(ddl)
	if err != nil {
		return nil, err
	}
	a.table = table
	return a, nil
}

// TableName 	生成表名
func (a *App) TableName() string {
	return a.table.tableName
}

// TableComment 	表注释
func (a *App) TableComment() string {
	return a.table.comment
}

func (a *App) Columns() []string {
	var columnsGroup = make([]string, 0, len(a.table.Columns))
	for _, col := range a.table.Columns {
		columnsGroup = append(columnsGroup, a.ColumnLine(col))
	}
	return columnsGroup
}

func (a *App) ColumnLine(col Column) string {
	type Col struct {
		Field   string
		Type    string
		Tag     []string
		Comment string
	}
	var column = Col{
		Field: a.conf.ToFieldStyle(col.Field),
		Type:  "any",
	}
	if a.conf.Comment {
		column.Comment = col.Comment
	}
	if val, ok := a.conf.MappingRela[utils.ColumnToKey(col)]; ok {
		column.Type = val
	}
	// tags
	for _, tagConf := range a.conf.Tags {
		switch strings.ToLower(tagConf.Name) {
		case "gorm":
			column.Tag = append(column.Tag, fmt.Sprintf(`gorm:"column:%s"`, col.Field))
		default:
			column.Tag = append(column.Tag, fmt.Sprintf(`%s:"%s"`, strings.ToLower(tagConf.Name), col.Field))
		}
	}
	// 生成
	str := fmt.Sprintf(`%s %s`, column.Type, column.Tag)
	if len(column.Tag) > 0 {
		str += fmt.Sprintf("`%s`", strings.Join(column.Tag, " "))
	}
	if column.Comment != "" {
		str += fmt.Sprintf(" // %s", column.Comment)
	}
	return str
}

func (a *App) BuildTpl() (string, error) {
	files, err := template.New("template").Funcs(map[string]any{
		"pascalCase": lo.PascalCase,
		"getWord": func(word string) string {
			var a = []rune(word)
			return string(a[0])
		},
	}).Parse(modelTpl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var builder strings.Builder
	err = files.Execute(&builder, a)
	source, _ := format.Source([]byte(builder.String()))
	return string(source), nil
}
