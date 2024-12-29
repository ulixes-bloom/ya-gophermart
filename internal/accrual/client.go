package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-gophermart/internal/config"
	appErrors "github.com/ulixes-bloom/ya-gophermart/internal/errors"
	"github.com/ulixes-bloom/ya-gophermart/internal/models"
	"github.com/ulixes-bloom/ya-gophermart/internal/workerpool"
)

type Client struct {
	http HTTPClient
	conf *config.Config
}

func NewClient(conf *config.Config) *Client {
	return &Client{
		conf: conf,
		http: &http.Client{},
	}
}

func (ac *Client) GetOrdersInfo(ctx context.Context, orders []models.Order) ([]models.Order, error) {
	wp := workerpool.New(ctx, ac.conf.AccrualRateLimit, ac.conf.AccrualRateLimit*2, ac.GetOrderInfo)

	for _, order := range orders {
		wp.Submit(&order)
	}
	wp.StopAndWait()

	resOrders := []models.Order{}
	for order := range wp.Results() {
		resOrders = append(resOrders, *order)
	}

	for err := range wp.Errors() {
		log.Warn().Msg(err.Error())
	}

	return resOrders, nil
}

func (ac *Client) GetOrderInfo(ctx context.Context, order *models.Order) (*models.Order, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		ac.conf.NormilizedAccrualSysAddr()+"/api/orders/"+order.Number,
		nil)
	if err != nil {
		return nil, fmt.Errorf("accrual.getOrderInfo: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := ac.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("accrual.getOrderInfo.doRequest: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		accrualResp := &models.AccrualResponse{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(accrualResp); err != nil {
			return nil, fmt.Errorf("accrual.getOrderInfo: %w", err)
		}

		order := &models.Order{
			UserID:  order.UserID,
			Number:  accrualResp.OrderNumber,
			Status:  mapAccrualResponseStatus(accrualResp.AccrualStatus),
			Accrual: accrualResp.Accrual,
		}
		return order, nil
	case http.StatusNoContent:
		return nil, errors.Join(appErrors.ErrAccrualOrderNotRegistered, fmt.Errorf("accrual.getOrderInfo: order %s", order.Number))
	case http.StatusTooManyRequests:
		return nil, errors.Join(appErrors.ErrAccrualTooManyRequests, fmt.Errorf("accrual.getOrderInfo: order %s", order.Number))
	default:
		return nil, fmt.Errorf("accrual.getOrderInfo: %d", resp.StatusCode)
	}
}

func mapAccrualResponseStatus(accrualStatus models.AccrualStatus) models.OrderStatus {
	switch accrualStatus {
	case models.AccrualStatusProcessing:
		return models.OrderStatusProcessing
	case models.AccrualStatusRegistered:
		return models.OrderStatusNew
	case models.AccrualStatusProcessed:
		return models.OrderStatusProcessed
	default:
		return models.OrderStatusInvalid
	}
}
