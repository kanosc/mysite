package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()
var m = melody.New()
var ServerMode string
var defaultExistTime = 72 * time.Hour
var maxMsgCountOfRoom int64 = 500

type ChatInfo struct {
	RoomName string   `json:"roomname"`
	Messages []string `json:"messages"`
}

type RoomInfo struct {
	RoomName string `form:"roomname"`
	MaxUser  string `form:"maxuser"`
	Password string `form:"password"`
}

const tokenSecret = "everyting is ok"

func getMsgKey(rn string) string {
	return fmt.Sprintf("chat:roomname:%s:messages", rn)
}

func getPwdKey(rn string) string {
	return fmt.Sprintf("chat:roomname:%s:password", rn)
}

func getAllRoomKey() string {
	return "chat:rooms"
}

func init() {
	m.HandleMessage(rcvMessage)
}

func RedisInit(mode string) {
	if mode == "debug" {
		ServerMode = mode
		redisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			//Password: "929319", // no password set
			Password: "myredis6379", // no password set
			DB:       0,
		})
		return
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "929319", // no password set
		//Password: "myredis6379", // no password set
		DB: 0,
	})

}

func rcvMessage(s *melody.Session, msg []byte) {
	log.Println("handling msg start")
	roomname, _ := s.Keys["roomname"].(string)
	msgkey := getMsgKey(roomname)
	pwdkey := getPwdKey(roomname)

	// limit max number of message
	err := limitMaxMsgNum(ctx, maxMsgCountOfRoom, msgkey, redisClient)
	if err != nil {
		log.Println("limit max msg failed", err.Error())
	}

	// insert message to message list
	err = insertMsg(ctx, msg, msgkey, redisClient)
	if err != nil {
		log.Println("insert msg failed", err.Error())
	}

	// refresh the expire time of message key and password key(if exist)
	err = refreshKeyExpire(ctx, defaultExistTime, msgkey, redisClient)
	if err != nil {
		log.Println("refresh key expire failed", err.Error())
	}

	if _, e := redisClient.Get(ctx, pwdkey).Result(); e != redis.Nil {
		err = refreshKeyExpire(ctx, defaultExistTime, pwdkey, redisClient)
		if err != nil {
			log.Println("refresh key expire failed", err.Error())
		}
	}

	// broadcast message to matched URL
	// m.Broadcast(msg)
	m.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

func limitMaxMsgNum(ctx context.Context, maxMsgs int64, msgkey string, rdc *redis.Client) error {
	msgCount, err := rdc.LLen(ctx, msgkey).Result()
	if err != nil {
		return err
	}
	if msgCount >= maxMsgCountOfRoom {
		log.Printf("length of messages is out of limit %v\n", maxMsgs)
		discardMsg, err := rdc.RPop(ctx, msgkey).Result()
		log.Printf("Msg full, remove oldest msg[%v]\n", discardMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertMsg(ctx context.Context, msg []byte, msgkey string, rdc *redis.Client) error {
	_, err := rdc.LPush(ctx, msgkey, string(msg)).Result()
	if err != nil {
		log.Println("add msg failed")
		return err
	}
	log.Printf("inserting msg[%v]\n", string(msg))
	return nil
}

func refreshKeyExpire(ctx context.Context, d time.Duration, key string, rdc *redis.Client) error {
	_, err := rdc.Expire(ctx, key, d).Result()
	if err != nil {
		log.Printf("Set expire time for [%v] fail\n", key)
		return err
	}
	return nil
}

func HandlePostMsg(c *gin.Context) {
	var roomname interface{} = c.Param("roomname")
	if roomname != "" {
		k := map[string]interface{}{"roomname": roomname}
		m.HandleRequestWithKeys(c.Writer, c.Request, k)
	}
}

func QueryRooms(c *gin.Context) {
	rooms, err := redisClient.LRange(ctx, getAllRoomKey(), 0, -1).Result()
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "internal error")
		c.Abort()
		return
	}
	for _, r := range rooms {
		_, err := redisClient.Get(ctx, getMsgKey(r)).Result()
		if err == redis.Nil {
			log.Println("room expires, start to delete ", r)
			_ = redisClient.LRem(ctx, getAllRoomKey(), 1, r)
		}
	}
	c.JSON(200, gin.H{
		"rooms": rooms,
	})

}

func CreateRoom(c *gin.Context) {
	var roomInfo RoomInfo
	c.ShouldBind(&roomInfo)
	if len(strings.Trim(roomInfo.RoomName, " ")) == 0 {
		c.String(http.StatusForbidden, "invalid roomname")
		c.Abort()
		return
	}

	log.Printf("creating new room, roomname[%v], maxuser[%v], password[%v]\n", roomInfo.RoomName, roomInfo.MaxUser, roomInfo.Password)
	roomk := getMsgKey(roomInfo.RoomName)
	roomp := getPwdKey(roomInfo.RoomName)
	rooms := getAllRoomKey()

	_, err := redisClient.Get(ctx, roomk).Result()
	if err != redis.Nil {
		log.Println("room already exist, return")
		c.String(http.StatusForbidden, "room has existed")
		c.Abort()
		return
	}

	_, err = redisClient.LPush(ctx, roomk, "room created success").Result()
	if err != nil {
		log.Println("create room", roomInfo.RoomName, "failed", roomInfo.MaxUser)
	}

	_, err = redisClient.Expire(ctx, roomk, defaultExistTime).Result()
	if err != nil {
		log.Println("Set expire time for room failed")
	}

	if roomInfo.Password != "" {
		_, err = redisClient.Set(ctx, roomp, roomInfo.Password, defaultExistTime).Result()
		if err != nil {
			log.Println("set password", roomp, "failed")
		}

		log.Println("set password", roomp, "success")
	}

	_, err = redisClient.LPush(ctx, rooms, roomInfo.RoomName).Result()

	c.String(http.StatusOK, "create room seccessful")
}

func DeleteRoom(c *gin.Context) {
	roomname := c.Param("roomname")

	msgkey := getMsgKey(roomname)
	pwdkey := getPwdKey(roomname)
	rooms := getAllRoomKey()

	// refresh the expire time of message key and password key(if exist)
	err := refreshKeyExpire(ctx, 0*time.Second, msgkey, redisClient)
	if err != nil {
		log.Println("refresh key expire failed", err.Error())
	}

	if _, e := redisClient.Get(ctx, pwdkey).Result(); e != redis.Nil {
		err = refreshKeyExpire(ctx, 0*time.Second, pwdkey, redisClient)
		if err != nil {
			log.Println("refresh key expire failed", err.Error())
		}
	}

	redisClient.LRem(ctx, rooms, 1, roomname)

	c.String(http.StatusOK, "delete room seccessful")
}

func GetRoomMessages(c *gin.Context) {
	// roomname, _ := c.MustGet("roomname").(string)
	roomname := c.Param("roomname")

	historyMsg, err := redisClient.LRange(ctx, getMsgKey(roomname), 0, -1).Result()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	response := ChatInfo{
		RoomName: roomname,
		Messages: historyMsg,
	}

	c.JSON(200, response)
}

type MyCustomClaims struct {
	Roomname string `json:"roomname"`
	jwt.StandardClaims
}

func PasswordChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomname := c.Param("roomname")
		_, err := redisClient.Get(ctx, getPwdKey(roomname)).Result()

		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			c.Abort()
			return
		}

		if err == redis.Nil {
			c.Next()
			return
		}

		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "no token in request",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			log.Printf("roomname in token[%v]\n", claims["roomname"])
			if claims["roomname"] != roomname {
				c.JSON(http.StatusUnauthorized, gin.H{
					"err": "roomname not matched",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "invalid token claim type",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Printf("invalid token for [%v], reason [%v]\n", tokenString, err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": "Invalid token " + err.Error(),
			})
			c.Abort()
			return
		}
		log.Printf("token validate successful[%v], for room [%v]\n", tokenString, roomname)
		c.Next()

	}
}

func AuthToken(c *gin.Context) {
	roomname := c.Param("roomname")
	pwdInReq := c.PostForm("password")
	if pwdInReq == "" || roomname == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "invalid auth request",
		})
		return
	}
	pwdInDB, err := redisClient.Get(ctx, getPwdKey(roomname)).Result()
	if err != nil {
		log.Println(err.Error(), getPwdKey(roomname))
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "fail to check pwd",
		})
		return
	}

	if pwdInReq != pwdInDB {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err": "password is wrong",
		})
		return
	}

	mySigningKey := []byte(tokenSecret)

	// Create the Claims
	claims := MyCustomClaims{
		roomname,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
			Issuer:    "server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	log.Printf("%v %v\n", ss, err)
	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}
