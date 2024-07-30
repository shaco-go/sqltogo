package modelparse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func NewModelParse(filePath string, tableName string) (*ModelParse, error) {
	var o = &ModelParse{
		tableName: tableName,
		filePath:  filePath,
	}
	// 解析文件
	err := o.astParse()
	if err != nil {
		return nil, err
	}
	// 解析结构体
	err = o.parseStructName()
	if err != nil {
		return nil, err
	}
	// 解析字段
	o.parseStructField()
	return o, nil
}

type ModelParse struct {
	tableName   string
	filePath    string
	structName  string
	structField []*ast.Field
	node        *ast.File
}

// AstParse 解析文件
func (m *ModelParse) astParse() error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, m.filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	m.node = node
	return nil
}

// GetStructName 获取对应的结构体名
func (m *ModelParse) parseStructName() error {
	for _, decl := range m.node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok && funcDecl.Recv != nil && funcDecl.Name.Name == "TableName" {
			// 查看return返回的表名
			for _, stmt := range funcDecl.Body.List {
				returnStmt, ok := stmt.(*ast.ReturnStmt)
				if ok {
					val, ok := returnStmt.Results[0].(*ast.BasicLit)
					if ok && val.Value == fmt.Sprintf("\"%s\"", m.tableName) && len(funcDecl.Recv.List) > 0 {
						expr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
						if ok {
							ident, ok := expr.X.(*ast.Ident)
							if ok {
								m.structName = ident.Name
								return nil
							}
						}
					}
				}
			}
		}
	}
	return NOT_FOUND_STRUCT
}

// GetStructField 获取结构体的字段
func (m *ModelParse) parseStructField() {
	for _, decl := range m.node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			if spec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
				if structType, ok := spec.Type.(*ast.StructType); ok && spec.Name.Name == m.structName {
					m.structField = structType.Fields.List
				}
			}
		}
	}
}
