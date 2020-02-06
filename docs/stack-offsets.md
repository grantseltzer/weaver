# Stack offsets

<i>Disclaimer: this doc is based on observations using the cmd/print-stack tool in go 1.13 on linux x86_64. This is subject to change with different versions of Go.</i>


- Look at size of largest data type that's being passed, that sets the window size. The maximum value for the window is 8.

- Each element added is limited by whether or not it will fit into that window.

- If it would go over a limit window then pad until back at 0, add it, then continue.
