package service

import (
	"analytics/internal/repository"
	"context"
	"fmt"
)

type AverageType struct {
	TypeID   int
	TypeName string
	Average  float64
}

type TypeService struct {
	typeRepo        repository.TypeRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
	recurringRepo   repository.RecurringTransactionRepositoryInterface
}

func NewTypeService(typeRepo repository.TypeRepositoryInterface, transactionRepo repository.TransactionRepositoryInterface, recurringRepo repository.RecurringTransactionRepositoryInterface) *TypeService {
	return &TypeService{typeRepo: typeRepo, transactionRepo: transactionRepo, recurringRepo: recurringRepo}
}

func (r *TypeService) GetAverageByType(ctx context.Context) ([]AverageType, error) {
	transactions, err := r.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	recurringTransactions, err := r.recurringRepo.GetAllRecurringTransactions(ctx)
	if err != nil {
		return nil, err
	}

	types, err := r.typeRepo.GetAllTypes(ctx)
	if err != nil {
		return nil, err
	}

	typeMap := make(map[int]string)
	for _, t := range types {
		typeMap[t.ID] = string(t.Name)
	}

	monthlySumsByType := make(map[string]float64)

	for _, tx := range transactions {
		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.TypeID,
			tx.Date.Year(),
			tx.Date.Month())

		monthlySumsByType[monthKey] += tx.Amount
	}

	for _, tx := range recurringTransactions {
		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.TypeID,
			tx.StartDate.Year(),
			tx.StartDate.Month())

		monthlySumsByType[monthKey] += tx.Amount
	}

	monthlySumsByTypeID := make(map[int][]float64)

	for key, sum := range monthlySumsByType {
		var typeID, year, month int
		fmt.Sscanf(key, "%d-%d-%d", &typeID, &year, &month)

		monthlySumsByTypeID[typeID] = append(monthlySumsByTypeID[typeID], sum)
	}

	var result []AverageType
	for typeID, monthlySums := range monthlySumsByTypeID {
		var total float64
		for _, sum := range monthlySums {
			total += sum
		}
		average := total / float64(len(monthlySums))

		result = append(result, AverageType{
			TypeID:   typeID,
			TypeName: typeMap[typeID],
			Average:  average,
		})
	}

	return result, nil
}
