package main

import (
	"time"
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewGorm(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open("postgres", db)
}

// Типы выдачи займов
type PaymentType struct {
	gorm.Model

	// Идентификатор типа выдачи
	ID uint `gorm:"AUTO_INCREMENT"`
	// Название типа
	Title string
}

// Валюты
type Currency struct {
	gorm.Model

	// Идентификатор валюы
	ID uint `gorm:"AUTO_INCREMENT"`
	// Название валюты
	Title string
	// ISO код валюты
	ISO string
	// Лигатура
	Ligature string
	// Текущий внутренний курс
	CrossRate float64
}

// Страны
type Country struct {
	gorm.Model
	
	// Идентификатор страны
	ID uint `gorm:"AUTO_INCREMENT"`
	// Название страны
	Title string
}

// Города
type City struct {
	gorm.Model

	// Идентификатор города
	ID uint `gorm:"AUTO_INCREMENT"`
	// Название города
	Title string
	// Страна
	Country Country
}

// Заказчики (пользователи, которые запрашивают кредит)
type Customer struct {
	gorm.Model

	// Идентификатор заказчика
	ID uint `gorm:"AUTO_INCREMENT"`
	// Имя
	FirstName string
	// Фамилия
	LastName string
	// Отчество
	MiddleName string
	// Дата рождения
	BirthDate time.Time
	// Место рождения
	BirthPlace string
	// Вид документа
	DocumentType DocumentType
	// Дата выдачи документа
	DocumentIssueDate time.Time
	// Кем выдан документ
	DocumentIssuedBy string
	// Телефон
	Phone string
	// Email
	Email string
	// Адрес
	Address string
	// Почтовый индекс
	Zip string
	// Город
	City City
}

// Запросы на выдачу кредита (заказы)
type Order struct {
	gorm.Model

	// Идентификатор заказа
	ID uint `gorm:"AUTO_INCREMENT"`
	// Желаемая сумма
	DesiredAmount float64
	// Заказчик
	Customer Customer
	// Валюа
	Currency Currency
	// Статус заказа
	Status OrderStatus
	// Дата-время создания заказа
	CreatedAt time.Time
	// Дата-время завершения размещения заказа
	PlacedAt time.Time
	// Дата-время завершения обработки заказа НФО
	CompletedAt time.Time
}

// Клиенты (НФО)
type Client struct {
	gorm.Model

	// Идентификатор клиента
	ID uint `gorm:"AUTO_INCREMENT"`
	// Название клиента
	Title string
	// Типы выдачи займов
	PaymentTypes []PaymentType
	// Города, с которыми работает клиент
	ServiceCities []City
	// Обработчик API-ответов
	Handler ClientHandler `gorm:"-"`
}

// Заявка в отдельное НФО
type Bid struct {
	gorm.Model

	// Идентификатор заявки
	ID uint `gorm:"AUTO_INCREMENT"`
	// Заказ
	Order Order
	// НФО
	Client Client
	// Хэш-код заказа на стороне НФО
	Reference string
	// Комментарий к заявки от НФО
	Notice string
	// Статус заявки
	Status BidStatus
}

// Предложенные пакеты (заказчику от НФО)
type OfferedPackage struct {
	gorm.Model

	// Идентификатор предложения
	ID uint `gorm:"AUTO_INCREMENT"`
	// Одобренная сумма займа
	ApprovedAmount float64
	// Валюта займа
	Currency Currency
	// Название пакета
	Title string
	// Тип выдачи займа
	PaymentType PaymentType
	// Заявка на кредит
	Bid Bid
	// Город выдачи займа (может быть nil)
	City *City
}