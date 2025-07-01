package repository

import (
	"context"

	"prediapp.local/db/model"
	e "prediapp.local/posts/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) e.ApiError
	GetPostByID(ctx context.Context, id int) (*model.Post, e.ApiError)
	GetPosts(ctx context.Context, offset, limit int) ([]*model.Post, e.ApiError)
	GetPostsByUserID(ctx context.Context, userID int) ([]*model.Post, e.ApiError)
	DeletePostByID(ctx context.Context, id int) e.ApiError // Eliminar el userID como argumento
	SearchPosts(ctx context.Context, query string, offset, limit int) ([]*model.Post, e.ApiError)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) e.ApiError {
	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		return e.NewInternalServerApiError("error creating post", err)
	}
	return nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*model.Post, e.ApiError) {
	var post model.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("post not found")
		}
		return nil, e.NewInternalServerApiError("error finding post", err)
	}
	post.Children = r.getChildPosts(ctx, post.ID)
	return &post, nil
}

func (r *postRepository) GetPosts(ctx context.Context, offset, limit int) ([]*model.Post, e.ApiError) {
	var posts []*model.Post
	if err := r.db.WithContext(ctx).
		Where("parent_post_id IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, e.NewInternalServerApiError("error finding posts", err)
	}
	for _, post := range posts {
		post.Children = r.getChildPosts(ctx, post.ID)
	}
	return posts, nil
}

func (r *postRepository) GetPostsByUserID(ctx context.Context, userID int) ([]*model.Post, e.ApiError) {
	var posts []*model.Post
	if err := r.db.WithContext(ctx).Where("user_id = ? AND parent_post_id IS NULL", userID).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, e.NewInternalServerApiError("error finding posts by user", err)
	}
	for _, post := range posts {
		post.Children = r.getChildPosts(ctx, post.ID)
	}
	return posts, nil
}

func (r *postRepository) DeletePostByID(ctx context.Context, id int) e.ApiError {
	var post model.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return e.NewNotFoundApiError("post not found")
		}
		return e.NewInternalServerApiError("error finding post", err)
	}

	// Aquí no validamos user_id porque viene del cuerpo en el controller
	if err := r.db.WithContext(ctx).Delete(&post).Error; err != nil {
		return e.NewInternalServerApiError("error deleting post", err)
	}
	return nil
}

func (r *postRepository) SearchPosts(ctx context.Context, query string, offset, limit int) ([]*model.Post, e.ApiError) {
	var posts []*model.Post
	db := r.db.WithContext(ctx)

	// WHERE fulltext
	db = db.Where("MATCH(body) AGAINST(? IN BOOLEAN MODE)", query)

	// ORDER BY relevancia usando clause.Expr
	db = db.Order(clause.Expr{
		SQL:  "MATCH(body) AGAINST(? IN BOOLEAN MODE) DESC",
		Vars: []interface{}{query},
	})

	// Paginación y ejecución
	if err := db.
		Offset(offset).
		Limit(limit).
		Find(&posts).
		Error; err != nil {
		return nil, e.NewInternalServerApiError("error searching posts", err)
	}

	return posts, nil
}

func (r *postRepository) getChildPosts(ctx context.Context, parentID int) []*model.Post {
	var children []*model.Post
	r.db.WithContext(ctx).Where("parent_post_id = ?", parentID).Order("created_at ASC").Find(&children)
	for _, child := range children {
		child.Children = r.getChildPosts(ctx, child.ID)
	}
	return children
}
