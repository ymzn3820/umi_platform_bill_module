package gredis

import (
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"strconv"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func ChooseDB(dbNum int, callback func(conn redis.Conn) error) error {
	conn := RedisConn.Get()
	defer conn.Close()

	// Select the specified database.
	if _, err := conn.Do("SELECT", dbNum); err != nil {
		return err
	}

	// Execute the callback with the selected database.
	err := callback(conn)
	if err != nil {
		return err
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

type HashRateEntry struct {
	HashRate            int64 `json:"hash_rate"`
	ExpirationTimestamp int64 `json:"expiration_timestamp"`
}

// UserHashRates 结构体用于分类统计通用和定向算力值及其有效期
type UserHashRates struct {
	UniversalHashRates []HashRateEntry
	DirectedHashRates  []HashRateEntry
}

// AddHashRateUniversal 会员通用算力部分写入
func AddHashRateUniversal(userId string, hashrate int64, duration time.Duration) error {

	conn := RedisConn.Get()
	defer conn.Close()

	if hashrate <= 0 {
		return nil
	}
	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	expirationTime := time.Now().Add(duration)
	expirationTimestamp := expirationTime.Unix()

	key := fmt.Sprintf("user:%s:hashrate:universal", userId)

	// Create an instance of the HashRateEntry struct with the hashrate and timestamp.
	entry := HashRateEntry{
		HashRate:            hashrate,
		ExpirationTimestamp: expirationTimestamp,
	}

	// Marshal the HashRateEntry struct into a JSON string.
	memberJSON, err := json.Marshal(entry)
	if err != nil {
		logging.Error("failed to marshal hashrate universal entry to JSON: %v", err)
		return fmt.Errorf("failed to marshal hashrate universal entry to JSON: %v", err)

	}

	_, err = conn.Do("ZADD", key, "NX", expirationTimestamp, memberJSON)
	if err != nil {
		logging.Error("failed to add hashrate directed to sorted set: %v", err)
		return fmt.Errorf("failed to add hashrate directed to sorted set: %v", err)

	}
	return nil
}

// AddHashRatePackage 会员通用算力部分写入
func AddHashRatePackage(userId string, hashrate int64, duration time.Duration) error {

	conn := RedisConn.Get()
	defer conn.Close()

	if hashrate <= 0 {
		return nil
	}

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	expirationTime := time.Now().Add(duration)
	expirationTimestamp := expirationTime.Unix()

	key := fmt.Sprintf("user:%s:hashrate:package", userId)

	// Create an instance of the HashRateEntry struct with the hashrate and timestamp.
	entry := HashRateEntry{
		HashRate:            hashrate,
		ExpirationTimestamp: expirationTimestamp,
	}

	// Marshal the HashRateEntry struct into a JSON string.
	memberJSON, err := json.Marshal(entry)
	if err != nil {
		logging.Error("failed to marshal hashrate package entry to JSON: %v", err)
		return fmt.Errorf("failed to marshal hashrate package entry to JSON: %v", err)
	}

	_, err = conn.Do("ZADD", key, "NX", expirationTimestamp, memberJSON)
	if err != nil {
		logging.Error("failed to add hashrate package to sorted set: %v", err)
		return fmt.Errorf("failed to add hashrate package to sorted set: %v", err)
	}
	return nil
}

// AddHashRateDirected 定向算力写入
func AddHashRateDirected(userId string, hashrate int64, duration time.Duration) error {

	if hashrate <= 0 {
		return nil
	}

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	expirationTime := time.Now().Add(duration)
	expirationTimestamp := expirationTime.Unix()

	key := fmt.Sprintf("user:%s:hashrate:directed", userId)

	// Create an instance of the HashRateEntry struct with the hashrate and timestamp.
	entry := HashRateEntry{
		HashRate:            hashrate,
		ExpirationTimestamp: expirationTimestamp,
	}

	// Marshal the HashRateEntry struct into a JSON string.
	memberJSON, err := json.Marshal(entry)
	if err != nil {
		logging.Error("failed to marshal hashrate directed entry to JSON: %v", err)
		return fmt.Errorf("failed to marshal hashrate directed entry to JSON: %v", err)
	}

	_, err = conn.Do("ZADD", key, "NX", expirationTimestamp, memberJSON)
	if err != nil {
		logging.Error("failed to add hasrate directed to sorted set: %v", err)
		return fmt.Errorf("failed to add hasrate directed to sorted set: %v", err)

	}
	return nil
}

// CheckHashRateValidityDirected 检查通用定向是否在有效期内
func CheckHashRateValidityDirected(userId string) (int, error) {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, err
	}

	key := fmt.Sprintf("user:%s:hashrate:directed", userId)
	currentTime := time.Now().Unix()

	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))

	if err != nil {
		logging.Error("failed to get hashrate directed from sorted set: %v", err)
		return e.ERROR_CHECK_HASHRATE, err
	}

	for _, m := range members {
		var entry HashRateEntry

		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			logging.Error("failed to unmarshal hashrate directed entry from JSON:%v", err)
			continue
		}
		// 可以减为负数，只判断当前余额是否大于0
		if entry.HashRate >= 0 && entry.ExpirationTimestamp > currentTime {
			return e.SUCCESS, nil
		}

	}
	return e.SUCCESS, nil
}

