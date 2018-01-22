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


func Run() {
    // Here we initialize the db and then assign it to a global var of RedisClient which is of type *gorm.DB
    // as defined in models.go
    RedisClient = InitializeDatabase()

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
    StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)

    coins := []string{"Bitcoin [BTC]", "Bitcoin Cash [BCH]", "Bitcoin Gold [BTG]", "Litecoin [LTC]", "Ethereum [ETH]", "Ethereum Classic [ETC]", "Ripple [XRP]", "Lumens [XLM]", "Waves [WAVES]", "NXT [NXT]", "ZCoin [XZC]"}

    for _, coin := range coins {
        message.Reply(coin)
    }
}

func RetrieveIdrTicker(message *tbot.Message) {
    StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)
    
    vipPublicAPI := os.Getenv("MARKET_API_URL")
    coinTicker := strings.ToLower(message.Vars["coin"])
    upCoinTicker := strings.ToUpper(coinTicker)

    timeNow := time.Now().UTC()
    timestampNow := timeNow.Format("200601021504")
    shouldGetLatest := !CheckIfTimestampIsCurrent(coinTicker, timestampNow)

    if shouldGetLatest {
        log.Println("Sending enquiry to Market.....")

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

            // check if the stat is not equal to new empty struct of Stat
            if stat != (Stat{}) {
                StoreMarketStat(coinTicker, &stat, timestampNow)
                SetMarketTimestamp(coinTicker, timestampNow)

                message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", upCoinTicker)
        
                time.Sleep(2 * time.Second)
        
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
    } else {
        high := RetrieveMarketStat(coinTicker, "high_price")
        low  := RetrieveMarketStat(coinTicker, "low_price")
        last := RetrieveMarketStat(coinTicker, "latest_price")
        buy  := RetrieveMarketStat(coinTicker, "buy_price")
        sell := RetrieveMarketStat(coinTicker, "sell_price")

        message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", upCoinTicker)

        time.Sleep(2 * time.Second)

        message.Replyf("Harga Tertinggi (24 jam): %s", high)
        message.Replyf("Harga Terendah (24 jam): %s", low)
        message.Replyf("Harga Terakhir: %s", last)
        message.Replyf("Harga Beli #1: %s", buy)
        message.Replyf("Harga Jual #1: %s", sell)
        message.Reply("==========")        
    }
}

func UnkownHandler(message *tbot.Message) {
    StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)
        
    message.Reply("Maaf, kami tidak mengerti perintah yang baru saja kamu ketik")
}