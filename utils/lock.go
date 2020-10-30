package utils

import (
	"time"

	"github.com/juju/fslock"
	"github.com/spf13/viper"
)

// Possibility for the future https://medium.com/@gdm85/distributed-locking-for-pennies-distrilock-967347e7f2dd

var noOp = func() error { return nil }

// Lock resource if viper returns it should be locked.
func Lock(resource string) (func() error, error) {
	shouldLock := viper.GetBool("file.enableLocks")
	lockTimeoutMs := viper.GetInt("file.lockTimeoutMs")

	if !shouldLock {
		// No OP
		return noOp, nil
	}

	lock := fslock.New(resource)
	err := lock.LockWithTimeout(time.Millisecond * time.Duration(lockTimeoutMs))
	if err != nil {
		return noOp, err
	}
	return lock.Unlock, nil
}
