package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/koding/multiconfig"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gitlab.ucloudadmin.com/graphql-example/model"
	"gitlab.ucloudadmin.com/graphql-example/object"
	_ "gitlab.ucloudadmin.com/graphql-example/util/loghelper"
	log "gitlab.ucloudadmin.com/wu/logrus"
)

type ServerCfg struct {
	Addr      string
	MysqlAddr string
}

func main() {
	//load config info
	m := multiconfig.NewWithPath("config.toml")
	svrCfg := new(ServerCfg)
	m.MustLoad(svrCfg)
	//new graphql schema
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    object.QueryType,
			Mutation: object.MutationType,
		},
	)
	if err != nil {
		log.WithError(err).Error("[main] invoke graphql.NewSchema() failed")
		return
	}

	model.InitSqlxClient(svrCfg.MysqlAddr)
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		//read user_id from gateway
		userIDStr := r.Header.Get("user_id")
		if len(userIDStr) > 0 {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			ctx = context.WithValue(ctx, "ContextUserIDKey", userID)
		}
		h.ContextHandler(ctx, w, r)

	})
	log.Fatal(http.ListenAndServe(svrCfg.Addr, nil))
}
