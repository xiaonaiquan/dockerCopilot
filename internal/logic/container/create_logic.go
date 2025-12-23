// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package container

import (
	"context"
	"fmt"

	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"
	"dockerCopilot/internal/utiles"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateReq) (resp *types.Resp, err error) {
	if req.Yaml == "" {
		return nil, fmt.Errorf("yaml内容为空")
	}

	ret, err := utiles.CreateContainerAsync(l.svcCtx, req.Yaml)
	if err != nil {
		return nil, err
	}

	return &types.Resp{
		Code: 0,
		Data: ret,
		Msg:  "所有服务处理完毕",
	}, nil
}
