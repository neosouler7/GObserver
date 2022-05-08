# GObserver
Project GObserver is made up of 3 main and subsidiary parts.
- main       : cpu / exchanges / db / tg 
- subsidiary : utils, config etc 

# Main Functions for each packages
### collectors
collector.go        // collects exchange's datas

### processors
processor.go        // calculate collector's datas & save hit count

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