package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Data    Pincode `json:"data"`
	OK      bool    `json:"ok"`
	Message string  `json:"message"`
}

type Pincode struct {
	Pincode  int    `json:"pincode"`
	District string `json:"district"`
	Taluk    string `json:"taluk"`
	State    string `json:"state"`
}

func GetPincode() []Pincode {
	pincodes := []Pincode{}

	content, err := ioutil.ReadFile("pincode.json")
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
		res, _ := json.Marshal(Response{OK: false, Message: "Invalid Pincode"})
		return c.Send(res)
	}

	i, err := strconv.Atoi(_pincode)
	if err != nil {
		fmt.Println(err)
		res, _ := json.Marshal(Response{OK: false, Message: "Invalid Pincode"})
		return c.Send(res)
	}

	pincodes := GetPincode()
	idx := IndexOf(pincodes, i)
	if idx == -1 {
		res, _ := json.Marshal(Response{OK: false, Message: "Record Not Found"})
		return c.Send(res)
	}

	res, _ := json.Marshal(Response{Data: pincodes[idx], OK: true, Message: "Success"})

	return c.Send(res)
}

func main() {
	app := fiber.New(fiber.Config{
		GETOnly:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Indian Pincode",
		Immutable:     true,
	})

	api := app.Group("/api")

	api.Get("/pincode/:pincode", PincodeHandler)

	app.Listen(":3000")
}
