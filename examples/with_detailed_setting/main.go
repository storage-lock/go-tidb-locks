package main

import (
	"context"
	"fmt"
	tidb_locks "github.com/storage-lock/go-tidb-locks"
	storage_lock "github.com/storage-lock/go-storage-lock"
)

func main() {

	// 第一步：创建一把分布式锁
	tidbDsn := "root:@tcp(127.0.0.1:4000)/storage_lock_test"
	lockId := "test-lock"
	// 可以在这个地方设置锁相关的选项
	options := storage_lock.NewStorageLockOptions().SetLockId(lockId)
	lock, err := tidb_locks.NewTidbLockByDsnWithOptions(context.Background(), tidbDsn, options)
	if err != nil {
		panic(err)
	}

	// 第二步：尝试加锁
	ownerId := "owner-A"
	err = lock.Lock(context.Background(), ownerId)
	if err != nil {
		panic(err)
	}
	// 加锁成功之后要记得释放锁
	defer func() {
		err := lock.UnLock(context.Background(), ownerId)
		if err != nil {
			panic(err)
		}
	}()

	// 临界区，操作临界资源
	fmt.Println("Lock success")

}