// CheckHashRateValidityUniversal 检查通用算力是否在有效期内
func CheckHashRateValidityUniversal(userId string) (int, error) {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, err
	}
	key := fmt.Sprintf("user:%s:hashrate:universal", userId)
	currentTime := time.Now().Unix()

	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))

	if err != nil {
		logging.Error("failed to get hashrate universal from sorted set: %v", err)
		return e.ERROR_CHECK_HASHRATE, err
	}

	for _, m := range members {
		var entry HashRateEntry

		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			logging.Error("failed to unmarshal hashrate universal entry from JSON:%v", err)
			continue
		}

		// 可以减为负数，只判断当前余额是否大于0
		if entry.HashRate >= 0 && entry.ExpirationTimestamp > currentTime {
			return e.SUCCESS, nil
		}

	}
	return e.SUCCESS, nil
}

// CheckHashRateValidityPackage 检查通用算力是否在有效期内
func CheckHashRateValidityPackage(userId string) (int, error) {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, err
	}
	key := fmt.Sprintf("user:%s:hashrate:package", userId)
	currentTime := time.Now().Unix()

	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))

	if err != nil {
		logging.Error("failed to get hashrate package from sorted set: %v", err)
		return e.ERROR_CHECK_HASHRATE, err
	}

	for _, m := range members {
		var entry HashRateEntry

		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			logging.Error("failed to unmarshal hashrate package entry from JSON:%v", err)
			continue
		}

		// 可以减为负数，只判断当前余额是否大于0
		if entry.HashRate >= 0 && entry.ExpirationTimestamp > currentTime {
			return e.SUCCESS, nil
		}

	}
	return e.SUCCESS, nil
}

