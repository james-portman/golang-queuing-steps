// probably need to use address on struct instance and functions
package main
import (
    "fmt"
    // "time"
    "strings"
    "reflect"
    "os"
)

type Steps struct {
    steps []Step
    input chan string
    output chan string
}

type Step struct {
    function string
    input chan string
    output chan string
}


func (s *Step) run() {
    fmt.Println("running step (function: "+s.function+")")
    go reflect.ValueOf(s).MethodByName(s.function).Call([]reflect.Value{})
}


func (s *Step) Lowercase() {
    for msg := range s.input {
        s.output <- strings.ToLower(msg)
    }
    fmt.Println("finished lowercase step")
    close(s.output)
}

func (s *Step) Uppercase() {
    for msg := range s.input {
        s.output <- strings.ToUpper(msg)
    }
    fmt.Println("finished uppercase step")
    close(s.output)
}


func (s *Steps) addStep(function string) {

    step := Step{function, nil, s.output}

    // check function exists before adding step
    if (reflect.ValueOf(&step).MethodByName(function) == reflect.Value{}) {
        fmt.Println("Couldnt find function "+function+" for step, quitting")
        os.Exit(1)
    }

    // if it's the first step then take input from the whole set of steps
    if len(s.steps) == 0 {
        step.input = s.input

    // else if already some steps then take input from the last
    } else {
        numsteps := len(s.steps)
        last := numsteps - 1

        // make last steps output a new channel instead of using whole steps output
        s.steps[last].output = make(chan string)
        // set this ones input to last ones output
        step.input = s.steps[numsteps-1].output
    }

    s.steps = append(s.steps, step)
}

func (s *Steps) run() {
    for i := range(s.steps) {
        s.steps[i].run()
    }
}



func main() {

    steps := Steps{}
    steps.input = make(chan string)
    steps.output = make(chan string)

    steps.addStep("Lowercase")
    steps.addStep("Uppercase")
    steps.addStep("Lowercase")
    steps.addStep("Uppercase")

    steps.run()

    // put some stuff in
    steps.input <-"test"
    steps.input <-"test2"
    steps.input <-"test3"
    close(steps.input)

    for msg := range steps.output {
        fmt.Println("got output: "+msg)
    }

}
