package main

import (
	"context"
	"dogovor.alif.tj/routes"
	"dogovor.alif.tj/utils"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var storagePathPtr = flag.String("storage", "exel-reports", "folder for storing files")

func main() {

	r := routes.Router()

	createStorageFolder()
	svc := http.Server{
		Handler: utils.AddCors(r),
		Addr:    ":3949",
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		exitCh := make(chan os.Signal, 1)
		signal.Notify(exitCh, os.Interrupt, os.Kill)
		<-exitCh
		err := svc.Shutdown(context.Background())
		if err != nil {
			log.Println("ERROR server shutdown: ", err)
		}
		close(exitCh)
	}()

	fmt.Println("Server is start to listening on port: ", svc.Addr)
	err := svc.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
			log.Fatal("ERROR in running process server: ", err)
		}
	}

	wg.Wait()
	fmt.Println("Server gracefully shut downed")
}

func createStorageFolder() {
	err := os.Mkdir(*storagePathPtr, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("ERROR can't create directory: %s", err)
		}
	}
}
