package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"github.com/pluveto/flystatic/cmd/flystatic/conf"
	"github.com/pluveto/flystatic/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthService struct {
	UserMap           map[string]conf.User
	DefaultSpeedLimit uint64 // 0 means no limit
}

func NewBasicAuthService(users []conf.User, defaultSpeedLimit uint64) *BasicAuthService {
	ret := &BasicAuthService{}
	ret.UserMap = make(map[string]conf.User)
	for _, user := range users {
		ret.UserMap[user.Username] = user
	}
	ret.DefaultSpeedLimit = defaultSpeedLimit
	return ret
}

var (
	ErrCrendential           = errors.New("invalid username or password")
	ErrUnsupportedHashMethod = errors.New("unsupported hash method")
)

func (s *BasicAuthService) IsEnabled() bool {
	return len(s.UserMap) > 0
}

func (s *BasicAuthService) Check(fsDir string) {
	for _, user := range s.UserMap {
		if user.PasswordHash == "" {
			logger.Fatal("password hash is empty for user: ", user.Username)
		}
		if user.PasswordCrypt == "" {
			logger.Fatal("password crypt method is empty for user: ", user.Username)
		}

		userFullFsDir := filepath.Join(fsDir, user.SubFsDir)

		if _, err := os.Stat(userFullFsDir); err != nil {
			logger.Fatal("user dir not exists: ", userFullFsDir)
		}
	}
}

func (s *BasicAuthService) Authenticate(username, password string) error {
	user, ok := s.UserMap[username]
	if !ok {
		logger.Debug("no such user: ", username)
		return ErrCrendential
	}
	hashMethod := user.PasswordCrypt
	// todo: sha256
	switch hashMethod {
	case conf.BcryptHash:
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if user.Username == username && err == nil {
			return nil
		}
		logger.Debug("bcrypt compare error: ", err)
	case conf.SHA256Hash:
		gen := sha256.New()
		gen.Write([]byte(password))
		expectedHash := hex.EncodeToString(gen.Sum(nil))
		if user.Username == username && user.PasswordHash == expectedHash {
			return nil
		}
		logger.Debug("sha256 compare error, expected hash: ", user.PasswordHash, ", actual hash: ", expectedHash)
	default:
		return ErrUnsupportedHashMethod
	}
	return ErrCrendential

}

func (s *BasicAuthService) GetAuthorizedSubDir(username string) (string, error) {
	user, ok := s.UserMap[username]
	if !ok {
		return "", errors.New("no such user")
	}
	return user.SubFsDir, nil
}

func (s *BasicAuthService) GetPathPrefix(username string) (string, error) {
	user, ok := s.UserMap[username]
	if !ok {
		return "", errors.New("no such user")
	}
	return user.SubPath, nil
}

func (s *BasicAuthService) GetSpeedLimit(username string) (uint64, error) {
	user, ok := s.UserMap[username]
	if !ok {
		return 0, errors.New("no such user")
	}
	if user.SpeedLimit == 0 {
		return (s.DefaultSpeedLimit), nil
	}
	return user.SpeedLimit, nil
}