// ConsumeHashRateDirected 消耗用户的定向算力
func ConsumeHashRateDirected(userId string, hashrateToConsume int64) (int, bool) {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, true
	}
	responseCode, checkErr := CheckHashRateValidityDirected(userId)
	if checkErr != nil {
		logging.Error("failed to check if hasrate directed is valid: %v", checkErr)
		return e.ERROR_CHECK_HASHRATE, true
	}

	if responseCode != 20000 {
		logging.Info("no valid hasrate directed available for consumption")
		return responseCode, true
	}
	key := fmt.Sprintf("user:%s:hashrate:directed", userId)
	currentTime := time.Now().Unix()

	// WATCH key to start the optimistic locking
	_, err := conn.Do("WATCH", key)
	if err != nil {
		logging.Error("failed to WATCH hashrate directed key: %v", err)
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Retrieving all the members that are valid
	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))
	if err != nil {
		logging.Error("failed to retrieve hashrate directed from sorted set: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Start a transaction with MULTI
	_, err = conn.Do("MULTI")
	if err != nil {
		logging.Error("failed to start MULTI transaction: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_START_MULTI, true
	}

	isUpdated := false

	// Iterate over the members
	for _, m := range members {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			conn.Do("UNWATCH") // Unwatch the key if there is an error
			logging.Error("failed to unmarshal hashrate universal entry from JSON: %v", err)
			continue
		}

		// Check rest hashrate is usable
		if entry.HashRate >= 0 {
			// Reduce the hashrate
			entry.HashRate -= hashrateToConsume

			// Serialize the updated hashrate data
			memberJSON, err := json.Marshal(entry)
			if err != nil {
				conn.Do("UNWATCH") // Unwatch the key if there is an error
				logging.Error("failed to marshal updated hashrate directed entry to JSON: %v", err)
				return e.ERROR_CONSUME_HASHRATE, true
			}

			// Remove the old member and add the updated one within the transaction
			conn.Send("ZREM", key, m.([]byte))
			conn.Send("ZADD", key, entry.ExpirationTimestamp, memberJSON)

			isUpdated = true
			break // Break after the first successful update
		}
	}

	// If no hashrate was updated, we can safely UNWATCH and abort the transaction
	if !isUpdated {
		conn.Do("UNWATCH")
		return e.INSUFFICIENT_HASHRATE, false
	}

	// Execute the transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		logging.Error("failed to EXEC transaction: %v", err)
		return e.ERROR_CONSUME_HASHRATE, true
	}
	fmt.Println(isUpdated)
	fmt.Println("isUpdatedisUpdated")
	// Consumption was successful
	return e.SUCCESS, false
}

// ConsumeHashRateUniversal 消耗用户的通用算力
func ConsumeHashRateUniversal(userId string, hashrateToConsume int64) (int, bool) {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, true
	}

	responseCode, checkErr := CheckHashRateValidityUniversal(userId)
	if checkErr != nil {
		logging.Error("failed to check if hasrate directed is valid: %v", checkErr)
		return e.ERROR_CHECK_HASHRATE, true
	}
	if responseCode != 20000 {
		logging.Info("no valid hasrate directed available for consumption")
		return responseCode, true
	}

	key := fmt.Sprintf("user:%s:hashrate:universal", userId)
	currentTime := time.Now().Unix()

	// WATCH key to start the optimistic locking
	_, err := conn.Do("WATCH", key)
	if err != nil {
		logging.Error("failed to WATCH hashrate universal key: %v", err)
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Retrieving all the members that are valid
	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))
	if err != nil {
		logging.Error("failed to retrieve hashrate universal from sorted set: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Start a transaction with MULTI
	_, err = conn.Do("MULTI")
	if err != nil {
		logging.Error("failed to start MULTI transaction: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_START_MULTI, true
	}

	isUpdated := false

	// Iterate over the members
	for _, m := range members {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			conn.Do("UNWATCH") // Unwatch the key if there is an error
			logging.Error("failed to unmarshal hashrate universal entry from JSON: %v", err)
			continue
		}

		// Check if the hashrate is enough for consumption
		if entry.HashRate > 0 {
			// Reduce the hashrate
			entry.HashRate -= hashrateToConsume

			// Serialize the updated hashrate data
			memberJSON, err := json.Marshal(entry)
			if err != nil {
				conn.Do("UNWATCH") // Unwatch the key if there is an error
				logging.Error("failed to marshal updated hashrate universal entry to JSON: %v", err)
				return e.ERROR_START_MULTI, true
			}

			// Remove the old member and add the updated one within the transaction
			conn.Send("ZREM", key, m.([]byte))
			conn.Send("ZADD", key, entry.ExpirationTimestamp, memberJSON)

			isUpdated = true
			break // Break after the first successful update
		}
	}

	// If no hashrate was updated, we can safely UNWATCH and abort the transaction
	if !isUpdated {
		conn.Do("UNWATCH")
		return e.INSUFFICIENT_HASHRATE, false
	}

	// Execute the transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		logging.Error("failed to EXEC transaction: %v", err)
		return e.ERROR_CONSUME_HASHRATE, true
	}

	// Consumption was successful
	return e.SUCCESS, false
}

