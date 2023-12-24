package main

import "proxy-engineering-thesis/server"

func main() {
	//agent := aws.NewCloudWatchConfiguration("GRUPA", "STREAM", "")
	//agent.InitLogStore()
	//err := agent.SendLog("cale te sie robi o i drugi raz")
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//}
	server.StartServer()
}
