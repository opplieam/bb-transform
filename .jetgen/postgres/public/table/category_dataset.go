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

var CategoryDataset = newCategoryDatasetTable("public", "category_dataset", "")

type categoryDatasetTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnInteger
	L1In        postgres.ColumnString
	L2In        postgres.ColumnString
	L3In        postgres.ColumnString
	L4In        postgres.ColumnString
	L5In        postgres.ColumnString
	L6In        postgres.ColumnString
	L7In        postgres.ColumnString
	L8In        postgres.ColumnString
	FullPathOut postgres.ColumnString
	NameOut     postgres.ColumnString
	Version     postgres.ColumnString
	Label       postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CategoryDatasetTable struct {
	categoryDatasetTable

	EXCLUDED categoryDatasetTable
}

// AS creates new CategoryDatasetTable with assigned alias
func (a CategoryDatasetTable) AS(alias string) *CategoryDatasetTable {
	return newCategoryDatasetTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CategoryDatasetTable with assigned schema name
func (a CategoryDatasetTable) FromSchema(schemaName string) *CategoryDatasetTable {
	return newCategoryDatasetTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CategoryDatasetTable with assigned table prefix
func (a CategoryDatasetTable) WithPrefix(prefix string) *CategoryDatasetTable {
	return newCategoryDatasetTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CategoryDatasetTable with assigned table suffix
func (a CategoryDatasetTable) WithSuffix(suffix string) *CategoryDatasetTable {
	return newCategoryDatasetTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCategoryDatasetTable(schemaName, tableName, alias string) *CategoryDatasetTable {
	return &CategoryDatasetTable{
		categoryDatasetTable: newCategoryDatasetTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newCategoryDatasetTableImpl("", "excluded", ""),
	}
}

func newCategoryDatasetTableImpl(schemaName, tableName, alias string) categoryDatasetTable {
	var (
		IDColumn          = postgres.IntegerColumn("id")
		L1InColumn        = postgres.StringColumn("l1_in")
		L2InColumn        = postgres.StringColumn("l2_in")
		L3InColumn        = postgres.StringColumn("l3_in")
		L4InColumn        = postgres.StringColumn("l4_in")
		L5InColumn        = postgres.StringColumn("l5_in")
		L6InColumn        = postgres.StringColumn("l6_in")
		L7InColumn        = postgres.StringColumn("l7_in")
		L8InColumn        = postgres.StringColumn("l8_in")
		FullPathOutColumn = postgres.StringColumn("full_path_out")
		NameOutColumn     = postgres.StringColumn("name_out")
		VersionColumn     = postgres.StringColumn("version")
		LabelColumn       = postgres.StringColumn("label")
		allColumns        = postgres.ColumnList{IDColumn, L1InColumn, L2InColumn, L3InColumn, L4InColumn, L5InColumn, L6InColumn, L7InColumn, L8InColumn, FullPathOutColumn, NameOutColumn, VersionColumn, LabelColumn}
		mutableColumns    = postgres.ColumnList{L1InColumn, L2InColumn, L3InColumn, L4InColumn, L5InColumn, L6InColumn, L7InColumn, L8InColumn, FullPathOutColumn, NameOutColumn, VersionColumn, LabelColumn}
	)

	return categoryDatasetTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		L1In:        L1InColumn,
		L2In:        L2InColumn,
		L3In:        L3InColumn,
		L4In:        L4InColumn,
		L5In:        L5InColumn,
		L6In:        L6InColumn,
		L7In:        L7InColumn,
		L8In:        L8InColumn,
		FullPathOut: FullPathOutColumn,
		NameOut:     NameOutColumn,
		Version:     VersionColumn,
		Label:       LabelColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}