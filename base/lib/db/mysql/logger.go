package mysql

import (
	"context"
	"time"

	"gorm.io/gorm/logger"

	localLogger "go-web-template/base/lib/logger"
)

var _ logger.Interface = (*GormLogger)(nil)

type GormLogger struct {
	SlowThreshold time.Duration
	level         logger.LogLevel
}

func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		SlowThreshold: slowThreshold,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	g.level = level
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.level == logger.Silent {
		return
	}

	localLogger.SInfoF(ctx, s, i)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.level == logger.Silent {
		return
	}

	localLogger.SWarnF(ctx, s, i)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.level == logger.Silent {
		return
	}

	localLogger.SErrorF(ctx, s, i)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsed := time.Since(begin)

	//执行SQL
	sql, rows := fc()

	//打印错误日志
	if err != nil {
		localLogger.SErrorF(ctx, "seq error, | sql=%v, rows=%v, elapsed=%v", sql, rows, elapsed)
	}

	// 打印慢查询日志
	if g.SlowThreshold != 0 && elapsed > g.SlowThreshold {
		localLogger.SWarnF(ctx, "database slow Log, | sql=%v, rows=%v, elapsed=%v", sql, rows, elapsed)
	}
}
