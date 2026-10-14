package finder

import (
	"fmt"
	"slices"
	"strings"
)

func stableSortColumn(
	tableAlias string,
	idColumn *string,
	columns *[]string,
) (string, string) {
	if columns != nil && slices.Contains(*columns, "id") {
		return "id", tableAlias + ".id"
	}
	if idColumn == nil {
		return "", ""
	}

	column := strings.TrimSpace(*idColumn)
	if column == "" {
		return "", ""
	}
	return strings.ToLower(column), tableAlias + "." + column
}

func appendStableSort(
	orderBy string,
	tableAlias string,
	idColumn *string,
	columns *[]string,
) string {
	_, qualifiedID := stableSortColumn(tableAlias, idColumn, columns)
	if qualifiedID == "" {
		return orderBy
	}

	normalizedOrder := strings.ToLower(strings.Join(strings.Fields(orderBy), ""))
	normalizedID := strings.ToLower(qualifiedID)
	if strings.Contains(normalizedOrder, normalizedID+"asc") ||
		strings.Contains(normalizedOrder, normalizedID+"desc") {
		return orderBy
	}
	if strings.TrimSpace(orderBy) == "" {
		return qualifiedID + " DESC"
	}
	return orderBy + "," + qualifiedID + " DESC"
}

// sortBuilder returns a string containing confirmed sorts only.
func sortBuilder(
	tableAlias, urlSorts string,
	idColumn *string,
	columns *[]string,
) string {
	var sortsString string

	sortsString = strings.ReplaceAll(urlSorts, " ", "")
	sortsString = strings.ReplaceAll(sortsString, ";", "")
	sortsString = strings.ReplaceAll(sortsString, "\\", "")

	sortsArr := strings.Split(sortsString, ",")

	var verifiedSortsArr []string

	stableColumn, _ := stableSortColumn(tableAlias, idColumn, columns)
	hasStableSort := false
	// In case no sort is provided it defaults to -created_at or -id
	if len(sortsArr) == 1 && sortsArr[0] == "" {
		var orderBy string
		if slices.Contains(*columns, "created_at") {
			orderBy = tableAlias + ".created_at DESC"
		} else if slices.Contains(*columns, "id") {
			orderBy = tableAlias + ".id DESC"
		} else if idColumn != nil {
			orderBy = fmt.Sprintf("%s.%s DESC", tableAlias, *idColumn)
		}
		verifiedSortsArr = append(verifiedSortsArr, orderBy)
		hasStableSort = stableColumn != "" && !slices.Contains(*columns, "created_at")
	}

	for _, value := range sortsArr {

		var orderBy string
		col := strings.ToLower(value)

		if strings.HasPrefix(value, "-") {
			col = strings.TrimPrefix(value, "-")
			orderBy = tableAlias + "." + col + " DESC"
		} else {
			orderBy = tableAlias + "." + col + " ASC"
		}
		pass := slices.Contains(*columns, col)

		if pass {
			if col == stableColumn {
				hasStableSort = true
			}
			verifiedSortsArr = append(verifiedSortsArr, orderBy)
		}
	}
	sortsString = strings.Join(verifiedSortsArr, ",")
	if !hasStableSort {
		sortsString = appendStableSort(sortsString, tableAlias, idColumn, columns)
	}

	return sortsString
}
