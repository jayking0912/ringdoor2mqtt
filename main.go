package main

import (
	"github.com/kataras/iris/v12"
	"master/mqtt"
)

func main(){


	//启动网页
	app := iris.Default()
	//app.Favicon("./assets/favicon.ico")



	//主页
	app.Get("/bell/attr",  func(ctx iris.Context) {

		//发送到mqtt客户端
		mqtt.Mqpublic("node-red/ringdoor",1,false,"{state:true}")
	})

	//mqtt 初始化
	go mqtt.MqInit()

	app.Run(iris.Addr(":80"),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}