// ConsumeHashRatePackage 消耗用户的通用算力
func ConsumeHashRatePackage(userId string, hashrateToConsume int64) (int, bool) {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return e.ERROR_SELECT_DB, true
	}

	responseCode, checkErr := CheckHashRateValidityPackage(userId)
	if checkErr != nil {
		logging.Error("failed to check if hasrate directed is valid: %v", checkErr)
		return e.ERROR_CHECK_HASHRATE, true
	}

	if responseCode != 20000 {
		logging.Info("no valid hasrate package available for consumption")
		return responseCode, true
	}

	key := fmt.Sprintf("user:%s:hashrate:package", userId)
	currentTime := time.Now().Unix()

	// WATCH key to start the optimistic locking
	_, err := conn.Do("WATCH", key)
	if err != nil {
		logging.Error("failed to WATCH hashrate package key: %v", err)
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Retrieving all the members that are valid
	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))
	if err != nil {
		logging.Error("failed to retrieve hashrate package from sorted set: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_CHECK_HASHRATE, true
	}

	// Start a transaction with MULTI
	_, err = conn.Do("MULTI")
	if err != nil {
		logging.Error("failed to start MULTI transaction: %v", err)
		conn.Do("UNWATCH") // Unwatch the key if there is an error
		return e.ERROR_START_MULTI, true
	}

	isUpdated := false

	// Iterate over the members
	for _, m := range members {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			conn.Do("UNWATCH") // Unwatch the key if there is an error
			logging.Error("failed to unmarshal hashrate package entry from JSON: %v", err)
			continue
		}

		// Check if the hashrate is enough for consumption

		if entry.HashRate > 0 {
			// Reduce the hashrate
			entry.HashRate -= hashrateToConsume

			// Serialize the updated hashrate data
			memberJSON, err := json.Marshal(entry)
			if err != nil {
				conn.Do("UNWATCH") // Unwatch the key if there is an error
				logging.Error("failed to marshal updated hashrate package entry to JSON: %v", err)
				return e.ERROR_START_MULTI, true
			}

			// Remove the old member and add the updated one within the transaction
			conn.Send("ZREM", key, m.([]byte))
			conn.Send("ZADD", key, entry.ExpirationTimestamp, memberJSON)

			isUpdated = true
			break // Break after the first successful update
		}
	}

	// If no hashrate was updated, we can safely UNWATCH and abort the transaction
	if !isUpdated {
		conn.Do("UNWATCH")
		return e.INSUFFICIENT_HASHRATE, false
	}

	// Execute the transaction
	_, err = conn.Do("EXEC")
	if err != nil {
		logging.Error("failed to EXEC transaction: %v", err)
		return e.ERROR_CONSUME_HASHRATE, true
	}

	// Consumption was successful
	return e.SUCCESS, false
}

type HashRateEntryConsolidate struct {
	HashRate int64 `json:"hash_rate"`
}

