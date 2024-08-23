package db

import (
	"database/sql"
	"errors"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	qb "github.com/jiveio/fluentsql" // Query builder
)

// ====================================================================
// ==================== Generic repository methods ====================
// ====================================================================

type Filter core.Data

// GetModel takes a model, query field, and condition value.
// query for getting one model by a specific field condition
func GetModel[T any](m *T, field string, value any, filters ...Filter) error {
	// DB Model instance
	db := Instance()
	var err error
	// Try/catch block
	try.Perform(func() {
		db.Where(field, qb.Eq, value)

		for field, value := range filters {
			db.Where(field, qb.Eq, value)
		}
		// Get first record then assign to `m` WHERE field = value
		err = db.First(m)
	}).Catch(func(e try.E) {
		// Log unexpected error!
		err = e.(error)
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error(err)
		}
	})
	return err
}

// ListModels is a query that returns slice of models, paginate by limit and offset
func ListModels[T any](page, limit int, field string, value any, filters ...Filter) ([]*T, int, error) {
	var models []*T
	var total int
	var err error

	db := Instance()

	var offset = 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	try.Perform(func() {
		db.Where(field, qb.Eq, value)
		for field, value := range filters {
			db.Where(field, qb.Eq, value)
		}

		total, err = db.Limit(limit, offset).Find(models)
		if err != nil {
			try.Throw(err)
		}
	}).Catch(func(e try.E) {
		err = e.(error)
	})

	return models, total, err
}

// CreateModel a query that creating a model by given model's data
func CreateModel[T any](m *T) error {
	db := Instance()
	var err error

	try.Perform(func() {
		// Begin transaction
		db.Begin()
		// Trying to create an instance
		if e := db.Create(m); e != nil {
			try.Throw(e)
		}
		// Commit transaction
		err = db.Commit()
	}).Catch(func(e try.E) {
		err = e.(error)
		// Rollback transaction
		_ = db.Rollback()
	})

	return err
}

// UpdateModel a query that updating a model by given model's data
func UpdateModel[T any](m *T) error {
	db := Instance()
	var err error

	try.Perform(func() {
		db.Begin()

		// Trying to update an instance
		if e := db.Update(m); e != nil {
			try.Throw(e)
		}

		err = db.Commit()
	}).Catch(func(e try.E) {
		err = e.(error)
		_ = db.Rollback()
	})

	return err
}

// DeleteModel a query that deleting a model by given model data.
func DeleteModel[T any](m *T) error {
	db := Instance()
	var err error

	try.Perform(func() {
		db.Begin()

		if e := db.Delete(m); e != nil {
			try.Throw(e)
		}

		err = db.Commit()
	}).Catch(func(e try.E) {
		err = e.(error)
		_ = db.Rollback()
	})

	return err
}
