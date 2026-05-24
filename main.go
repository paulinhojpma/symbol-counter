package response_analyser

// package main

// import (
// 	"context"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/paulinhojpma/symbol-counter/response_analyser"
// )

// func main() {
// 	// fmt.Println("GENERATOR - ", string(symbolGenerator()))
// 	// fmt.Println("------------------------- ")
// 	ctx, cancel := context.WithCancel(context.Background())
// 	m1, m2, m3 := &MyCustomWC{}, &MyCustomWC{}, &MyCustomWC{}
// 	fmt.Println("initializing data analyzer")
// 	da := response_analyser.NewDataAnalyser()

// 	fmt.Println("starting sending")

// 	runWorkers(worker, ctx, da, 100)
// 	fmt.Println("Subscribing")
// 	da.Subscribe(m1, time.Second)
// 	da.Subscribe(m2, time.Second*5)
// 	da.Subscribe(m3, time.Second*20)
// 	time.Sleep(time.Second * 30)
// 	fmt.Printf("MAPS - %+v\n", da.GetCurrentCounts())
// 	fmt.Println("Canceling")
// 	go cancel()
// 	da.Shutdown(ctx)
// 	time.Sleep(time.Second * 5)

// }

// func symbolGenerator() []byte {
// 	simbols := []string{"AU_ieu13", "103956", "ghgqb", "10002", "a012ne"}
// 	simbolsMap := map[string]string{}
// 	size := rand.Intn(20) + 1
// 	simStr := ""
// 	for i := 0; i < size; i++ {
// 		simbolsMap[simbols[rand.Intn(len(simbols))]] = simbols[rand.Intn(len(simbols))]
// 	}
// 	for _, v := range simbolsMap {
// 		if simStr != "" {
// 			simStr = fmt.Sprintf("%s\n%s", simStr, v)

// 		} else {
// 			simStr = v
// 		}
// 	}

// 	return []byte(simStr)

// }
// func worker(ctx context.Context, da response_analyser.DataExporter, workerNum int) {
// 	i := 0
// 	fmt.Printf("starting worker - [%d]\n", workerNum)
// 	for {

// 		i++

// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("closing senders")
// 			fmt.Println("num requests - ", i)
// 			return
// 		default:
// 			symbols := symbolGenerator()
// 			now := time.Now()
// 			da.Analyse(ctx, symbols)
// 			result := time.Since(now)
// 			if result > time.Millisecond*7 {
// 				fmt.Println("limit reach ", result)
// 			}

// 			time.Sleep(time.Microsecond)
// 			//time.Sleep(time.Second)
// 		}

// 	}
// }
// func runWorkers(worker func(ctx context.Context, da response_analyser.DataExporter, workerNum int), ctx context.Context, da response_analyser.DataExporter, num int) {
// 	for i := 0; i < num; i++ {
// 		go worker(ctx, da, i)
// 	}
// }
