package service

import (
	"context"
	"groups/internal/dto"
	"groups/internal/model"
	"groups/internal/repository"
	e "groups/pkg/utils"
)


type groupService struct {
	groupRepo repository.GroupRepository
}

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, group dto.CreateGroupRequestDTO) (dto.CreateGroupRequestDTO, e.ApiError)
	GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError)
	GetGroups(ctx context.Context) ([]dto.GroupResponseDTO, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError
}

func NewGroupService(groupRepo repository.GroupRepository) GroupServiceInterface {
	return &groupService{groupRepo: groupRepo}
}

func (s *groupService) CreateGroup(ctx context.Context, group dto.CreateGroupRequestDTO) (dto.CreateGroupRequestDTO, e.ApiError) {
	// Check if the groupName already exists
	if _, err := s.groupRepo.GetGroupByGroupName(ctx, group.GroupName); err == nil {
		return dto.CreateGroupRequestDTO{}, e.NewBadRequestApiError("group already exists")
	}

	newGroup := &model.Group{
		GroupName: group.GroupName,
		Description: group.Description,
	}

	if err := s.groupRepo.CreateGroup(ctx, newGroup); err != nil {
		return dto.CreateGroupRequestDTO{}, e.NewInternalServerApiError("error creating group", err)
	}

	response := dto.CreateGroupRequestDTO{
		GroupName: newGroup.GroupName,
		Description: newGroup.Description,
	}

	return response, nil
}

func (s *groupService) GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError) {
	group, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return dto.GroupResponseDTO{}, err
	}

	response := dto.GroupResponseDTO{
		ID: group.ID,
		GroupName: group.GroupName,
		Description: group.Description,
	}

	return response, nil
}

func (s *groupService) GetGroups(ctx context.Context) ([]dto.GroupResponseDTO, e.ApiError) {
	groups, err := s.groupRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.GroupResponseDTO
	for _, group := range groups {
		response = append(response, dto.GroupResponseDTO{
			ID: group.ID,
			GroupName: group.GroupName,
			Description: group.Description,
		})
	}

	return response, nil
}

func (s *groupService) DeleteGroupByID(ctx context.Context, id int) e.ApiError {
	_, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return e.NewBadRequestApiError("group not found")
	}

	if err := s.groupRepo.DeleteGroupByID(ctx, id); err != nil {
		return e.NewInternalServerApiError("error deleting group", err)
	}
	return nil
}

func (s *groupService) DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError {
	_, err := s.groupRepo.GetGroupByGroupName(ctx, groupName)
	if err != nil {
		return e.NewBadRequestApiError("group not found")
	}

	if err := s.groupRepo.DeleteGroupByGroupName(ctx, groupName); err != nil {
		return e.NewInternalServerApiError("error deleting group", err)
	}
	return nil
}