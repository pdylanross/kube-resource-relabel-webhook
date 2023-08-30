package config

import (
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	ListenAddress string
	TLS           *ServerTLSConfig

	FriendlyName string
	SetupFunc    func(router *gin.Engine) error
}

type ServerTLSConfig struct {
	TLSCertFilePath string
	TLSKeyFilePath  string
}

type LoggerConfig struct {
	Level  string
	Format string
}
