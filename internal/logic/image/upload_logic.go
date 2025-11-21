package image

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(r *http.Request, req *types.UploadReq) (resp *types.Resp, err error) {
	l.Infof("Upload request: %v", r)
	resp = &types.Resp{}
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		resp.Code = 400
		resp.Msg = "请求不合法"
		resp.Data = map[string]interface{}{}
		return resp, err
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		resp.Code = 400
		resp.Msg = "未选择文件"
		resp.Data = map[string]interface{}{}
		return resp, err
	}
	defer file.Close()

	if header == nil || header.Filename == "" {
		resp.Code = 400
		resp.Msg = "文件名无效"
		resp.Data = map[string]interface{}{}
		return resp, errors.New("invalid filename")
	}

	maxSize := int64(10 << 20)
	if header.Size > 0 && header.Size > maxSize {
		resp.Code = 400
		resp.Msg = "文件过大"
		resp.Data = map[string]interface{}{}
		return resp, errors.New("file too large")
	}

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	allowedTypes := map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
		"image/gif":  {},
		"image/webp": {},
		"image/bmp":  {},
	}
	if _, ok := allowedTypes[contentType]; !ok {
		resp.Code = 400
		resp.Msg = "仅允许上传图片"
		resp.Data = map[string]interface{}{}
		return resp, errors.New("unsupported content type")
	}

	basePath := os.Getenv("UPLOAD_DIR")
	if basePath == "" {
		basePath = "/data/uploads"
	}
	if err = os.MkdirAll(basePath, 0755); err != nil {
		resp.Code = 500
		resp.Msg = "服务器错误"
		resp.Data = map[string]interface{}{}
		return resp, err
	}

	name := req.Name
	if name == "" {
		name = "image"
	}
	ext := strings.ToLower(filepath.Ext(name))
	allowedExt := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
		".webp": {},
		".bmp":  {},
	}
	if _, ok := allowedExt[ext]; !ok {
		resp.Code = 400
		resp.Msg = "仅允许上传图片"
		resp.Data = map[string]interface{}{}
		return resp, errors.New("unsupported extension")
	}

	dstPath := filepath.Join(basePath, name)

	out, err := os.Create(dstPath)
	if err != nil {
		resp.Code = 500
		resp.Msg = "服务器错误"
		resp.Data = map[string]interface{}{}
		return resp, err
	}
	defer out.Close()

	if _, err = out.Write(buf[:n]); err != nil {
		resp.Code = 500
		resp.Msg = "服务器错误"
		resp.Data = map[string]interface{}{}
		return resp, err
	}
	if _, err = io.Copy(out, file); err != nil {
		resp.Code = 500
		resp.Msg = "服务器错误"
		resp.Data = map[string]interface{}{}
		return resp, err
	}

	resp.Code = 200
	resp.Msg = "success"
	resp.Data = map[string]interface{}{"filename": name, "path": dstPath, "contentType": contentType}
	return resp, nil
}
