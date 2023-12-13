package main

import "proxy-engineering-thesis/server"

func main() {
	server.StartServer()
	//p := sql.NewProxy("proxy_cale_te",
	//	model.NewAddress("127.0.0.1", "4444"),
	//	datasource.NewSqlDataSource("127.0.0.1", "5432", "postgres", "postgres", "library-db"),
	//)
	//p.Start()
}
