package testing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	etcdctl "go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdV3ConfigLoader struct {
	Config     *etcdctl.Config
	Key        string
	ConfigType string
}

func (l *EtcdV3ConfigLoader) LoadConfig(vc *viper.Viper) error {
	cli, err := etcdctl.New(*l.Config)
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, l.Key)
	cancel()
	if err != nil {
		return err
	}
	if len(resp.Kvs) <= 0 {
		return fmt.Errorf("Config key %s  is not setted in etcd server", l.Key)
	}
	configBytes := resp.Kvs[0].Value
	if l.ConfigType != "" {
		vc.SetConfigType(l.ConfigType)
	} else {
		vc.SetConfigType("toml")
	}
	vc.ReadConfig(bytes.NewReader(configBytes))
	return nil
}

func NewEtcdV3ConfigLoader(endpoint, key string) *EtcdV3ConfigLoader {
	return &EtcdV3ConfigLoader{
		Config: &etcdctl.Config{
			Endpoints:   []string{endpoint},
			DialTimeout: 5 * time.Second,
		},
		Key:        key,
		ConfigType: "toml",
	}
}
