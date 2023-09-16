package mysql_locks

import (
	"context"
	"database/sql"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
	tidb_storage "github.com/storage-lock/go-tidb-storage"
)

var dsnStorageLockFactoryBeanFactory *storage_lock_factory.StorageLockFactoryBeanFactory[string, *sql.DB] = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *sql.DB]()

func NewTidbLockByDsn(ctx context.Context, dsn string, lockId string) (*storage_lock.StorageLock, error) {
	factory, err := GetTidbLockFactoryByDsn(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return factory.CreateLock(lockId)
}

func NewTidbLockByDsnWithOptions(ctx context.Context, uri string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	factory, err := GetTidbLockFactoryByDsn(ctx, uri)
	if err != nil {
		return nil, err
	}
	return factory.CreateLockWithOptions(options)
}

func GetTidbLockFactoryByDsn(ctx context.Context, uri string) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return dsnStorageLockFactoryBeanFactory.GetOrInit(ctx, uri, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := tidb_storage.NewTidbConnectionManagerFromDsn(uri)
		options := tidb_storage.NewTidbStorageOptions().SetConnectionManager(connectionManager)
		storage, err := tidb_storage.NewTidbStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		factory := storage_lock_factory.NewStorageLockFactory(storage, options.ConnectionManager)
		return factory, nil
	})
}
