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

var vipPublicAPI = os.Getenv("MARKET_API_URL")

func Run() {
    // Here we initialize the db and then assign it to a global var of RedisClient which is of type *gorm.DB
    // as defined in models.go
    models.RedisClient = models.InitializeDatabase()

    bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    go bot.HandleFunc("/koin rupiah", ListAllIdrCoins)
    go bot.HandleFunc("/harga {coin}", RetrieveIdrTicker)

    // Set default handler if you want to process unmatched input
    go bot.HandleDefault(UnkownHandler)

    bot.ListenAndServe()
}

func ListAllIdrCoins(message *tbot.Message) {
    go models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)

    coins := []string{"Bitcoin [btc]", "Bitcoin Cash [bch]", "Bitcoin Gold [btg]", "Litecoin [ltc]", "Ethereum [eth]", "Ethereum Classic [etc]", "Ripple [xrp]", "Lumens [xlm]", "Waves [waves]", "NXT [nxt]", "ZCoin [xzc]"}

    message.Reply("Berikut adalah daftar koin yang nilai tukarnya dengan rupiah dapat kamu tanya ke saya")
    for _, coin := range coins {
        message.Reply(coin)
    }
    message.Reply("==========")
}

func RetrieveIdrTicker(message *tbot.Message) {
    coinTicker := strings.ToLower(message.Vars["coin"])
    timestampNow := getTimestampNow()

    shouldGetLatest := !models.CheckIfTimestampIsCurrent(coinTicker, timestampNow)

    go models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)

    if shouldGetLatest {
        resp := sendRequestToFetchTicker(coinTicker, message)
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            stat := parseResponseAsTicker(resp)

            // check if the stat is not equal to new empty struct of Stat
            if stat != (models.Stat{}) {
                go models.StoreMarketStat(coinTicker, &stat, timestampNow)
                go models.SetMarketTimestamp(coinTicker, timestampNow)

                pseudoTicker := stat.ConvertToPseudoTicker()

                relayStats(coinTicker, message, pseudoTicker)
            } else {
                // The endpoint always return 200 no matter what, so this is basically the handler in case no Ticker was found
                message.Replyf("Maaf, saya tidak bisa mendapatkan info mengenai aktivitas perdagangan IDR-%s", strings.ToUpper(coinTicker))
            }
        }
    } else {
        stat := models.RetrieveMarketStats(coinTicker)

        relayStats(coinTicker, message, stat)
    }
}

func UnkownHandler(message *tbot.Message) {
    go models.StoreUser(message.From.UserName, message.From.FirstName, message.From.LastName, message.From.ID)
    message.Replyf("Maaf %s, saya tidak mengerti perintah yang baru saja kamu ketik", message.From)
}



//// Private methods ////

func relayStats(ticker string, message *tbot.Message, stat *models.PseudoTicker) {
    message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", strings.ToUpper(ticker))

    stat.DisplayAsMoney()

    message.Replyf("Harga Terakhir: %s", stat.Last)
    message.Replyf("Harga Beli #1: %s", stat.Buy)
    message.Replyf("Harga Jual #1: %s", stat.Sell)
    message.Replyf("Harga Tertinggi (24 jam): %s", stat.High)
    message.Replyf("Harga Terendah (24 jam): %s", stat.Low)
    message.Reply("==========")
}

func getTimestampNow() string {
    timeNow := time.Now().UTC()
    return timeNow.Format("200601021504")
}

// HTTP-related methods
func sendRequestToFetchTicker(coinTicker string, message *tbot.Message) *http.Response {
    resp, err := http.Get(vipPublicAPI + coinTicker + "_idr/ticker")
    if err != nil {
        message.Reply("Maaf, saya gagal mendapatkan data terbaru")        
        log.Fatal(err)
    }

    return resp
}

func parseResponseAsTicker(resp *http.Response) models.Stat {
    log.Println("Sending enquiry to Market.....")

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    stat := models.Stat{}
    json.Unmarshal([]byte(body), &stat)

    return stat
}
