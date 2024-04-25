package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yuxing.code/bookstore/internal/store"
	"yuxing.code/bookstore/server"
	"yuxing.code/bookstore/store/factory"
)

func main() {
	// 创建图书数据存储模块实例
	s, err := factory.New(store.MemStoreProviderName)
	if err != nil {
		panic(err)
	}

	// 创建http服务
	srv := server.NewBookStoreServer(":8080", s)

	errChan, err := srv.ListenAndServe()
	if err != nil {
		log.Println("web server start failed:", err)
		return
	}

	log.Println("web server start success")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select { // 监听来自errChan及c的事件
	case err := <-errChan:
		log.Println("web server run failed:", err)
		return
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx) // 优雅关闭http服务器
	}

	if err != nil {
		log.Println("bookstore program exit error:", err)
		return
	}

	log.Println("bookstore program exit success")
}
