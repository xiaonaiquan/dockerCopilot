// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package container

import (
	"context"

	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"dockerCopilot/internal/utiles"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.IdReq) (resp *types.Resp, err error) {
	// 调用工具函数删除容器
	err = utiles.DeleteContainer(l.svcCtx, req.Id, true)
	if err != nil {
		l.Logger.Errorf("delete container failed: %v", err)
		return &types.Resp{
			Code: 1,
			Msg:  "删除容器失败: " + err.Error(),
			Data: nil,
		}, err
	}
	return &types.Resp{
		Code: 0,
		Msg:  "删除容器成功",
		Data: nil,
	}, nil
}