// ConsolidateHashRate 处理算力余额低于阈值，进行相关处理
func ConsolidateHashRate(userId string, cate int) error {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	var hashRateKey string

	switch cate {
	case 1:
		hashRateKey = fmt.Sprintf("user:%s:hashrate:universal", userId)
	case 2:
		hashRateKey = fmt.Sprintf("user:%s:hashrate:directed", userId)
	case 3:
		hashRateKey = fmt.Sprintf("user:%s:hashrate:package", userId)
	}

	scoresAndMembers, err := redis.Strings(conn.Do("ZRANGE", hashRateKey, 0, -1, "WITHSCORES"))
	if err != nil {
		return err
	}

	if len(scoresAndMembers)/2 == 1 {
		return nil
	}

	var maxScore int64 = -1
	var maxScoreMember string
	var totalHashRateToAdd int64 = 0

	for i := 1; i < len(scoresAndMembers); i += 2 {
		score, err := strconv.ParseInt(scoresAndMembers[i], 10, 64)
		if err != nil {
			return err
		}
		if score > maxScore {
			maxScore = score
			maxScoreMember = scoresAndMembers[i-1]
		}
	}

	for i := 0; i < len(scoresAndMembers); i += 2 {
		member := scoresAndMembers[i]

		var entry HashRateEntry
		err = json.Unmarshal([]byte(member), &entry)
		if err != nil {
			return err
		}

		if entry.HashRate < 10 {
			totalHashRateToAdd += entry.HashRate

			_, err = conn.Do("ZREM", hashRateKey, member)
			if err != nil {
				return err
			}
		}
	}

	if totalHashRateToAdd > 0 {
		var maxMemberEntry HashRateEntry
		err = json.Unmarshal([]byte(maxScoreMember), &maxMemberEntry)
		if err != nil {
			return err
		}

		maxMemberEntry.HashRate += totalHashRateToAdd

		updatedMaxScoreMemberData, err := json.Marshal(maxMemberEntry)
		if err != nil {
			return err
		}

		// 删除当前的最大 score 成员
		_, err = conn.Do("ZREM", hashRateKey, maxScoreMember)
		if err != nil {
			return err
		}

		// 重新添加更新后的成员数据
		_, err = conn.Do("ZADD", hashRateKey, maxScore, string(updatedMaxScoreMemberData))
		if err != nil {
			return err
		}
	}

	return nil
}

// ClearExpiredHashRate 清除过期算力值
func ClearExpiredHashRate(userId string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	directedKey := fmt.Sprintf("user:%s:hashrate:directed", userId)
	universalKey := fmt.Sprintf("user:%s:hashrate:universal", userId)
	packageKey := fmt.Sprintf("user:%s:hashrate:package", userId)
	currentTime := time.Now().Unix()

	// 删除有序集合中过期的算力值
	if _, err := conn.Do("ZREMRANGEBYSCORE", directedKey, "-inf", currentTime); err != nil {
		return fmt.Errorf("error removing expired directed hashrate: %v", err)
	}

	if _, err := conn.Do("ZREMRANGEBYSCORE", universalKey, "-inf", currentTime); err != nil {
		return fmt.Errorf("error removing expired universal hashrate: %v", err)
	}

	if _, err := conn.Do("ZREMRANGEBYSCORE", packageKey, "-inf", currentTime); err != nil {
		return fmt.Errorf("error removing expired package hashrate: %v", err)
	}

	return nil
}

// HasRemainingHashRate 检查是否还有剩余算力值
func HasRemainingHashRate(userId string, prodCate int64) (bool, int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return false, 0, err
	}

	var key string

	switch prodCate {
	case 6:
		key = fmt.Sprintf("user:%s:hashrate:package", userId)
	case 3:
		key = fmt.Sprintf("user:%s:hashrate:universal", userId)

	}
	currentTime := time.Now().Unix()

	members, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, currentTime, "+inf"))

	if err != nil {
		return false, 0, fmt.Errorf("error retrieving remaining hashrate: %v", err)
	}

	var totalRemaining int64

	for _, m := range members {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			logging.Error("error when ummarshal members")
			continue
		}

		totalRemaining += entry.HashRate
	}
	fmt.Println(totalRemaining)
	return totalRemaining > 0, totalRemaining, nil
}

