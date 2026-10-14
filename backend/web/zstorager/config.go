package zstorager

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/denghuo98/zfast/internal/dao"
	webSchema "github.com/denghuo98/zfast/zschema/zweb"
)

var config *webSchema.UploadConfig

func SetConfig(c *webSchema.UploadConfig) {
	config = c
}

func GetConfig() *webSchema.UploadConfig {
	return config
}

func GetModel(ctx context.Context) *gdb.Model {
	return dao.SysAttachment.Ctx(ctx)
}
