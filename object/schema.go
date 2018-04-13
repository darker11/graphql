package object

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"gitlab.ucloudadmin.com/graphql-example/model"
	log "gitlab.ucloudadmin.com/wu/logrus"
)

func GenerateID() string {
	return strings.Split(uuid.NewV4().String(), "-")[0] // (Id)  8
}

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "[用户管理] 创建用户",
			Args: graphql.FieldConfigArgument{
				"userName": &graphql.ArgumentConfig{
					Description: "用户名称",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Description: "用户邮箱",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"pwd": &graphql.ArgumentConfig{
					Description: "用户密码",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"phone": &graphql.ArgumentConfig{
					Description: "用户联系方式",
					Type:        graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId, _ := strconv.Atoi(GenerateID())
				user := &model.User{
					Name: p.Args["userName"].(string),
					Email: sql.NullString{
						String: p.Args["email"].(string),
						Valid:  true,
					},
					Pwd:    p.Args["pwd"].(string),
					Phone:  int64(p.Args["phone"].(int)),
					UserID: uint64(userId),
					Status: int64(model.EnableStatus),
				}
				if err := model.InsertUser(user); err != nil {
					log.WithError(err).Error("[mutaition.createUser] invoke InserUser() failed")
					return false, err
				}
				return true, nil

			},
		},
		"changeUserName": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "[用户管理] 修改用户名称",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Description: "用户ID",
					Type:        graphql.NewNonNull(graphql.Int),
				},
				"userName": &graphql.ArgumentConfig{
					Description: "用户名称",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Args["userId"].(uint64)
				name := p.Args["userName"].(string)
				if err := model.ChangeUserName(userId, name); err != nil {
					log.WithError(err).Error("[mutaition.changeUserName] invoke InserUser() failed")
					return false, err
				}
				return true, nil

			},
		},
		"deleteUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "[用户管理] 删除用户",
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Description: "用户ID",
					Type:        graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Args["userId"].(uint64)
				if err := model.DeleteUser(userId, model.DisableStatus); err != nil {
					log.WithError(err).Error("[mutaition.deleteUser] invoke InserUser() failed")
					return false, err
				}
				return true, nil

			},
		},
	},
})

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"UserInfo": &graphql.Field{
			Description: "[用户管理] 获取指定用户的信息",
			Type:        UserInfoType,
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Description: "用户ID",
					Type:        graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Args["userId"].(uint64)
				user, err := model.GetUser(userId)
				if err != nil {
					log.WithError(err).Error("[query.UserInfo] invoke InserUser() failed")
					return false, err
				}
				return UserInfo{
					Name:   user.Name,
					UserID: user.UserID,
					Email:  user.Email.String,
					Phone:  user.Phone,
					Pwd:    user.Pwd,
					Status: model.UserStatusType(user.Status),
				}, nil

			},
		},
		"UserListInfo": &graphql.Field{
			Description: "[用户管理] 获取指定用户的信息",
			Type:        graphql.NewNonNull(graphql.NewList(UserInfoType)),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				users, err := model.GetUsers()
				if err != nil {
					log.WithError(err).Error("[query.UserInfo] invoke InserUser() failed")
					return false, err
				}
				usersList := make([]*UserInfo, 0)
				for _, v := range users {
					userInfo := new(UserInfo)
					userInfo.Name = v.Name
					userInfo.Email = v.Email.String
					userInfo.Phone = v.Phone
					userInfo.Pwd = v.Pwd
					userInfo.Status = model.UserStatusType(v.Status)
					usersList = append(usersList, userInfo)

				}
				return usersList, nil

			},
		},
	},
})
