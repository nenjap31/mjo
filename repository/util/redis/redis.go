package redis

import (
	"time"

	"github.com/go-redis/redis"

	log "github.com/sirupsen/logrus"
)

// LOGPREFIX For logging purpose
const LOGPREFIX = "[UTILS REDIS REPOSITORY]"

type IRepository interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (*string, error)
	Keys(pattern string) ([]string, error)
	Delete(key string) error
}

type Repository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *Repository {
	return &Repository{
		client,
	}
}

func (r *Repository) Set(key string, value string, expiration time.Duration) error {
	command := r.client.Set(key, value, expiration)
	if command.Err() != nil {
		log.Error(LOGPREFIX, " Set - ", command.Err().Error())
		return command.Err()
	}
	return nil
}

func (r *Repository) Get(key string) (*string, error) {
	command := r.client.Get(key)
	if command.Err() != nil {
		if command.Err() == redis.Nil {
			return nil, nil
		}
		log.Error(LOGPREFIX, " Get - ", command.Err().Error())
		return nil, command.Err()
	}
	value := command.Val()
	return &value, nil
}

func (r *Repository) Keys(pattern string) ([]string, error) {
	command := r.client.Keys(pattern)
	if command.Err() != nil {
		if command.Err() == redis.Nil {
			return nil, nil
		}
		log.Error(LOGPREFIX, " Keys - ", command.Err().Error())
		return nil, command.Err()
	}

	value := command.Val()
	return value, nil
}

func (r *Repository) Delete(key string) error {
	command := r.client.Del(key)
	if command.Err() != nil {
		log.Error(LOGPREFIX, " Set - ", command.Err().Error())
		return command.Err()
	}
	return nil
}