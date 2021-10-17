package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	Name, Port   string
	ResponseTime int
}

func (s Service) HandleRequest(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(s.ResponseTime) * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Response from %v \n", s.Name)))
}

var (
	serviceName  = flag.String("svc", "", "Name of service")
	port         = flag.String("p", "", "HTTP port to listen to")
	responseTime = flag.Int("rt", 0, "Time in ms to wait before response")
)

func main() {
	flag.Parse()
	svc := Service{
		Name:         *serviceName,
		Port:         *port,
		ResponseTime: *responseTime,
	}
	http.HandleFunc("/", svc.HandleRequest)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", svc.Port), nil)
		if err != nil {
			fmt.Printf("Error listening on server : %v", svc.Name)
			return
		}
	}()
	fmt.Printf("%v Listening on port : %v , Press <Enter> to exit \n", svc.Name, svc.Port)
	fmt.Scanln()
}
