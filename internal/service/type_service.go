package service

import (
	"analytics/internal/repository"
	"context"
	"fmt"
	"log"
)

type AverageType struct {
	Type    string
	Average float64
}

type TypeService struct {
	transactionRepo repository.TransactionRepositoryInterface
}

func NewTypeService(transactionRepo repository.TransactionRepositoryInterface) *TypeService {
	return &TypeService{transactionRepo: transactionRepo}
}

func (r *TypeService) GetAverageByType(ctx context.Context) ([]AverageType, error) {
	transactions, err := r.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		log.Printf("[TypeService.GetAverageByType] ERROR: Failed to fetch transactions: %v", err)
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	monthlySumsByType := make(map[string]float64)

	for _, tx := range transactions {
		monthKey := fmt.Sprintf("%s-%d-%d",
			tx.Type,
			tx.Date.Year(),
			tx.Date.Month())

		monthlySumsByType[monthKey] += tx.Amount
	}

	monthlySumsByTypeID := make(map[string][]float64)

	for key, sum := range monthlySumsByType {
		var typeName string
		var year, month int
		fmt.Sscanf(key, "%s-%d-%d", &typeName, &year, &month)

		monthlySumsByTypeID[typeName] = append(monthlySumsByTypeID[typeName], sum)
	}

	var result []AverageType
	for typeName, monthlySums := range monthlySumsByTypeID {
		var total float64
		for _, sum := range monthlySums {
			total += sum
		}
		average := total / float64(len(monthlySums))

		result = append(result, AverageType{
			Type:    typeName,
			Average: average,
		})
	}

	return result, nil
}
