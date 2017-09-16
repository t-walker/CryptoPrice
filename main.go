package main

import (
  "fmt"
  "time"
  "github.com/jinzhu/gorm" // ORM
    _ "github.com/jinzhu/gorm/dialects/sqlite"
  "net/http"
  "log"
  "encoding/json"
)

type Tick struct {
  gorm.Model
  Name string
  DateTime time.Time
  Price float32
}

type Price struct {
  USD float32 `json:"USD"`
}


func getCurrentPrice(db *gorm.DB) {
  ethPoint := &Price{}
  btcPoint := &Price{}

  makeGetRequest("https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD", ethPoint)
  makeGetRequest("https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD", btcPoint)

  ethTick := Tick{Name: "ETH", DateTime: time.Now(), Price: ethPoint.USD}
  btcTick := Tick{Name: "BTC", DateTime: time.Now(), Price: btcPoint.USD}

  fmt.Printf("eth: %f\n", ethPoint.USD)
  fmt.Printf("btc: %f\n", btcPoint.USD)

  db.Create(&ethTick)
  db.Create(&btcTick)
}


func makeGetRequest(url string, price *Price) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("NewRequest", err)
  }

  client := &http.Client{}

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Do: ", err)
  }

  defer resp.Body.Close()

  json.NewDecoder(resp.Body).Decode(&price)
}


func main() {
  db, err := gorm.Open("sqlite3", "test.db")

  db.AutoMigrate(&Tick{})

  if err != nil {
    panic("failed to connect to database")
  }


  fmt.Printf("App is running\n")
  for true {
    getCurrentPrice(db)
    time.Sleep(30 * time.Second)
  }

  defer db.Close()
}
