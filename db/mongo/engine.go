package mongo

import (
	"context"
	"github.com/xieqiaoyu/xin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

//Config config provide mongo connection
type Config interface {
	GetMongoURI(id string) (source string, err error)
}

//Service mongo connect service
type Service struct {
	instances *sync.Map
	config    Config
}

//NewService create a new mongo connect service
func NewService(config Config) *Service {
	return &Service{
		instances: new(sync.Map),
		config:    config,
	}
}

//Engine get redis client by id
func (s *Service) Engine(id string) (*mongo.Client, error) {
	instance, exists := s.instances.Load(id)
	if exists {
		return instance.(*mongo.Client), nil
	}
	mongoURI, err := s.config.GetMongoURI(id)
	if err != nil {
		return nil, xin.NewTracedE(err)
	}
	connectOption := options.Client()
	connectOption.ApplyURI(mongoURI)

	client, err := mongo.NewClient(connectOption)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to Create mongo client [%s] :%s", id, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, xin.WrapEf(&xin.InternalError{}, "Fail to Connect mongo  [%s] :%s", id, err)
	}
	instance, loaded := s.instances.LoadOrStore(id, client)
	if loaded {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.Disconnect(ctx)
	}
	return instance.(*mongo.Client), nil
}

//Close close  all connections
func (s *Service) Close() error {
	var err error
	s.instances.Range(func(key, engine interface{}) bool {
		s.instances.Delete(key)
		client := engine.(*mongo.Client)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Disconnect(ctx)
		if err != nil {
			return false
		}
		return true
	})
	return err
}
