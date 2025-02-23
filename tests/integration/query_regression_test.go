//go:build integration
// +build integration

package integration

import (
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/query/tx"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/version"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func TestUUIDSerializationQueryServiceIssue1501(t *testing.T) {
	// https://github.com/ydb-platform/ydb-go-sdk/issues/1501
	// test with special uuid - all bytes are different for check any byte swaps

	t.Run("old-send-with-force-wrapper", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		expectedResultWithBug := "2d9e498b-b746-9cfb-084d-de4e1cb4736e"
		id := uuid.MustParse(idString)
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS UUID;

SELECT CAST($val AS Utf8)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").UUIDWithIssue1501Value(id).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res string

		err = row.Scan(&res)
		require.NoError(t, err)
		require.Equal(t, expectedResultWithBug, res)
	})
	t.Run("old-receive-to-bytes", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS Text;

SELECT CAST($val AS UUID)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Text(idString).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res [16]byte

		err = row.Scan(&res)
		require.ErrorIs(t, err, types.ErrIssue1501BadUUID)
	})
	t.Run("old-receive-to-bytes-with-force-wrapper", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		expectedResultWithBug := "8b499e2d-46b7-fb9c-4d08-4ede6e73b41c"
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS Text;

SELECT CAST($val AS UUID)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Text(idString).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res types.UUIDBytesWithIssue1501Type

		err = row.Scan(&res)
		require.NoError(t, err)

		resUUID := uuid.UUID(res.AsBytesArray())
		require.Equal(t, expectedResultWithBug, resUUID.String())
	})
	t.Run("old-receive-to-string", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS Text;

SELECT CAST($val AS UUID)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Text(idString).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res string

		err = row.Scan(&res)
		require.ErrorIs(t, err, types.ErrIssue1501BadUUID)
	})
	t.Run("old-receive-to-string-with-force-wrapper", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		expectedResultWithBug := []byte{0x8b, 0x49, 0x9e, 0x2d, 0x46, 0xb7, 0xfb, 0x9c, 0x4d, 0x8, 0x4e, 0xde, 0x6e, 0x73, 0xb4, 0x1c}
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS Text;

SELECT CAST($val AS UUID)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Text(idString).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var resFromDB types.UUIDBytesWithIssue1501Type

		err = row.Scan(&resFromDB)
		require.NoError(t, err)

		res := resFromDB.AsBrokenString()
		require.Equal(t, expectedResultWithBug, []byte(res))
	})
	t.Run("old-send-receive-with-force-wrapper", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		id := uuid.MustParse(idString)
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS UUID;

SELECT $val`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").UUIDWithIssue1501Value(id).Build()),
			query.WithTxControl(tx.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var resWrapper types.UUIDBytesWithIssue1501Type
		err = row.Scan(&resWrapper)
		require.NoError(t, err)

		resUUID := uuid.UUID(resWrapper.AsBytesArray())

		require.Equal(t, id, resUUID)
	})
	t.Run("good-send", func(t *testing.T) {
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		id := uuid.MustParse(idString)
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS UUID;

SELECT CAST($val AS Utf8)`,
			query.WithIdempotent(),
			query.WithTxControl(query.SerializableReadWriteTxControl()),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Uuid(id).Build()),
		)

		require.NoError(t, err)

		var res string
		err = row.Scan(&res)
		require.NoError(t, err)
		require.Equal(t, id.String(), res)
	})
	t.Run("good-receive", func(t *testing.T) {
		// test old behavior - for test way of safe work with data, written with bagged API version
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS Utf8;

SELECT CAST($val AS UUID)`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Text(idString).Build()),
			query.WithTxControl(query.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res uuid.UUID

		err = row.Scan(&res)
		require.NoError(t, err)

		resString := strings.ToUpper(res.String())
		require.Equal(t, idString, resString)
	})
	t.Run("good-send-receive", func(t *testing.T) {
		var (
			scope = newScope(t)
			ctx   = scope.Ctx
			db    = scope.Driver()
		)

		idString := "6E73B41C-4EDE-4D08-9CFB-B7462D9E498B"
		id := uuid.MustParse(idString)
		row, err := db.Query().QueryRow(ctx, `
DECLARE $val AS UUID;

SELECT $val`,
			query.WithIdempotent(),
			query.WithParameters(ydb.ParamsBuilder().Param("$val").Uuid(id).Build()),
			query.WithTxControl(query.SerializableReadWriteTxControl()),
		)

		require.NoError(t, err)

		var res uuid.UUID
		err = row.Scan(&res)
		require.NoError(t, err)
		require.Equal(t, id.String(), res.String())
	})
}

// https://github.com/ydb-platform/ydb-go-sdk/issues/1506
func TestIssue1506TypedNullPushdown(t *testing.T) {
	if version.Lt(os.Getenv("YDB_VERSION"), "24.1") {
		t.Skip("query service not allowed in YDB version '" + os.Getenv("YDB_VERSION") + "'")
	}

	scope := newScope(t)
	ctx := scope.Ctx
	db := scope.Driver()

	val := int64(123)
	row, err := db.Query().QueryRow(ctx, `
DECLARE $arg1 AS Int64?;
DECLARE $arg2 AS Int64?;

SELECT CAST($arg1 AS Utf8) AS v1, CAST($arg2 AS Utf8) AS v2
`, query.WithParameters(
		ydb.ParamsBuilder().
			Param("$arg1").BeginOptional().Int64(nil).EndOptional().
			Param("$arg2").BeginOptional().Int64(&val).EndOptional().
			Build(),
	),
	)
	require.NoError(t, err)

	var res struct {
		V1 *string `sql:"v1"`
		V2 *string `sql:"v2"`
	}
	err = row.ScanStruct(&res)
	require.NoError(t, err)

	require.Nil(t, res.V1)
	require.Equal(t, "123", *res.V2)
}
