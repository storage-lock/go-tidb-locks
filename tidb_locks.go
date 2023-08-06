package locks

import (
	"context"
	mysql_storage "github.com/storage-lock/go-mysql-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	tidb_storage "github.com/storage-lock/go-tidb-storage"
)

// NewTidbStorageLock 高层API，使用默认配置快速创建基于TiDB的分布式锁
// lockId: 要锁住的资源的唯一ID
// dsn: 用来存储锁的TiDB数据库连接方式
func NewTidbStorageLock(ctx context.Context, lockId string, dsn string) (*storage_lock.StorageLock, error) {
	connectionProvider := tidb_storage.NewTidbConnectionProviderFromDSN(dsn)
	storageOptions := &tidb_storage.TidbStorageOptions{
		MySQLStorageOptions: &mysql_storage.MySQLStorageOptions{
			ConnectionManager: connectionProvider,
			TableName:         storage.DefaultStorageTableName,
		},
	}

	s, err := tidb_storage.NewTidbStorage(ctx, storageOptions)
	if err != nil {
		return nil, err
	}

	lockOptions := &storage_lock.StorageLockOptions{
		LockId:               lockId,
		LeaseExpireAfter:     storage_lock.DefaultLeaseExpireAfter,
		LeaseRefreshInterval: storage_lock.DefaultLeaseRefreshInterval,
		//VersionMissRetryTimes: storage_lock.DefaultVersionMissRetryTimes,
	}
	return storage_lock.NewStorageLockWithOptions(s, lockOptions)
}