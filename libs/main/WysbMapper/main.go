package WysbMapper

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type Scanner struct {
	IPRange string
	Results map[string][]string
	mu      sync.Mutex
}

func NewScanner(ipRange string) *Scanner {
	return &Scanner{
		IPRange: ipRange,
		Results: make(map[string][]string),
	}
}

func (s *Scanner) Scan() {
	ipList := s.getIPList(s.IPRange)
	var wg sync.WaitGroup

	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			s.scanIP(ip)
		}(ip)
	}

	wg.Wait()
}

func (s *Scanner) getIPList(ipRange string) []string {
	var ipList []string
	_, ipNet, err := net.ParseCIDR(ipRange)
	if err != nil {
		return ipList 
	}

	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		ipList = append(ipList, ip.String())
	}

	return ipList
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		if ip[j]++; ip[j] > 0 {
			break
		}
	}
}

func (s *Scanner) scanIP(ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for port := 1; port <= 1024; port++ {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			s.Results[ip] = append(s.Results[ip], fmt.Sprintf("Porta %d aberta", port))
			conn.Close()
		}
	}

	s.scanHTTP(ip)
}

func (s *Scanner) scanHTTP(ip string) {
	for _, protocol := range []string{"http", "https"} {
		url := fmt.Sprintf("%s://%s", protocol, ip)
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				s.Results[ip] = append(s.Results[ip], fmt.Sprintf("%s acessível", url))
			}
		}
	}
}

func (s *Scanner) PrintResults() {
	for ip, results := range s.Results {
		fmt.Printf("Resultados para %s:\n", ip)
		for _, result := range results {
			fmt.Println(result)
		}
	}
}

func Disclaimer() {
	fmt.Println("This library is for educational and ethical purposes. The creator is not responsible for misuse.")
}