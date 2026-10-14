package zstorager

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// CosDrive 腾讯云 COS 驱动
type CosDrive struct{}

// Upload 上传到腾讯云 COS 对象存储
func (d *CosDrive) Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error) {
	if config.CosPath == "" {
		err = gerror.New("COS存储驱动必须配置存储路径!")
		return
	}
	if config.CosBucketURL == "" {
		err = gerror.New("COS存储驱动必须配置Bucket访问地址!")
		return
	}

	reader, err := file.Open()
	if err != nil {
		return
	}
	defer func() { _ = reader.Close() }()

	bucketURL, err := url.Parse(config.CosBucketURL)
	if err != nil {
		return "", gerror.Wrap(err, "COS Bucket访问地址格式错误")
	}
	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CosSecretId,
			SecretKey: config.CosSecretKey,
		},
	})

	fullPath = GenFullPath(config.CosPath, gfile.Ext(file.Filename))
	_, err = client.Object.Put(ctx, fullPath, reader, nil)
	return
}

// CreateMultipart 创建分片事件
func (d *CosDrive) CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}

// UploadPart 上传分片
func (d *CosDrive) UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}
