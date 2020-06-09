package modules

import (
	"github.com/go-redis/redis"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("RedisClient", newRedisClient)
}

type redisClient struct {
	client *redis.Client
}

func newRedisClient(options *redis.Options) (*redisClient, error) {
	client := redis.NewClient(options)
	return &redisClient{client: client}, nil
}

func (self *redisClient) Dox(args ...interface{}) (interface{}, error) {
	cmd := self.client.Do(args...)
	err := cmd.Err()
	data, err := cmd.Result()
	if err != nil {
		if err == redis.Nil  {
			err = nil
		}
	}

	return data, err
}

func (self *redisClient) Test() () {

}
