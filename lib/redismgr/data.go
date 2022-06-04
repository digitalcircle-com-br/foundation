package redismgr

// import (
// 	"context"
// 	"time"
// )

// func DataSet(k string, v interface{}, to int) error {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.Set(context, k, v, time.Duration(to)*time.Second).Err()
// }

// func DataGet(k string) (string, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.Get(context, k).Result()
// }

// func DataDel(k string) (int64, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.Del(context, k).Result()
// }

// func DataHSet(k string, v ...interface{}) (int64, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.HSet(context, k, v...).Result()
// }

// func DataHGet(k string, v string) (string, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.HGet(context, k, v).Result()
// }

// func DataHDel(k string, v ...string) (int64, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.HDel(context, k, v...).Result()
// }

// func DataHGetAll(k string) (map[string]string, error) {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.HGetAll(context, k).Result()
// }

// func Enqueue(q string, v ...interface{}) error {
// 	context, cancel := Ctx()
// 	defer cancel()
// 	return rediscli.LPush(context, q, v...).Err()
// }

// func Dequeue(q string, to int) ([]string, error) {
// 	context := context.Background()
// 	return rediscli.BRPop(context, time.Second*time.Duration(to), q).Result()
// }
