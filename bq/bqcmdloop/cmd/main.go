package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	n := time.Now().Unix()
	f, err := os.Create(fmt.Sprintf("../output/%d.txt", n))
	if err != nil {
		panic(err)
	}

	for {
		cmd := strings.Join(os.Args[1:], " ")
		fmt.Println(cmd)
		fmt.Println()
		_, err = f.WriteString(fmt.Sprintf("%s %s\n", time.Now(), cmd))
		if err != nil {
			fmt.Printf("failed write output: %s\n", err)
		}

		out, err := exec.Command(os.Args[1], os.Args[2:]...).Output()
		if err != nil {
			fmt.Printf("failed command %s\n%s\n", out, err)
			os.Exit(1)
		}

		_, err = f.Write(out)
		if err != nil {
			fmt.Printf("failed write output: %s\n", err)
		}
		_, err = f.WriteString("\n\n")
		if err != nil {
			fmt.Printf("failed write output: %s\n", err)
		}
		time.Sleep(5 * time.Minute)
	}
}
