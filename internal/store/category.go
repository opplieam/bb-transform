package store

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/opplieam/bb-transform/.jetgen/postgres/public/model"
	. "github.com/opplieam/bb-transform/.jetgen/postgres/public/table"
	"github.com/opplieam/bb-transform/internal/transform"
)

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

type AllCategoryResult struct {
	ID       int32  `sql:"primary_key" alias:"category.id" `
	Name     string `alias:"category.name" `
	HasChild bool   `alias:"category.has_child" `
	Path     string
}

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

func (c *CategoryStore) CleanUp(version string) error {
	stmt := CategoryDataset.DELETE().WHERE(CategoryDataset.Version.EQ(String(version)))
	_, err := stmt.Exec(c.db)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryStore) InsertDataset(dataset []model.CategoryDataset) error {
	stmt := CategoryDataset.INSERT(CategoryDataset.AllColumns.Except(CategoryDataset.ID)).MODELS(dataset)
	_, err := stmt.Exec(c.db)
	if err != nil {
		return err
	}
	return nil
}