// 判断当前用户是否还有可使用的算力， 判断最大socre和当前时间的关系
func CheckIsActive(userId string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	// Select the appropriate database
	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return false, err
	}

	// Current timestamp
	currentTime := time.Now().Unix()

	// Define the keys
	hashRateKeys := []string{
		fmt.Sprintf("user:%s:hashrate:universal", userId),
		fmt.Sprintf("user:%s:hashrate:directed", userId),
		fmt.Sprintf("user:%s:hashrate:package", userId),
	}

	// Check the max score for each key
	for _, key := range hashRateKeys {
		// Retrieve the maximum score for the key
		values, err := redis.Values(conn.Do("ZREVRANGE", key, 0, 0, "WITHSCORES"))
		if err != nil {
			return false, fmt.Errorf("error retrieving values for key %s: %v", key, err)
		}

		// If no scores were returned, continue to the next key
		if len(values) == 0 {
			continue
		}

		// Extract the member and score from the values
		var member string
		var score int64
		if _, err := redis.Scan(values, &member, &score); err != nil {
			return false, fmt.Errorf("error scanning values for key %s: %v", key, err)
		}

		// If the max score is greater than the current time, the user is active
		if score > currentTime {
			return true, nil
		}
	}

	// If none of the max scores are greater than the current time, the user is not active
	return false, nil
}

// RenewHashRate 续费
func RenewHashRate(userId string, hashrateDirected int64, hashrateUniversal int64, duration time.Duration) error {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return err
	}

	// Check hashrate less than threshold
	mapCate := []int{1, 2, 3}
	for _, value := range mapCate {
		err := ConsolidateHashRate(userId, value)
		if err != nil {
			return err
		}
	}

	if err := ClearExpiredHashRate(userId); err != nil {

		return fmt.Errorf("error clearing expired hashrate: %v", err)
	}

	// TODO 增加判断当前的产品是流量包还是会员
	isActive, err := CheckIsActive(userId)

	if err != nil {
		return fmt.Errorf("error user %s CheckIsActive: %v", userId, err)
	}

	if isActive {
		if hashrateUniversal > 0 && hashrateDirected > 0 {
			if err := AddHashRateDirected(userId, hashrateDirected, duration); err != nil {
				return fmt.Errorf("error adding directed hashrate: %v", err)
			}
		}

		if hashrateUniversal > 0 && hashrateDirected > 0 {
			if err := AddHashRateUniversal(userId, hashrateUniversal, duration); err != nil {
				return fmt.Errorf("error adding universal hashrate: %v", err)
			}
		}

		if hashrateDirected == 0 && hashrateUniversal > 0 {
			if err := AddHashRatePackage(userId, hashrateUniversal, duration); err != nil {
				return fmt.Errorf("error adding package hashrate: %v", err)
			}
		}
	} else {
		//hasRemaining, remainingHashrate, err := HasRemainingHashRate(userId, prodCate)
		//fmt.Println(hasRemaining)
		//fmt.Println(remainingHashrate)
		//fmt.Println("remainingHashrateremainingHashrate")
		//
		//if err != nil {
		//	return fmt.Errorf("error checking remaining hashrate: %v", err)
		//}

		if hashrateDirected > 0 {
			if err := AddHashRateDirected(userId, hashrateDirected, duration); err != nil {
				return fmt.Errorf("error adding directed hashrate: %v", err)
			}
		}

		if hashrateUniversal > 0 && hashrateDirected > 0 {
			if err := AddHashRateUniversal(userId, hashrateUniversal, duration); err != nil {
				return fmt.Errorf("error adding universal hashrate: %v", err)
			}
		}
		if hashrateDirected == 0 && hashrateUniversal > 0 {
			if err := AddHashRatePackage(userId, hashrateUniversal, duration); err != nil {
				return fmt.Errorf("error adding package hashrate: %v", err)
			}
		}
	}

	return nil
}

type HashRateRulesFields struct {
	Model         string  `json:"model"`
	Unit          string  `json:"unit"`
	ConsumePoints float64 `json:"consume_points"`
}

func HashRateRules() ([]HashRateRulesFields, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	// 选择正确的数据库
	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return nil, err
	}

	// 获取存储规则的 JSON 字符串
	jsonString, err := redis.String(conn.Do("HGET", "hashrateRules", "pricing"))
	if err != nil {
		// 键可能不存在或者发生了其他错误
		return nil, err
	}

	// 初始化 HashRateRulesFields 切片用于存放解码后的数据
	var rules []HashRateRulesFields

	// 解码 JSON 到切片
	err = json.Unmarshal([]byte(jsonString), &rules)
	if err != nil {
		// 处理 JSON 解码错误
		return nil, err
	}

	// 返回解码后的规则和 nil 错误
	return rules, nil
}

