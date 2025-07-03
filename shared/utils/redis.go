package utils

import (
	"context"
	"fmt"
	"hupu/shared/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 封装了 Redis 客户端
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient 创建一个新的 RedisClient 实例
func NewRedisClient(ctx context.Context) (*RedisClient, error) {
	cfg := config.GlobalConfig.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// --- String Commands ---

// Set 设置键值对，expiration 为过期时间，0 表示不过期
func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(rc.ctx, key, value, expiration).Err()
}

// Get 获取键的值
func (rc *RedisClient) Get(key string) (string, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key 不存在时返回空字符串和 nil error
	}
	return val, err
}

// GetBytes 获取键的值 (byte slice)
func (rc *RedisClient) GetBytes(key string) ([]byte, error) {
	val, err := rc.client.Get(rc.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Key 不存在时返回 nil slice 和 nil error
	}
	return val, err
}

// Del 删除一个或多个键
func (rc *RedisClient) Del(keys ...string) (int64, error) {
	return rc.client.Del(rc.ctx, keys...).Result()
}

// Incr 原子性地将 key 中储存的数字值增一
func (rc *RedisClient) Incr(key string) (int64, error) {
	return rc.client.Incr(rc.ctx, key).Result()
}

// IncrBy 原子性地将 key 所储存的值加上指定的增量值
func (rc *RedisClient) IncrBy(key string, value int64) (int64, error) {
	return rc.client.IncrBy(rc.ctx, key, value).Result()
}

// Decr 原子性地将 key 中储存的数字值减一
func (rc *RedisClient) Decr(key string) (int64, error) {
	return rc.client.Decr(rc.ctx, key).Result()
}

// DecrBy 原子性地将 key 所储存的值减去指定的减量值
func (rc *RedisClient) DecrBy(key string, value int64) (int64, error) {
	return rc.client.DecrBy(rc.ctx, key, value).Result()
}

// SetNX (SET if Not eXists) 只在键不存在时，才对键进行设置操作
func (rc *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return rc.client.SetNX(rc.ctx, key, value, expiration).Result()
}

// --- Hash Commands ---

