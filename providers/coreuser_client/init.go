package article_client

import (
	"log"

	// 初始化coreuser

	"code.byted.org/kite/kitc"
)

func initClient() (*kitc.KitcClient, error) {
	var err error
	opts := []kitc.Option{}
	// 自定义下游服务的地址
	client, err := kitc.NewClient("toutiao.learning.coreuser", opts...)
	if err != nil {
		log.Fatalf("init client error: %s", err)
	}
	return client, err
}
