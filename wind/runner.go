package wind

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"golang.org/x/time/rate"
	"strings"
	"golang.org/x/net/http2"
	"crypto/tls"
	"fmt"
	"sync/atomic"
	"strconv"
	"runtime"
	"context"
	"time"
	"net/http"
	"net"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Run() {
	channel := make(chan string, *client)
	done := make(chan int, 0)
	l := rate.NewLimiter(rate.Limit(*r), *burst)
	c, _ := context.WithCancel(context.TODO())
	go func(limitrate *rate.Limiter, can *context.Context) {
		if *timeLimit != 0 {
			tt := time.After(time.Second * time.Duration(*timeLimit))
			endFlag := false
			for !endFlag {
				select {
				case <-tt:
					endFlag = true
				default:
					limitrate.Wait(*can)
					channel <- *url
				}
			}
			done <- 1
		} else {
			for i := int64(0); i < *requests; i++ {
				limitrate.Wait(*can)
				channel <- *url
			}
		}
	}(l, &c)
	total := *requests
	var process int64 = 0
	var qps int64 = 0
	var failedReq int64 = 0
	var totalFailedReq int64 = 0
	var fourXXReq int64 = 0
	var fiveXXReq int64 = 0
	var successReq int64 = 0
	var totalTransferred int64 = 0
	var count int64
	var ticker = time.NewTicker(1 * time.Second)
	t1 := time.Now()
	var alive = *client
	cipher, ok := cipherDic[*cipherSuite]
	var cipherSuites []uint16
	if ok {
		cipherSuites = []uint16{cipher}
	} else {
		cipherSuites = nil
	}
	headMap := make(map[string]string)
	if *heard != "" {
		for _, t := range strings.Split(*heard, ",") {
			idx := strings.Index(t, ":")
			if strings.Index(t, ":") > 0 {
				headMap[string(t[0:idx])] = string(t[idx+1:])
			}
		}
	}
	for count = 0; count <= *client; count++ {
		go func() {
			var c http.Client
			if *protocol == "h2" {
				c = http.Client{
					Transport: &http2.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
							CipherSuites:       cipherSuites,
						},
						DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
							return tls.DialWithDialer(&net.Dialer{Timeout: time.Duration(*timeout) * time.Second}, network, addr, cfg)
						},
					},
					Timeout: time.Second * time.Duration(*timeout),
				}

			} else {
				c = http.Client{
					Transport: &http.Transport{
						Dial: func(nettw, addr string) (net.Conn, error) {
							c, err := net.DialTimeout("tcp4", *proxy, time.Duration(*timeout)*time.Second)
							if err != nil {
								if *proxy != "" {
									addr = *proxy
								}
								fmt.Printf("[Q-wind] h1 dial %v proxy error,error info:%v", addr, err)
								return nil, err
							}
							return c, nil
						},
						DisableKeepAlives: !*keepAlive,
					},
				}
			}
		forEnd:
			for {
				select {
				case l := <-channel:
					req, err := http.NewRequest(*method, l, nil)
					if err != nil {
						fmt.Printf("new request error:%v\n", err)
					}
					for key, value := range headMap {
						req.Header.Set(key, value)
					}
					resp, err := c.Do(req)
					if err != nil {
						atomic.AddInt64(&failedReq, 1)
						atomic.AddInt64(&totalFailedReq, 1)
						fmt.Printf("do request error:%v\n", err)
					} else {
						if resp != nil {
							if resp.ContentLength > 0 {
								atomic.AddInt64(&totalTransferred, resp.ContentLength)
							}
							if resp.StatusCode == 200 {
								atomic.AddInt64(&successReq, 1)
							} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
								atomic.AddInt64(&fourXXReq, 1)
							} else if resp.StatusCode >= 500 && resp.StatusCode < 600 {
								atomic.AddInt64(&fiveXXReq, 1)
							}
						}
						atomic.AddInt64(&qps, 1)
						defer resp.Body.Close()
					}
					atomic.AddInt64(&process, 1)

				case <-time.After(time.Second * time.Duration(*timeout)):
					break forEnd
				}
			}
			atomic.AddInt64(&alive, -1)
		}()
	}

	//table := tablewriter.NewWriter(os.Stdout)
	//table.SetAlignment(tablewriter.ALIGN_LEFT)
	//table.SetHeader([]string{"Elapsed", "total", "Completed", "Requests per second [sed]", "Failed requests", "Clients", "Duration"})
	fmt.Printf("|%-10s|%-10s|%-10s|%-30s|%-20s|%-10s|%-10s|\n", "Elapsed", "total", "Completed", "Requests per second [sed]", "Failed requests", "Clients", "Duration")
	go func() {
		for t := range ticker.C {
			elapsed := time.Since(t1)
			_ = t
			reqPerSec := atomic.LoadInt64(&qps)
			atomic.StoreInt64(&qps, 0)
			failedReqPerSec := atomic.LoadInt64(&failedReq)
			atomic.StoreInt64(&failedReq, 0)
			doneReq := atomic.LoadInt64(&process)
			if reqPerSec < 1 {
				reqPerSec = 1
			}
			remainReq := total - doneReq
			sample := remainReq / reqPerSec
			if remainReq != 0 && remainReq < reqPerSec {
				sample = 1
			}
			rDurTime := "+" + strconv.FormatInt(sample, 10) + "s"
			dur, _ := time.ParseDuration(rDurTime)
			fmt.Printf("|%-10s|%-10s|%-10s|%-30s|%-20s|%-10s|%-10s|\n", fmt.Sprintf("%0.1fs", elapsed.Seconds()), fmt.Sprintf("%v", total),
				fmt.Sprintf("%v", doneReq), fmt.Sprintf("%v", reqPerSec), fmt.Sprintf("%v", failedReqPerSec),
				fmt.Sprintf("%v", alive), fmt.Sprintf("%v", dur))
			if process == total {
				done <- 1
				break
			}
			//table.Append([]string{fmt.Sprintf("%0.1fs", elapsed.Seconds()), fmt.Sprintf("%v", total),
			//	fmt.Sprintf("%v", doneReq), fmt.Sprintf("%v", reqPerSec), fmt.Sprintf("%v", failedReqPerSec),
			//	fmt.Sprintf("%v", alive), fmt.Sprintf("%v", dur)})
			//table.Render()
		}
	}()
	<-done
	elapsed := time.Since(t1)
	time.Sleep(1 * time.Second)
	resultTable := tablewriter.NewWriter(os.Stdout)
	resultTable.SetAlignment(tablewriter.ALIGN_LEFT)
	resultTable.SetHeader([]string{"Item", "Value"})
	resultTable.Append([]string{"Time taken for stress", fmt.Sprintf("%0.1f seconds", elapsed.Seconds())})
	resultTable.Append([]string{"Complete requests", fmt.Sprintf("%v", process)})
	resultTable.Append([]string{"Failed requests", fmt.Sprintf("%v", totalFailedReq)})
	resultTable.Append([]string{"Success requests", fmt.Sprintf("%v", successReq)})
	resultTable.Append([]string{"4xx requests", fmt.Sprintf("%v", fourXXReq)})
	resultTable.Append([]string{"5xx requests", fmt.Sprintf("%v", fiveXXReq)})
	resultTable.Append([]string{"Total transferred", fmt.Sprintf("%v bytes", totalTransferred)})
	resultTable.Append([]string{"Transfer rate", fmt.Sprintf("%v [bytes/sec] received", totalTransferred/int64(elapsed.Seconds()))})
	resultTable.Render()
	fmt.Println(Binary)

}
