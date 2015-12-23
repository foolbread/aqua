//@auther foolbread
//@time 2015-12-22
//@file aqua/common/storage/redis_handler.go
package storage

import (
	"github.com/mediocregopher/radix.v2/redis"
	"sync"
	"time"
)

var connect_timeout time.Duration = 2 * time.Minute

type RedisHandler struct {
	server string
	pwd    string
	client *redis.Client
	lo     *sync.Mutex
}

func NewRedisHandler(ser string, pwd string) (*RedisHandler, error) {
	ret := new(RedisHandler)
	ret.server = ser
	ret.pwd = pwd
	ret.lo = new(sync.Mutex)

	return ret, ret.connect()
}

func (s *RedisHandler) SetKey(key, value string) error {
	rsp := s.redisCmd("SET", key, value)
	return rsp.Err
}

func (s *RedisHandler) GetKey(key string) (string, error) {
	rsp := s.redisCmd("GET", key)
	if rsp.Err != nil {
		return "", rsp.Err
	}

	return rsp.Str()
}

func (s *RedisHandler) DelKey(key string) error {
	rsp := s.redisCmd("DEL", key)
	if rsp.Err != nil {
		return rsp.Err
	}

	return nil
}

func (s *RedisHandler) ExistsKey(key string) (bool, error) {
	rsp := s.redisCmd("EXISTS", key)
	if rsp.Err != nil {
		return false, rsp.Err
	}

	v, _ := rsp.Int()
	return v > 0, nil
}

//////////////////////////////////////////////////////////////////////////
func (s *RedisHandler) authHandler() error {
	rsp := s.redisCmd("AUTH", s.pwd)
	if rsp.Err != nil {
		return rsp.Err
	}

	return nil
}

func (s *RedisHandler) redisCmd(cmd string, args ...interface{}) *redis.Resp {
	s.lo.Lock()
	rsp := s.client.Cmd(cmd, args)
	if rsp.Err != nil {
		err := s.connect()
		if err != nil {
			rsp.Err = err
			s.lo.Unlock()
			return rsp
		}

		s.lo.Unlock()
		return s.client.Cmd(cmd, args)
	}

	s.lo.Unlock()
	return rsp
}

//////////////////////////////////////////////////////////////////////////
func (s *RedisHandler) connect() error {
	if len(s.pwd) > 0 {
		return s.initHandlerWithAuth()
	}

	return s.initHandler()
}

func (s *RedisHandler) initHandler() error {
	var err error
	s.client, err = redis.DialTimeout("tcp", s.server, connect_timeout)
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisHandler) initHandlerWithAuth() error {
	var err error
	s.client, err = redis.DialTimeout("tcp", s.server, connect_timeout)
	if err != nil {
		return err
	}

	return s.authHandler()
}

func (s *RedisHandler) closeHandler() {
	if s.client != nil {
		s.client.Close()
	}
}
