package mysql

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"code.byted.org/gopkg/context"
	"code.byted.org/gopkg/env"

	"code.byted.org/golf/ssconf"
	gormConf "code.byted.org/gopkg/dbutil/conf"
	"code.byted.org/gopkg/dbutil/gormdb"
	"code.byted.org/gopkg/gorm"
	"code.byted.org/gopkg/logs"
	"code.byted.org/gopkg/metrics"
	mysqldriver "code.byted.org/gopkg/mysql-driver"
	"code.byted.org/learning_fe/go_modules/utils"
	"gopkg.in/yaml.v2"
)

const (
	MetricsNamePrefix         = "toutiao.service.thrift.db."
	MetricsStatusSuccess      = "success"
	MetricsStatusError        = "error"
	MetricsStatusMiss         = "miss"
	MetricsStatusMissNoAffect = "noaffect"
)

var (
	metricsClient          *metrics.MetricsClient
	MetricsCallDbMethodKey = utils.CtxKey{}
)

func init() {
	metricsClient = metrics.NewDefaultMetricsClient(MetricsNamePrefix+env.PSM(), true)
}

type DBWrapper struct {
	readHandler  *gormdb.DBHandler
	writeHandler *gormdb.DBHandler
}

func NewDbWrapperWithSsConf(
	confPath string,
	dbName string,
	wrapConfFunc func(*gormConf.DBOptional, *gormConf.DBOptional),
) (*DBWrapper, error) {
	confPath, _ = filepath.Abs(confPath)
	ssConf, _ := ssconf.LoadSsConfFile(confPath)

	var readDbConf gormConf.DBOptional
	var writeDbConf gormConf.DBOptional
	if !env.IsProduct() {
		readDbConf = gormConf.GetDbConf(ssConf, dbName, gormConf.Offline)
		writeDbConf = gormConf.GetDbConf(ssConf, dbName, gormConf.Offline)
	} else {
		readDbConf = gormConf.GetDbConf(ssConf, dbName, gormConf.Read)
		writeDbConf = gormConf.GetDbConf(ssConf, dbName, gormConf.Write)
	}

	readDbConf.DriverName = "mysql2"
	writeDbConf.DriverName = "mysql2"

	if wrapConfFunc != nil {
		wrapConfFunc(&readDbConf, &writeDbConf)
	}

	wrapper := &DBWrapper{
		readHandler:  gormdb.NewDBHandlerWithOptional(&readDbConf),
		writeHandler: gormdb.NewDBHandlerWithOptional(&writeDbConf),
	}

	return wrapper, nil
}

func NewDBWrapperWithYamlFile(dbname, yamlFile string) (*DBWrapper, error) {
	yamlBytes, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	return NewDBWrapper(dbname, yamlBytes)
}

func NewDBWrapper(dbname string, yamlBytes []byte) (*DBWrapper, error) {
	mysqldriver.SetPSMCluster(env.PSM(), env.Cluster())
	dbOptionals := make(map[string]gormConf.DBOptional)
	if err := yaml.Unmarshal(yamlBytes, &dbOptionals); err != nil {
		return nil, err
	}

	writeDBOptional, ok := dbOptionals[dbname+"_write"]
	if !ok {
		return nil, errors.New("No Write DB config")
	}

	readDBOptional, ok := dbOptionals[dbname+"_read"]
	if !ok {
		readDBOptional = writeDBOptional
	}

	// Dev environment
	if testDBOptional, ok := dbOptionals[dbname+"_test"]; ok && !env.IsProduct() {
		readDBOptional = testDBOptional
		writeDBOptional = testDBOptional
	}

	wrapper := &DBWrapper{
		readHandler:  gormdb.NewDBHandlerWithOptional(&readDBOptional),
		writeHandler: gormdb.NewDBHandlerWithOptional(&writeDBOptional),
	}

	return wrapper, nil
}

func NewDBWrapperRaw(
	readHandler *gormdb.DBHandler,
	writeHandler *gormdb.DBHandler,
) (*DBWrapper, error) {
	return &DBWrapper{
		readHandler:  readHandler,
		writeHandler: writeHandler,
	}, nil
}

func (dw *DBWrapper) Init() {
	RegisterMetricsCallback(dw.Read())
	RegisterMetricsCallback(dw.Write())
}

// GetReadDBConnection get read connection.
func (dw *DBWrapper) GetReadDBConnection() (*gorm.DB, error) {
	return dw.readHandler.GetConnection()
}

// GetWriteDBConnection get write connection.
func (dw *DBWrapper) GetWriteDBConnection() (*gorm.DB, error) {
	return dw.writeHandler.GetConnection()
}

// Read get read connection.
func (dw *DBWrapper) Read() *gorm.DB {
	conn, _ := dw.GetReadDBConnection()
	return conn
}

// Write get write connection.
func (dw *DBWrapper) Write() *gorm.DB {
	conn, _ := dw.GetWriteDBConnection()
	return conn
}

// ReadWithMethod get read connection and set method into context.
func (dw *DBWrapper) ReadWithMethod(ctx context.Context, method string) *gorm.DB {
	if len(method) != 0 {
		ctx = context.WithValue(ctx, MetricsCallDbMethodKey, method)
	}
	return dw.Read().Context(ctx)
}

// WriteWithMethod get write connection and set method into context.（不推荐用）
func (dw *DBWrapper) WriteWithMethod(ctx context.Context, method string) *gorm.DB {
	if len(method) != 0 {
		ctx = context.WithValue(ctx, MetricsCallDbMethodKey, method)
	}
	return dw.Write().Context(ctx)
}

func FormatExpr(exp string) interface{} {
	return gorm.Expr(exp)
}

