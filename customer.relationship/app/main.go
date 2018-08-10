package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"customer.relationship/app/cr"
	"github.com/garyburd/redigo/redis"
	"google.golang.org/grpc"
)

const (
	hostip = "localhost:6379"
	host   = "0.0.0.0"
	port   = "8080"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cr.RegisterSearchServer(s, &SearchServer{})
	s.Serve(lis)
}

//GetFriendsIDList user by search from redis
func GetFriendsIDList(UserIDList []string) (string, error) {
	useridlist := make([]interface{}, len(UserIDList))
	for i, v := range UserIDList {
		useridlist[i] = v
	}

	log.Println(useridlist)
	conn, err := ConnectRedis()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	friendinfor, err := Searchredis(conn, useridlist)
	if err != nil {
		panic(err)
	}
	friendlist, err := json.Marshal(friendinfor)
	return string(friendlist), err
}

//Searchredis use for search from reids
func Searchredis(conn redis.Conn, userid []interface{}) (map[string]string, error) {
	FriendsList, err := redis.Strings(conn.Do("MGET", userid...))
	if err != nil {
		fmt.Println("search error")
	}
	UserInfor := make(map[string]string)
	for i := 0; i < len(FriendsList); i++ {
		UserInfor[userid[i].(string)] = FriendsList[i]
	}
	return UserInfor, err
}

//ConnectRedis use for connet redis
func ConnectRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", hostip)
	if err != nil {
		fmt.Println("connect error.")
	}
	return conn, err
}

//SearchServer use by reload method
type SearchServer struct{}

//SearchFromRedis use by reload method
func (s *SearchServer) SearchFromRedis(ctx context.Context, in *cr.UseridRequest) (*cr.FriendListReply, error) {
	useridlist := in.UserIdList
	log.Printf("%+v", in)
	friendslist, err := GetFriendsIDList(useridlist)
	if err != nil {
		//dlog.Error("Msg", "SearchFromRedis", "Error", err.Error())
		return nil, err
	}
	return &cr.FriendListReply{
		FriendsList: friendslist}, nil
}
