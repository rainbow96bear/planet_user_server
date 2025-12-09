package service

// type ThemeService struct {
// 	ProfilesRepo *repository.ProfilesRepository
// }

// // 테마 조회
// func (s *ThemeService) GetTheme(ctx context.Context, UserID uuid.UUID) (string, error) {
// 	theme, err := s.ProfilesRepo.GetTheme(ctx, UserID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get theme: %w", err)
// 	}
// 	return theme, nil
// }

// // 테마 설정
// func (s *ThemeService) SetTheme(ctx context.Context, userID uuid.UUID, theme string) error {
// 	if err := s.ProfilesRepo.SetTheme(ctx, userID, theme); err != nil {
// 		return fmt.Errorf("failed to set theme: %w", err)
// 	}
// 	return nil
// }
