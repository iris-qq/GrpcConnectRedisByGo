package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"customer.relationship/app/cr"
	"google.golang.org/grpc"
)

const (
	hostip = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(hostip, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := cr.NewSearchClient(conn)
	useridlist := []string{"name", "name1", "193283405", "193276594"}
	response, err := c.SearchFromRedis(context.Background(), &cr.UseridRequest{
		UserIdList: useridlist})
	if err != nil {
		log.Fatalf("could not greet:%v", err)
	}
	friendslistinfor := response.FriendsList
	var FriendListMap map[string]string
	json.Unmarshal([]byte(friendslistinfor), &FriendListMap)
	for k, v := range FriendListMap {
		fmt.Println(k, v)
	}
}
