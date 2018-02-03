package outbound

import(
    "os"    
    "log"
    "net/http"
    "io/ioutil"
)


var vipPublicAPI = os.Getenv("MARKET_API_URL")


func FetchMarketPrice(ticker string) ([]byte, error) {
    resp, err := http.Get(vipPublicAPI + ticker + "_idr/ticker")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("Sending enquiry to Market.....")

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

        return body, err
    } else {
      body := []byte{}

      return body, err
    }
}