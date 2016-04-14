package main

import (
	"time"
)

// Запрос на эксресс-кредит
type ExpressRequest struct {
	// Заказчик
	Customer Customer
	// Заказ
	Order Order
	// НФО, которым будут разосланы запросы
	Clients []Client
	// Дата-время создания запроса
	CreatedAt time.Time
}

// Ответ НФО на запрос экспресс-кредита
type ExpressResponse struct {
	// Заказ
	Order Order
	// НФО
	Client Client
	// Заказчик
	Customer Customer
	// Запрос успешно выполнен
	IsSuccess bool
	// Запрос завершен по таймауту
	IsTimedOut bool
	// Количество попыток выполения запроса
	Times int
	// Дата-время начала формирования ответа
	CreatedAt time.Time
	// Дата-время завершения формирования ответа
	CompletedAt time.Time
	// Номер заказа на стороне НФО
	Reference string
	// Комментарий к заявки от НФО
	Notice string
	// Модель заявки
	Bid *Bid
}

// Запрос на выдачу экспресс-кредита
func (service *ServiceExpress) OrderExpress(customer Customer, order Order) (result []*ExpressResponse) {
	request := &ExpressRequest{
		Customer: customer,
		Order: order,
		Clients: service.Clients,
		CreatedAt: time.Now(),
	}

	result = make([]*ExpressResponse, len(service.Clients))

	for _, response := range request.SendVia(service) {
		service.PersistResponse(response)
		result = append(result, response)
	}

	return result
}

// Сохранить модель ответа (заявка на кредит)
func (service *ServiceExpress) PersistResponse(response *ExpressResponse) bool {
	if response.Bid == nil {
		response.Bid = &Bid{
			Status: Pending,
			Order: response.Order,
			Client: response.Client,
		}
	} else {
		bid := response.Bid
		bid.CompletedAt = time.Now()
		bid.Reference = response.Reference
		bid.Notice = response.Notice

		status := Success

		switch (true) {
		case response.IsTimedOut:
			status = TimedOut
		case !response.IsTimedOut && !response.IsSuccess:
			status = Failure
		}

		bid.Status = status
	}

	service.DB.Save(response.Bid)

	return response.Bid.ID > 0
}

// Выполнен ли запроса на экспресс-кредит
func (request *ExpressRequest) IsCompleted(results map[uint]*ExpressResponse) bool {
	count := 0

	for _, result := range results {
		if result.IsSuccess || result.IsTimedOut {
			count++
		}
	}

	return count == len(request.Clients)
}

// Отправить запрос
func (request *ExpressRequest) SendVia(service *ServiceExpress) (responses map[uint]*ExpressResponse) {
	ch := make(chan *ExpressResponse, len(request.Clients))
	
	for _, client := range request.Clients {
		go func() {
			response := &ExpressResponse{
				Order: request.Order,
				Customer: request.Customer,
				Client: client,
				CreatedAt: time.Now(),
			}

			if service.PersistResponse(response) {
				responses[client.ID] = response
				ch <- client.Handler.SendOrder(request.Customer, request.Order)
			}
		}()
	}

	for {
		select {
		case r := <- ch:
			responses[r.Client.ID].Reference = r.Reference
			responses[r.Client.ID].IsSuccess = r.IsSuccess
			responses[r.Client.ID].Notice = r.Notice
			responses[r.Client.ID].CompletedAt = time.Now()

			if request.IsCompleted(responses) {
				return responses
			}
		case <-time.After(100 * time.Millisecond):
			for _, client := range request.Clients {
				if !responses[client.ID].IsSuccess && !responses[client.ID].IsTimedOut {
					responses[client.ID].Times++
				}

				if responses[client.ID].Times >= 300 {
					responses[client.ID].IsTimedOut = true
					responses[client.ID].CompletedAt = time.Now()
				}
			}
		}
	}

	return responses
}
