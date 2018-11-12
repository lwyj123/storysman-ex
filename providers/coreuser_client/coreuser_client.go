package article_client

import (
	"context"
	"sync"

	_ "code.byted.org/learning/learning_open_api/clients/toutiao/learning/coreuser"

	"code.byted.org/learning/learning_open_api/thrift_gen/toutiao/learning/coreuser"

	"code.byted.org/kite/kitc"
)

// TODO: client做一层封装，不然这样只能通过client.Call传个字符串很蛋疼的

// CoreuserClient coreuser的thrift Client包装
type CoreuserClient struct {
	KitcClient *kitc.KitcClient
}

var _coreuserClient *ArticleClient
var _coreuserClientOnce sync.Once

// CoreuserClientInstance 返回 CoreuserClient 单例对象
func CoreuserClientInstance() *CoreuserClient {
	_coreuserClientOnce.Do(func() {
		kitcClient, _ := initClient()
		_coreuserClient = &CoreuserClient{
			KitcClient: kitcClient,
			// cache??
		}
	})
	return _coreuserClient
}

// HelloRPC RPC demo方法
func (client *CoreuserClient) HelloRPC(ctx context.Context) (string, error) {
	req := &coreuser.GetCoreUserRequest{
		UserId: 86224290727,
	}
	resp, err := client.KitcClient.Call("GetCoreUser", ctx, req)
	if err != nil {
		return "", err
	}

	realResp := resp.RealResponse().(*coreuser.GetCoreUserResponse)
	return realResp.UserInfo.Name, nil
}
