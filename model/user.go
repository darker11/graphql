package model

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	log "gitlab.ucloudadmin.com/wu/logrus"
)

type CommonInfo struct {
	ID          uint64         `db:"id"`
	CreatedTime time.Time      `db:"created_time"`
	DeletedTime mysql.NullTime `db:"deleted_time"`
	UpdatedTime time.Time      `db:"updated_time"`
}

/*
请按照如下内容在数据库中建表
CREATE TABLE `t_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) NOT NULL,
  `name` char(100) NOT NULL,
  `pwd` char(128) NOT NULL,
  `email` char(50) NOT NULL,
  `phone` bigint(20) NOT NULL,
  `status` int(10) NOT NULL,
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8
*/

type User struct {
	CommonInfo
	UserID uint64         `db:"user_id"`
	Name   string         `db:"name"`
	Pwd    string         `db:"pwd"`
	Email  sql.NullString `db:"email"`
	Phone  int64          `db:"phone"`
	Status int64          `db:"status"`
}

func GetUser(userId uint64) (*User, error) {
	sql := `SELECT * FROM t_user where user_id = ?`
	var user User
	if err := dbx.Get(&user, sql, userId); err != nil {
		log.WithError(err).Error("[model.GetUser] invoke mysql failed")
		return nil, err
	}
	return &user, nil
}

func GetUsers() ([]User, error) {
	var users []User
	sql := `SELECT * FROM t_user`
	if err := dbx.Select(&users, sql); err != nil {
		log.WithError(err).Error("[model.GetUsers] invoke mysql failed")
		return nil, err
	}
	return users, nil

}

func InsertUser(user *User) error {
	sql := `INSERT into t_user (user_id, name, pwd, email, phone, status) VALUES (:user_id, :name, :pwd, :email, :phone, :status)`
	if _, err := dbx.NamedExec(sql, user); err != nil {
		log.WithError(err).Error("[model.InserUser] invoke mysql failed")
		return err
	}
	return nil
}

//标记删除
func DeleteUser(userId uint64, status UserStatusType) error {
	sql := `UPDATE t_user status = :status WHERE user_id = :user_id`
	user := User{
		UserID: userId,
		Status: int64(status),
		CommonInfo: CommonInfo{
			DeletedTime: mysql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
	}
	if _, err := dbx.NamedExec(sql, user); err != nil {
		log.WithError(err).Error("[model.DeleteUser] invoke mysql failed")
		return err

	}
	return nil
}

func ChangeUserName(userId uint64, name string) error {
	sql := `UPDATE t_user name = :name WHERE user_id = :user_id`
	user := User{
		UserID: userId,
		Name:   name,
		CommonInfo: CommonInfo{
			UpdatedTime: time.Now(),
		},
	}
	if _, err := dbx.NamedExec(sql, user); err != nil {
		log.WithError(err).Error("[model.ChangeUserName] invoke mysql failed")
		return err
	}
	return nil
}
