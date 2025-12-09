package main

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/check", checkEmail)
	fmt.Println("DIRECT SMTP RCPT TO VERIFIER → http://localhost:8081/check")
	r.Run(":8081")

}

func checkEmail(c *gin.Context) {

	var req struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil || req.Email == "" {

		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return

	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	result := verifyEmailDirect(email)

	c.JSON(200, gin.H{

		"email":  email,
		"result": result,
	})
}

// THIS IS THE REAL DIRECT RCPT TO METHOD — WORKS INSTANTLY
func verifyEmailDirect(email string) string {

	parts := strings.Split(email, "@")

	if len(parts) != 2 {

		return "Invalid format"

	}

	domain := parts[1]

	mx, err := net.LookupMX(domain)

	if err != nil || len(mx) == 0 {

		return "No mail server (invalid domain)"

	}

	for _, server := range mx {

		host := strings.TrimSuffix(server.Host, ".")

		conn, err := net.DialTimeout("tcp", host+":25", 5*time.Second)
		if err != nil {

			fmt.Println(err)
			fmt.Println(45)

			continue

		}

		defer conn.Close()

		client, err := smtp.NewClient(conn, host)

		if err != nil {
			fmt.Println(err)
			fmt.Println(50)

			continue
		}

		defer client.Quit()

		if err := client.Hello("checker.local"); err != nil {

			continue

		}

		if err := client.Mail("test@checker.local"); err != nil {
			fmt.Println(err)
			fmt.Println(55)

			continue

		}

		//THIS IS RCPT TO
		err = client.Rcpt(email)
		fmt.Println(err)
		fmt.Println(60)

		if err == nil {
			fmt.Println(err)

			return "EXISTS (ACTIVE & DELIVERABLE)"

		}

		if strings.Contains(err.Error(), "550") ||
			strings.Contains(err.Error(), "user unknown") ||
			strings.Contains(err.Error(), "no such user") {
			return "DOES NOT EXIST"

		}

	}

	return "BLOCKED or TEMPORARY ERROR"

}