func emit(method string, status string, latency int64) {
	tags := map[string]string{
		"method": method,
	}
	metricsClient.EmitCounter("call."+status+".throughput", 1, "", tags)
	metricsClient.EmitTimer("call."+status+".latency.us", latency, "", tags)
}

// Commit commit write action (update、delete...) to db with transaction.
func (dw *DBWrapper) Commit(ctx context.Context, method string, entries ...interface{}) (err error) {
	conn := dw.WriteWithMethod(ctx, method)
	conn = conn.Context(ctx)

	tx := conn.Begin()
	// tx = tx.LogMode(true)
	for _, ientry := range entries {
		switch entry := ientry.(type) {
		case *CreatedEntry:
			r := tx.Create(entry.Object)
			err = r.Error
		case *UpdatedEntry:
			r := tx.Save(entry.Object)
			err = r.Error
		case *DeletedEntry:
			r := tx.Delete(entry.Object, entry.Where...)
			err = r.Error
		case *UpdatedWithColumnsEntry:
			if len(entry.Table) != 0 {
				tx = tx.Table(entry.Table)
			}
			if entry.Model != nil {
				tx = tx.Model(entry.Model)
			}
			if len(entry.Where) != 0 {
				for _, searchOpt := range entry.Where {
					tx = tx.Where(searchOpt.Format, searchOpt.Args...)
				}
			}
			var r *gorm.DB
			if entry.Updates != nil {
				r = tx.Updates(entry.Updates)
			} else if entry.Update != nil {
				r = tx.Update(entry.Update...)
			}
			var rowsAffected int64
			rowsAffected, err = r.RowsAffected, r.Error
			if err == nil && entry.StrickMode && rowsAffected == 0 {
				err = DBNoRowBeAffected
			}
		case *ExecEntry:
			r := tx.Exec(entry.Format, entry.Args...)
			if r.Error != nil {
				err = r.Error
			}
		default:
			err = errors.New("Unknown Entry")
		}

		if err != nil {
			logs.CtxError(ctx, "DBWrapper commit error:%s", err)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

// Load query db.
func (dw *DBWrapper) Load(ctx context.Context, method string, obj interface{}, opt QueryOption) (err error) {
	return dw.LoadWithOption(ctx, method, obj, opt, &LoadOption{
		IgnoreNotFoundError: false,
	})
}

// LoadWithOption query db with config option.
func (dw *DBWrapper) LoadWithOption(ctx context.Context, method string, obj interface{}, opt QueryOption, loadOpt *LoadOption) (err error) {
	conn := dw.ReadWithMethod(ctx, method)
	conn = conn.Context(ctx)

	if opt.OrderBy != "" {
		conn = conn.Order(opt.OrderBy)
	}
	if len(opt.Where) != 0 {
		for _, searchOpt := range opt.Where {
			conn = conn.Where(searchOpt.Format, searchOpt.Args...)
		}
	}

	if opt.Limit != 0 {
		conn = conn.Offset(opt.Offset).Limit(opt.Limit)
	}

	if len(opt.Select) > 0 {
		conn = conn.Select(opt.Select)
	}

	if opt.GetAll {
		conn = conn.Find(obj)
	} else if opt.Count {
		if len(opt.Table) > 0 {
			conn = conn.Table(opt.Table)
		}
		if opt.CountModel != nil {
			conn = conn.Model(opt.CountModel)
		}
		conn = conn.Count(obj)
	} else {
		conn = conn.First(obj)
	}

	if conn.Error != nil {
		if conn.RecordNotFound() {
			if loadOpt == nil || !loadOpt.IgnoreNotFoundError {
				return DBNotFoundError
			} else {
				return DBNotFoundError
			}
		}
		logs.CtxError(ctx, "Load error err:%s opt:%v", conn.Error, opt)
		return conn.Error
	}

	return nil
}

// LoadAndCount load and count.
func (dw *DBWrapper) LoadAndCount(
	ctx context.Context,
	actionName string,
	opt QueryOption,
	queryResult interface{},
	countResult *int64,
) error {
	opt.GetAll = true
	countOpt := opt.Clone()
	countOpt.Count = true
	countOpt.GetAll = false
	countOpt.Limit = 0

	loadErrorChan := make(chan error, 2)

	loadFn := func(dw *DBWrapper, ctx context.Context, method string, opt *QueryOption, queryResult interface{}) {
		loadErrorChan <- dw.Load(ctx, method, queryResult, *opt)
	}

	go loadFn(dw, ctx, "GET"+actionName, &opt, queryResult)
	go loadFn(dw, ctx, "Load"+actionName, countOpt, countResult)

	return utils.WaitForError(loadErrorChan, 2)
}

// GetByKey get data from db with key.
// Eg. select * from a where b = xx;
func (dw *DBWrapper) GetByKey(
	ctx context.Context,
	idColName string,
	id interface{},
	data interface{},
	actionName string,
) error {
	return dw.Load(ctx, "Get"+actionName, data, NewEqualsQueryOption(idColName, id))
}

// GetByID get data from db by id.
// Eg. select * from a where id = xx;
func (dw *DBWrapper) GetByID(
	ctx context.Context,
	id interface{},
	data interface{},
	actionName string,
) error {
	return dw.GetByKey(ctx, "id", id, data, actionName)
}

// SaveDbObject insert or update a row.
func (dw *DBWrapper) SaveDbObject(
	ctx context.Context,
	data interface{},
	isCreate bool,
	actionName string,
) error {
	var entry interface{}
	var method string
	if isCreate {
		method = "Create" + actionName
		entry = &CreatedEntry{
			Object: data,
		}
	} else {
		method = "Update" + actionName
		entry = &UpdatedEntry{
			Object: data,
		}
	}
	return dw.Commit(ctx, method, entry)
}
