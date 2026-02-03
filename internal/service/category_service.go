package service

import (
	"analytics/internal/repository"
	"context"
	"fmt"
)

type AverageCategory struct {
	CategoryID   int
	CategoryName string
	Average      float64
}

type CategoryService struct {
	categoryRepo    repository.CategoryRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
}

func NewCategoryService(categoryRepo repository.CategoryRepositoryInterface, transactionRepo repository.TransactionRepositoryInterface) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo, transactionRepo: transactionRepo}
}

func (r *CategoryService) GetAverageByCategory(ctx context.Context) ([]AverageCategory, error) {
	transactions, err := r.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	categories, err := r.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[int]string)
	for _, c := range categories {
		categoryMap[c.ID] = c.Name
	}

	monthlySumsByCategory := make(map[string]float64)

	for _, tx := range transactions {
		monthKey := fmt.Sprintf("%d-%d-%d",
			tx.CategoryID,
			tx.Date.Year(),
			tx.Date.Month())

		monthlySumsByCategory[monthKey] += tx.Amount
	}

	monthlySumsByCategoryID := make(map[int][]float64)

	for key, sum := range monthlySumsByCategory {
		var categoryID, year, month int
		fmt.Sscanf(key, "%d-%d-%d", &categoryID, &year, &month)

		monthlySumsByCategoryID[categoryID] = append(monthlySumsByCategoryID[categoryID], sum)
	}

	var result []AverageCategory
	for categoryID, monthlySums := range monthlySumsByCategoryID {
		var total float64
		for _, sum := range monthlySums {
			total += sum
		}
		average := total / float64(len(monthlySums))

		result = append(result, AverageCategory{
			CategoryID:   categoryID,
			CategoryName: categoryMap[categoryID],
			Average:      average,
		})
	}

	return result, nil
}
