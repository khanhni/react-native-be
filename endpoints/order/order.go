package order

import (
	"context"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-kit/kit/endpoint"
	"github.com/Lycheeeeeee/react-native-be/domain"
	"github.com/Lycheeeeeee/react-native-be/service"
)

func GetAllByDate(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.OrderDate)
		res, err := s.OrderService.GetAllByDate(ctx, req)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return make([]domain.Order, 0), nil
			}

			return nil, err
		}

		return res, nil
	}
}

func Create(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.Order)
		now := time.Now()
		if req.ReceiveTime != nil {
			receiveTimeStr := strconv.Itoa(*req.ReceiveTime) + "m"
			receiveMinute, err := time.ParseDuration(receiveTimeStr)
			if err != nil {
				return nil, err
			}

			req.OrderTime = now.Add(receiveMinute * time.Nanosecond)
		}
		res, err := s.OrderService.Create(ctx, &req)
		if err != nil {
			return nil, err
		}
		for i := range req.Details {
			detail, err := s.DetailService.Create(ctx, &req.Details[i])
			if err != nil {
				return nil, err
			}
			req.Details[i] = *detail
		}

		return res, nil
	}
}

func Get(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		res, err := s.OrderService.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func Update(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.Order)
		oldOrder, err := s.OrderService.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		if !req.OrderTime.IsZero() {
			oldOrder.OrderTime = req.OrderTime
		}
		if req.ReceiveTime != nil {
			oldOrder.ReceiveTime = req.ReceiveTime
		}
		if !req.AccountID.IsZero() {
			oldOrder.AccountID = req.AccountID
		}
		if !req.DetailID.IsZero() {
			oldOrder.DetailID = req.DetailID
		}
		if !req.ShopID.IsZero() {
			oldOrder.ShopID = req.ShopID
		}
		if req.Status != nil {
			oldOrder.Status = req.Status
		}

		res, err := s.OrderService.Update(ctx, oldOrder)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func GetByShopID(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		return s.OrderService.GetByShopID(ctx, req.ID)
	}
}

func GetByMonth(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MonthRequest)
		return s.OrderService.GetByMonth(ctx, req.Month)
	}
}
