package user_repo_mapping

import (
	user_domain "github.com/098765432m/grpc-kafka/user/internal/domain"
	user_repo "github.com/098765432m/grpc-kafka/user/internal/infrastructure/repository/sqlc/user"
)

func FromUserRepoToUserDomain(userRepo user_repo.User) user_domain.User {
	return user_domain.User{
		Id:          userRepo.ID.String(),
		Username:    userRepo.Username,
		Password:    userRepo.Password,
		Address:     userRepo.Address,
		Email:       userRepo.Email,
		PhoneNumber: userRepo.PhoneNumber,
		FullName:    userRepo.FullName,
		Role:        string(userRepo.Role),
		HotelId:     userRepo.HotelID.String(),
	}
}

func FromUsersRepoToUsersDomain(usersRepo []user_repo.User) []user_domain.User {

	users := make([]user_domain.User, 0, len(usersRepo))
	for _, u := range usersRepo {
		users = append(users, user_domain.User{
			Id:          u.ID.String(),
			Username:    u.Username,
			Password:    u.Password,
			Address:     u.Address,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			FullName:    u.FullName,
			Role:        string(u.Role),
			HotelId:     u.HotelID.String(),
		})
	}

	return users
}
