package appservice

import (
	"context"
	"go-mall/api/reply"
	"go-mall/api/request"
	"go-mall/common/errcode"
	"go-mall/common/util"
	"go-mall/logic/do"
	"go-mall/logic/domainservice"
)

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

func (das *DemoAppSvc) GetDemoIdentities() ([]int64, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]int64, 0, len(demos))

	for _, demo := range demos {
		identities = append(identities, demo.Id)
	}
	return identities, nil
}

func (das *DemoAppSvc) CreateDemoOrder(orderRequesst *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	err := util.CopyProperties(demoOrderDo, orderRequesst)
	if err != nil {
		errcode.Wrap("请求转换成demoOrderDo失败", err)
		return nil, err
	}

	demoOrderDo, err = das.demoDomainSvc.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// 做一些创建订单成功后的其他业务逻辑
	// 比如异步发送创建订单成功的通知

	replyDemoOrder := new(reply.DemoOrder)
	err = util.CopyProperties(replyDemoOrder, demoOrderDo)
	if err != nil {
		errcode.Wrap("demoOrderDo转换成replyDemoOrder失败", err)
		return nil, err
	}

	return replyDemoOrder, nil
}
