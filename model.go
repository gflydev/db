/*
Convert from model struct to table structure.
*/

package db

import (
	"fmt"
	"github.com/gflydev/core/errors"
	"reflect"
	"regexp"
	"strings"
)

// ====================================================================
//                         Struct and Data
// ====================================================================

// Constants for model tag attributes and processing
const (
	MODEL     = "model"   // Tag `model` used to store metadata for struct fields
	TABLE     = "table"   // Table name in the database
	TYPE      = "type"    // Column types for database representation
	REFERENCE = "ref"     // Column reference to another table
	CASCADE   = "cascade" // Cascade rules for DELETE and UPDATE
	RELATION  = "rel"     // Relation to another table
	NAME      = "name"    // Column name in the database
)

// MetaData represents metadata string information
type MetaData string

// Table structure that maps a Go struct to a database table
type Table struct {
	Name      string         // Name of the table
	Columns   []Column       // List of columns in the table
	Primaries []Column       // List of primary key columns
	Values    map[string]any // Values of the table columns
	Relation  []*Table       // Related tables for relational mapping
	HasData   bool           // Indicates if the table has valid data
}

// Column structure that maps a struct field to a database table column
type Column struct {
	Key      string // Name of the struct field the column maps to
	Name     string // Name of the database column
	Primary  bool   // Indicates if the column is a primary key
	Types    string // Data type of the column
	Ref      string // Reference to another table column
	Relation string // Relation to another table
	IsZero   bool   // Indicates if the column value is the zero value for its type
	HasValue bool   // Indicates if the column has a valid (non-zero) value
}

// isNotData determines if the column is not valid data for the table
//
// Returns:
//
//	bool - true if the column is not associated with valid data
func (c *Column) isNotData() bool {
	return !c.HasValue || c.Relation != "" || c.Ref != ""
}

// NewTable initializes a new Table instance
//
// Returns:
//
//	*Table - A pointer to the initialized Table struct
func NewTable() *Table {
	tbl := new(Table)
	tbl.Values = make(map[string]any)

	return tbl
}

// ModelData converts a Go struct into a Table representation
//
// Parameters:
//
//	model (any): The struct or pointer to struct to be converted
//
// Returns:
//
//	*Table - The table structure representing the input struct
//	error  - An error if the input is not a struct or pointer to struct
func ModelData(model any) (*Table, error) {
	// Get the type and value of the input model
	typ := reflect.TypeOf(model)
	value := reflect.ValueOf(model)

	// Create a new Table structure
	tbl := NewTable()

	// Convert the struct name to snake_case for table name
	tbl.Name = toSnakeCase(typ.Name())

	// Process the model fields if it's a struct or a pointer to struct
	if typ.Kind() == reflect.Struct {
		return processModel(typ, value, tbl), nil
	} else if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		return processModel(typ.Elem(), value, tbl), nil
	}

	// Return an error if the input is not a struct
	return nil, errors.New("input param should be a struct")
}

// ====================================================================
//                         Process methods
// ====================================================================

