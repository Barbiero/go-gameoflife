package main

import (
	"log"
	"os"
	"syscall/js"
)

func main() {
	globalThis := js.Global()

	w, h := 200, 200
	gol := &GameOfLife{
		width:      w,
		height:     h,
		board:      [][]State{},
		boardReady: make(chan []byte, 2),
	}
	gol.Init()

	js.CopyBytesToJS(globalThis.Get("golBuffer"), gol.flatBoard())
	globalThis.Call("drawGameOfLife")

	// Allows js to stop the go program
	stopCh := make(chan struct{}, 1)

	globalThis.Set("stopGo", js.FuncOf(func(this js.Value, args []js.Value) any {
		stopCh <- struct{}{}
		return js.Undefined()
	}))
	defer globalThis.Set("stopGo", js.Undefined())

	// ticker := time.Tick(time.Millisecond * 1000)
	gol.RunStep()
	promiseCh := make(chan []js.Value, 1)
	for {
		select {
		case <-stopCh:
			log.Println("stopping due to signal")
			os.Exit(0)
		case board := <-gol.boardReady:
			gol.RunStep()

			js.CopyBytesToJS(globalThis.Get("golBuffer"), board)
			globalThis.Call("drawGameOfLife")

			// await animation frame
			var promise = globalThis.Get("Promise").New(js.FuncOf(func(_ js.Value, args []js.Value) any {
				res := args[0]

				globalThis.Call("requestAnimationFrame", res)
				return js.Undefined()
			}))
			promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
				promiseCh <- args
				return js.Undefined()
			}))
			<-promiseCh
		}
	}
}
