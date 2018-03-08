# Q-wind

## Description

S-wing is a request pressure test , written in golang.
It can customize the number of requests, pressure measurement time, header, proxy, protocols and suites.

## Usage of Q-wind

```bash
  -Z string
    	Specify SSL/TLS cipher suite (See openssl ciphers)
  -b int
    	burst limit,0 for not limit.Default 1 (default 1)
  -c int
    	Number of multiple requests to make at a time.Default 1 (default 1)
  -k	Use HTTP KeepAlive feature,argument invalid in h2.Default false
  -m string
    	Request Method,[GET,HEAD,POST],Default GET (default "GET")
  -n int
    	Number of requests to perform.Default 1 (default 1)
  -p string
    	Request Method,[h1,h2],Default  h1 (default "h1")
  -r int
    	Number of requests each second.Default 1 (default 1)
  -s int
    	Seconds to max. wait for each response. Default 30 seconds (default 30)
  -t int
    	 Seconds to max. to spend on benchmarking,This implies -b 0
  -url string
    	stress url, Default  https://127.0.0.1 (default "https://127.0.0.1")
  -x string
    	proxy:port ,proxy server and port number to use
```

### Ciphers Suite Support
    AES128-SHA
    AES256-SHA
 	AES128-SHA256
 	AES128-GCM-SHA256
 	AES256-GCM-SHA384
 	ECDHE-ECDSA-AES128-SHA
 	ECDHE-ECDSA-AES256-SHA
 	ECDHE-RSA-AES128-SHA
 	ECDHE-RSA-AES256-SHA
 	ECDHE-ECDSA-AES128-SHA256
 	ECDHE-RSA-AES128-SHA256
 	ECDHE-RSA-AES128-GCM-SHA256
 	ECDHE-ECDSA-AES128-GCM-SHA256
 	ECDHE-RSA-AES256-GCM-SHA384
 	ECDHE-ECDSA-AES256-GCM-SHA384
 	ECDHE-RSA-CHACHA20-POLY1305
 	ECDHE-ECDSA-CHACHA20-POLY1305

### Protocols Support
    http1
    http2