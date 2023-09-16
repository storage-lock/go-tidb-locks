package mysql_locks

import (
	"context"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewTidbLockByDsn(t *testing.T) {
	envName := "STORAGE_LOCK_TIDB_DSN"
	tidbDsn := os.Getenv(envName)
	assert.NotEmpty(t, tidbDsn)

	factory, err := GetTidbLockFactoryByDsn(context.Background(), tidbDsn)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 20
	storage_lock_test_helper.EveryOnePlayTimes = 100
	storage_lock_test_helper.TestStorageLock(t, factory)
}
