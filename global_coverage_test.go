package tidb_locks

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	storage_lock "github.com/storage-lock/go-storage-lock"
	tidb_storage "github.com/storage-lock/go-tidb-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewTidbLockBySqlDb(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	// Init 建表（TiDB 复用 MySQL 建表 SQL）
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	factory, err := GetTidbLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetTidbLockFactoryBySqlDb(ctx, mockDB)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewTidbLockBySqlDb(ctx, mockDB, "test-tidb-lock-sqldb")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-tidb-lock-sqldb-opt")
	lock2, err := NewTidbLockBySqlDbWithOptions(ctx, mockDB, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestNewTidbLockBySqlDb_Errors(t *testing.T) {
	old := newTidbStorage
	newTidbStorage = func(ctx context.Context, options *tidb_storage.TidbStorageOptions) (*tidb_storage.TidbStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newTidbStorage = old }()

	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	_, err = GetTidbLockFactoryBySqlDb(ctx, mockDB)
	assert.NotNil(t, err)
	_, err = NewTidbLockBySqlDb(ctx, mockDB, "x")
	assert.NotNil(t, err)
	_, err = NewTidbLockBySqlDbWithOptions(ctx, mockDB, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewTidbLockByDsn_Local(t *testing.T) {
	// TiDB 协议兼容 MySQL，用本机 sl-mysql 验证成功路径
	dsn := "root:rootpass@tcp(127.0.0.1:13307)/storage_lock_test"
	ctx := context.Background()

	probe, err := GetTidbLockFactoryByDsn(ctx, dsn)
	if err != nil {
		t.Skipf("本机 sl-mysql 不可用，跳过: %v", err)
	}
	assert.NotNil(t, probe)

	// 缓存命中
	factory2, err := GetTidbLockFactoryByDsn(ctx, dsn)
	assert.Nil(t, err)
	assert.Equal(t, probe, factory2)

	lock, err := NewTidbLockByDsn(ctx, dsn, "test-tidb-lock-dsn")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-tidb-lock-dsn-opt")
	lock2, err := NewTidbLockByDsnWithOptions(ctx, dsn, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}
