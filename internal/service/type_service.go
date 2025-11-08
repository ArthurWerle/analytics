package service

import (
	"analytics/internal/repository"
	"context"
)

type AverageType struct {
	TypeID   int
	TypeName string
	Average  float64
}

type TypeService struct {
	typeRepo        *repository.TypeRepository
	transactionRepo *repository.TransactionRepository
	recurringRepo   *repository.RecurringTransactionRepository
}

func NewTypeService(typeRepo *repository.TypeRepository, transactionRepo *repository.TransactionRepository, recurringRepo *repository.RecurringTransactionRepository) *TypeService {
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

	spendByType := make(map[int]struct {
		Total float64
		Count int
	})

	for _, tx := range transactions {
		stats := spendByType[tx.TypeID]
		stats.Total += tx.Amount
		stats.Count++
		spendByType[tx.TypeID] = stats
	}

	for _, tx := range recurringTransactions {
		stats := spendByType[tx.TypeID]
		stats.Total += tx.Amount
		stats.Count++
		spendByType[tx.TypeID] = stats
	}

	var result []AverageType
	for typeID, stats := range spendByType {
		average := stats.Total / float64(stats.Count)

		result = append(result, AverageType{
			TypeID:   typeID,
			TypeName: typeMap[typeID],
			Average:  average,
		})
	}

	return result, nil
}
