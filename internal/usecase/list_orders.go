package usecase

import "github.com/willcmarques/goexpert-desafio-clean-arch/internal/entity"

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.List()
	if err != nil {
		return []OrderOutputDTO{}, err
	}

	var ordersOutput []OrderOutputDTO
	for _, order := range orders {
		orderOutput := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		ordersOutput = append(ordersOutput, orderOutput)
	}
	return ordersOutput, nil
}
