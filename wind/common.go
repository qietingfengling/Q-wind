package wind

import (
	"crypto/tls"
	"flag"
	"os"
	"fmt"
	"strings"
)

var RunFlag = true

var cipherDic = map[string]uint16{
	"AES128-SHA":                    tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"AES256-SHA":                    tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"AES128-SHA256":                 tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	"AES128-GCM-SHA256":             tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	"AES256-GCM-SHA384":             tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-ECDSA-AES128-SHA":        tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"ECDHE-ECDSA-AES256-SHA":        tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"ECDHE-RSA-AES128-SHA":          tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"ECDHE-RSA-AES256-SHA":          tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"ECDHE-ECDSA-AES128-SHA256":     tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	"ECDHE-RSA-AES128-SHA256":       tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	"ECDHE-RSA-AES128-GCM-SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-ECDSA-AES128-GCM-SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-RSA-AES256-GCM-SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-ECDSA-AES256-GCM-SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	"ECDHE-RSA-CHACHA20-POLY1305":   tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	"ECDHE-ECDSA-CHACHA20-POLY1305": tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}
var (
	flagSet     = flag.NewFlagSet("Q-wind", flag.ExitOnError)
	requests    = flagSet.Int64("n", 1, "Number of requests to perform.Default 1")
	r           = flagSet.Int64("r", 1, "Number of requests each second.Default 1")
	burst       = flagSet.Int("b", 1, "burst limit,0 for not limit.Default 1")
	client      = flagSet.Int64("c", 1, "Number of multiple requests to make at a time.Default 1")
	timeLimit   = flagSet.Int64("t", 0, " Seconds to max. to spend on benchmarking,This implies -b 0")
	timeout     = flagSet.Int("s", 30, "Seconds to max. wait for each response. Default 30 seconds")
	keepAlive   = flagSet.Bool("k", false, "Use HTTP KeepAlive feature,argument invalid in h2.Default false")
	cipherSuite = flagSet.String("Z", "", "Specify SSL/TLS cipher suite (See openssl ciphers)")
	url         = flagSet.String("url", "https://127.0.0.1", "stress url, Default  https://127.0.0.1")
	method      = flagSet.String("m", "GET", "Request Method,[GET,HEAD,POST],Default GET")
	protocol    = flagSet.String("p", "h1", "Request Method,[h1,h2],Default  h1")
	proxy       = flagSet.String("x", "", "proxy:port ,proxy server and port number to use")
	heard       = flagSet.String("H", "", "header")
)

func init() {
	flagSet.Parse(os.Args[1:])
	_, ok := cipherDic[*cipherSuite]
	if !ok && *cipherSuite != "" {
		fmt.Printf("[Q-wind] cipher error\n")
		RunFlag = false
	}
	if !strings.Contains("GET,HEAD,POST", strings.ToUpper(*method)) {
		fmt.Printf("[Q-wind] method error,pealse select GET,HEAD,POST\n")
		RunFlag = false
	}
	//timeLimit not 0 ,burst default 0
	if *timeLimit != 0 {
		*burst = 0
	}
}
