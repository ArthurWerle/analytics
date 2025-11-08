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

	spendByTypeAndMonth := make(map[string]struct {
		Total float64
		Count int
	})

	for _, tx := range transactions {
		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.TypeID,
			tx.Date.Year(),
			tx.Date.Month())

		stats := spendByTypeAndMonth[monthKey]
		stats.Total += tx.Amount
		stats.Count++
		spendByTypeAndMonth[monthKey] = stats
	}

	for _, tx := range recurringTransactions {
		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.TypeID,
			tx.StartDate.Year(),
			tx.StartDate.Month())

		stats := spendByTypeAndMonth[monthKey]
		stats.Total += tx.Amount
		stats.Count++
		spendByTypeAndMonth[monthKey] = stats
	}

	monthlyAveragesByType := make(map[int][]float64)

	for key, stats := range spendByTypeAndMonth {
		var typeID, year, month int
		fmt.Sscanf(key, "%d-%d-%d", &typeID, &year, &month)

		monthlyAverage := stats.Total / float64(stats.Count)
		monthlyAveragesByType[typeID] = append(monthlyAveragesByType[typeID], monthlyAverage)
	}

	var result []AverageType
	for typeID, monthlyAverages := range monthlyAveragesByType {
		var sum float64
		for _, avg := range monthlyAverages {
			sum += avg
		}
		overallAverage := sum / float64(len(monthlyAverages))

		result = append(result, AverageType{
			TypeID:   typeID,
			TypeName: typeMap[typeID],
			Average:  overallAverage,
		})
	}

	return result, nil
}
