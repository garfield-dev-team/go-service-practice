package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %q RequestURI: %q", r.Method, r.RequestURI)
	if _, err := w.Write([]byte("===测试内容2333")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 注意如果在同一个协程中接受
	// 则必须使用 buffered channel
	done := make(chan os.Signal)
	// 监听系统信号转发给 channel
	// 如果不指定信号，则转发所有信号
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// 阻塞直到接受到系统信号
		<-done
		// 也可以通过 for 循环遍历 channel
		//for o := range done {
		//	switch o {
		//	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
		//		log.Println("Server shutdown")
		//	default:
		//		log.Println("Unknown error")
		//	}
		//}
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server: ", err)
		}
	}()

	log.Println("Starting server at :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
