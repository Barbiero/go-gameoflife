# Game of Life with WASM + Go

this is an experiment to understand WASM with GO, where I implement [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) in GO while rendering it on a canvas in HTML.

The rendering itself is done from javascript side as I couldn't figure out a way to do it with better performance from the GO side. I use GO channels to calculate an array buffer that will sent the state of the game of life board to JS, and then I call a `drawGameOfLife` function defined in `index.html` that will take the information from the buffer and draw it on the canvas.

From GO's side, I call js's `requestAnimationFrame` in order to wait until I can try to render again; meanwhile, a different goroutine has already calculated the next step on the board. This means that calculating the board and drawing on the canvas are separate routines which will wait for the next one, whichever is ready first.

# Running

simply `make serve`. If you don't have npm, you'll receive an error when trying to serve the contents, but you can simply serve the contents of this folder with any other method you want.
