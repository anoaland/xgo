package repository

import "strings"

func SetKeywordLikeVarsByTotalExpr(keyword string, total int) (vars []interface{}) {
	for i := 0; i < total; i++ {
		vars = append(vars, "%"+strings.ToLower(keyword)+"%")
	}

	return vars
}

func LowerLikeQuery(field string) string {
	return "LOWER(" + field + ") LIKE ?"
}
