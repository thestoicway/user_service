package service

import "context"

func (s *userServiceImpl) SignOut(ctx context.Context, refreshToken string) (err error) {
	// Decode token
	claims, err := s.JwtManager.DecodeToken(refreshToken)

	if err != nil {
		return err
	}

	// Delete session
	if err = s.Session.DeleteSession(ctx, claims.ID); err != nil {
		return err
	}

	return
}
