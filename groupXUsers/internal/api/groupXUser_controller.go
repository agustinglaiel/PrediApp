package api

import (
	"groupXUsers/internal/service"
)


type groupXUserController struct{
	groupXUserService service.GroupXUserServiceInterface
}

func NewGroupXUserController(groupXUserService service.GroupXUserServiceInterface) *groupXUserController {
	return &groupXUserController{
		groupXUserService: groupXUserService,
	}
}

// func (sc *groupXUserController) CreateGroupXUser(ctx context.Context, groupXUser dto.GroupXUserCreateDTO) (dto.GroupXUserCreateDTO, e.ApiError) {

// }