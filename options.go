package sqlquery

// Supported SQL database flavors for query compilation.
const (
	MySQLFlavor Flavor = iota + 1
	PostgreSQLFlavor
	SQLiteFlavor
)

// Flavor represents the SQL dialect used for query compilation.
// It controls the format of compiled SQL to match database-specific syntax.
type Flavor int

// FindOptions configures the behavior of FindQuery for building SELECT queries.
//
// Fields:
//   - Flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - Fields: Column names to select (defaults to "*" for all columns)
//   - Filters: WHERE conditions using field names and optional operators (see package docs for filter syntax)
//   - ForUpdate: Whether to add FOR UPDATE clause for row locking
//   - ForUpdateMode: Optional mode for FOR UPDATE (e.g., "NOWAIT", "SKIP LOCKED")
type FindOptions struct {
	Flavor        Flavor
	Fields        []string
	Filters       map[string]interface{}
	ForUpdate     bool
	ForUpdateMode string
}

// WithFields returns a new FindOptions with the specified field list.
// Use this to select specific columns instead of "*".
func (f *FindOptions) WithFields(fields []string) *FindOptions {
	copy := *f
	copy.Fields = fields
	return &copy
}

// WithFilter returns a new FindOptions with an additional filter condition.
// Supports special operators via dot notation (e.g., "age.gte", "status.in").
// See package documentation for complete filter syntax.
func (f *FindOptions) WithFilter(field string, value interface{}) *FindOptions {
	copy := *f
	copy.Filters[field] = value
	return &copy
}

// WithForUpdate returns a new FindOptions with FOR UPDATE clause enabled.
// The mode parameter can specify locking behavior (e.g., "NOWAIT", "SKIP LOCKED").
// Pass an empty string for default FOR UPDATE behavior.
func (f *FindOptions) WithForUpdate(mode string) *FindOptions {
	copy := *f
	copy.ForUpdate = true
	copy.ForUpdateMode = mode
	return &copy
}

// NewFindOptions creates a new FindOptions with default values.
// The fields list defaults to ["*"] (all columns), and an empty filters map is initialized.
func NewFindOptions(flavor Flavor) *FindOptions {
	return &FindOptions{
		Fields:  []string{"*"},
		Flavor:  flavor,
		Filters: make(map[string]interface{}),
	}
}

// FindAllOptions configures the behavior of FindAllQuery for building paginated SELECT queries.
//
// Fields:
//   - Flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - Fields: Column names to select (defaults to "*" for all columns)
//   - Filters: WHERE conditions using field names and optional operators (see package docs for filter syntax)
//   - Limit: Maximum number of rows to return
//   - Offset: Number of rows to skip before returning results
//   - OrderBy: ORDER BY clause (e.g., "created_at DESC", "name ASC, id DESC")
//   - ForUpdate: Whether to add FOR UPDATE clause for row locking
//   - ForUpdateMode: Optional mode for FOR UPDATE (e.g., "NOWAIT", "SKIP LOCKED")
type FindAllOptions struct {
	Flavor        Flavor
	Fields        []string
	Filters       map[string]interface{}
	Limit         int
	Offset        int
	OrderBy       string
	ForUpdate     bool
	ForUpdateMode string
}

// WithFields returns a new FindAllOptions with the specified field list.
// Use this to select specific columns instead of "*".
func (f *FindAllOptions) WithFields(fields []string) *FindAllOptions {
	copy := *f
	copy.Fields = fields
	return &copy
}

// WithFilter returns a new FindAllOptions with an additional filter condition.
// Supports special operators via dot notation (e.g., "age.gte", "status.in").
// See package documentation for complete filter syntax.
func (f *FindAllOptions) WithFilter(field string, value interface{}) *FindAllOptions {
	copy := *f
	copy.Filters[field] = value
	return &copy
}

