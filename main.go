// package main
package main

import (
	ctrlr "LeafMS-BackEnd/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	//uncomment the following method only when you want to insert public holidays of a country for a specific year
	// utils.PersistPublicHolidays(2024, "IN")

	routes := mux.NewRouter()
	routes.HandleFunc("/login", ctrlr.HandleLogin).Methods("GET")
	authRoute := routes.NewRoute().Subrouter()
	authRoute.Use(ctrlr.HandleAuth)
	authRoute.HandleFunc("/apply", ctrlr.HandleApply).Methods("PUT")
	authRoute.HandleFunc("/leaves", ctrlr.HandleViewLeaves).Methods("GET")
	// authROute.HandleFunc("/team-leaves", ctrlr.Ha)
	authRoute.HandleFunc("/applications", ctrlr.HandleViewLeaveApplications).Methods(("GET"))
	authRoute.HandleFunc("/approval-based-applications", ctrlr.HandleViewLeavesBasedOnApproval).Methods(("GET"))
	authRoute.HandleFunc("/approve", ctrlr.HandleLeaveApproval).Methods("PATCH")
	authRoute.HandleFunc("/holidays", ctrlr.HandleViewHolidays).Methods("GET")

	http.ListenAndServe(":8080", routes)

	// var sli []int = []int{1, 2, 3}
	// var channel = make(chan int)

	// // var wg *sync.WaitGroup
	// // wg.Add(1)
	// go addMeth(sli, channel)

	// extractedVal <- channel
	// fmt.Println(extractedVal)

	// var val = 5
	// var stringVal = "stringify132"
	// me(val)
	// me(stringVal)

	// wg.Wait()

	// var sli1 = []int{1, 2, 3, 4, 5}
	// var sli2 = []int{3, 4, 5, 6, 1, 2}

	// var sli3 = []int{8, 7, 6, 5, 4, 3}
	// var sli4 = []int{4, 3, 8, 7, 6, 5}
	// var sli5 = []int{4, 3, 7, 8, 6, 5}
	// fmt.Println(checkRotation(sli1, sli2))
	// fmt.Println(checkRotation(sli3, sli4))
	// fmt.Println(checkRotation(sli3, sli5))
	// fmt.Println(checkRotation(sli5, sli4))
}

// func me(val interface{}) {

// 	switch val.(type) {
// 	case string:
// 		fmt.Println("string")
// 	case int:
// 		fmt.Println("int")
// 	default:
// 		fmt.Println("default")
// 	}
// }

// func addMeth(sli []int, ch chan int) {
// 	var res = 0
// 	for _, v := range sli {
// 		res += v
// 	}

// 	fmt.Println(ch)
// 	ch <- res
// 	fmt.Println(ch)
// 	// wg.Done()
// }

// func checkRotation(sli1 []int, sli2 []int) bool {
// 	var start = sli1[0]
// 	var rotationLength = 0
// 	for index, val := range sli2 {
// 		if start == val {
// 			rotationLength = index
// 		}
// 	}

// 	for index, val := range sli1 {
// 		var maxIndexRange = len(sli2) - 1
// 		var updatedIndex = index + rotationLength
// 		if updatedIndex > maxIndexRange {
// 			updatedIndex = updatedIndex - maxIndexRange - 1
// 		}
// 		if val != sli2[updatedIndex] {
// 			return false
// 		}
// 	}
// 	return true

// }

// 1,2.3.4.5
// 3,4,5,1,2

//8,7,6,5,4,3

// import "fmt"

// func generator(nums ...int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		for _, n := range nums {
// 			out <- n
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// func square(in <-chan int) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		for n := range in {
// 			out <- n * n
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// func main() {
// 	nums := generator(1, 2, 3, 4, 5)
// 	squares := square(nums)

// 		for n := range squares {
// 			fmt.Println(n)
// 		}
// 	}
// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type PubSub struct {
// 	subscribers []chan string
// 	mu          sync.Mutex
// }

// func (ps *PubSub) Subscribe() <-chan string {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()

// 	ch := make(chan string, 5) // Buffered channel
// 	ps.subscribers = append(ps.subscribers, ch)
// 	return ch
// }

// func (ps *PubSub) Publish(msg string) {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()

// 	for _, sub := range ps.subscribers {
// 		sub <- msg
// 	}
// }

// func main() {
// 	ps := &PubSub{}

// 	// Subscribe
// 	sub1 := ps.Subscribe()
// 	sub2 := ps.Subscribe()

// 	// Publisher sends messages
// 	go func() {
// 		for i := 1; i <= 3; i++ {
// 			ps.Publish(fmt.Sprintf("Message %d", i))
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	// Subscribers receive messages
// 	go func() {
// 		for msg := range sub1 {
// 			fmt.Println("Subscriber 1 received:", msg)
// 		}
// 	}()
// 	go func() {
// 		for msg := range sub2 {
// 			fmt.Println("Subscriber 2 received:", msg)
// 		}
// 	}()

// 	time.Sleep(5 * time.Second)
// }
