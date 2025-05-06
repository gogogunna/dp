package main

import "fmt"

type s struct {
	g int
}

func main() {
	b := []s{
		{
			g: 1,
		},
		{
			g: 2,
		},
	}

	for i := range b {
		b[i].g = 120
	}

	fmt.Println(b)
}
