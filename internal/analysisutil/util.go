// Package analysisutil предоставляет вспомогательные функции для анализа AST.
package analysisutil

import "go/ast"

// IsStringLiteral проверяет, является ли выражение строковым литералом.
func IsStringLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind.String() == "STRING"
}

// GetLoggerName извлекает имя логгера из SelectorExpr.
func GetLoggerName(sel *ast.SelectorExpr) string {
	if ident, ok := sel.X.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// GetMethodName извлекает имя метода из SelectorExpr.
func GetMethodName(sel *ast.SelectorExpr) string {
	if sel != nil && sel.Sel != nil {
		return sel.Sel.Name
	}
	return ""
}
