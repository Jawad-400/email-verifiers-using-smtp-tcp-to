package main

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	

	"github.com/gin-gonic/gin"  // ← JUST THIS LINE
)

func main() {
	r := gin.Default()
	r.POST("/check", checkEmail)

	fmt.Println("Server running → http://localhost:8081/check")
	fmt.Println("Send: {\"email\": \"test@aol.com\"}")
	r.Run(":8081")
}

func checkEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BindJSON(&req); err != nil || req.Email == "" {
		c.JSON(400, gin.H{"result": "Invalid request"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	if !strings.HasSuffix(email, "@aol.com") {
		c.JSON(200, gin.H{"email": email, "result": "Not an AOL email"})
		return
	}

	exists, err := verifyAOL(email)

	if err != nil {
		c.JSON(200, gin.H{
			"email":  email,
			"result": "Cannot verify (blocked by AOL)",
		})
		return
	}

	if exists {
		c.JSON(200, gin.H{
			"email":  email,
			"result": "Exists",
		})
	} else {
		c.JSON(200, gin.H{
			"email":  email,
			"result": "Does NOT exist",
		})
	}
}

func verifyAOL(email string) (bool, error) {
	mx, _ := net.LookupMX("aol.com")
	if len(mx) == 0 {
		return false, fmt.Errorf("no MX")
	}

	conn, err := smtp.Dial(mx[0].Host + ":25")
	if err != nil {
		return false, err
	}
	defer conn.Close()

	conn.Mail("test@localhost")
	err = conn.Rcpt(email)

	if err != nil {
		if strings.Contains(err.Error(), "550") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}