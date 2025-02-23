package connector

import "context"

var (
	_ interface {
		GetDatabaseName() string
	} = (*connWrapper)(nil)

	_ interface {
		Version(ctx context.Context) (version string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		IsTableExists(ctx context.Context, tableName string) (tableExists bool, err error)
	} = (*connWrapper)(nil)

	_ interface {
		IsColumnExists(ctx context.Context, tableName string, columnName string) (columnExists bool, err error)
	} = (*connWrapper)(nil)

	_ interface {
		IsPrimaryKey(ctx context.Context, tableName string, columnName string) (ok bool, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetColumns(ctx context.Context, tableName string) (columns []string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetColumnType(ctx context.Context, tableName string, columnName string) (dataType string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetPrimaryKeys(ctx context.Context, tableName string) (pkCols []string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetTables(ctx context.Context, folder string, recursive bool, excludeSysDirs bool) (tables []string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetIndexes(ctx context.Context, tableName string) (indexes []string, err error)
	} = (*connWrapper)(nil)

	_ interface {
		GetIndexColumns(ctx context.Context, tableName string, indexName string) (columns []string, err error)
	} = (*connWrapper)(nil)
)
