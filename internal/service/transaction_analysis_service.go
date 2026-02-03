package service

import (
	"analytics/internal/domain"
	"analytics/internal/repository"

	"context"
	"fmt"
	"time"
)

type AverageCategorySpendByMonth struct {
	CategoryID   int
	CategoryName string
	Month        time.Time
	AverageSpend float64
}

type TransactionAnalysisService struct {
	transactionRepo *repository.TransactionRepository
	categoryRepo    *repository.CategoryRepository
}

func NewTransactionAnalysisService(
	transactionRepo *repository.TransactionRepository,
	categoryRepo *repository.CategoryRepository,
) *TransactionAnalysisService {
	return &TransactionAnalysisService{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

func (r *TransactionAnalysisService) GetAverageSpendByCategory(ctx context.Context) ([]AverageCategorySpendByMonth, error) {
	transactions, err := r.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := r.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[int]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	spendByMonth := make(map[string]struct {
		Total float64
		Count int
	})

	processTransaction(transactions, spendByMonth)

	var result []AverageCategorySpendByMonth
	for key, value := range spendByMonth {
		var categoryID int
		var year, month int
		fmt.Sscanf(key, "%d-%d-%d", &categoryID, &year, &month)

		date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		averageSpend := value.Total / float64(value.Count)

		result = append(result, AverageCategorySpendByMonth{
			CategoryID:   categoryID,
			CategoryName: categoryMap[categoryID],
			Month:        date,
			AverageSpend: averageSpend,
		})
	}

	return result, nil
}

func processTransaction(transactions []domain.Transaction, spendByMonth map[string]struct {
	Total float64
	Count int
}) {
	for _, tx := range transactions {
		if tx.Type != domain.Expense {
			continue
		}

		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.CategoryID,
			tx.Date.Year(),
			tx.Date.Month())

		monthly := spendByMonth[monthKey]
		monthly.Total += tx.Amount
		monthly.Count++
		spendByMonth[monthKey] = monthly
	}
}
