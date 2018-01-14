package main
import (
    "log"
    "os"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"

    "github.com/yanzay/tbot"
    _ "github.com/joho/godotenv/autoload"        
)

// The Stat struct should be able to decode a JSON object like this:
// {
//     "ticker": {
//         "high": "16500000",
//         "low": "13611000",
//         "vol_eth": "4453.35873651",
//         "vol_idr": "66475714301",
//         "last": "14301000",
//         "buy": "14301000",
//         "sell": "14375000",
//         "server_time": 1515107696
//     }
// }
type Stat struct {
    Ticker struct {
        High            string `json:"high"`
        Low             string `json:"low"`
        VolEth          string `json:"vol_eth"`
        VolIdr          string `json:"vol_idr"`
        Last            string `json:"last"`
        Buy             string `json:"buy"`
        Sell            string `json:"sell"`
        ServerTime      int    `json:"server_time"`
    }
}

func Run() {
    bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    bot.HandleFunc("/idrs", ListAllIdrCoins)
    bot.HandleFunc("/idr {coin}", RetrieveIdrTicker)

    // Set default handler if you want to process unmatched input
    bot.HandleDefault(UnkownHandler)

    bot.ListenAndServe()
}

func ListAllIdrCoins(message *tbot.Message) {
    coins := []string{"Bitcoin [BTC]", "Bitcoin Cash [BCH]", "Bitcoin Gold [BTG]", "Litecoin [LTC]", "Ethereum [ETH]", "Ethereum Classic [ETC]", "Ripple [XRP]", "Lumens [XLM]", "Waves [WAVES]", "NXT [NXT]", "ZCoin [XZC]"}
    
    for _, coin := range coins {
        message.Reply(coin)
    }
}

func RetrieveIdrTicker(message *tbot.Message) {
    message.Reply("Tunggu sebentar ya")

    vipPublicAPI := os.Getenv("MARKET_API_URL")
    coinTicker := strings.ToLower(message.Vars["coin"])
    upCoinTicker := strings.ToUpper(coinTicker)

    resp, err := http.Get(vipPublicAPI + coinTicker + "_idr/ticker")
    if err != nil {
        message.Reply("Maaf, ada sesuatu yang salah")
        log.Fatal(err)
    }

    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

        stat := Stat{}
        json.Unmarshal([]byte(body), &stat)

        // check if the stat is equal to new empty struct of Stat
        if stat != (Stat{}) {
            message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", upCoinTicker)
    
            time.Sleep(1 * time.Second)
    
            message.Replyf("Harga Tertinggi (24 jam): %s", stat.Ticker.High)
            message.Replyf("Harga Terendah (24 jam): %s", stat.Ticker.Low)
            message.Replyf("Harga Terakhir: %s", stat.Ticker.Last)
            message.Replyf("Harga Beli #1: %s", stat.Ticker.Buy)
            message.Replyf("Harga Jual #1: %s", stat.Ticker.Sell)
            message.Replyf("Volume (24 Jam): %s", stat.Ticker.VolIdr)
            message.Reply("==========")
        } else {
            // The endpoint always return 200 no matter what, so this is basically the handler in case no Ticker was found
            message.Replyf("Maaf, saya tidak bisa mendapatkan info mengenai aktivitas perdagangan IDR-%s", upCoinTicker)
        }
    }
}

func UnkownHandler(message *tbot.Message) {
    message.Reply("Maaf, kami tidak mengerti perintah yang baru saja kamu ketik")
}