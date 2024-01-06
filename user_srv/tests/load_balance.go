// 测试负载均衡
package main

import (
	"context"
	"fmt"
	"log"
	"mxshop_srvs/user_srv/proto"
	"strconv"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
)

func main() {
	consulConfig := api.DefaultConfig()
	consulClient, _ := api.NewClient(consulConfig)
	services, _, _ := consulClient.Health().Service("user-srv", "srv", true, nil)
	addr := services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)
	fmt.Println(addr)
	conn, err := grpc.Dial("consul://127.0.0.1:8500/user-srv?wait=14s&tag=srv",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	userSrvClient := proto.NewUserClient(conn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}

	for index, data := range rsp.Data {
		fmt.Println(index, data)
	}
}
