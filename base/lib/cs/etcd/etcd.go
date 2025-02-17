package etcd

import (
	"context"
	"fmt"
	"time"

	"go-web-template/base/common/utils"

	etcdCliend "go.etcd.io/etcd/client/v3"
)

func Init() {
	config := etcdCliend.Config{
		Endpoints:        []string{"http://127.0.0.1:2379"},
		AutoSyncInterval: time.Second,
	}

	client, err := etcdCliend.New(config)
	utils.PanicAndPrintIfNotNil(err)

	do, err := client.Do(context.Background(), etcdCliend.OpPut("hello", "world"))
	utils.PanicAndPrintIfNotNil(err)
	put := do.Put()
	fmt.Printf("%#v\n", put)

	getDo, err := client.Do(context.Background(), etcdCliend.OpGet("hello"))
	utils.PanicAndPrintIfNotNil(err)
	get := getDo.Get()

	fmt.Printf("%#v\n", get)
}
