# GObserver
Project GObserver is made up of 3 main and subsidiary parts.
    - main       : cpu / db / tg 
    - subsidiary : utils, asyncLogging etc 

# Main Functions for each packages
### tg
func Start          // package start, init tg config/bot
func listenMsg      // listens msg channel
func SendMsg        // sends tg msg
func HandleErr      // sends error as tg msg

### db
func Start          // package start, init bolt DB
func createBucket   // create bucket
func clearBucket    // clear bucket
func SaveCheckPoint // saves timestamp 
func GetCheckPoint  // loads timestamp 
func UpdateMold     // update

### cpu
collector.go
func upbit        // collect upbit ticker
func korbit       // collect korbit ticker
func bithumb      // collect bithumb ticker

processor.go
type dddd struct {

}

updater.go (let's use cron expression)
func updater      // updates ddd struct to db

### useful commands
boltbrowswer <filename.db>
godoc -http=localhost:6060 → http://localhost:6060/pkg/github.com/neosouler7/GObserver