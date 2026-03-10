package auth

import (
	"context"
)

type SessionRepository interface {
	Create(ctx context.Context, session *JwtSession) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (JwtSession, error)
	GetByUserAndRefreshToken(ctx context.Context, userID uint, token string) (JwtSession, error)
	GetByUserID(ctx context.Context, userID uint) ([]JwtSession, error)
	DeleteUserSession(ctx context.Context, userID uint, sessionID string) error
	DeleteAllUserSessions(ctx context.Context, userID uint) error
}
