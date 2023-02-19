package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pluveto/flystatic/cmd/flystatic/conf"
	"github.com/pluveto/flystatic/cmd/flystatic/middleware"
	"github.com/pluveto/flystatic/pkg/logger"
)

type AuthService interface {
	IsEnabled() bool
	Authenticate(username, password string) error
	GetAuthorizedSubDir(username string) (string, error)
	GetPathPrefix(username string) (string, error)
	GetSpeedLimit(username string) (uint64, error)
}

type StaticServer struct {
	AuthService AuthService
	Host        string
	Port        int
	Path        string
	FsDir       string
	TlsConf     conf.SSL
	Middlewares []func(http.HandlerFunc) http.HandlerFunc
}

func NewStaticServer(authService AuthService, serverConf conf.Server, tlsConf conf.SSL) *StaticServer {
	return &StaticServer{
		AuthService: authService,
		Host:        serverConf.Host,
		Port:        serverConf.Port,
		Path:        serverConf.Path,
		FsDir:       serverConf.FsDir,
		TlsConf:     tlsConf,
	}
}

func (s *StaticServer) AddMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func (s *StaticServer) check() {
	if nil == s.AuthService {
		logger.Fatal("AuthService is nil")
	}
	if s.FsDir == "" {
		logger.Fatal("FsDir is empty")
	}
	// if !path.IsAbs(s.FsDir) {
	// logger.Fatal("FsDir is not an absolute path")
	// }

	var err error
	s.FsDir, err = filepath.Abs(s.FsDir)
	if err != nil {
		logger.Fatal("FsDir is not a valid path", err)
	}

	// must exists
	if _, err := os.Stat(s.FsDir); os.IsNotExist(err) {
		logger.Fatal("FsDir does not exist", err)
	}

	if !isUnderHomeDir(s.FsDir) {
		tmpDir := filepath.Join(os.TempDir(), "flystatic")
		if !strings.HasPrefix(s.FsDir, tmpDir) {
			logger.Warn("You're using a path which isn't under home dir as mapped directory. This may cause security issues.")
		}
	}

	logger.Debug("FsDir: ", s.FsDir)
}

func isUnderHomeDir(s string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal("Error when getting home dir", err)
	}
	return strings.HasPrefix(s, homeDir)
}

func (s *StaticServer) wrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range s.Middlewares {
			h = middleware(h)
		}
		h(w, r)
	}
}

func (s *StaticServer) Listen() {
	s.check()

	http.HandleFunc("/", s.wrapHandler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request: ", r.Method, r.URL.Path)

		if s.AuthService.IsEnabled() {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			}

			err := s.AuthService.Authenticate(username, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				logger.Error("Unauthorized: ", err)
				return
			}

			subFsDir, err := s.AuthService.GetAuthorizedSubDir(username)
			if err != nil {
				http.Error(w, "Internal Error.", http.StatusInternalServerError)
				logger.Errorf("Error when getting authorized sub dir for user %s: %s", username, err)
				return
			}

			userPrefix, err := s.AuthService.GetPathPrefix(username)
			if err != nil {
				http.Error(w, "Internal Error.", http.StatusInternalServerError)
				logger.Errorf("Error when getting path prefix for user %s: %s", username, err)
			}

			fs := http.Dir(filepath.Join(s.FsDir, subFsDir))
			h := http.FileServer(fs)
			h = s.ApplyUserSpeedLimitMiddleware(username, h)
			h = http.StripPrefix(buildHttpPathPrefix(s.Path, userPrefix), h)
			h.ServeHTTP(w, r)
		} else {
			fs := http.Dir(s.FsDir)
			h := http.FileServer(fs)
			h = http.StripPrefix(s.Path, h)
			h.ServeHTTP(w, r)
		}

	}))

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	var err error
	if s.TlsConf.Enabled {
		err = http.ListenAndServeTLS(addr, s.TlsConf.Cert, s.TlsConf.Key, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}

	logger.Fatal("failed to listen and serve on", addr, ":", err)
}

var speedLimiterMap map[string]http.Handler

func (s *StaticServer) ApplyUserSpeedLimitMiddleware(username string, h http.Handler) http.Handler {
	speedLimiter, ok := speedLimiterMap[username]
	if ok {
		return speedLimiter
	}

	userSpeedLimit, err := s.AuthService.GetSpeedLimit(username)
	if err != nil {
		logger.Fatal("Error when getting speed limit for user", username, ":", err)
	}

	speedLimiter = middleware.NewSpeedLimiter(float64(userSpeedLimit), 100)
	speedLimiterMap[username] = speedLimiter

	return speedLimiter
}

func buildHttpPathPrefix(path, userPrefix string) string {
	if userPrefix == "" {
		return path
	}
	joined := filepath.Join(path, userPrefix)
	// use unix path separator
	joined = strings.ReplaceAll(joined, "\\", "/")
	return strings.TrimSuffix(joined, "/")
}
