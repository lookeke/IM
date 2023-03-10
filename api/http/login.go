package http

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"im/db"
	"im/schema"
	"log"
	"net/http"
)

func Login(account string, password string, user *schema.UserBasic) schema.Status {
	collection := db.Mongo.Database("im").Collection("user_basic")
	if account != "" || password != "" {
		filter := bson.D{{"account", account}}
		err := collection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {

				return schema.Status{
					Code:    http.StatusForbidden,
					Message: "数据库不存在该账号",
					Body:    "该用户未注册",
				}
			}
			log.Println(err)
		}
	} else if account != user.Password && password != user.Password {
		return schema.Status{
			Code:    http.StatusUnauthorized,
			Message: "账号密码错误",
			Body:    "账号密码错误",
		}
	}
	return schema.Status{
		Code:    http.StatusOK,
		Message: "登陆成功",
		Body:    user,
	}
}
