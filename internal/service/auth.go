package service

type AuthService struct {
	UsersRepo *repository.UsersRepository
}

func (s *AuthService)VerifyUser(ctx  context.Context, nickname, userUuid string)(bool, error){
	dbUuid , err := s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
	if err != nil {
		return false, err
	}

	return dbUuid == userUuid, nil
}