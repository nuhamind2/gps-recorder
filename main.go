package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/nats-io/go-nats"
	"encoding/json"
	"time"
	"strings"
)

type Config struct {
	Nats struct {
		Host string
		Port int32
	}
	Pg struct {
		Host     string
		Port     int32
		Username string
		Password string
		Database string
	}
}

func main() {
	fmt.Println("hello world")
	fmt.Println(os.Getenv("GOPATH"))
	data, err := ioutil.ReadFile("config.yml")
	if (err != nil) {
		fmt.Println("Config not found")
		os.Exit(1)
	}
	fmt.Println(string(data))
	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if (err != nil) {
		fmt.Println("config error")
		fmt.Println(err)
		os.Exit(1)
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d", c.Pg.Username, c.Pg.Database, c.Pg.Password, c.Pg.Host, c.Pg.Port)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)

	nc, err := nats.Connect(nats.DefaultURL)
	if (err != nil) {
		fmt.Println("nats connect failed")
		os.Exit(1)
	}

	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	nc.Subscribe("tenant1.clientPublish.*", func(msg *nats.Msg) {
		var result map[string]interface{}
		json.Unmarshal(msg.Data, &result)
		asset_id, ok := result["clientId"].(string)
		if (!ok) {
			fmt.Println("no clientId field")
			os.Exit(1)
		}
		fmt.Println(asset_id)
		gps_data, err := json.Marshal(result["payload"])
		if (err != nil) {
			fmt.Println()
			os.Exit(1)
		}
		//s := `{"topic":"gps","payload":{"lat":-7.794176162141601,"lon":110.20834654044077,"time":1557372081},"clientId":"DEMO1-5054","username":"name5054"}`
		event_id := strings.Split(msg.Subject,".")[2]
		res, err := db.Exec("INSERT INTO gpstrace.gpstrace (event_id,asset_id,gps_data) VALUES($1,$2,$3)", event_id,asset_id,string(gps_data))
		fmt.Println("res")
		fmt.Println(res)
		fmt.Println("err")
		fmt.Println(err)

	})

	for {
		time.Sleep(1)
	}

	nc.Close()
}
