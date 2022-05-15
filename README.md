# GObserver
Project GObserver is made up of 3 main and subsidiary parts.
- main       : cpu / exchanges / db / tg 
- subsidiary : utils, config etc 

# Main Functions for each packages
### collectors
collector.go        // collects exchange's datas

* orderbook
used in trade TAKER-TAKER.
then, we should store orderbook price as TAKER struct.

type taker struct {
    exchange string
    market string
    symbol string
    askPrice string // ask[0]
    askVolume string // ask[0]
    bidPrice string // bid[0]
    bidVolume string // bid[0]
    targetAskPrice string // sumOver target
    targetAskVolume string // sumOver target
    targetBidPrice string // sumOver target
    targetBidVolume string // sumOver target
    timestamp string // returned timestamp
}

obMap[obKey("upb:krw:btc")] = taker(struct)

* transaction
used in trade MAKER-TAKER.
we already have TAKER data, so save transaction data as MAKER struct.

type maker struct {
    exchange string
    market string
    symbol string
    askPrice string
    askVolume string // ignore if less then minAmount & should be bigger than taker's bidVolume
    bidPrice string
    bidVolume string // ignore if less then minAmount & should be bigger than taker's askVolume
    timestamp string
}

txMap[obKey("upb:krw:btc")] = maker(struct)

### processors
processor.go        // calculate collector's datas & save hit count

loop possible combinations. (tradeType x exchange(nC2) x ASK/BID)

* TAKER-TAKER
compare TAKER's targetPrice(ASK/BID) of 2 exchanges.

* MAKER-TAKER
compare 2 exchanges MAKER vs TAKER ASK/BID price.

### updaters
updater.go          // saves hit count map to db

### exchanges
upb.go              // collect & save upbit ob & tx
kbt.go              // collect & save korbit ob & tx

### db
func Start          // package start, init bolt DB
func createBucket   // create bucket
func clearBucket    // clear bucket
func SaveCheckPoint // saves timestamp 
func GetCheckPoint  // loads timestamp 
func UpdateMold     // update

### tg
func Start          // package start, init tg config/bot
func listenMsg      // listens msg channel
func SendMsg        // sends tg msg
func HandleErr      // sends error as tg msg

# useful commands
boltbrowswer <filename.db>
godoc -http=localhost:6060 → http://localhost:6060/pkg/github.com/neosouler7/GObserver