// WithLimit returns a new FindAllOptions with the specified LIMIT value.
func (f *FindAllOptions) WithLimit(limit int) *FindAllOptions {
	copy := *f
	copy.Limit = limit
	return &copy
}

// WithOffset returns a new FindAllOptions with the specified OFFSET value.
func (f *FindAllOptions) WithOffset(offset int) *FindAllOptions {
	copy := *f
	copy.Offset = offset
	return &copy
}

// WithOrderBy returns a new FindAllOptions with the specified ORDER BY clause.
// The orderBy parameter should be a valid SQL ORDER BY expression (e.g., "created_at DESC").
func (f *FindAllOptions) WithOrderBy(orderBy string) *FindAllOptions {
	copy := *f
	copy.OrderBy = orderBy
	return &copy
}

// WithForUpdate returns a new FindAllOptions with FOR UPDATE clause enabled.
// The mode parameter can specify locking behavior (e.g., "NOWAIT", "SKIP LOCKED").
// Pass an empty string for default FOR UPDATE behavior.
func (f *FindAllOptions) WithForUpdate(mode string) *FindAllOptions {
	copy := *f
	copy.ForUpdate = true
	copy.ForUpdateMode = mode
	return &copy
}

// NewFindAllOptions creates a new FindAllOptions with default values.
// The fields list defaults to ["*"] (all columns), and an empty filters map is initialized.
func NewFindAllOptions(flavor Flavor) *FindAllOptions {
	return &FindAllOptions{
		Fields:  []string{"*"},
		Flavor:  flavor,
		Filters: make(map[string]interface{}),
	}
}

// UpdateOptions configures the behavior of UpdateWithOptionsQuery for building UPDATE queries.
//
// Fields:
//   - Flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - Assignments: Field-value pairs to set (e.g., {"status": "active", "updated_at": time.Now()})
//   - Filters: WHERE conditions using field names and optional operators (see package docs for filter syntax)
type UpdateOptions struct {
	Flavor      Flavor
	Assignments map[string]interface{}
	Filters     map[string]interface{}
}

// WithAssignment returns a new UpdateOptions with an additional field assignment.
// Use this to specify which fields to update and their new values.
func (u *UpdateOptions) WithAssignment(field string, value interface{}) *UpdateOptions {
	copy := *u
	copy.Assignments[field] = value
	return &copy
}

// WithFilter returns a new UpdateOptions with an additional filter condition.
// Supports special operators via dot notation (e.g., "age.gte", "status.in").
// See package documentation for complete filter syntax.
func (u *UpdateOptions) WithFilter(field string, value interface{}) *UpdateOptions {
	copy := *u
	copy.Filters[field] = value
	return &copy
}

// NewUpdateOptions creates a new UpdateOptions with default values.
// Empty maps are initialized for both assignments and filters.
func NewUpdateOptions(flavor Flavor) *UpdateOptions {
	return &UpdateOptions{
		Flavor:      flavor,
		Assignments: make(map[string]interface{}),
		Filters:     make(map[string]interface{}),
	}
}

// DeleteOptions configures the behavior of DeleteWithOptionsQuery for building DELETE queries.
//
// Fields:
//   - Flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - Filters: WHERE conditions using field names and optional operators (see package docs for filter syntax)
type DeleteOptions struct {
	Flavor  Flavor
	Filters map[string]interface{}
}

// WithFilter returns a new DeleteOptions with an additional filter condition.
// Supports special operators via dot notation (e.g., "age.gte", "status.in").
// See package documentation for complete filter syntax.
func (d *DeleteOptions) WithFilter(field string, value interface{}) *DeleteOptions {
	copy := *d
	copy.Filters[field] = value
	return &copy
}

// NewDeleteOptions creates a new DeleteOptions with default values.
// An empty filters map is initialized.
func NewDeleteOptions(flavor Flavor) *DeleteOptions {
	return &DeleteOptions{
		Flavor:  flavor,
		Filters: make(map[string]interface{}),
	}
}
