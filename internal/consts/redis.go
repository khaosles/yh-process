package consts

import "time"

/*
   @File: redis.go
   @Author: khaosles
   @Time: 2023/8/6 10:01
   @Desc:
*/

const (
	ExpireTime  = time.Hour * 5
	RedisPrefix = "yh:"
)

const (
	RedisRelease  = RedisPrefix + "release:"
	RedisCallback = RedisPrefix + "callback:"

	RedisUpdate = RedisPrefix + "data-update"
)

const (
	RedisUri         = RedisPrefix + "%s:uri"          // {dag}:uri
	RedisDataType    = RedisPrefix + "%s:datatype"     // {dag}:datatype
	RedisNodata      = RedisPrefix + "%s:nodata"       // {dag}:nodata
	RedisSourceDir   = RedisPrefix + "%s:dir:source"   // {dag}:dir:source
	RedisTmpDir      = RedisPrefix + "%s:dir:tmp"      // {dag}:dir:tmp
	RedisTifDir      = RedisPrefix + "%s:dir:tif"      // {dag}:dir:tmp
	RedisBinDir      = RedisPrefix + "%s:dir:bin"      // {dag}:dir:bin
	RedisUpdateTime  = RedisPrefix + "%s:time:update"  // {dag}:time:update
	RedisReportTime  = RedisPrefix + "%s:time:report"  // {dag}:time:report
	RedisDataVersion = RedisPrefix + "%s:time:version" // {dag}:time:version
	RedisFieldBase   = RedisPrefix + "%s:field:base"   // {dag}:field:base
	RedisFieldUV     = RedisPrefix + "%s:field:uv"     // {dag}:field:uv
)
