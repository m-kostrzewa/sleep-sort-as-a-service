package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sort/{csv_list}", sortHandler)
	http.Handle("/", r)

	logger := logrus.New()
	loggedR := handlers.LoggingHandler(logger.Writer(), r)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      loggedR,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	list := vars["csv_list"]
	listSplit := strings.Split(list, ",")

	nums := []int{}
	for _, v := range listSplit {
		num, err := strconv.Atoi(v)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		nums = append(nums, num)
	}

	sorted := sort(nums)
	w.WriteHeader(http.StatusOK)

	js, err := json.Marshal(sorted)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func sort(nums []int) []int {
	sortedCh := make(chan int, len(nums))
	for _, num := range nums {
		go func(n int) {
			maxTime := 1.0
			sleepTimeSecs := maxTime - maxTime/(float64(n))
			inMilis := sleepTimeSecs * 1000
			time.Sleep(time.Duration(inMilis) * time.Millisecond)
			sortedCh <- n
		}(num)
	}

	sorted := []int{}
	for i := 0; i < len(nums); i++ {
		sorted = append(sorted, <-sortedCh)
	}
	return sorted
}
