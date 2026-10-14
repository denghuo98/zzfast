package zstorager

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// Attachment 附件记录（与数据访问层解耦的通用模型）
// 由 AttachmentRepository 负责持久化，字段与 entity.SysAttachment 保持一致以便转换
type Attachment struct {
	Id        int64       `json:"id"`
	AppId     string      `json:"appId"`
	MemberId  int64       `json:"memberId"`
	CateId    uint64      `json:"cateId"`
	Drive     string      `json:"drive"`
	Name      string      `json:"name"`
	Kind      string      `json:"kind"`
	MimeType  string      `json:"mimeType"`
	NaiveType string      `json:"naiveType"`
	Path      string      `json:"path"`
	FileUrl   string      `json:"fileUrl"`
	Size      int64       `json:"size"`
	Ext       string      `json:"ext"`
	Md5       string      `json:"md5"`
	Status    int         `json:"status"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// AttachmentRepository 附件记录存储接口
// 由使用方（service 层）注入实现，避免 zstorager 直接依赖数据访问层（internal/dao、internal/model）
type AttachmentRepository interface {
	// HasFile 按 md5 查询附件记录，不存在返回 nil
	HasFile(ctx context.Context, md5 string) (*Attachment, error)
	// Create 写入附件记录并返回自增 ID
	Create(ctx context.Context, a *Attachment) (int64, error)
}

// FileMeta 文件元数据
type FileMeta struct {
	Filename  string `json:"filename"`  // 文件名称
	Size      int64  `json:"size"`      // 文件大小
	Kind      string `json:"kind"`      // 文件上传类型
	MimeType  string `json:"mimeType"`  // 文件扩展类型
	NaiveType string `json:"naiveType"` // NaiveUI类型
	Ext       string `json:"ext"`       // 文件扩展名
	Md5       string `json:"md5"`       // 文件hash
}

// MultipartProgress 分片进度
type MultipartProgress struct {
	UploadId      string      `json:"uploadId"`      // 上传事件ID
	ThirdUploadId string      `json:"thirdUploadId"` // 第三方上传事件ID
	Meta          *FileMeta   `json:"meta"`          // 文件元数据
	ShardCount    int         `json:"shardCount"`    // 分片数量
	UploadedIndex []int       `json:"uploadedIndex"` // 已上传的分片索引
	CreatedAt     *gtime.Time `json:"createdAt"`     // 创建时间
}

// CheckMultipartParams 检查文件分片
type CheckMultipartParams struct {
	UploadType string `json:"uploadType"  dc:"文件类型"`
	FileName   string `json:"fileName"    dc:"文件名称"`
	Size       int64  `json:"size"        dc:"文件大小"`
	Md5        string `json:"md5"         dc:"文件md5值"`
	ShardCount int    `json:"shardCount"  dc:"分片数量"`
	meta       *FileMeta
}

type CheckMultipartModel struct {
	UploadId        string      `json:"uploadId"        dc:"上传事件ID"`
	Attachment      *Attachment `json:"attachment"      dc:"附件"`
	WaitUploadIndex []int       `json:"waitUploadIndex" dc:"等待上传的分片索引"`
	Progress        float64     `json:"progress"        dc:"上传进度"`
	SizeFormat      string      `json:"sizeFormat"      dc:"文件大小"`
}

// UploadPartParams 分片上传
type UploadPartParams struct {
	UploadId   string            `json:"uploadId"    dc:"上传事件ID"`
	UploadType string            `json:"uploadType"  dc:"文件类型"`
	FileName   string            `json:"fileName"    dc:"文件名称"`
	Size       int64             `json:"size"        dc:"文件大小"`
	Md5        string            `json:"md5"         dc:"文件md5值"`
	Index      int               `json:"index"       dc:"分片索引"`
	File       *ghttp.UploadFile `json:"file" type:"file" dc:"分片文件"`
	mp         *MultipartProgress
}

type UploadPartModel struct {
	Attachment *Attachment `json:"attachment" dc:"附件"`
	Progress   float64     `json:"progress"   dc:"上传进度"`
	Finish     bool        `json:"finish"     dc:"是否完成"`
}