// SummarizeUserHashRates 统计用户下所有可用的算力值
func SummarizeUserHashRates(userId string) (map[string]int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return nil, err
	}

	universalKey := fmt.Sprintf("user:%s:hashrate:universal", userId)
	directedKey := fmt.Sprintf("user:%s:hashrate:directed", userId)
	packageKey := fmt.Sprintf("user:%s:hashrate:package", userId)
	currentTime := time.Now().Unix()

	// Initialize sums for universal and directed hash rates
	var universalHashRateSum, directedHashRateSum, packageHashRateSum int64

	// Sum universal hash rates
	universalMembers, err := redis.Values(conn.Do("ZRANGEBYSCORE", universalKey, currentTime, "+inf"))
	if err != nil {
		return nil, fmt.Errorf("error retrieving universal hashrates: %v", err)
	}
	for _, m := range universalMembers {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			continue // Error handling: log the error or return it
		}
		universalHashRateSum += entry.HashRate
	}
	// Sum directed hash rates
	directedMembers, err := redis.Values(conn.Do("ZRANGEBYSCORE", directedKey, currentTime, "+inf"))
	if err != nil {
		return nil, fmt.Errorf("error retrieving directed hashrates: %v", err)
	}

	for _, m := range directedMembers {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			continue // Error handling: log the error or return it
		}
		directedHashRateSum += entry.HashRate
	}

	// Sum directed hash rates
	packageHashrate, err := redis.Values(conn.Do("ZRANGEBYSCORE", packageKey, currentTime, "+inf"))
	if err != nil {
		return nil, fmt.Errorf("error retrieving package hashrates: %v", err)
	}

	for _, m := range packageHashrate {
		var entry HashRateEntry
		if err := json.Unmarshal(m.([]byte), &entry); err != nil {
			continue // Error handling: log the error or return it
		}
		packageHashRateSum += entry.HashRate
	}

	// Create the desired data structure with the summed hashrate values
	summedHashRates := map[string]int64{
		"universal":    universalHashRateSum,
		"directed":     directedHashRateSum,
		"package":      packageHashRateSum,
		"member_total": directedHashRateSum + universalHashRateSum,
		"total":        directedHashRateSum + universalHashRateSum + packageHashRateSum,
	}
	return summedHashRates, nil
}

type ComplimentaryHashRates struct {
	CreatedAt string `json:"created_at"`
	HashRates int64  `json:"hash_rates"`
	Reason    string `json:"reason"`
	ExpireAt  string `json:"expire_at"`
}

// ComplimentaryHashRate redis 存储用户赠送的算力
func ComplimentaryHashRate(userId string, hashRate int64, reason string) (*ComplimentaryHashRates, error) {

	conn := RedisConn.Get()
	defer conn.Close()

	if _, err := conn.Do("SELECT", setting.UsingDBSetting.HashRate); err != nil {
		return nil, err
	}

	validPeriod := time.Duration(365) * 24 * time.Hour

	createdAt := time.Now()
	expireAt := createdAt.Add(validPeriod)

	createdAtStr := createdAt.Format("2006-01-02 15:04:05")
	expireAtStr := expireAt.Format("2006-01-02 15:04:05")

	switch hashRate {
	case 50:
		reason = "邀请好友注册赠送"
	case 100:
		reason = "注册赠送"
	}
	complimentaryHashRates := ComplimentaryHashRates{
		CreatedAt: createdAtStr,
		HashRates: hashRate,
		Reason:    reason,
		ExpireAt:  expireAtStr,
	}
	key := fmt.Sprintf("complimentary:user:%s:hashrate", userId)

	jsonData, err := json.Marshal(complimentaryHashRates)

	if err != nil {
		return nil, err
	}

	_, err = conn.Do("ZADD", key, expireAt.Unix(), jsonData)

	if err != nil {
		return nil, err
	}
	return &complimentaryHashRates, nil
}
