package main

// test out the goopt package...

import (
	"fmt"
	goopt "github.com/droundy/goopt"
)

var amVerbose = goopt.Bool("--verbose", false, "output verbosely")
var amHappy = goopt.Flag([]string{"-h", "--happy"},
	[]string{"-u", "--unhappy", "--sad"}, "be happy", "be unhappy")

var foo = goopt.String([]string{"--name"}, "anonymous", "pick your name")
var bar = goopt.String([]string{"-b"}, "BOO!", "pick your scary sound")
var speed = goopt.Alternatives([]string{"--speed","--velocity"},
	                             []string{"slow","medium","fast"},
	                             "set the speed")

func main() {
	goopt.Summary = "silly test program"
	goopt.Parse(nil)
	if *amVerbose {
		fmt.Println("I am verbose.")
	}
	if *amHappy {
		fmt.Println("I am happy")
	} else {
		fmt.Println("I am unhappy")
	}
	fmt.Println("Your name is", *foo)
	fmt.Println(*bar, "... Did I scare you?")
	fmt.Println("I am going so very", *speed,"!!!")
}
