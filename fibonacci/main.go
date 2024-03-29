package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Fib sends the Fibonacci sequence on demand.
func Fib() <-chan *big.Int {
	num := make(chan *big.Int)
	a := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(1)
	go func() {
		for {
			num <- a
			a, b = b, new(big.Int).Add(a, b)
		}
	}()
	return num
}

// FibNth computes the nth Fibonacci number using fewer allocations.
func FibNth(n uint64) *big.Int {
	a := new(big.Int).SetUint64(0)
	b := new(big.Int).SetUint64(1)
	for i := uint64(0); i < n; i++ {
		a, b = b, a.Add(a, b)
	}
	return a
}

// FibCached returns a func that computes n Fibonacci numbers and caches results.
func FibCached() func(uint64) []*big.Int {
	var nums []*big.Int
	num := Fib()
	return func(n uint64) []*big.Int {
		for i := uint64(len(nums)); i < n; i++ {
			nums = append(nums, <-num)
		}
		return nums[:n]
	}
}

func serveFib() func(http.ResponseWriter, *http.Request, httprouter.Params) {
	fib := FibCached()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		nParam := p.ByName("n")
		n, err := strconv.ParseInt(nParam, 10, 64)
		if err != nil {
			writeResponse(w, nil, fmt.Errorf("n must be integer, got %s", nParam))
			return
		}
		if n < 0 {
			writeResponse(w, nil, fmt.Errorf("n must be positive, got %d", n))
			return
		}
		writeResponse(w, fib(uint64(n)), nil)
	}
}

type fibResponse struct {
	Nums []*big.Int `json:"nums"`
	Err  error      `json:"error"`
}

func writeResponse(w http.ResponseWriter, nums []*big.Int, err error) {
	b, err := json.Marshal(&fibResponse{nums, err})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

func main() {
	router := httprouter.New()
	router.GET("/api/fibonacci/:n", serveFib())
	log.Fatal(http.ListenAndServe(":8080", router))
}
