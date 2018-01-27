package models
import (
    "log"
    "os"
    "strconv"
    "github.com/go-redis/redis"
    "github.com/leekchan/accounting"
)

var RedisClient *redis.Client

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
        VolIdr          string `json:"vol_idr"`
        Last            string `json:"last"`
        Buy             string `json:"buy"`
        Sell            string `json:"sell"`
        ServerTime      int    `json:"server_time"`
    }
}

func (s *Stat) ConvertToPseudoTicker() *PseudoTicker {
    pseudoTicker        := new(PseudoTicker)
    pseudoTicker.High   =   s.Ticker.High
    pseudoTicker.Low    =   s.Ticker.Low
    pseudoTicker.Buy    =   s.Ticker.Buy
    pseudoTicker.Sell   =   s.Ticker.Sell
    pseudoTicker.Last   =   s.Ticker.Last

    return pseudoTicker
}

type PseudoTicker struct {
    High            string
    Low             string
    Last            string
    Buy             string
    Sell            string
}

func (p *PseudoTicker) DisplayAsMoney() {
    highInt, _ := strconv.Atoi(p.High)
    lowInt, _ := strconv.Atoi(p.Low)
    lastInt, _ := strconv.Atoi(p.Last)
    buyInt, _ := strconv.Atoi(p.Buy)
    sellInt, _ := strconv.Atoi(p.Sell)

    ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}
    
    p.High  = ac.FormatMoneyInt(highInt)
    p.Low   = ac.FormatMoneyInt(lowInt)
    p.Last  = ac.FormatMoneyInt(lastInt)
    p.Buy   = ac.FormatMoneyInt(buyInt)
    p.Sell  = ac.FormatMoneyInt(sellInt)
}

func InitializeDatabase() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    pong, err := client.Ping().Result()
    log.Println(pong, err)

    return client
}

func StoreUser(username string, firstName string, lastName string, telegramUID int) {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":telegram:user:" + strconv.Itoa(telegramUID)

    // Check if user already stored, using the key
    val, userCheckingErr := RedisClient.Exists(recordKey).Result()
    if userCheckingErr != nil {
        panic(userCheckingErr)
    }

    // If the user is not found, indicated by val is equal to 0, then store the user
    if val != 1 {
        fullName := firstName + " " + lastName

        userNameStoringErr := RedisClient.HSet(recordKey, "username", username).Err()
        if userNameStoringErr != nil {
            panic(userNameStoringErr)
        }    
        fullNameStoringErr := RedisClient.HSet(recordKey, "fullname", fullName).Err()
        if fullNameStoringErr != nil {
            panic(fullNameStoringErr)
        }

        log.Println("Successfully stored the user")
    }
}

func CheckIfTimestampIsCurrent(ticker string, timestampNow string) bool {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:timestamp"
    
    statPresence, statCheckingErr := RedisClient.HExists(recordKey, ticker).Result()
    if statCheckingErr != nil {
        panic(statCheckingErr)
    }

    timestampNowInt, _ := strconv.Atoi(timestampNow)
    log.Println(statPresence)

    if statPresence {
        existingTimestamp, _ := RedisClient.HGet(recordKey, ticker).Result()
        existingTimestampInt, _ := strconv.Atoi(existingTimestamp)

        log.Println(timestampNowInt)
        log.Println(existingTimestampInt)
        return existingTimestampInt + 5 > timestampNowInt
    } else {
        return false
    }
}

func getMarketStat(ticker string, statAttr string) string {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:" + ticker

    attr, err := RedisClient.HGet(recordKey, statAttr).Result()
    if err != nil {
        panic(err)
    }

    return attr
}

func RetrieveMarketStats(ticker string) *PseudoTicker {
    high := getMarketStat(ticker, "high_price")
    low  := getMarketStat(ticker, "low_price")
    last := getMarketStat(ticker, "latest_price")
    buy  := getMarketStat(ticker, "buy_price")
    sell := getMarketStat(ticker, "sell_price")

    pseudoTicker := new(PseudoTicker)
    pseudoTicker.High =   high
    pseudoTicker.Low =   low
    pseudoTicker.Buy =   buy
    pseudoTicker.Sell =   sell
    pseudoTicker.Last =   last

    return pseudoTicker
}

func SetMarketTimestamp(ticker string, timestampNow string) {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:timestamp"
    timestampInt, _ := strconv.Atoi(timestampNow)

    err := RedisClient.HSet(recordKey, ticker, timestampInt).Err()
    if err != nil {
        panic(err)
    }

    log.Printf("Successfully updated market timestamp for %s", ticker)
}

func StoreMarketStat(ticker string, stat *Stat, timestampNow string) {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:" + ticker
    
    statCheckingErr := RedisClient.Exists(recordKey).Err()
    if statCheckingErr != nil {
        panic(statCheckingErr)
    }

    newHighReplyErr := RedisClient.HSet(recordKey, "high_price", stat.Ticker.High).Err()
    newLowReplyErr := RedisClient.HSet(recordKey, "low_price", stat.Ticker.Low).Err()
    newLatestReplyErr := RedisClient.HSet(recordKey, "latest_price", stat.Ticker.Last).Err()
    newBuyReplyErr := RedisClient.HSet(recordKey, "buy_price", stat.Ticker.Buy).Err()
    newSellReplyErr := RedisClient.HSet(recordKey, "sell_price", stat.Ticker.Sell).Err()

    if newHighReplyErr != nil || newLowReplyErr != nil || newLatestReplyErr != nil || newBuyReplyErr != nil || newSellReplyErr != nil {
        log.Fatal("Failed to store the latest market stat")
    }

    log.Printf("Successfully updated market stat at %s", timestampNow)
}