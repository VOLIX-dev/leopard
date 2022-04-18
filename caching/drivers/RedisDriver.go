package drivers

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/volix-dev/leopard/caching"
	"strconv"
	"time"
)

type RedisDriver struct {
	client *redis.Client
}

type RedisSettings struct {
	Host     string
	Port     int
	Password string
	Database int
}

func init() {
	caching.Register("redis", func(config any) (caching.Driver, error) {
		conf, ok := config.(RedisSettings)

		if !ok {
			return nil, errors.New("invalid redis settings")
		}

		return newRedisDriver(conf)
	})
}

func newRedisDriver(config RedisSettings) (*RedisDriver, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.Database,
	})
	return &RedisDriver{
		client,
	}, nil
}

func (r *RedisDriver) Get(key string, target any) (bool, error) {
	if _, stru := target.(*struct{}); stru {
		_, marshaler := target.(encoding.BinaryUnmarshaler)
		if !marshaler {

			data, err := r.client.Get(context.TODO(), key).Bytes()
			if err != nil {
				return false, err
			}

			err = json.Unmarshal(data, target)
			if err != nil {
				return false, err
			}
		}
	}
	return true, r.client.Get(context.TODO(), key).Scan(target)
}

func (r *RedisDriver) Set(key string, value any) error {
	return r.SetTTL(key, value, 0)
}

func (r *RedisDriver) SetTTL(key string, value any, ttl int) error {
	if _, stru := value.(struct{}); stru {
		_, marshaler := value.(encoding.BinaryMarshaler)
		if !marshaler {
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}

			return r.client.Set(context.TODO(), key, data, time.Duration(ttl)*time.Second).Err()
		}
	}

	return r.client.Set(context.TODO(), key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisDriver) Delete(key string) error {
	return r.client.Del(context.TODO(), key).Err()
}

func (r *RedisDriver) Close() error {
	return r.client.Close()
}

func (r *RedisDriver) Open() error {
	// Empty because redis driver opens every request
	return nil
}
