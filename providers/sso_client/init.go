package sso_client

import (

	// 初始化coreuser

	"net/url"
	"time"

	"code.byted.org/learning/learning_open_api/providers/redis"
	"code.byted.org/learning_fe/go_modules/sso"
	"github.com/gin-gonic/gin"
)

func ssoRedirectFn(ctx *gin.Context, url string) {
	isAjax := ginUtils.IsAjax(ctx.Request)
	if isAjax {
		ginUtils.ReturnsRedirect(ctx, url, "跳转至登录")
	} else {
		ctx.Redirect(302, url)
	}
}

func ssoHandleRedirectURL(ctx *gin.Context, urlv *url.Values) {
	isAjax := ginUtils.IsAjax(ctx.Request)
	if isAjax {
		urlv.Del("url")
	} else {
	}
}

func initSSO() (*sso.SSO, error) {
	var err error
	redisClient := redis.RedisInstance().GetClient()
	store, err := sso.NewRedisStore(redisClient, "learning_open:sso:internal:", 30*24*time.Hour)
	if err != nil {
		return nil, err
	}
	bytedanceSSO = sso.NewBytedanceSSO(
		"/learning/open/sso/auth",
		Conf.URLPrefix+"/learning/open",
		store,
		"__BYTEDANCE_SSO__",
		&sso.SSOCallback{
			Redirect:          ssoRedirectFn,
			HandleRedirectURL: ssoHandleRedirectURL,
		},
	)
	return bytedanceSSO, nil
}
