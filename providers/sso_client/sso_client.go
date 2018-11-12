package sso_client

import (
	"sync"

	_ "code.byted.org/learning/learning_open_api/clients/toutiao/learning/coreuser"
	"code.byted.org/learning_fe/go_modules/sso"
)

// BytedanceSSO BytedanceSSO用于SSO
type BytedanceSSO struct {
	bytedanceSSO *sso.SSO
}

var _bytedanceSSO *BytedanceSSO
var _bytedanceSSOOnce sync.Once

// BytedanceSSOInstance 返回 BytedanceSSO 单例对象
func BytedanceSSOInstance() *BytedanceSSO {
	_bytedanceSSOOnce.Do(func() {
		bytedanceSSO, _ := initSSO()
		_bytedanceSSO = &BytedanceSSO{
			bytedanceSSO: bytedanceSSO,
			// cache??
		}
	})
	return _bytedanceSSO
}