// processModel converts a struct's fields into a Table representation.
//
// Parameters:
//   - typ (reflect.Type): The type of the struct.
//   - value (reflect.Value): The value of the struct.
//   - tbl (*Table): The Table instance to populate.
//
// Returns:
//   - *Table: The populated Table instance.
func processModel(typ reflect.Type, value reflect.Value, tbl *Table) *Table {
	// Pointer case: Dereference pointers to get the underlying value.
	value = reflect.Indirect(value)

	// Loop through the fields of the struct.
	for i := 0; i < typ.NumField(); i++ {
		var col Column // Define a Column to store field metadata.

		// Get type field (information about the struct field).
		typeField := typ.Field(i)

		// Get the value of the field.
		var valueField reflect.Value
		// FIX: Handle panic caused by field index mismatch.
		if value.NumField() == typ.NumField() {
			valueField = value.Field(i)
		} else {
			valueField = reflect.ValueOf(nil)
		}

		// Extract attributes from the `model` tag.
		attr := readTags(typeField.Tag.Get(MODEL))

		// Check if the field is a primary column by analyzing its `type` attribute.
		isPrimaryColumn := isPrimary(attr[TYPE])

		// Process special MetaData type fields for table-specific settings.
		if typeField.Type == reflect.TypeOf(MetaData("")) {
			if slice, tableOk := attr[TABLE]; tableOk && len(slice) > 0 {
				tbl.Name = slice[0]
			}
			continue
		}

		// Set the database column name.
		if slice, nameOk := attr[NAME]; nameOk && len(slice) > 0 {
			col.Name = slice[0]
		} else {
			// Default to snake_case conversion of the struct field name.
			col.Name = toSnakeCase(typeField.Name)
		}
		col.Key = typeField.Name // Store the original struct field name.

		// Check if the field has a valid value.
		validValue := valueField.IsValid()

		// Handle zero values for primary columns or other fields.
		validValueType := (isPrimaryColumn && valueField.CanInt() && !valueField.IsZero()) || !isPrimaryColumn

		// Store the value in the Table and update column metadata.
		if validValue && validValueType {
			tbl.Values[col.Name] = valueField.Interface()
			tbl.HasData = true
			col.HasValue = true
			col.IsZero = valueField.IsZero()
		} else {
			tbl.Values[col.Name] = nil
			col.HasValue = false
			col.IsZero = true
		}

		// Process types for the column from the `type` attribute.
		if slice, typeOk := attr[TYPE]; typeOk {
			col.Types = getTypes(slice)
		}

		// Process references if the column has a `ref` attribute.
		if slice, refOk := attr[REFERENCE]; refOk {
			col.Ref = getReferences(slice[0], col.Name)
		}

		// Process cascade rules if the column has a `cascade` attribute.
		if slice, casOk := attr[CASCADE]; casOk {
			col.Ref += getCascade(slice)
		}

		// Process relationships if the column has a `rel` attribute.
		if slice, relOk := attr[RELATION]; relOk && validValue {
			col.Relation = toSnakeCase(slice[0])

			// Handle relationships for arrays (slices) or single structs.
			if typeField.Type.Kind() == reflect.Slice {
				for n := 0; n < valueField.Len(); n++ {
					elemVal := valueField.Index(n)
					_tbl, _ := ModelData(elemVal.Interface()) // Convert slice element to Table.
					tbl.Relation = append(tbl.Relation, _tbl)
				}
			} else {
				_tbl, _ := ModelData(valueField.Interface()) // Convert single field value to Table.
				tbl.Relation = append(tbl.Relation, _tbl)
			}
		}

		// Mark the column as a primary key if applicable.
		col.Primary = isPrimaryColumn
		if col.Primary {
			tbl.Primaries = append(tbl.Primaries, col)
		}

		// Add the column to the list of table columns.
		tbl.Columns = append(tbl.Columns, col)
	}

	// Return the populated Table instance.
	return tbl
}

// readTags parses the `model` tag string into a map of attributes and their corresponding values.
//
// Parameters:
//   - tags (string): A semicolon-separated string of attributes in the format 'key:value1,value2,...'.
//
// Returns:
//   - map[string][]string: A map where keys are attribute names (e.g., "type") and values are slices
//     of strings representing the attribute values (e.g., []string{"BOOLEAN"}).
func readTags(tags string) map[string][]string {
	// If the tags string is empty, return a default "type:BOOLEAN" attribute.
	if tags == "" {
		return map[string][]string{TYPE: {"BOOLEAN"}}
	}

	// Remove all spaces from the tags string.
	tags = strings.ReplaceAll(tags, " ", "")

	// Split the tags string by semicolons to get individual attributes.
	attributes := strings.Split(tags, ";")

	// Map to store parsed attributes and their values.
	var vals = make(map[string][]string)

	// Iterate through each attribute string.
	for i := 0; i < len(attributes); i++ {
		// Split each attribute string into the attribute name (key) and values.
		pre := strings.SplitN(attributes[i], ":", 2)
		// Assign the attribute values (split by commas) to the corresponding key in the map.
		vals[pre[0]] = strings.Split(pre[1], ",")
	}

	// Return the parsed attributes as a map.
	return vals
}

