package ent_util

import (
	"context"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

type DriverX struct {
	*sql.Driver
	explainer Explainer
}

func NewDriverX(drv *sql.Driver, dialect string) dialect.Driver {
	return &DriverX{
		Driver:    drv,
		explainer: newExplainer(dialect),
	}
}

func (d *DriverX) Exec(ctx context.Context, query string, args, v any) error {
	//TODO implement me
	panic("implement me")
}

func (d *DriverX) Query(ctx context.Context, query string, args, v any) error {
	raw := d.explainer.Explain(query, args.([]interface{})...)
	slog.Info("sql", slog.Any("sql", raw))

	return d.Driver.Query(ctx, query, args, v)
}

// Explainer 基于 [GORM DryRun](https://gorm.io/zh_CN/docs/session.html#DryRun)
// 此处, `gorm.io/driver/sqlite` 依赖 CGO, 项目依赖一个Pure Go 的 sqlite3 驱动 `github.com/logoove/sqlite`
// 同时使用会抛出错误 `sql: Register called twice for driver sqlite3`
// 注意：SQL 并不总是能安全地执行，GORM 仅将其用于日志，它可能导致会 SQL 注入
type Explainer interface {
	Explain(sql string, vars ...interface{}) string
}

func newExplainer(dialect string) Explainer {
	switch dialect {
	case "mysql":
		return mysql.Open("")
	case "postgres":
		return postgres.Open("")
	case "sqlite3":
		return &sqlite3Explainer{}
	default:
		panic("implement me")
	}
}

type sqlite3Explainer struct{}

func (s *sqlite3Explainer) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}
