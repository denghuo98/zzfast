package zstorager

import (
	"context"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
)

// OssDrive 阿里云 OSS 驱动
type OssDrive struct{}

// Upload 上传到阿里云 OSS
func (d *OssDrive) Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error) {
	if config.OssPath == "" {
		err = gerror.New("OSS存储驱动必须配置存储路径!")
		return
	}
	if config.OssEndpoint == "" || config.OssBucket == "" {
		err = gerror.New("OSS存储驱动必须配置Endpoint和Bucket!")
		return
	}

	reader, err := file.Open()
	if err != nil {
		return
	}
	defer func() { _ = reader.Close() }()

	client, err := oss.New(config.OssEndpoint, config.OssSecretId, config.OssSecretKey)
	if err != nil {
		return
	}

	bucket, err := client.Bucket(config.OssBucket)
	if err != nil {
		return
	}

	fullPath = GenFullPath(config.OssPath, gfile.Ext(file.Filename))
	err = bucket.PutObject(fullPath, reader)
	return
}

// CreateMultipart 创建分片事件
func (d *OssDrive) CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}

// UploadPart 上传分片
func (d *OssDrive) UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}
