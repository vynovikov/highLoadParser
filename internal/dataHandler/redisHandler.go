package dataHandler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/vynovikov/highLoadParser/internal/logger"
)

type redisHandler struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisHandler() *redisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB})
	})
	ctx := context.Background()
	return &redisHandler{
		client: client,
		ctx:    ctx,
	}
}

func (r *redisHandler) Create(d DataHandlerDTO, bou Boundary) (ProducerUnit, error) {
	_, _, resTT := newKeyGeneralFromDTO(d), newKeyDetailed(d), newResult(d)

	val, err := newValue(d, bou)
	if err != nil &&
		!errors.Is(err, errHeaderNotFull) &&
		!errors.Is(err, errHeaderEnding) {

		return resTT, err
	}
	marshalledValue, _ := json.Marshal(val)

	_, err = r.client.Set(r.ctx, d.TS(), marshalledValue, 0).Result()

	if err != nil {
		logger.L.Errorf("in redisHandler.Create error :%v\n", err)
		return resTT, err
	}

	return resTT, nil
}

func (r *redisHandler) Set(k KeyDetailed, v Value) error {

	keyBytes, err := json.Marshal(k)
	if err != nil {
		logger.L.Errorf("in redisHandler.Set unable to marshal :%v\n", err)
		return err
	}

	valueBytes, err := json.Marshal(v)
	if err != nil {
		logger.L.Errorf("in redisHandler.Set unable to marshal :%v\n", err)
		return err
	}

	_, err = r.client.Set(r.ctx, string(keyBytes), valueBytes, 0).Result()

	if err != nil {
		logger.L.Errorf("in redisHandler.Set error while set to redis: %v\n", err)
		return err
	}
	return nil
}

func (r *redisHandler) Get(k KeyDetailed) (Value, error) {

	res := Value{}

	keyBytes, err := json.Marshal(k)
	if err != nil {
		logger.L.Errorf("in redisHandler.Get unable to marshal :%v\n", err)
		return Value{}, err
	}

	val, err := r.client.Get(r.ctx, string(keyBytes)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {

			return Value{}, errors.Join(ErrKeyNotFound, err)
		} else {

			logger.L.Errorf("in redisHandler.Get error while get from redis: %v\n", err)
			return Value{}, err
		}

	}
	err = json.Unmarshal(val, &res)
	if err != nil {
		logger.L.Errorf("in redisHandler.Get unable to unmarshal :%v\n", err)
		return Value{}, err
	}

	return res, nil
}

func (r *redisHandler) Read(DataHandlerDTO) (Value, error) {
	return Value{}, nil
}

func (r *redisHandler) Updade(DataHandlerDTO, Boundary) (ProducerUnit, error) {
	return &ProducerUnitStruct{}, nil
}

func (r *redisHandler) Delete(KeyDetailed) error {
	return nil
}
