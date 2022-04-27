# GObserver
Project GObserver is made up of 3 main and subsidiary parts.
    - main       : cpu / db / tg 
    - subsidiary : utils, asyncLogging etc 

# Main Functions for each packages
### tg
func Start        // package start
func initBot      // init tg Bot
func listenMsg    // listens msg channel
func SendMsg      // sends tg msg

### db
func Start        // package start
func initDB       // init bolt DB
func createBucket // create bucket
func clearBucket  // clear bucket
func update       // update

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