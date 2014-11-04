package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Conf struct {
	Subdomain string
	Token     string
	IsMember  bool
	Emails    []string
}

type Profiles struct {
	Status int32 "json:`status`"
	Data   []struct {
		ProfileName string `json:"profileName"`
		Meta        struct {
			Links []struct {
				Href string `json:"href"`
			} `json:"links"`
		} `json:"meta"`
	} "json:`data`"
	Error string "json:`error`"
}

func main() {
	// Load config file
	f, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Config file error: %v\n", err)
		os.Exit(1)
	}
	d := json.NewDecoder(f)
	conf := Conf{}
	err = d.Decode(&conf)
	if err != nil {
		fmt.Printf("Config parsing error: %v\n", err)
		f.Close()
		os.Exit(1)
	}
	f.Close()

	if conf.Subdomain == "" {
		fmt.Errorf("Please set the 'Subdomain' in the config file to the subdomain of your Microcosm site.")
		os.Exit(1)
	}

	if conf.Token == "" {
		fmt.Errorf("Please set the 'Token' in the config file to the 'access_token' cookie of your signed in admin user or moderator.")
		os.Exit(1)
	}

	if len(conf.Emails) == 0 {
		fmt.Errorf("Please provide one or more 'Emails' in the config file.")
		os.Exit(1)
	}

	fmt.Println("config.json loaded")

	// Fetch (or create) profile info for each email address
	var emails string
	for i, e := range conf.Emails {
		if i > 0 {
			emails += ","
		}
		emails += fmt.Sprintf(`{"email":"%s"}`, e)

	}
	emails = "[" + emails + "]"

	req, err := http.NewRequest(
		"POST",
		"https://"+conf.Subdomain+".microco.sm/api/v1/users?access_token="+conf.Token,
		bytes.NewBufferString(emails),
	)
	if err != nil {
		fmt.Printf("Get user info request error: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Get user info error: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Get users error = %s\n", resp.Status)
		os.Exit(1)
	}

	// Load profiles
	d = json.NewDecoder(resp.Body)
	p := Profiles{}
	err = d.Decode(&p)
	if err != nil {
		fmt.Printf("Profiles parsing error: %v\n", err)
		os.Exit(1)
	}

	// Set or unset the given var
	for _, m := range p.Data {
		url := "https://" + conf.Subdomain + ".microco.sm" + m.Meta.Links[0].Href + "/attributes/is_member?access_token=" + conf.Token

		if conf.IsMember {
			fmt.Printf("Adding: %s ", m.ProfileName)
			req, err = http.NewRequest("PUT", url, bytes.NewBufferString(`{"value": true}`))
		} else {
			fmt.Printf("Removing: %s ", m.ProfileName)
			req, err = http.NewRequest("DELETE", url, nil)
		}
		if err != nil {
			fmt.Printf("\nRequest error: %v\n", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("\nClient error: %v\n", err)
			os.Exit(1)
		} else {
			if resp.StatusCode == http.StatusOK {
				fmt.Printf("OK\n")
			} else {
				if !conf.IsMember && resp.StatusCode == http.StatusNotFound {
					fmt.Printf("OK\n")
				}
			}
		}
	}
}
