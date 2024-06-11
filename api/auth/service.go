package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	userSchema "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	IsEmailRegisted(email string) bool
	SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error)
	GenerateToken(user *userSchema.User) (string, string, error)
	CreateKeystore(client *userSchema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	SignToken(claims jwt.RegisteredClaims) (string, error)
	FindApiKey(key string) (*schema.ApiKey, error)
}

type service struct {
	network.BaseService
	keystoreQuery mongo.Query[schema.Keystore]
	apikeyQuery   mongo.Query[schema.ApiKey]
	userService   user.UserService
	// token
	rsaPrivateKey        *rsa.PrivateKey
	rsaPublicKey         *rsa.PublicKey
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	tokenIssuer          string
	tokenAudience        string
}

func NewAuthService(
	db mongo.Database,
	dbQueryTimeout time.Duration,
	env *config.Env,
	userService user.UserService,
) AuthService {
	privatePem, err := utils.LoadPEMFileInto(env.RSAPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		panic(err)
	}

	publicPem, err := utils.LoadPEMFileInto(env.RSAPublicKeyPath)
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPem)
	if err != nil {
		panic(err)
	}

	s := service{
		BaseService:   network.NewBaseService(dbQueryTimeout),
		userService: userService,
		keystoreQuery: mongo.NewQuery[schema.Keystore](db, schema.KeystoreCollectionName),
		apikeyQuery:   mongo.NewQuery[schema.ApiKey](db, schema.CollectionName),
		// token key
		rsaPrivateKey: rsaPrivateKey,
		rsaPublicKey:  rsaPublicKey,
		// token claim
		accessTokenValidity:  time.Duration(env.AccessTokenValiditySec),
		refreshTokenValidity: time.Duration(env.RefreshTokenValiditySec),
		tokenIssuer:          env.TokenIssuer,
		tokenAudience:        env.TokenAudience,
	}
	return &s
}

func (s *service) IsEmailRegisted(email string) bool {
	user, _ := s.userService.FindUserByEmail(email)
	return user != nil
}

func (s *service) SignUpBasic(signupDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	role, err := s.userService.FindRoleByCode(userSchema.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]userSchema.Role, 1)
	roles[0] = *role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signupDto.Password), 5)
	if err != nil {
		return nil, err
	}

	user, err := userSchema.NewUser(signupDto.Email, string(hashed), &signupDto.Name, signupDto.ProfilePicUrl, roles)
	if err != nil {
		return nil, err
	}

	user, err = s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewUserToken(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) GenerateToken(user *userSchema.User) (string, string, error) {
	primaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	secondaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	_, err = s.CreateKeystore(user, primaryKey, secondaryKey)
	if err != nil {
		return "", "", err
	}

	now := jwt.NewNumericDate(time.Now())

	accessTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenValidity * time.Second)),
		ID:        primaryKey,
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenValidity * time.Second)),
		ID:        secondaryKey,
	}

	accessToken, err := s.SignToken(accessTokenClaims)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.SignToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) CreateKeystore(client *userSchema.User, primaryKey string, secondaryKey string) (*schema.Keystore, error) {
	ctx, cancel := s.Context()
	defer cancel()

	doc, err := schema.NewKeystore(client.ID, primaryKey, secondaryKey)
	if err != nil {
		return nil, err
	}

	id, err := s.keystoreQuery.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc.ID = *id
	return doc, nil
}

func (s *service) SignToken(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *service) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.RegisteredClaims); ok {
			return &claims, nil
		}
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *service) DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})

	if token.Valid {
		if claims, ok := token.Claims.(jwt.RegisteredClaims); ok {
			return &claims, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *service) FindApiKey(key string) (*schema.ApiKey, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"key": key, "status": true}

	apikey, err := s.apikeyQuery.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}
