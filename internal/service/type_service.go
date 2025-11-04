package service

import (
	"analytics/internal/repository"
	"context"
	"fmt"
	"time"
)

type TypeService struct {
	typeRepo *repository.TypeRepository
}

func NewTypeService(typeRepo *repository.TypeRepository) *TypeService {
	return &TypeService{typeRepo: typeRepo}
}

func (r *TypeService) GetAverageSpendByType(ctx context.Context) ([]AverageCategorySpendByMonth, error) {

}
