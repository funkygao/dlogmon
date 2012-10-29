package main

import (
	"os"
	"secondbit.org/pastry"
)

func main() {
	hostname, _ := os.Hostname()
	println(hostname)

	id, _ := pastry.NodeIDFromBytes([]byte(hostname))
	node := pastry.NewNode(id, "121.121.12.12", "1.1.1.1", "region", 9000)
	credentials := pastry.Passphrase("I <3 Gophers.")
	cluster := pastry.NewCluster(node, credentials)
	go func() {
		defer cluster.Stop()
		cluster.Listen()
	}()

	cluster.Join("2.2.23.3", 8080)
	select {}

}
