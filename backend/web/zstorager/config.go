package zstorager

import (
	webSchema "github.com/denghuo98/zzfast/zschema/zweb"
)

var (
	config     *webSchema.UploadConfig
	repository AttachmentRepository
)

func SetConfig(c *webSchema.UploadConfig) {
	config = c
}

func GetConfig() *webSchema.UploadConfig {
	return config
}

// SetRepository 注入附件记录存储实现（由 service 层注册，避免 zstorager 依赖数据访问层）
func SetRepository(r AttachmentRepository) {
	repository = r
}
