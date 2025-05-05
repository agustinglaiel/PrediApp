package service

import (
	"context"
	"posts/internal/dto"
	"posts/internal/model"
	"posts/internal/repository"
	e "posts/pkg/utils"
	"time"
)

type PostService interface {
    CreatePost(ctx context.Context, request dto.PostCreateRequestDTO) (dto.PostResponseDTO, e.ApiError)
    GetPostByID(ctx context.Context, id int) (dto.PostResponseDTO, e.ApiError)
    GetPosts(ctx context.Context) ([]dto.PostResponseDTO, e.ApiError)
    GetPostsByUserID(ctx context.Context, userID int) ([]dto.PostResponseDTO, e.ApiError)
    DeletePostByID(ctx context.Context, id int) e.ApiError
}

type postService struct {
    postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
    return &postService{postRepo: postRepo}
}

func (s *postService) CreatePost(ctx context.Context, request dto.PostCreateRequestDTO) (dto.PostResponseDTO, e.ApiError) {
    if request.UserID <= 0 {
        return dto.PostResponseDTO{}, e.NewBadRequestApiError("user_id is required and must be greater than 0")
    }

    // Validar que el parent_post_id exista, si se proporciona
    if request.ParentPostID != nil {
        _, apiErr := s.postRepo.GetPostByID(ctx, *request.ParentPostID)
        if apiErr != nil {
            return dto.PostResponseDTO{}, e.NewBadRequestApiError("parent post not found")
        }
    }

    newPost := &model.Post{
        UserID:       request.UserID,
        ParentPostID: request.ParentPostID,
        Body:         request.Body,
        CreatedAt:    time.Now(),
    }

    if err := s.postRepo.CreatePost(ctx, newPost); err != nil {
        return dto.PostResponseDTO{}, err
    }

    response := dto.PostResponseDTO{
        ID:           newPost.ID,
        UserID:       newPost.UserID,
        ParentPostID: newPost.ParentPostID,
        Body:         newPost.Body,
        CreatedAt:    newPost.CreatedAt.Format(time.RFC3339),
    }

    return response, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (dto.PostResponseDTO, e.ApiError) {
    post, apiErr := s.postRepo.GetPostByID(ctx, id)
    if apiErr != nil {
        return dto.PostResponseDTO{}, apiErr
    }

    response := s.mapPostToResponseDTO(post)
    return response, nil
}

func (s *postService) GetPosts(ctx context.Context) ([]dto.PostResponseDTO, e.ApiError) {
    posts, apiErr := s.postRepo.GetPosts(ctx)
    if apiErr != nil {
        return nil, apiErr
    }

    var response []dto.PostResponseDTO
    for _, post := range posts {
        response = append(response, s.mapPostToResponseDTO(post))
    }
    return response, nil
}

func (s *postService) GetPostsByUserID(ctx context.Context, userID int) ([]dto.PostResponseDTO, e.ApiError) {
    posts, apiErr := s.postRepo.GetPostsByUserID(ctx, userID)
    if apiErr != nil {
        return nil, apiErr
    }

    var response []dto.PostResponseDTO
    for _, post := range posts {
        response = append(response, s.mapPostToResponseDTO(post))
    }
    return response, nil
}

func (s *postService) DeletePostByID(ctx context.Context, id int) e.ApiError {
    return s.postRepo.DeletePostByID(ctx, id) // El userID ya no se pasa desde el contexto
}

// mapPostToResponseDTO convierte un modelo Post a un DTO de respuesta
func (s *postService) mapPostToResponseDTO(post *model.Post) dto.PostResponseDTO {
    var children []dto.PostResponseDTO
    for _, child := range post.Children {
        children = append(children, s.mapPostToResponseDTO(child))
    }

    return dto.PostResponseDTO{
        ID:           post.ID,
        UserID:       post.UserID,
        ParentPostID: post.ParentPostID,
        Body:         post.Body,
        CreatedAt:    post.CreatedAt.Format(time.RFC3339),
        Children:     children,
    }
}