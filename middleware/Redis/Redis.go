package Redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"golandproject/common"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetDownloadKey(key string) (bool, string) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, ""
	}
	return true, val
}

func AddDownloadKey(value string) string {
	var key string
	for {
		key = common.RandStringRunes(64)
		if result, _ := rdb.Exists(ctx, key).Result(); result == 0 {
			break
		}
	}
	err := rdb.Set(ctx, key, value, 60*time.Second).Err()
	if err != nil {
		return ""
	}
	return key
}

func AddUploadKey() string {
	var key string
	for {
		key = common.RandStringRunes(64)
		if result, _ := rdb.Exists(ctx, key).Result(); result == 0 {
			break
		}
	}
	err := rdb.Set(ctx, key, "1", 12*60*60*time.Second).Err() //12小时失效
	if err != nil {
		return ""
	}
	return key
}

func DelUploadKey(key string) int64 {
	result, _ := rdb.Del(ctx, key).Result()
	return result
}

func GetUploadKey(key string) bool {
	if result, _ := rdb.Exists(ctx, key).Result(); result == 0 {
		return false
	}
	return true
}

func NewShare(fileid, deadtime int64, uid, password string) (key string) {
	for {
		key = common.RandStringRunes(16)
		if result, _ := rdb.Exists(ctx, key).Result(); result == 0 {
			break
		}
	}
	if password == "" {
		rdb.HMSet(ctx, key, "fileid", fileid, "uid", uid)
	} else {
		rdb.HMSet(ctx, key, "fileid", fileid, "password", password, "uid", uid)
	}
	if deadtime != 3 {
		if deadtime == 1 {
			rdb.Expire(ctx, key, 7*24*60*60*time.Second)
		} else {
			rdb.Expire(ctx, key, 30*24*60*60*time.Second)
		}
	}
	return key
}

func GetShare(key string) (isexist bool, fileid int64, password, uid string) {
	if result, _ := rdb.Exists(ctx, key).Result(); result == 0 {
		return false, 0, "", ""
	}
	rdb.HGet(ctx, key, "fileid").Scan(&fileid)
	rdb.HGet(ctx, key, "password").Scan(&password)
	rdb.HGet(ctx, key, "uid").Scan(&uid)
	return true, fileid, password, uid
}

func DeleteShare(shareid string) bool {
	num := rdb.HDel(ctx, shareid, "fileid", "password", "uid")
	return num.Val() != 0
}
