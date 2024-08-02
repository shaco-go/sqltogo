package sqlparse

import (
	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/format"
	"github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/pingcap/tidb/pkg/parser/types"
	sqlerror "github.com/shaco-go/sqltogo/internal/sqlparse/errors"
	"strings"
)

type Table struct {
	tableName string
	Columns   []Column
	comment   string
}

type Column struct {
	Field    string
	Type     string
	Null     bool
	Unsigned bool
	Comment  string
	Sql      string
}

// 解析
func (t *Table) parse(sql string) (*Table, error) {
	p := parser.New()
	stmt, _, err := p.ParseSQL(sql)
	if err != nil {
		return nil, err
	}
	if len(stmt) == 0 {
		return nil, sqlerror.DDLParseFail
	}
	dom := stmt[0]
	var table Table
	_, ok := dom.Accept(&table)
	if !ok {
		return nil, sqlerror.DDLParseFail
	}
	return &table, nil
}

// 解析列
func (t *Table) parseColumn(dom *ast.ColumnDef) {
	c := Column{
		Field: dom.Name.String(),
		Type:  types.TypeStr(dom.Tp.GetType()),
	}
	for _, item := range dom.Options {
		switch item.Tp {
		case ast.ColumnOptionNull:
			c.Null = true
		case ast.ColumnOptionComment:
			if val, ok := item.Expr.(*test_driver.ValueExpr); ok {
				c.Comment = val.Datum.GetString()
			}
		}
	}

	var sb strings.Builder
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)
	err := dom.Restore(ctx)
	if err == nil {
		c.Sql = sb.String()
	}
	t.Columns = append(t.Columns, c)
}

func (t *Table) Enter(in ast.Node) (ast.Node, bool) {
	// 表名
	if val, ok := in.(*ast.TableName); ok {
		t.tableName = val.Name.String()
	}
	// 表注释
	if val, ok := in.(*ast.TableOption); ok {
		if val.Tp == ast.TableOptionComment {
			t.comment = val.StrValue
		}
	}
	// 解析列
	if val, ok := in.(*ast.ColumnDef); ok {
		t.parseColumn(val)
	}
	return in, false
}

func (t *Table) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