// HSet 将哈希表 key 中的字段 field 的值设为 value
func (rc *RedisClient) HSet(key, field string, value interface{}) error {
	return rc.client.HSet(rc.ctx, key, field, value).Err()
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中
func (rc *RedisClient) HMSet(key string, values map[string]interface{}) error {
	return rc.client.HMSet(rc.ctx, key, values).Err()
}

// HGet 获取存储在哈希表中指定字段的值
func (rc *RedisClient) HGet(key, field string) (string, error) {
	val, err := rc.client.HGet(rc.ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// HGetBytes 获取存储在哈希表中指定字段的值 (byte slice)
func (rc *RedisClient) HGetBytes(key, field string) ([]byte, error) {
	val, err := rc.client.HGet(rc.ctx, key, field).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// HMGet 获取所有给定字段的值
func (rc *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	return rc.client.HMGet(rc.ctx, key, fields...).Result()
}

// HGetAll 获取在哈希表中指定 key 的所有字段和值
func (rc *RedisClient) HGetAll(key string) (map[string]string, error) {
	return rc.client.HGetAll(rc.ctx, key).Result()
}

// HDel 删除一个或多个哈希表字段
func (rc *RedisClient) HDel(key string, fields ...string) (int64, error) {
	return rc.client.HDel(rc.ctx, key, fields...).Result()
}

// HExists 查看哈希表的指定字段是否存在
func (rc *RedisClient) HExists(key, field string) (bool, error) {
	return rc.client.HExists(rc.ctx, key, field).Result()
}

// HKeys 获取哈希表中的所有字段
func (rc *RedisClient) HKeys(key string) ([]string, error) {
	return rc.client.HKeys(rc.ctx, key).Result()
}

// HLen 获取哈希表中字段的数量
func (rc *RedisClient) HLen(key string) (int64, error) {
	return rc.client.HLen(rc.ctx, key).Result()
}

// --- List Commands ---

// LPush 将一个或多个值插入到列表头部
func (rc *RedisClient) LPush(key string, values ...interface{}) (int64, error) {
	return rc.client.LPush(rc.ctx, key, values...).Result()
}

// RPush 将一个或多个值插入到列表尾部
func (rc *RedisClient) RPush(key string, values ...interface{}) (int64, error) {
	return rc.client.RPush(rc.ctx, key, values...).Result()
}

// LPop 移出并获取列表的第一个元素
func (rc *RedisClient) LPop(key string) (string, error) {
	val, err := rc.client.LPop(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// RPop 移出并获取列表的最后一个元素
func (rc *RedisClient) RPop(key string) (string, error) {
	val, err := rc.client.RPop(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// LRange 获取列表指定范围内的元素
func (rc *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return rc.client.LRange(rc.ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (rc *RedisClient) LLen(key string) (int64, error) {
	return rc.client.LLen(rc.ctx, key).Result()
}

// LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除
func (rc *RedisClient) LTrim(key string, start, stop int64) error {
	return rc.client.LTrim(rc.ctx, key, start, stop).Err()
}

// --- Set Commands ---

// SAdd 向集合添加一个或多个成员
func (rc *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	return rc.client.SAdd(rc.ctx, key, members...).Result()
}

// SRem 移除集合中一个或多个成员
func (rc *RedisClient) SRem(key string, members ...interface{}) (int64, error) {
	return rc.client.SRem(rc.ctx, key, members...).Result()
}

// SMembers 返回集合中的所有成员
func (rc *RedisClient) SMembers(key string) ([]string, error) {
	return rc.client.SMembers(rc.ctx, key).Result()
}

// SIsMember 判断 member 元素是否是集合 key 的成员
func (rc *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	return rc.client.SIsMember(rc.ctx, key, member).Result()
}

// SCard 获取集合的成员数
func (rc *RedisClient) SCard(key string) (int64, error) {
	return rc.client.SCard(rc.ctx, key).Result()
}

// SPop 移除并返回集合中的一个随机元素
func (rc *RedisClient) SPop(key string) (string, error) {
	val, err := rc.client.SPop(rc.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// SRandMember 返回集合中一个或多个随机数
func (rc *RedisClient) SRandMemberN(key string, count int64) ([]string, error) {
	return rc.client.SRandMemberN(rc.ctx, key, count).Result()
}

// --- Sorted Set Commands ---

// ZAdd 向有序集合添加一个或多个成员，或者更新已存在成员的分数
func (rc *RedisClient) ZAdd(key string, members ...redis.Z) (int64, error) {
	return rc.client.ZAdd(rc.ctx, key, members...).Result()
}

// ZRem 移除有序集合中的一个或多个成员
func (rc *RedisClient) ZRem(key string, members ...interface{}) (int64, error) {
	return rc.client.ZRem(rc.ctx, key, members...).Result()
}

// ZRange 通过索引区间返回有序集合成指定区间内的成员 (按分数从小到大排序)
func (rc *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	return rc.client.ZRange(rc.ctx, key, start, stop).Result()
}

// ZRevRange 通过索引区间返回有序集合成指定区间内的成员 (按分数从大到小排序)
func (rc *RedisClient) ZRevRange(key string, start, stop int64) ([]string, error) {
	return rc.client.ZRevRange(rc.ctx, key, start, stop).Result()
}

// ZRangeWithScores 通过索引区间返回有序集合成指定区间内的成员和分数 (按分数从小到大排序)
func (rc *RedisClient) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return rc.client.ZRangeWithScores(rc.ctx, key, start, stop).Result()
}

// ZRevRangeWithScores 通过索引区间返回有序集合成指定区间内的成员和分数 (按分数从大到小排序)
func (rc *RedisClient) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return rc.client.ZRevRangeWithScores(rc.ctx, key, start, stop).Result()
}

// ZRangeByScore 通过分数区间返回有序集合的成员 (按分数从小到大排序)
func (rc *RedisClient) ZRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return rc.client.ZRangeByScore(rc.ctx, key, opt).Result()
}

// ZRevRangeByScore 通过分数区间返回有序集合的成员 (按分数从大到小排序)
func (rc *RedisClient) ZRevRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error) {
	return rc.client.ZRevRangeByScore(rc.ctx, key, opt).Result()
}

// ZCard 获取有序集合的成员数
func (rc *RedisClient) ZCard(key string) (int64, error) {
	return rc.client.ZCard(rc.ctx, key).Result()
}

// ZScore 返回有序集中，成员的分数值
func (rc *RedisClient) ZScore(key, member string) (float64, error) {
	return rc.client.ZScore(rc.ctx, key, member).Result()
}

// ZRank 返回有序集合中指定成员的排名 (按分数从小到大排序)
func (rc *RedisClient) ZRank(key, member string) (int64, error) {
	rank, err := rc.client.ZRank(rc.ctx, key, member).Result()
	if err == redis.Nil {
		return -1, nil // 成员不存在时返回 -1
	}
	return rank, err
}

// ZRevRank 返回有序集合中指定成员的排名 (按分数从大到小排序)
func (rc *RedisClient) ZRevRank(key, member string) (int64, error) {
	rank, err := rc.client.ZRevRank(rc.ctx, key, member).Result()
	if err == redis.Nil {
		return -1, nil // 成员不存在时返回 -1
	}
	return rank, err
}

// --- Key Commands ---

// Exists 检查给定 key 是否存在
func (rc *RedisClient) Exists(keys ...string) (int64, error) {
	return rc.client.Exists(rc.ctx, keys...).Result()
}

// Expire 为给定 key 设置过期时间
func (rc *RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	return rc.client.Expire(rc.ctx, key, expiration).Result()
}

// ExpireAt 设置 key 的过期时间戳 (Unix time)
func (rc *RedisClient) ExpireAt(key string, tm time.Time) (bool, error) {
	return rc.client.ExpireAt(rc.ctx, key, tm).Result()
}

// TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)
func (rc *RedisClient) TTL(key string) (time.Duration, error) {
	return rc.client.TTL(rc.ctx, key).Result()
}

// PTTL 以毫秒为单位返回 key 的剩余的过期时间
func (rc *RedisClient) PTTL(key string) (time.Duration, error) {
	return rc.client.PTTL(rc.ctx, key).Result()
}

// Persist 移除给定 key 的过期时间，使得 key 永不过期
func (rc *RedisClient) Persist(key string) (bool, error) {
	return rc.client.Persist(rc.ctx, key).Result()
}

// Keys 查找所有符合给定模式 pattern 的 key
func (rc *RedisClient) Keys(pattern string) ([]string, error) {
	return rc.client.Keys(rc.ctx, pattern).Result()
}

// --- Pub/Sub Commands ---

// Publish 将信息发送到指定的频道
func (rc *RedisClient) Publish(channel string, message interface{}) (int64, error) {
	return rc.client.Publish(rc.ctx, channel, message).Result()
}

// Subscribe 订阅给定的一个或多个频道的信息
// 返回一个 PubSub 对象，需要调用其 ReceiveMessage 或 Channel 方法来接收消息
func (rc *RedisClient) Subscribe(channels ...string) *redis.PubSub {
	return rc.client.Subscribe(rc.ctx, channels...)
}

// --- Transaction Commands ---

// TxPipeline 创建一个事务管道
// 你可以在回调函数中执行多个命令，这些命令会被原子性地执行
func (rc *RedisClient) TxPipeline(fn func(pipe redis.Pipeliner) error) ([]redis.Cmder, error) {
	return rc.client.TxPipelined(rc.ctx, fn)
}

// Pipeline 创建一个普通管道
// 命令会被打包发送，但不是原子性的
func (rc *RedisClient) Pipeline() redis.Pipeliner {
	return rc.client.Pipeline()
}

// --- Scripting Commands ---

// Eval 执行 Lua 脚本
func (rc *RedisClient) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return rc.client.Eval(rc.ctx, script, keys, args...).Result()
}

// EvalSha 执行缓存在服务器中的 Lua 脚本 (通过 SHA1 校验和)
func (rc *RedisClient) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return rc.client.EvalSha(rc.ctx, sha1, keys, args...).Result()
}

// ScriptLoad 将脚本添加到脚本缓存中，但并不立即执行它
func (rc *RedisClient) ScriptLoad(script string) (string, error) {
	return rc.client.ScriptLoad(rc.ctx, script).Result()
}

// ScriptExists 检查指定的脚本是否已经被保存在缓存当中
func (rc *RedisClient) ScriptExists(hashes ...string) ([]bool, error) {
	return rc.client.ScriptExists(rc.ctx, hashes...).Result()
}

// ScriptFlush 清空服务器的所有 Lua 脚本缓存
func (rc *RedisClient) ScriptFlush() error {
	return rc.client.ScriptFlush(rc.ctx).Err()
}

// --- Other Commands ---

// Ping 测试服务器是否连接正常
func (rc *RedisClient) Ping() (string, error) {
	return rc.client.Ping(rc.ctx).Result()
}

// UnderlyingClient 返回底层的 go-redis 客户端，以便进行更高级的操作
func (rc *RedisClient) UnderlyingClient() *redis.Client {
	return rc.client
}
