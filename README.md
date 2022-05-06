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
upb.go              // collect & save upbit ob & tx
kbt.go              // collect & save korbit ob & tx

cpu.go              // main function of package

processor.go        // calculate & save hit count

updater.go          // saves to db every second

### useful commands
boltbrowswer <filename.db>
godoc -http=localhost:6060 → http://localhost:6060/pkg/github.com/neosouler7/GObserver