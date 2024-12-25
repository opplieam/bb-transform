//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var MatchCategory = newMatchCategoryTable("public", "match_category", "")

type matchCategoryTable struct {
	postgres.Table

	// Columns
	ID      postgres.ColumnInteger
	L1      postgres.ColumnString
	L2      postgres.ColumnString
	L3      postgres.ColumnString
	L4      postgres.ColumnString
	L5      postgres.ColumnString
	L6      postgres.ColumnString
	L7      postgres.ColumnString
	L8      postgres.ColumnString
	MatchID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MatchCategoryTable struct {
	matchCategoryTable

	EXCLUDED matchCategoryTable
}

// AS creates new MatchCategoryTable with assigned alias
func (a MatchCategoryTable) AS(alias string) *MatchCategoryTable {
	return newMatchCategoryTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MatchCategoryTable with assigned schema name
func (a MatchCategoryTable) FromSchema(schemaName string) *MatchCategoryTable {
	return newMatchCategoryTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MatchCategoryTable with assigned table prefix
func (a MatchCategoryTable) WithPrefix(prefix string) *MatchCategoryTable {
	return newMatchCategoryTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MatchCategoryTable with assigned table suffix
func (a MatchCategoryTable) WithSuffix(suffix string) *MatchCategoryTable {
	return newMatchCategoryTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMatchCategoryTable(schemaName, tableName, alias string) *MatchCategoryTable {
	return &MatchCategoryTable{
		matchCategoryTable: newMatchCategoryTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newMatchCategoryTableImpl("", "excluded", ""),
	}
}

func newMatchCategoryTableImpl(schemaName, tableName, alias string) matchCategoryTable {
	var (
		IDColumn       = postgres.IntegerColumn("id")
		L1Column       = postgres.StringColumn("l1")
		L2Column       = postgres.StringColumn("l2")
		L3Column       = postgres.StringColumn("l3")
		L4Column       = postgres.StringColumn("l4")
		L5Column       = postgres.StringColumn("l5")
		L6Column       = postgres.StringColumn("l6")
		L7Column       = postgres.StringColumn("l7")
		L8Column       = postgres.StringColumn("l8")
		MatchIDColumn  = postgres.IntegerColumn("match_id")
		allColumns     = postgres.ColumnList{IDColumn, L1Column, L2Column, L3Column, L4Column, L5Column, L6Column, L7Column, L8Column, MatchIDColumn}
		mutableColumns = postgres.ColumnList{L1Column, L2Column, L3Column, L4Column, L5Column, L6Column, L7Column, L8Column, MatchIDColumn}
	)

	return matchCategoryTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:      IDColumn,
		L1:      L1Column,
		L2:      L2Column,
		L3:      L3Column,
		L4:      L4Column,
		L5:      L5Column,
		L6:      L6Column,
		L7:      L7Column,
		L8:      L8Column,
		MatchID: MatchIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