// getTypes processes a slice of strings representing data types and formats them
// into a single string that combines these types for database schema definitions.
//
// Parameters:
//   - slice ([]string): A slice of strings where each string is a type attribute.
//
// Returns:
//   - string: A formatted string of type attributes, concatenated and ready for use
//     in database schema definitions (e.g., "PRIMARY KEY VARCHAR").
func getTypes(slice []string) (out string) {
	for i := 0; i < len(slice); i++ {
		var t string // Temporary variable to store the formatted type.
		switch slice[i] {
		case "primary": // Convert "primary" to "PRIMARY KEY".
			t = "PRIMARY KEY"
		default: // Convert other types to uppercase.
			t = strings.ToUpper(slice[i])
		}
		out += t + " " // Append the type to the output with a trailing space.
	}

	return // Return the concatenated string of types.
}

// isPrimary determines if a slice of strings contains the keyword "primary",
// which designates a database column as the primary key.
//
// Parameters:
//   - slice ([]string): A slice of strings containing column attributes.
//
// Returns:
//   - bool: Returns true if the slice includes the "primary" attribute, false otherwise.
func isPrimary(slice []string) bool {
	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case "primary": // Check if the current attribute is "primary".
			return true
		}
	}

	// Return false if no "primary" attribute is found.
	return false
}

// getReferences constructs a SQL REFERENCES clause for foreign key constraints.
//
// Parameters:
//   - item (string): The name of the referenced table.
//   - colName (string): The name of the current column in the table.
//
// Returns:
//   - string: A SQL string defining a REFERENCES clause for the foreign key,
//     formatted as "REFERENCES <table_name> (<column_name>) ".
func getReferences(item, colName string) string {
	// Convert the table name to snake_case.
	tName := toSnakeCase(item)

	// Split the column name to extract the referenced column (assumes column names are in "table_column" format).
	refColum := strings.SplitN(colName, "_", 2)

	// Return the formatted REFERENCES clause.
	return fmt.Sprintf("REFERENCES %s (%s) ", tName, refColum[1])
}

// getCascade generates SQL Cascade rules based on the provided slice of keywords.
//
// Parameters:
//   - slice ([]string): A slice of strings where each string represents a cascade rule,
//     such as "delete" or any other keyword for update cascades.
//
// Returns:
//   - string: A concatenated string of SQL Cascade rules derived from the input slice.
func getCascade(slice []string) (out string) {
	for i := 0; i < len(slice); i++ {
		switch slice[i] {
		case "delete": // If the cascade rule is "delete", append "ON DELETE CASCADE".
			out += "ON DELETE CASCADE "
		default: // For any other rule, append "ON UPDATE CASCADE".
			out += "ON UPDATE CASCADE "
		}
	}

	// Return the final concatenated string of cascade rules.
	return
}

// matchFirstCap is a regular expression used to match a transition from
// any character to an uppercase letter followed by lowercase letters.
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")

// matchAllCap is a regular expression used to match a transition between
// lowercase digits or letters and an uppercase letter.
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase converts a given string from camelCase or PascalCase to snake_case.
//
// Parameters:
//   - str (string): The input string in camelCase or PascalCase format.
//
// Returns:
//   - string: The converted string in snake_case format.
func toSnakeCase(str string) string {
	// Replace matches for transitions from any character to an uppercase
	// followed by lowercase with an underscore-separated format.
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")

	// Replace matches for transitions between lowercase digits or letters
	// and uppercase letters with an underscore-separated format.
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	// Convert the result to lowercase and return.
	return strings.ToLower(snake)
}
