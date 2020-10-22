package utils

import (
	"time"

	"github.com/juju/fslock"
	"github.com/spf13/viper"
)

// Possibility for the future https://medium.com/@gdm85/distributed-locking-for-pennies-distrilock-967347e7f2dd

var NoOp = func() error { return nil }

func Lock(resource string) (func() error, error) {
	viper.SetDefault("file.enableLocks", true)
	viper.SetDefault("file.lockTimeoutMs", 100)
	shouldLock := viper.GetBool("file.enableLocks")
	lockTimeoutMs := viper.GetInt("file.lockTimeoutMs")

	if !shouldLock {
		// No OP
		return NoOp, nil
	}

	lock := fslock.New(resource)
	err := lock.LockWithTimeout(time.Millisecond * time.Duration(lockTimeoutMs))
	if err != nil {
		return NoOp, err
	}
	return lock.Unlock, nil
}
