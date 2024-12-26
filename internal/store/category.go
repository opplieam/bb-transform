// Package store provides data access objects and methods for interacting with the database.
// It includes functionalities for fetching and manipulating category data used in machine learning dataset generation.
package store

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/opplieam/bb-transform/.jetgen/postgres/public/model"
	. "github.com/opplieam/bb-transform/.jetgen/postgres/public/table"
	"github.com/opplieam/bb-transform/internal/transform"
)

// CategoryStore provides methods for interacting with the category related tables in the database.
// It uses a sql.DB connection to execute queries and manage category data.
type CategoryStore struct {
	db *sql.DB
}

// NewCategoryStore creates a new instance of CategoryStore.
// It takes a *sql.DB connection as input and returns a pointer to a CategoryStore.
func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

// AllCategoryResult represents the result structure for a query that fetches all categories,
// including their ID, name, whether they have children, and their hierarchical path.
// The structure is designed to be used with a recursive CTE query to extract hierarchical data.
type AllCategoryResult struct {
	ID       int32  `sql:"primary_key" alias:"category.id" `
	Name     string `alias:"category.name" `
	HasChild bool   `alias:"category.has_child" `
	Path     string
}

// OriginalCategory retrieves the original category structure from the database.
// It constructs a recursive Common Table Expression (CTE) to traverse the category hierarchy,
// starting from categories with no parent (root categories) and recursively joining
// child categories until it reaches the deepest level. The function then filters
// the result to include only the deepest categories (categories without children),
// which represent the end of each branch in the hierarchy. The result is a map
// where the key is the category ID and the value is a struct containing the category's
// name and its full hierarchical path, represented as a string concatenated with " > ".
// The function returns an error if any issues occur during the database query.
func (c *CategoryStore) OriginalCategory() (transform.Category, error) {
	cr := CTE("CategoryRecursive")
	pathCol := StringColumn("AllCategoryResult.Path").From(cr)
	stmt := WITH_RECURSIVE(
		cr.AS(
			SELECT(
				Category.ID, Category.Name, Category.HasChild, CAST(Category.Name).AS_TEXT().AS(pathCol.Name()),
			).FROM(
				Category,
			).WHERE(
				Category.ParentID.IS_NULL(),
			).UNION(
				SELECT(
					Category.ID, Category.Name, Category.HasChild, CAST(pathCol.CONCAT(String(" > ")).CONCAT(Category.Name)).AS_TEXT(),
				).FROM(
					Category.
						INNER_JOIN(cr, Category.ParentID.EQ(Category.ID.From(cr))),
				),
			),
		),
	)(
		SELECT(
			cr.AllColumns(),
		).FROM(
			cr,
		).WHERE(
			Category.HasChild.From(cr).IS_FALSE(),
		),
	)
	var dest []AllCategoryResult
	if err := stmt.Query(c.db, &dest); err != nil {
		return nil, err
	}

	catResult := make(transform.Category)
	for _, v := range dest {
		catResult[v.ID] = transform.CategoryDeepest{
			Name: v.Name,
			Path: v.Path,
		}
	}
	return catResult, nil
}

// MatchedCategory retrieves all matched categories from the 'match_category' table where 'match_id' is not null.
// It returns a slice of model.MatchCategory representing the matched categories or an error if the query fails.
func (c *CategoryStore) MatchedCategory() ([]model.MatchCategory, error) {
	stmt := SELECT(
		MatchCategory.AllColumns,
	).FROM(
		MatchCategory,
	).WHERE(
		MatchCategory.MatchID.IS_NOT_NULL(),
	)

	var dest []model.MatchCategory
	if err := stmt.Query(c.db, &dest); err != nil {
		return nil, err
	}
	return dest, nil
}

// CleanUp removes all category datasets from the 'category_dataset' table that match a specific version.
// It takes a version string as input and returns an error if the deletion fails.
func (c *CategoryStore) CleanUp(version string) error {
	stmt := CategoryDataset.DELETE().WHERE(CategoryDataset.Version.EQ(String(version)))
	_, err := stmt.Exec(c.db)
	if err != nil {
		return err
	}
	return nil
}

// InsertDataset inserts multiple category dataset records into the 'category_dataset' table.
// It takes a slice of model.CategoryDataset as input, excluding the 'id' column which is assumed to be auto-generated.
// It returns an error if the insertion fails.
func (c *CategoryStore) InsertDataset(dataset []model.CategoryDataset) error {
	stmt := CategoryDataset.INSERT(CategoryDataset.AllColumns.Except(CategoryDataset.ID)).MODELS(dataset)
	_, err := stmt.Exec(c.db)
	if err != nil {
		return err
	}
	return nil
}
