package contact

import (
	"io"
)

// UploadKey 代表 上传的服务
const UploadKey = "hade:upload"

// UploadService 表示传入的参数
type UploadService interface {
	Upload(profile, fileName string, reader *io.Reader) error
}
