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

    "./models"
)


func Run() {
    // Here we initialize the db and then assign it to a global var of RedisClient which is of type *gorm.DB
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()

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
    models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)

    coins := []string{"Bitcoin [BTC]", "Bitcoin Cash [BCH]", "Bitcoin Gold [BTG]", "Litecoin [LTC]", "Ethereum [ETH]", "Ethereum Classic [ETC]", "Ripple [XRP]", "Lumens [XLM]", "Waves [WAVES]", "NXT [NXT]", "ZCoin [XZC]"}

    for _, coin := range coins {
        message.Reply(coin)
    }
}

func RetrieveIdrTicker(message *tbot.Message) {
    models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)
    
    vipPublicAPI := os.Getenv("MARKET_API_URL")
    coinTicker := strings.ToLower(message.Vars["coin"])
    upCoinTicker := strings.ToUpper(coinTicker)

    timeNow := time.Now().UTC()
    timestampNow := timeNow.Format("200601021504")
    shouldGetLatest := !models.CheckIfTimestampIsCurrent(coinTicker, timestampNow)

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

            stat := models.Stat{}
            json.Unmarshal([]byte(body), &stat)

            // check if the stat is not equal to new empty struct of Stat
            if stat != (models.Stat{}) {
                models.StoreMarketStat(coinTicker, &stat, timestampNow)
                models.SetMarketTimestamp(coinTicker, timestampNow)

                pseudoTicker := stat.ConvertToPseudoTicker()

                relayStats(upCoinTicker, message, pseudoTicker)
            } else {
                // The endpoint always return 200 no matter what, so this is basically the handler in case no Ticker was found
                message.Replyf("Maaf, saya tidak bisa mendapatkan info mengenai aktivitas perdagangan IDR-%s", upCoinTicker)
            }
        }
    } else {
        stat := models.RetrieveMarketStats(coinTicker)

        relayStats(upCoinTicker, message, stat)
    }
}

func UnkownHandler(message *tbot.Message) {
    models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)
        
    message.Reply("Maaf, kami tidak mengerti perintah yang baru saja kamu ketik")
}

// Private methods //
func relayStats(ticker string, message *tbot.Message, stat *models.PseudoTicker) {
    message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", ticker)

    time.Sleep(2 * time.Second)

    stat.DisplayAsMoney()

    message.Replyf("Harga Terakhir: %s", stat.Last)
    message.Replyf("Harga Beli #1: %s", stat.Buy)
    message.Replyf("Harga Jual #1: %s", stat.Sell)
    message.Replyf("Harga Tertinggi (24 jam): %s", stat.High)
    message.Replyf("Harga Terendah (24 jam): %s", stat.Low)
    message.Reply("==========")
}