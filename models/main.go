package models
import (
    "log"
    "os"
    "strconv"
    "github.com/go-redis/redis"
    "github.com/leekchan/accounting"
)

var RedisClient *redis.Client

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

type PseudoTicker struct {
    High            string
    Low             string
    Last            string
    Buy             string
    Sell            string
}

// Structs-related functions
func (s *Stat) ConvertToPseudoTicker() *PseudoTicker {
    pseudoTicker        := new(PseudoTicker)
    pseudoTicker.High   =   s.Ticker.High
    pseudoTicker.Low    =   s.Ticker.Low
    pseudoTicker.Buy    =   s.Ticker.Buy
    pseudoTicker.Sell   =   s.Ticker.Sell
    pseudoTicker.Last   =   s.Ticker.Last

    return pseudoTicker
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


// Redis-related functions
func StoreUser(username string, firstName string, telegramUID int) {
    recordKey := os.Getenv("REDIS_GLAD_NAMESPACE") + ":telegram:user:" + strconv.Itoa(telegramUID)

    // Check if user already stored, using the key
    val, userCheckingErr := RedisClient.Exists(recordKey).Result()
    if userCheckingErr != nil {
        panic(userCheckingErr)
    }

    // If the user is not found, indicated by val is equal to 0, then store the user
    if val != 1 {
        var user = make(map[string]interface{})
        user["username"] = username
        user["fullname"] = firstName

        storingUserErr := RedisClient.HMSet(recordKey, user).Err()
        if storingUserErr != nil {
            log.Fatal("Failed to store the user")
        }    

        log.Println("Successfully stored the user")
    }
}

func CheckIfTimestampIsCurrent(ticker string, timestampNow string) bool {
    recordKey := keyForTimestampRecord()
    
    statPresence, statCheckingErr := RedisClient.HExists(recordKey, ticker).Result()
    if statCheckingErr != nil {
        panic(statCheckingErr)
    }

    timestampNowInt, _ := strconv.Atoi(timestampNow)

    if statPresence {
        existingTimestamp, _ := RedisClient.HGet(recordKey, ticker).Result()
        existingTimestampInt, _ := strconv.Atoi(existingTimestamp)

        return existingTimestampInt + 5 > timestampNowInt
    } else {
        return false
    }
}

func SetMarketTimestamp(ticker string, timestampNow string) {
    recordKey := keyForTimestampRecord()
    timestampInt, _ := strconv.Atoi(timestampNow)

    err := RedisClient.HSet(recordKey, ticker, timestampInt).Err()
    if err != nil {
        panic(err)
    }

    log.Printf("Successfully updated market timestamp for %s", ticker)
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

func StoreMarketStat(ticker string, stat *Stat, timestampNow string) {
    recordKey := keyForStatRecord(ticker)

    var st = make(map[string]interface{})
    st["high_price"] = stat.Ticker.High
    st["low_price"] = stat.Ticker.Low
    st["latest_price"] = stat.Ticker.Last
    st["buy_price"] = stat.Ticker.Buy
    st["sell_price"] = stat.Ticker.Sell

    err := RedisClient.HMSet(recordKey, st).Err()
    if(err != nil){
        log.Fatal("Failed to store the latest market stat")
    }

    log.Printf("Successfully updated market stat at %s", timestampNow)
}


// private functions
func getMarketStat(ticker string, statAttr string) string {
    recordKey := keyForStatRecord(ticker)

    attr, err := RedisClient.HGet(recordKey, statAttr).Result()
    if err != nil {
        panic(err)
    }

    return attr
}

func keyForTimestampRecord() string {
    return os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:timestamp"   
}

func keyForStatRecord(coinName string) string {
    return os.Getenv("REDIS_GLAD_NAMESPACE") + ":vip:stat:" + coinName
}