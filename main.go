package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Response struct {
	Data    *Pincode `json:"data"`
	OK      bool     `json:"ok"`
	Message string   `json:"message"`
}

type Pincode struct {
	Pincode  int    `json:"pincode"`
	District string `json:"district"`
	Taluk    string `json:"taluk"`
	State    string `json:"state"`
}

func GetPincode() []Pincode {
	pincodes := []Pincode{}

	content, err := ioutil.ReadFile("./public/pincode.json")
	if err != nil {
		fmt.Println(err)
		return pincodes
	}

	err = json.Unmarshal(content, &pincodes)
	if err != nil {
		fmt.Println(err)
		return pincodes
	}

	return pincodes
}

func IndexOf(a []Pincode, x int) int {
	for i, n := range a {
		if x == n.Pincode {
			return i
		}
	}

	return -1
}

var InvalidPincodeMessage = Response{OK: false, Message: "Invalid Pincode"}
var RecordNotFoundMessage = Response{OK: false, Message: "Record Not Found"}

func PincodeHandler(c *fiber.Ctx) error {
	c.Response().Header.EnableNormalizing()

	c.Response().Header.SetContentType("application/json")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("referrer-policy", "no-referrer")
	c.Response().Header.Set("x-frame-options", "SAMEORIGIN")
	c.Response().Header.Set("vary", "origin")
	c.Response().Header.Set("x-content-type-options", "nosniff")
	c.Response().Header.Set("x-dns-prefetch-control", "off")
	c.Response().Header.Set("x-download-options", "noopen")
	c.Response().Header.Set("x-permitted-cross-domain-policies", "none")
	c.Response().Header.Set("x-xss-protection", "0")

	_pincode := c.Params("pincode")
	if len(_pincode) != 6 {
		res, _ := json.Marshal(InvalidPincodeMessage)
		return c.Send(res)
	}

	i, err := strconv.Atoi(_pincode)
	if err != nil {
		fmt.Println(err)
		res, _ := json.Marshal(InvalidPincodeMessage)
		return c.Send(res)
	}

	pincodes := GetPincode()
	idx := IndexOf(pincodes, i)
	if idx == -1 {
		res, _ := json.Marshal(RecordNotFoundMessage)
		return c.Send(res)
	}

	res, _ := json.Marshal(Response{Data: &pincodes[idx], OK: true, Message: "Success"})
	return c.Send(res)
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	app := fiber.New(fiber.Config{
		GETOnly:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Indian Pincode",
		Immutable:     true,
	})

	app.Static(`/static`, "./public")
	app.Use(etag.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        5,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			res, _ := json.Marshal(Response{OK: false, Message: "Rate Limit Reached"})
			c.SendStatus(fiber.StatusTooManyRequests)
			return c.Send(res)
		},
	}))

	app.Use(favicon.New(favicon.Config{File: "./public/favicon.ico"}))

	api := app.Group("/api")
	api.Get("/pincode/:pincode", PincodeHandler)

	app.Get("**", func(c *fiber.Ctx) error {
		res, _ := json.Marshal(Response{OK: false, Message: "Invalid Request"})
		c.SendStatus(fiber.StatusNotFound)
		return c.Send(res)
	})

	app.Listen(":" + port)
}
