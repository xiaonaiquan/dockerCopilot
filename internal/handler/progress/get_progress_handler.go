package progress

import (
	"net/http"

	"dockerCopilot/internal/logic/progress"
	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetProgressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetProgressReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := progress.NewGetProgressLogic(r.Context(), svcCtx)
		resp, err := l.GetProgress(&req)
		if err != nil {
			httpx.WriteJson(w, resp.Code, resp)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
