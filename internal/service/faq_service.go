package service

import (
	"context"
	"fmt"

	"github.com/lucsky/cuid"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type IFaqService interface {
	CreateFAQ(ctx context.Context, question, answer, category string, orderIndex int) (*dto.FAQItem, error)
	UpdateFAQ(ctx context.Context, id, question, answer, category string, orderIndex int, isActive bool) (*dto.FAQItem, error)
	DeleteFAQ(ctx context.Context, id string) error
	ListAllFAQs(ctx context.Context, page, limit int) ([]dto.FAQItem, int64, error)
	ListPublicFAQs(ctx context.Context) ([]dto.FAQItem, error)
	GetFAQByID(ctx context.Context, id string) (*dto.FAQItem, error)
}

type FaqService struct {
	faqRepo *repository.FAQRepository
}

func NewFaqService(faqRepo *repository.FAQRepository) IFaqService {
	return &FaqService{faqRepo: faqRepo}
}

func (s *FaqService) CreateFAQ(ctx context.Context, question, answer, category string, orderIndex int) (*dto.FAQItem, error) {
	if question == "" || answer == "" {
		return nil, domain.ErrInvalidInput
	}

	if category == "" {
		category = "general"
	}

	faq := &domain.FAQ{
		ID:         cuid.New(),
		Question:   question,
		Answer:     answer,
		Category:   category,
		OrderIndex: orderIndex,
		IsActive:   true,
	}

	if err := s.faqRepo.Create(ctx, faq); err != nil {
		return nil, fmt.Errorf("create faq: %w", err)
	}

	item := dto.ToFAQItem(*faq)
	return &item, nil
}

func (s *FaqService) UpdateFAQ(ctx context.Context, id, question, answer, category string, orderIndex int, isActive bool) (*dto.FAQItem, error) {
	faq, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if question != "" {
		faq.Question = question
	}
	if answer != "" {
		faq.Answer = answer
	}
	if category != "" {
		faq.Category = category
	}
	faq.OrderIndex = orderIndex
	faq.IsActive = isActive

	if err := s.faqRepo.Update(ctx, faq); err != nil {
		return nil, fmt.Errorf("update faq: %w", err)
	}

	item := dto.ToFAQItem(*faq)
	return &item, nil
}

func (s *FaqService) DeleteFAQ(ctx context.Context, id string) error {
	_, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}
	return s.faqRepo.Delete(ctx, id)
}

func (s *FaqService) ListAllFAQs(ctx context.Context, page, limit int) ([]dto.FAQItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	faqs, total, err := s.faqRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("list faqs: %w", err)
	}

	items := make([]dto.FAQItem, len(faqs))
	for i, f := range faqs {
		items[i] = dto.ToFAQItem(f)
	}
	return items, total, nil
}

func (s *FaqService) ListPublicFAQs(ctx context.Context) ([]dto.FAQItem, error) {
	faqs, err := s.faqRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list public faqs: %w", err)
	}

	items := make([]dto.FAQItem, len(faqs))
	for i, f := range faqs {
		items[i] = dto.ToFAQItem(f)
	}
	return items, nil
}

func (s *FaqService) GetFAQByID(ctx context.Context, id string) (*dto.FAQItem, error) {
	faq, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	item := dto.ToFAQItem(*faq)
	return &item, nil
}
