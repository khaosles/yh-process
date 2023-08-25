package redisutil

import (
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/khaosles/go-contrib/redis"
	"yh-process/internal/model"
	"yh-process/internal/vo"
)

/*
   @File: redis.go
   @Author: khaosles
   @Time: 2023/8/21 16:32
   @Desc:
*/

func Set(key string, val string, expire time.Duration) error {
	return redis.SetExpire(key, val, expire)
}

func GetOfString(key string) (string, error) {
	return redis.Get(key)
}

func GetOfTime(key string) (time.Time, error) {
	str, err := redis.Get(key)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.DateTime, str)
}

func GetOfI64(key string) (int64, error) {
	str, err := redis.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 0, 0)
}

func GetOfF32(key string) (float32, error) {
	str, err := redis.Get(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseFloat(str, 0)
	if err != nil {
		return 0, err
	}
	return float32(i), err
}

func GetOfURLInfo(key string) ([]*vo.URLInfo, error) {
	str, err := redis.Get(key)
	if err != nil {
		return nil, err
	}
	var uris []*vo.URLInfo
	err = sonic.UnmarshalString(str, &uris)
	if err != nil {
		return nil, err
	}
	return uris, nil
}

func SetOfURLInfo(key string, val []*vo.URLInfo, expire time.Duration) error {
	bs, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return redis.SetExpire(key, string(bs), expire)
}

func GetOfElements(key string) ([]*model.SysElementInfo, error) {
	str, err := redis.Get(key)
	if err != nil {
		return nil, err
	}
	var elements []*model.SysElementInfo
	err = sonic.UnmarshalString(str, &elements)
	if err != nil {
		return nil, err
	}
	return elements, nil
}

func SetOfElements(key string, val []*model.SysElementInfo, expire time.Duration) error {
	bs, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return redis.SetExpire(key, string(bs), expire)
}

func GetOfUV(key string) ([]*vo.UV, error) {
	str, err := redis.Get(key)
	if err != nil {
		return nil, err
	}
	var uvs []*vo.UV
	err = sonic.UnmarshalString(str, &uvs)
	if err != nil {
		return nil, err
	}
	return uvs, nil
}

func SetOfUV(key string, val []*vo.UV, expire time.Duration) error {
	bs, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return redis.SetExpire(key, string(bs), expire)
}
