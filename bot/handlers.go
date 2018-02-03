package bot

import (
    "strings"
    
    "github.com/yanzay/tbot"

    "../models"
)

func unkownHandler(message *tbot.Message) {
    message.Replyf("Maaf %s, saya tidak mengerti permintaanmu barusan", message.From.FirstName)
    message.Reply("Kamu dapat melihat perintah yang tersedia dengan mengetik /help")
}


func customHelpHandler(message *tbot.Message) {
    message.Reply("Inilah daftar perintah yang tersedia untuk kamu:")
    message.Reply("/koin - Untuk mengetahui semua koin yang dapat ditanyakan harganya")
    message.Reply("/harga NAMA_SINGKAT_KOIN || contohnya, ketik: /harga btc")
}


func listAllIdrCoins(message *tbot.Message) {
    coins := []string{"Bitcoin [btc]", "Bitcoin Cash [bch]", "Bitcoin Gold [btg]", "Litecoin [ltc]", "Ethereum [eth]", "Ethereum Classic [etc]", "Ripple [xrp]", "Lumens [xlm]", "Waves [waves]", "NXT [nxt]", "ZCoin [xzc]"}

    message.Reply("Berikut adalah daftar koin yang nilai tukarnya dengan rupiah dapat kamu tanya ke saya")
    for _, coin := range coins {
        message.Reply(coin)
    }
    message.Reply("==========")
}


func retrieveIdrTradeStat(message *tbot.Message) {
    coinTicker := strings.ToLower(message.Vars["coin"])

    message.Replyf("Berikut info mengenai aktivitas perdagangan IDR-%s", strings.ToUpper(coinTicker))

    stat := models.RetrieveMarketStats(coinTicker)
    stat.DisplayAsMoney()

    message.Replyf("Harga Terakhir: %s", stat.Last)
    message.Replyf("Harga Beli #1: %s", stat.Buy)
    message.Replyf("Harga Jual #1: %s", stat.Sell)
    message.Replyf("Harga Tertinggi (24 jam): %s", stat.High)
    message.Replyf("Harga Terendah (24 jam): %s", stat.Low)
    message.Reply("==========")
}
