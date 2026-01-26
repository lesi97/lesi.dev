package cache

func IsRedisNil(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "redis: nil"
}
