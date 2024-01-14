package service

import "context"

func (s *userServiceImpl) SignOut(ctx context.Context, refreshToken string) (err error) {

	// Decode token
	result, err := s.JwtManager.DecodeToken(refreshToken)

	if err != nil {
		return err
	}

	// Delete session
	if err = s.Session.DeleteSession(ctx, result.ID); err != nil {
		return err
	}
	
	return
}
