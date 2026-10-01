package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Rdb struct {
	User string
	Pass string
	Host string
	Port string
}

func NewRdb(user, pass, host, port string) *Rdb {
	return &Rdb{
		User: user,
		Pass: pass,
		Host: host,
		Port: port,
	}
}

func (r *Rdb) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Username: r.User,
		Password: r.Pass,
	})
}

//  rdb := redis.NewClient(&redis.Options{
//         Addr:     "localhost:6379",
//         Password: "", // no password set
//         DB:       0,  // use default DB
//     })
//     defer rdb.Close()
