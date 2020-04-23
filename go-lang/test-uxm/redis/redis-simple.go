
// redis simple example
// (c) 2020 unix-world.org

package main

import (
	"fmt"
	"log"
	"strconv"
	"github.com/gomodule/redigo/redis"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "6379"
	DB_NUM = 1
)

var (
	c redis.Conn
	err error
	reply interface{}
)

func init() {
	fmt.Println("===== Connecting to Redis: [Host:Port = " + CONN_HOST + ":" + CONN_PORT +  " / DB = " + strconv.Itoa(DB_NUM) + "] =====")
	c, err = redis.Dial("tcp", CONN_HOST + ":" + CONN_PORT, redis.DialDatabase(DB_NUM))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	defer c.Close()

	var redisKeyName = "hello"
	var redisKeyValue = "world"

	fmt.Println("Redis SET sample key: [" + redisKeyName + "] with value of: [" + redisKeyValue + "]")
	_, err := c.Do("SET", redisKeyName, redisKeyValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Redis GET the above sample key")
	s, err := redis.String(c.Do("GET", redisKeyName))
	if err != nil {
		log.Fatal(err)
	}
	if(s != redisKeyValue) {
		fmt.Println("ERROR: Redis Key GET failed. Supposed to be: [" + redisKeyValue + "] but is: [" + s + "]")
	} else {
		fmt.Println("OK: Key value in Redis is: [" + s + "]")
	}

	var redisCMD = "FLUSHDB"
	fmt.Println("Redis Command: " + redisCMD)
	c.Do(redisCMD)

	fmt.Println("... DONE ... program will exit now ...")

}
