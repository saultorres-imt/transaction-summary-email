package main

import (
	"fmt"
	"log"
)

const (
	subject = "Summary of account transactions"
)

func main() {
	if err := sendEmail("saultorres.imt@outlook.com", "saulftg.22@gmail.com", subject, "Testing SES"); err != nil {
		log.Fatalf("Fail to send email %v", err)
	}
	fmt.Print("Email sent\n")
}
