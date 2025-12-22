package domainservice

import (
	"context"
	"go-mall/common/errcode"
	"go-mall/dal/dao"
	"go-mall/logic/do"
)

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.DemoOrder, 0, len(demos))
	for _, demo := range demos {
		demoOrders = append(demoOrders, &do.DemoOrder{
			Id:           demo.Id,
			UserId:       demo.UserId,
			BillMoney:    demo.BillMoney,
			OrderNo:      demo.OrderNo,
			OrderGoodsId: demo.OrderGoodsId,
			State:        demo.State,
			PaidAt:       demo.PaidAt,
			CreatedAt:    demo.CreatedAt,
			UpdatedAt:    demo.UpdatedAt,
		})
	}

	return demoOrders, nil
}
