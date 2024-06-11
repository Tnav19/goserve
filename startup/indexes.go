package startup

import (
	auth "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	contact "github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	user "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
)

func EnsureDbIndexes(db mongo.Database) {
	go mongo.Schema[auth.Keystore](&auth.Keystore{}).EnsureIndexes(db)
	go mongo.Schema[auth.ApiKey](&auth.ApiKey{}).EnsureIndexes(db)
	go mongo.Schema[user.User](&user.User{}).EnsureIndexes(db)
	go mongo.Schema[user.Role](&user.Role{}).EnsureIndexes(db)
	go mongo.Schema[contact.Message](&contact.Message{}).EnsureIndexes(db)
}
