package csrf

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"sync"

	"github.com/Beretta350/authentication/pkg/util"
)

const (
	Secret  string = "csrfSecret"
	Session string = "csrfSession"
	Token   string = "csrfToken"
)

type CSRFWrapper interface {
	GenerateToken(secretSent string, session string) string
	ValidateToken(token string, session string) bool
	IsIgnoredPath(path string) bool
	GetSecret() string
}

var instance *csrfWrapper
var once sync.Once

type csrfWrapper struct {
	secretKey   string
	ignorePaths []string
}

// Singleton
func NewCSRFWrapper(secret string, ignore []string) *csrfWrapper {
	once.Do(func() {
		instance = &csrfWrapper{secretKey: secret, ignorePaths: ignore}
	})
	return instance
}

func GetCSRFWrapper() *csrfWrapper {
	return instance
}

func (wrap *csrfWrapper) GenerateToken(secretSent string, session string) string {
	h := sha1.New()
	io.WriteString(h, session+"-"+secretSent)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func (wrap *csrfWrapper) ValidateToken(token string, session string) bool {
	return wrap.GenerateToken(wrap.secretKey, session) == token
}

func (wrap *csrfWrapper) IsIgnoredPath(path string) bool {
	return util.InArray(wrap.ignorePaths, path)
}

func (wrap *csrfWrapper) GetSecret() string {
	return wrap.secretKey
}
