package service

type ProfileService struct {
	UsersRepo *repository.UsersRepository
}

func (s *ProfileService)GetProfilInfo(ctx  context.Context, nickname string)(*dto.profileInfo, error){
	profile, err := s.UsersRepo.GetProfileInfo(ctx, nickname)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("fail to get profile info ERR[%s]", err.Error())
	}

	return profile, nil
}

func (s *ProfileService)UpdateProfile(ctx  context.Context, profile *dto.profileInfo) error{
	err := s.UsersRepo.UpdateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return nil
}