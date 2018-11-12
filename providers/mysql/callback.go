package mysql

import (
	"time"

	"code.byted.org/gopkg/gorm"
)

const (
	ScopeKeyStartTime = "XStartTime"
)

// MetricsBeforeQuery add start time for metrics.
func MetricsBeforeQuery(scope *gorm.Scope) {
	now := time.Now()
	scope.Set(ScopeKeyStartTime, &now)
}

// MetricsAfterQuery emit metrics cost.
func MetricsAfterQuery(scope *gorm.Scope) {
	val, ok := scope.Get(ScopeKeyStartTime)
	if !ok {
		return
	}
	start, ok := val.(*time.Time)
	if !ok {
		return
	}
	latency := time.Since(*start).Nanoseconds() / 1000

	var status string
	switch scope.DB().Error {
	case nil:
		status = MetricsStatusSuccess
	case DBNotFoundError:
		status = MetricsStatusMiss
	case DBNoRowBeAffected:
		status = MetricsStatusMissNoAffect
	default:
		status = MetricsStatusError
	}

	method := "none"
	if v, ok := scope.DB().Ctx.Value(MetricsCallDbMethodKey).(string); ok {
		method = v
	}

	emit(method, status, latency)
}

// RegisterMetricsCallback register metrics callback.
func RegisterMetricsCallback(db *gorm.DB) {
	db.Callback().Query().Before("gorm:query").Register("before_query", MetricsBeforeQuery)
	db.Callback().Query().After("gorm:query").Register("after_query", MetricsAfterQuery)
	db.Callback().Create().Before("gorm:create").Register("before_create", MetricsBeforeQuery)
	db.Callback().Create().After("gorm:create").Register("after_create", MetricsAfterQuery)
	db.Callback().Delete().Before("gorm:delete").Register("before_delete", MetricsBeforeQuery)
	db.Callback().Delete().After("gorm:delete").Register("after_delete", MetricsAfterQuery)
	db.Callback().Update().Before("gorm:update").Register("before_update", MetricsBeforeQuery)
	db.Callback().Update().After("gorm:update").Register("after_update", MetricsAfterQuery)
	db.Callback().RowQuery().Before("gorm:row_query").Register("before_row_query", MetricsBeforeQuery)
	db.Callback().RowQuery().After("gorm:row_query").Register("after_row_query", MetricsAfterQuery)
}
