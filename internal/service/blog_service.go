package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type IBlogService interface {
	CreateBlog(ctx context.Context, userID string, req CreateBlogRequest) (*dto.BlogDetail, error)
	ListBlogs(ctx context.Context, page, limit int, search, categoryID, tag string) ([]dto.BlogListItem, *dto.Pagination, error)
	ListMyBlogs(ctx context.Context, userID string, page, limit int, status string) ([]dto.BlogListItem, *dto.Pagination, error)
	GetBlogByID(ctx context.Context, id string) (*dto.BlogDetail, error)
	UpdateBlog(ctx context.Context, blogID, userID string, req UpdateBlogRequest) error
	AdminUpdateBlog(ctx context.Context, blogID string, req UpdateBlogRequest) error
	DeleteBlog(ctx context.Context, blogID, userID string) error
	AdminDeleteBlog(ctx context.Context, blogID string) error
}

type CreateBlogRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	CoverImgURL string   `json:"cover_img_url"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	CategoryID  string   `json:"category_id"`
}

type UpdateBlogRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Content     *string  `json:"content,omitempty"`
	CoverImgURL *string  `json:"cover_img_url,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Status      *string  `json:"status,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

type BlogService struct {
	blogRepo repository.IBlogRepository
}

func NewBlogService(blogRepo repository.IBlogRepository) IBlogService {
	return &BlogService{blogRepo: blogRepo}
}

func makeSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 100 {
		slug = slug[:100]
	}
	return slug
}

func (s *BlogService) CreateBlog(ctx context.Context, userID string, req CreateBlogRequest) (*dto.BlogDetail, error) {
	blog := &domain.Blog{
		ID:          cuid.New(),
		Title:       req.Title,
		Slug:        makeSlug(req.Title),
		Description: req.Description,
		Content:     req.Content,
		CoverImgURL: req.CoverImgURL,
		Status:      req.Status,
		UserID:      userID,
		CategoryID:  req.CategoryID,
	}

	if blog.Status == "" {
		blog.Status = "draft"
	}

	if len(req.Tags) > 0 {
		tagsJSON, _ := json.Marshal(req.Tags)
		blog.Tags = tagsJSON
	}

	if err := s.blogRepo.Create(ctx, blog); err != nil {
		return nil, fmt.Errorf("create blog: %w", err)
	}

	created, err := s.blogRepo.FindByID(ctx, blog.ID)
	if err != nil {
		return nil, fmt.Errorf("find created blog: %w", err)
	}

	return dto.ToBlogDetail(created), nil
}

func (s *BlogService) ListBlogs(ctx context.Context, page, limit int, search, categoryID, tag string) ([]dto.BlogListItem, *dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	blogs, total, err := s.blogRepo.ListAll(ctx, page, limit, search, "published", categoryID, tag)
	if err != nil {
		return nil, nil, fmt.Errorf("list blogs: %w", err)
	}

	items := make([]dto.BlogListItem, len(blogs))
	for i, b := range blogs {
		items[i] = dto.ToBlogListItem(b)
	}

	return items, &dto.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (s *BlogService) ListMyBlogs(ctx context.Context, userID string, page, limit int, status string) ([]dto.BlogListItem, *dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	blogs, total, err := s.blogRepo.ListByUserID(ctx, userID, page, limit, status)
	if err != nil {
		return nil, nil, fmt.Errorf("list my blogs: %w", err)
	}

	items := make([]dto.BlogListItem, len(blogs))
	for i, b := range blogs {
		items[i] = dto.ToBlogListItem(b)
	}

	return items, &dto.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (s *BlogService) GetBlogByID(ctx context.Context, id string) (*dto.BlogDetail, error) {
	blog, err := s.blogRepo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get blog: %w", err)
	}
	return dto.ToBlogDetail(blog), nil
}

func (s *BlogService) UpdateBlog(ctx context.Context, blogID, userID string, req UpdateBlogRequest) error {
	blog, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("find blog for update: %w", err)
	}

	if blog.UserID != userID {
		return domain.ErrForbidden
	}

	applyUpdates(blog, req)
	return s.blogRepo.Update(ctx, blog)
}

// AdminUpdateBlog allows admin to update any blog without ownership check
func (s *BlogService) AdminUpdateBlog(ctx context.Context, blogID string, req UpdateBlogRequest) error {
	blog, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("find blog for admin update: %w", err)
	}

	applyUpdates(blog, req)
	return s.blogRepo.Update(ctx, blog)
}

func (s *BlogService) DeleteBlog(ctx context.Context, blogID, userID string) error {
	blog, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("find blog for delete: %w", err)
	}

	if blog.UserID != userID {
		return domain.ErrForbidden
	}

	return s.blogRepo.Delete(ctx, blogID)
}

// AdminDeleteBlog allows admin to delete any blog without ownership check
func (s *BlogService) AdminDeleteBlog(ctx context.Context, blogID string) error {
	_, err := s.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		return fmt.Errorf("find blog for admin delete: %w", err)
	}

	return s.blogRepo.Delete(ctx, blogID)
}

func applyUpdates(blog *domain.Blog, req UpdateBlogRequest) {
	if req.Title != nil {
		blog.Title = *req.Title
		blog.Slug = makeSlug(*req.Title)
	}
	if req.Description != nil {
		blog.Description = *req.Description
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.CoverImgURL != nil {
		blog.CoverImgURL = *req.CoverImgURL
	}
	if req.Status != nil {
		blog.Status = *req.Status
	}
	if req.CategoryID != nil {
		blog.CategoryID = *req.CategoryID
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		blog.Tags = tagsJSON
	}
}
