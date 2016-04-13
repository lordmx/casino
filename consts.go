package main

// Тип документа
type DocumentType string

// Статус заказа
type OrderStatus string

// Статус заявок
type BidStatus string

const (
	// Паспорт
	Passport DocumentType = "passport"

	// Заказ в обработке
	Processing OrderStatus = "processing"
	// Заказ завершен
	Completed OrderStatus = "completed"
	// Заказ отменен
	Canceled OrderStatus = "canceled"

	// Заявка ждет обработки
	Pending BidStatus = "pending"
	// Заявка успешна
	Success BidStatus = "success"
	// Заявка завершилась по таймауту
	TimedOut BidStatus = "timed_out"
	// НФО ответило ошибкой
	Failure BidStatus = "failure"
)