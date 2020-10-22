# Configuration

## Lock Configuration

In order to configure locking in the CAS, you can use two keys: `file.enableLocks` and `file.lockTimeoutMs`.

```yaml
file:
  enableLocks: true
  lockTimeoutMs: 100
```

### file.enableLocks

Use this configuration to disable locks on the filesystem entirely. This is useful in scenarios where you have a distributed file system or for any other reason locking on a file is not working.

hyper-cas should work properly anyway, since most operations tend to be idempotent and can be easily retried.

**Values**: `true`, `false`

### file.lockTimeoutMS

Use this configuration to configure how lock hyper-cas should wait for a file lock before returning an error.

**Values**: `the number of milliseconds to wait for a lock`
