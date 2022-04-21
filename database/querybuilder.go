package database

type DatabaseDriver interface {
	// GetName get the database driver name
	GetName() string

	// Insert a model into the database
	Insert(model *Model) error

	// Get a model from the database
	Get(model *Model) QueryBuilder

	// Update a model in the database
	Update(model *Model) error

	// Delete a model from the database
	Delete(model *Model) error
}

type QueryBuilder interface {
	// Select certain fields from the table
	Select(fields ...string) QueryBuilder

	// From the table
	From(table string) QueryBuilder

	// Where filter on a field
	Where(field string, operator Operator, value interface{}) QueryBuilder

	// Limit the number of results
	Limit(limit int) QueryBuilder

	// Offset the results
	Offset(offset int) QueryBuilder

	// Order the results
	Order(field string, direction string) QueryBuilder

	// Group the results
	Group(field string) QueryBuilder

	// Having filter on a field
	Having(field string, operator Operator, value interface{}) QueryBuilder

	// With join with a relationship
	With(model *Model) QueryBuilder

	// WithCustom join with a custom relation name
	// This is useful for joining if there are 2 relations with the same model.
	WithCustom(model *Model, relationName string) QueryBuilder
}
