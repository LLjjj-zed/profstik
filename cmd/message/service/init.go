package service

import (
	jwt "github.com/132982317/profstik/middleware"
	tool "github.com/132982317/profstik/pkg/utils/crypt"
)

var (
	Jwt        *jwt.JWT
	publicKey  string
	privateKey string
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
	publicKey, _ = tool.ReadKeyFromFile(tool.PublicKeyFilePath)
	privateKey, _ = tool.ReadKeyFromFile(tool.PrivateKeyFilePath)
}
