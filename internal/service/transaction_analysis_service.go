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
	recurringRepo   *repository.RecurringTransactionRepository
	categoryRepo    *repository.CategoryRepository
	typeRepo        *repository.TypeRepository
}

func NewTransactionAnalysisService(
	transactionRepo *repository.TransactionRepository,
	recurringRepo *repository.RecurringTransactionRepository,
	categoryRepo *repository.CategoryRepository,
	typeRepo *repository.TypeRepository,
) *TransactionAnalysisService {
	return &TransactionAnalysisService{
		transactionRepo: transactionRepo,
		recurringRepo:   recurringRepo,
		categoryRepo:    categoryRepo,
		typeRepo:        typeRepo,
	}
}

func (r *TransactionAnalysisService) GetAverageSpendByCategory(ctx context.Context) ([]AverageCategorySpendByMonth, error) {
	transactions, err := r.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	// recurringTransactions, err := r.recurringRepo.GetAllRecurringTransactions(ctx)
	// if err != nil {
	// 	return nil, err
	// }

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

	processRegularTransactions(transactions, spendByMonth)

	// Process recurring transactions
	// processRecurringTransactions(recurringTransactions, spendByMonth)

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

func processRegularTransactions(transactions []domain.Transaction, spendByMonth map[string]struct {
	Total float64
	Count int
}) {
	for _, tx := range transactions {
		if tx.TypeID != 3 {
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

func processRecurringTransactions(transactions []domain.RecurringTransaction, spendByMonth map[string]struct {
	Total float64
	Count int
}) {
	for _, tx := range transactions {
		if tx.TypeID != 3 {
			continue
		}

		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.CategoryID,
			time.Now().Year(),
			time.Now().Month())

		monthly := spendByMonth[monthKey]
		monthly.Total += tx.Amount
		monthly.Count++
		spendByMonth[monthKey] = monthly
	}
}
