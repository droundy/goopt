package main

// test out the goopt package...

import (
	"fmt"
	"strings"
	goopt "github.com/droundy/goopt"
)

var amHappy = goopt.Flag([]string{"-h", "--happy"}, []string{"-u", "--unhappy", "--sad"}, false, "be happy", "be unhappy")

var foo = goopt.String([]string{"--name"}, "anonymous", "pick your name")
var bar = goopt.String([]string{"-b"}, "BOO!", "pick your scary sound")
var speed = goopt.Alternatives([]string{"--speed", "--velocity"},
	[]string{"slow", "medium", "fast"},
	"set the speed")
var list = goopt.Strings([]string{"--list", "-l"}, "add", "Add words to the word list")

func main() {
	goopt.Summary = "silly test program"
	goopt.Parse(nil)
	if *amHappy {
		fmt.Println("I am happy")
	} else {
		fmt.Println("I am unhappy")
	}
	fmt.Println("Your name is", *foo)
	fmt.Println(*bar, "... Did I scare you?")
	fmt.Println("I am going so very", *speed, "!!!")
	fmt.Println("I really like words like: " + strings.Join(list.Data(), " "))
}
