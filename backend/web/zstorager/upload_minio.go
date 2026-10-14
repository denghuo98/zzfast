package zstorager

import (
	"context"
	"mime"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
)

// MinioDrive MinIO 对象存储驱动
type MinioDrive struct{}

// Upload 上传到 MinIO 对象存储
func (d *MinioDrive) Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error) {
	if config.MinioPath == "" {
		err = gerror.New("MinIO存储驱动必须配置存储路径!")
		return
	}
	if config.MinioEndpoint == "" || config.MinioBucket == "" {
		err = gerror.New("MinIO存储驱动必须配置Endpoint和Bucket!")
		return
	}

	client, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL == 1,
	})
	if err != nil {
		return "", err
	}

	if err = s3utils.CheckValidBucketName(config.MinioBucket); err != nil {
		return
	}

	fullPath = GenFullPath(config.MinioPath, gfile.Ext(file.Filename))
	if err = s3utils.CheckValidObjectName(fullPath); err != nil {
		return
	}

	reader, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()

	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = client.PutObject(ctx, config.MinioBucket, fullPath, reader, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return
}

// CreateMultipart 创建分片事件
func (d *MinioDrive) CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}

// UploadPart 上传分片
func (d *MinioDrive) UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error) {
	err = gerror.New("当前驱动暂不支持分片上传！")
	return
}
