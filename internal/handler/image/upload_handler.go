// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package image

import (
	"net/http"

	"dockerCopilot/internal/logic/image"
	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := image.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(r, &req)
		if err != nil {
			httpx.WriteJson(w, resp.Code, resp)
		} else {
			httpx.WriteJson(w, resp.Code, resp)
		}
	}
}
