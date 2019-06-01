package main

import (
    "fmt"
    "log"
    "errors"
    "context"
    "os"

    "github.com/apache/thrift/lib/go/thrift"
    "github.com/jamiees2/tinyscribe/scribe"
    "github.com/jamiees2/tinyscribe/fb303"
)

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

var (
    unsupportedError = errors.New("Unsupported action")
)

type ScribeHandler struct {}

func NewScribeHandler() *ScribeHandler {
    return &ScribeHandler{}
}

func (s *ScribeHandler) Log(ctx context.Context, messages []*scribe.LogEntry) (r scribe.ResultCode, err error)  {
    for _, message := range messages {
        log.Printf("[%s] %s", message.GetCategory(), message.GetMessage())
    }
    return scribe.ResultCode_OK, nil
}


  // Returns a descriptive name of the service
func (s *ScribeHandler) GetName(ctx context.Context) (r string, err error) {
  return "tinyscribe", nil
}
// Returns the version of the service
func (s *ScribeHandler) GetVersion(ctx context.Context) (r string, err error) {
    return "0.0.1", nil
}
// Gets the status of this service
func (s *ScribeHandler) GetStatus(ctx context.Context) (r fb303.FbStatus, err error) {
    return fb303.FbStatus_ALIVE, nil
}
// User friendly description of status, such as why the service is in
// the dead or warning state, or what is being started or stopped.
func (s *ScribeHandler) GetStatusDetails(ctx context.Context) (r string, err error) {
    return "We're alive! gowai!", nil
}
// Gets the counters for this service
func (s *ScribeHandler) GetCounters(ctx context.Context) (r map[string]int64, err error) {
    return map[string]int64{}, unsupportedError
}
// Gets the value of a single counter
// 
// Parameters:
//  - Key
func (s *ScribeHandler) GetCounter(ctx context.Context, key string) (r int64, err error) {
    return 0, unsupportedError
}
// Sets an option
// 
// Parameters:
//  - Key
//  - Value
func (s *ScribeHandler) SetOption(ctx context.Context, key string, value string) (err error) {
    return unsupportedError
}
// Gets an option
// 
// Parameters:
//  - Key
func (s *ScribeHandler) GetOption(ctx context.Context, key string) (r string, err error) {
    return "", unsupportedError
}
// Gets all options
func (s *ScribeHandler) GetOptions(ctx context.Context) (r map[string]string, err error) {
    return map[string]string{}, unsupportedError
}
// Returns a CPU profile over the given time interval (client and server
// must agree on the profile format).
// 
// Parameters:
//  - ProfileDurationInSec
func (s *ScribeHandler) GetCpuProfile(ctx context.Context, profileDurationInSec int32) (r string, err error) {
    return "", unsupportedError

}
// Returns the unix time that the server has been running since
func (s *ScribeHandler) AliveSince(ctx context.Context) (r int64, err error) {
    // No idea
    return 0, nil
}
// Tell the server to reload its configuration, reopen log files, etc
func (s *ScribeHandler) Reinitialize(ctx context.Context) (err error) {
  return nil
}
// Suggest a shutdown to the server
func (s *ScribeHandler) Shutdown(ctx context.Context) (err error) {
    os.Exit(0)
    return nil
}


func main() {
    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

    addr := getEnv("LISTEN_HOST", ":1463")
    transport, err := thrift.NewTServerSocket(addr)
    if err != nil {
        log.Printf("error: %v", err)
        os.Exit(1)
    }
    handler := NewScribeHandler()
    processor := scribe.NewScribeProcessor(handler)
    server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

    fmt.Println("Starting the scribe server... on ", addr)
    server.Serve()
}


