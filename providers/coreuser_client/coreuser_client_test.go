package article_client_test

import (
	"testing"

	"code.byted.org/learning_fe/go_modules/utils"

	_ "code.byted.org/learning/learning_open_api/clients/toutiao/learning/coreuser"

	"code.byted.org/learning/learning_open_api/thrift_gen/toutiao/learning/coreuser"
)

const (
	AppID     = 1345
	IsProduct = false
)

var (
	UserId = int64(86224290727)
)

// go test -run="IMBatchDeleteConversationParticipant"
func TestDemo(t *testing.T) {
	client := coreuser_client.CoreuserClientInstance()
	ctx := utils.CreateTestContext("toutiao.learning_fe.open")
	req := &coreuser.GetCoreUserRequest{
		UserId: UserId,
	}
	resp, err := client.KitcClient.Call("GetCoreUser", ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	realResp := resp.RealResponse().(*coreuser.GetCoreUserResponse)
	t.Log(realResp)
	// t.Log(resp)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
