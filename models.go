package main
import (
    "log"
    "os"
    "strconv"
    "github.com/go-redis/redis"
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
        VolEth          string `json:"vol_eth"`
        VolIdr          string `json:"vol_idr"`
        Last            string `json:"last"`
        Buy             string `json:"buy"`
        Sell            string `json:"sell"`
        ServerTime      int    `json:"server_time"`
    }
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
    
    statPresence, statCheckingErr := RedisClient.Exists(recordKey).Result()
    if statCheckingErr != nil {
        panic(statCheckingErr)
    }

    timestampNowInt, _ := strconv.Atoi(timestampNow)
    log.Println(statPresence)

    if statPresence == 1 {
        existingTimestamp, _ := RedisClient.HGet(recordKey, ticker).Result()
        existingTimestampInt, _ := strconv.Atoi(existingTimestamp)

        log.Println(timestampNowInt)
        log.Println(existingTimestampInt)
        return timestampNowInt < existingTimestampInt + 5
    } else {
        SetMarketTimestamp(recordKey, ticker, timestampNowInt)

        return false
    }
}

func SetMarketTimestamp(recordKey string, ticker string, timestampInt int) {
    err := RedisClient.HSet(recordKey, ticker, timestampInt).Err()
    if err != nil {
        panic(err)
    }
}