package infra

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConnector struct {
	Clinet *redis.Client
}

func (r *RedisConnector) CreateUser(ctx context.Context, username string, password string, db int) error {
	isExist, checkUserErr := r.IsUserExist(ctx, username)
	if checkUserErr != nil {
		return checkUserErr
	}

	if isExist {
		return fmt.Errorf(fmt.Sprintf("duplicate user: %s", username))
	}

	cmds := []interface{}{
		"acl",
		"setuser",
		username,
		fmt.Sprintf(">%s", password),
		"on",
		"+@all",
		"allchannels",
		"allkeys",
		"-acl",
		"-bgrewriteaof",
		"-bgsave",
		"-config",
		"-module",
	}

	if db >= 0 {
		cmds = append(cmds, "-select", fmt.Sprintf("+select|%d", db))
	}

	return r.Clinet.Do(ctx, cmds...).Err()
}

func (r *RedisConnector) UpdateUser(ctx context.Context, username string, password string, db int) error {
	cmds := []interface{}{
		"acl",
		"setuser",
		username,
		fmt.Sprintf(">%s", password),
		"on",
		"+@all",
		"allchannels",
		"allkeys",
		"-acl",
		"-bgrewriteaof",
		"-bgsave",
		"-config",
		"-module",
	}

	if db >= 0 {
		cmds = append(cmds, "-select", fmt.Sprintf("+select|%d", db))
	}

	pipe := r.Clinet.Pipeline()
	pipe.Do(ctx, "acl", "deluser", username)
	pipe.Do(ctx, cmds...)
	_, err := pipe.Exec(ctx)

	return err
}

func (r *RedisConnector) GetUserInfo(ctx context.Context, username string) (interface{}, error) {
	return r.Clinet.Do(ctx, "acl", "getuser", username).Result()
}

func (r *RedisConnector) IsUserExist(ctx context.Context, username string) (bool, error) {
	result, err := r.GetUserInfo(ctx, username)
	if err != nil {
		if result == nil {
			return false, nil
		}

		return false, err
	}

	if result == nil {
		return false, nil
	}
	userInfo, ok := result.([]interface{})
	if !ok {
		return false, nil
	}

	return len(userInfo) != 0, nil
}

func (r *RedisConnector) DeleteUser(ctx context.Context, username string) error {
	return r.Clinet.Do(ctx, "acl", "deluser", username).Err()
}
