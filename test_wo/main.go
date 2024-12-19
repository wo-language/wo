package main

import (
	"fmt"
	"os"
)

func main() {

	var f, err = os.Open("bang.wo")
	//var f, err = os!Open("bang.wo")

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)

}
