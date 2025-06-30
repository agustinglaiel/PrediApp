package service

import (
	"context"
	"groups/internal/client"
	"groups/internal/dto"
	"groups/internal/model"
	"groups/internal/repository"
	e "groups/pkg/utils"
	"time"
)


type groupService struct {
	groupRepo repository.GroupRepository
}

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, group dto.CreateGroupRequestDTO) (dto.GroupResponseDTO, e.ApiError)
	GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError)
	GetGroupsByUserId(ctx context.Context, userID int) ([]dto.GroupResponseDTO, e.ApiError) 
	GetGroups(ctx context.Context) ([]dto.GroupListResponseDTO, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	JoinGroup(ctx context.Context, request dto.RequestJoinGroupDTO) e.ApiError
	ManageGroupInvitation(ctx context.Context, request dto.ManageGroupInvitationDTO) e.ApiError
}

func NewGroupService(groupRepo repository.GroupRepository) GroupServiceInterface {
	return &groupService{groupRepo: groupRepo}
}

func (s *groupService) CreateGroup(ctx context.Context, request dto.CreateGroupRequestDTO) (dto.GroupResponseDTO, e.ApiError) {
	// Verificar si el nombre del grupo ya existe
	if _, err := s.groupRepo.GetGroupByGroupName(ctx, request.GroupName); err == nil {
		return dto.GroupResponseDTO{}, e.NewBadRequestApiError("Group name already exists")
	}

	newGroup := &model.Group{
		GroupName:   request.GroupName,
		Description: request.Description,
	}

	// Crear el grupo en la base de datos
	if err := s.groupRepo.CreateGroup(ctx, newGroup); err != nil {
		return dto.GroupResponseDTO{}, e.NewInternalServerApiError("Error creating group", err)
	}

	// Asociar al usuario creador en `groups_users` con el rol de `creator`
	if err := s.groupRepo.AddUserToGroup(ctx, newGroup.ID, request.UserID, "creator"); err != nil {
		return dto.GroupResponseDTO{}, e.NewInternalServerApiError("Error adding creator to group", err)
	}

	// Construir la respuesta con la lista de usuarios
	response := dto.GroupResponseDTO{
		ID:          newGroup.ID,
		GroupName:   newGroup.GroupName,
		Description: newGroup.Description,
		GroupCode:   newGroup.GroupCode,
		Users: []dto.GroupUserResponseDTO{
			{
				UserID: request.UserID,
				Role:   "creator",
				Score:  nil, // Se traerá desde el microservicio de `users` en otra función
			},
		},
		CreatedAt:   newGroup.CreatedAt.Format(time.RFC3339), // ✅ Convertimos a string ISO 8601
		UpdatedAt:   newGroup.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *groupService) GetGroupsByUserId(ctx context.Context, userID int) ([]dto.GroupResponseDTO, e.ApiError) {
	groups, apiErr := s.groupRepo.GetGroupsByUserId(ctx, userID)
	if apiErr != nil {
		return nil, apiErr
	}

	usersClient := client.NewUsersClient()
	var responses []dto.GroupResponseDTO

	for _, g := range groups {
		var users []dto.GroupUserResponseDTO

		for _, gu := range g.GroupUsers {
			score, err := usersClient.GetUserScore(gu.UserID)
			if err != nil {
				score = 0
			}
			users = append(users, dto.GroupUserResponseDTO{
				UserID: gu.UserID,
				Role:   gu.GroupRole,
				Score:  &score,
			})
		}

		responses = append(responses, dto.GroupResponseDTO{
			ID:          g.ID,
			GroupName:   g.GroupName,
			Description: g.Description,
			GroupCode:   g.GroupCode,
			Users:       users,
			CreatedAt:   g.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}

func (s *groupService) GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError) {
	group, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return dto.GroupResponseDTO{}, err
	}

	usersClient := client.NewUsersClient()

	// Mapear usuarios del grupo con sus puntajes
	var users []dto.GroupUserResponseDTO
	for _, groupUser := range group.GroupUsers {
		score, err := usersClient.GetUserScore(groupUser.UserID)
		if err != nil {
			score = 0 // Si falla, asignamos 0 por defecto
		}

		users = append(users, dto.GroupUserResponseDTO{
			UserID: groupUser.UserID,
			Role:   groupUser.GroupRole,
			Score:  &score,
		})
	}

	response := dto.GroupResponseDTO{
		ID:          group.ID,
		GroupName:   group.GroupName,
		Description: group.Description,
		GroupCode:   group.GroupCode,
		Users:       users,
		CreatedAt:   group.CreatedAt.Format(time.RFC3339), // Convertimos a string ISO 8601
		UpdatedAt:   group.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}


func (s *groupService) GetGroups(ctx context.Context) ([]dto.GroupListResponseDTO, e.ApiError) {
	groups, err := s.groupRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.GroupListResponseDTO
	for _, group := range groups {
		response = append(response, dto.GroupListResponseDTO{
			ID:        group.ID,
			GroupName: group.GroupName,
			GroupCode: group.GroupCode,
		})
	}

	return response, nil
}

func (s *groupService) DeleteGroupByID(ctx context.Context, id int) e.ApiError {
	_, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return e.NewBadRequestApiError("Group not found")
	}

	if err := s.groupRepo.DeleteGroupByID(ctx, id); err != nil {
		return e.NewInternalServerApiError("Error deleting group", err)
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

func (s *groupService) JoinGroup(ctx context.Context, request dto.RequestJoinGroupDTO) e.ApiError {
	group, err := s.groupRepo.GetGroupByCode(ctx, request.GroupCode)
	if err != nil {
		return err
	}

	// Verificar si el usuario ya está en el grupo
	role, err := s.groupRepo.GetUserRoleInGroup(ctx, group.ID, request.UserID)
	if err != nil {
		return err
	}
	if role == "member" {
		return e.NewBadRequestApiError("User is already a member of the group")
	}
	if role == "invited" {
		return e.NewBadRequestApiError("User has already requested to join the group and is awaiting approval")
	}

	// Agregar al usuario como "invited"
	if err := s.groupRepo.AddUserToGroup(ctx, group.ID, request.UserID, "invited"); err != nil {
		return e.NewInternalServerApiError("Error adding user to group", err)
	}

	return nil
}

//Un usuario solicita unirse a un grupo mediante la ruta POST /groups/join. 
///En este momento, su estado en la tabla group_x_users queda como "invited".
// Luego, el creador del grupo debe aceptar o rechazar la invitación usando POST /groups/manage-invitation.
//Si la acción es "accept", el usuario pasa a "member".
//Si la acción es "reject", el usuario es eliminado del grupo.
func (s *groupService) ManageGroupInvitation(ctx context.Context, request dto.ManageGroupInvitationDTO) e.ApiError {
	group, err := s.groupRepo.GetGroupByID(ctx, request.GroupID)
	if err != nil {
		return err
	}

	// Verificar si el usuario que realiza la acción (CreatorID) es el creator del grupo
	creatorRole, err := s.groupRepo.GetUserRoleInGroup(ctx, group.ID, request.CreatorID)
	if err != nil {
		return err
	}
	if creatorRole != "creator" {
		return e.NewForbiddenApiError("Only the group creator can manage invitations")
	}

	// Verificar si el usuario a gestionar está en el grupo y es "invited"
	targetRole, err := s.groupRepo.GetUserRoleInGroup(ctx, group.ID, request.TargetUserID)
	if err != nil {
		return err
	}
	if targetRole != "invited" {
		return e.NewBadRequestApiError("User is not in invited status")
	}

	// Si la acción es "accept", cambiamos el rol a "member"
	if request.Action == "accept" {
		if err := s.groupRepo.UpdateUserRoleInGroup(ctx, group.ID, request.TargetUserID, "member"); err != nil {
			return e.NewInternalServerApiError("Error accepting user into group", err)
		}
	}

	// Si la acción es "reject", eliminamos al usuario del grupo
	if request.Action == "reject" {
		if err := s.groupRepo.RemoveUserFromGroup(ctx, group.ID, request.TargetUserID); err != nil {
			return e.NewInternalServerApiError("Error rejecting user from group", err)
		}
	}

	return nil
}
