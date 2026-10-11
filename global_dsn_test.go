package tidb_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalDsnLockFactoryBeanFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, dsnStorageLockFactoryBeanFactory)
}

func TestNewTidbLockByDsn_InvalidDsn(t *testing.T) {
	_, err := NewTidbLockByDsn(context.Background(), "invalid-dsn", "test-lock")
	assert.Error(t, err)
}

func TestGetTidbLockFactoryByDsn_InvalidDsn(t *testing.T) {
	_, err := GetTidbLockFactoryByDsn(context.Background(), "invalid-dsn")
	assert.Error(t, err)
}

func TestNewTidbLockByDsnWithOptions_InvalidDsn(t *testing.T) {
	_, err := NewTidbLockByDsnWithOptions(context.Background(), "invalid-dsn", nil)
	assert.Error(t, err)
}

