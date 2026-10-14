package common

import (
	"context"

	"github.com/dustin/go-humanize"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/denghuo98/zfast/internal/dao"
	"github.com/denghuo98/zfast/internal/model/entity"
	"github.com/denghuo98/zfast/web/zstorager"
	"github.com/denghuo98/zfast/zconsts"
	commonSchema "github.com/denghuo98/zfast/zschema/common"
	"github.com/denghuo98/zfast/zservice"
)

type sCommonUpload struct{}

func NewCommonUpload() *sCommonUpload {
	return &sCommonUpload{}
}

func init() {
	zservice.RegisterCommonUpload(NewCommonUpload())
	// 注入附件记录存储实现，桥接 zstorager 与数据访问层
	zstorager.SetRepository(&attachmentRepository{})
}

// attachmentRepository 附件记录存储实现
type attachmentRepository struct{}

// HasFile 按 md5 查询附件记录，不存在返回 nil
func (r *attachmentRepository) HasFile(ctx context.Context, md5 string) (*zstorager.Attachment, error) {
	var res *entity.SysAttachment
	if err := dao.SysAttachment.Ctx(ctx).Where("md5", md5).Scan(&res); err != nil {
		return nil, gerror.Wrap(err, "检查文件hash时出现错误")
	}
	if res == nil {
		return nil, nil
	}

	// 附件存在则更新最后上传时间，保证上传列表更新显示在最前面
	if res.Id > 0 {
		update := g.Map{
			"status":     zconsts.StatusEnabled,
			"updated_at": gtime.Now(),
		}
		_, _ = dao.SysAttachment.Ctx(ctx).WherePri(res.Id).Data(update).Update()
	}

	var attachment zstorager.Attachment
	if err := gconv.Struct(res, &attachment); err != nil {
		return nil, err
	}
	return &attachment, nil
}

// Create 写入附件记录并返回自增 ID
func (r *attachmentRepository) Create(ctx context.Context, a *zstorager.Attachment) (int64, error) {
	var e entity.SysAttachment
	if err := gconv.Struct(a, &e); err != nil {
		return 0, err
	}
	id, err := dao.SysAttachment.Ctx(ctx).Data(e).OmitEmptyData().InsertAndGetId()
	if err != nil {
		return 0, gerror.Wrap(err, zconsts.ErrorORM)
	}
	return id, nil
}

// UploadFile 上传文件
func (s *sCommonUpload) UploadFile(ctx context.Context, uploadType string, file *ghttp.UploadFile) (res *commonSchema.AttachmentListModel, err error) {
	attachment, err := zstorager.DoUpload(ctx, uploadType, file)
	if err != nil {
		return
	}

	attachment.FileUrl = zstorager.LastUrl(ctx, attachment.FileUrl, attachment.Drive)

	var e entity.SysAttachment
	if err = gconv.Struct(attachment, &e); err != nil {
		return nil, err
	}
	res = &commonSchema.AttachmentListModel{
		SysAttachment: e,
		SizeFormat:    humanize.IBytes(gconv.Uint64(e.Size)),
	}
	return
}

// CheckMultipart 检查文件分片
func (s *sCommonUpload) CheckMultipart(ctx context.Context, in *commonSchema.CheckMultipartInp) (res *commonSchema.CheckMultipartModel, err error) {
	data, err := zstorager.CheckMultipart(ctx, in.CheckMultipartParams)
	if err != nil {
		return nil, err
	}
	res = new(commonSchema.CheckMultipartModel)
	res.CheckMultipartModel = data
	return
}

// UploadPart 上传分片
func (s *sCommonUpload) UploadPart(ctx context.Context, in *commonSchema.UploadPartInp) (res *commonSchema.UploadPartModel, err error) {
	data, err := zstorager.UploadPart(ctx, in.UploadPartParams)
	if err != nil {
		return nil, err
	}
	res = new(commonSchema.UploadPartModel)
	res.UploadPartModel = data
	return
}
