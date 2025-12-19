package db

import (
	"github.com/gflydev/core"
	"github.com/gflydev/utils/arr"
	"github.com/jivegroup/fluentsql"
)

type acceptedOrdersConstraint interface {
	[]string | core.Data
}

// ProcessOrderBy processes the order by parameter and returns the order key and direction.
//
// Parameters:
//   - orderBy (string): The order by parameter that may include a direction prefix (e.g. "-created_at" for descending)
//   - acceptedOrders (any): List of accepted order keys ([]string) or map of field to column mappings (core.Data)
//   - defaultOrderKey (...any): Optional default order key and direction. First argument is key (string), second is direction (fluentsql.OrderByDir)
//
// Examples:
//
//	mb.ProcessOrderBy(builder, "title", []string{"created_at", "title"}, "created_at", mb.Desc)
//	mb.ProcessOrderBy(builder, "-id", core.Data{"id": "categories.id", "name": "categories.name"})
func ProcessOrderBy[T acceptedOrdersConstraint](db *DBModel, orderBy string, acceptedOrders T, defaultOrderKey ...any) {
	direction := Asc
	orderKey := orderBy

	// Try to get input value
	if orderBy != "" {
		if orderBy[0] == '-' {
			orderKey = orderBy[1:]
			direction = Desc
		}

		// Handle different types of acceptedOrders
		switch v := any(acceptedOrders).(type) {
		case []string:
			if ok := arr.Contains[string](v, orderKey); !ok {
				orderKey = ""
			}
		case core.Data:
			if columnName, ok := v[orderKey]; ok {
				orderKey = columnName.(string)
			} else {
				orderKey = ""
			}
		}
	}

	// Try to get default value
	if orderKey == "" && len(defaultOrderKey) > 0 {
		for _, key := range defaultOrderKey {
			switch v := key.(type) {
			case string:
				orderKey = v
			case fluentsql.OrderByDir:
				direction = v
			}
		}
	}

	if orderKey != "" {
		db.OrderBy(orderKey, direction)
	}
}
