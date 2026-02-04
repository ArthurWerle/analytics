package service

import (
	"analytics/internal/repository"
	"context"
	"fmt"
	"log"
)

type AverageType struct {
	TypeName string
	Average  float64
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

	// Group monthly sums by type
	// Key: "type-year-month", Value: sum for that month
	monthlySumsByTypeMonth := make(map[string]float64)

	for _, tx := range transactions {
		if tx.Date == nil {
			continue
		}
		monthKey := fmt.Sprintf("%s-%d-%d",
			tx.Type,
			tx.Date.Year(),
			tx.Date.Month())

		monthlySumsByTypeMonth[monthKey] += tx.Amount
	}

	// Group monthly sums by type only (to calculate average across months)
	monthlySumsByType := make(map[string][]float64)

	for key, sum := range monthlySumsByTypeMonth {
		// Extract type name (everything before the first dash followed by a digit)
		typeName := extractTypeName(key)
		monthlySumsByType[typeName] = append(monthlySumsByType[typeName], sum)
	}

	var result []AverageType
	for typeName, monthlySums := range monthlySumsByType {
		var total float64
		for _, sum := range monthlySums {
			total += sum
		}
		average := total / float64(len(monthlySums))

		result = append(result, AverageType{
			TypeName: typeName,
			Average:  average,
		})
	}

	return result, nil
}

// extractTypeName extracts the type name from a key like "income-2025-10" or "expense-2025-12"
func extractTypeName(key string) string {
	// Find the position where the year starts (first digit after a dash)
	for i := 0; i < len(key); i++ {
		if key[i] == '-' && i+1 < len(key) && key[i+1] >= '0' && key[i+1] <= '9' {
			return key[:i]
		}
	}
	return key
}
