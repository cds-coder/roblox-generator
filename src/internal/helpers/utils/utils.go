package utils

import (
	"RobloxRegister/src/internal/helpers/class"
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	randa "crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	proxies []string
	mu      sync.Mutex
)

type ECDSASignature struct {
	R, S *big.Int
}

func init() {

	file, err := os.Open("input/proxies.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	proxies = lines
}

func GetProxy() string {
	return proxies[rand.Intn(len(proxies))]
}

func GenerateSecureAuth(serverNonce string) (*class.SecureAuth, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), randa.Reader)
	if err != nil {
		return nil, err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	clientPublicKey := base64.StdEncoding.EncodeToString(pubBytes)

	clientEpochTimestamp := time.Now().Unix()

	payload := clientPublicKey + "|" +
		strconv.FormatInt(clientEpochTimestamp, 10) + "|" +
		serverNonce

	hash := sha256.Sum256([]byte(payload))

	r, s, err := ecdsa.Sign(randa.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	sigStruct := ECDSASignature{R: r, S: s}

	derSignature, err := asn1.Marshal(sigStruct)
	if err != nil {
		return nil, err
	}

	saiSignature := base64.StdEncoding.EncodeToString(derSignature)

	return &class.SecureAuth{
		ClientPublicKey:      clientPublicKey,
		ClientEpochTimestamp: clientEpochTimestamp,
		ServerNonce:          serverNonce,
		SaiSignature:         saiSignature,
	}, nil
}

func ParseArkoseHeader(headerVal string) (*class.ArkoseResponse, error) {
	decoded, err := base64.StdEncoding.DecodeString(headerVal)
	if err != nil {
		return nil, err
	}

	var ark class.ArkoseResponse
	if err := json.Unmarshal(decoded, &ark); err != nil {
		return nil, err
	}

	return &ark, nil
}

func SaveAccount(user, pass, cookie string) error {
	line := user + ":" + pass + ":" + cookie + "\n"

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile("output/accounts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line)
	return err
}

func rgbGradient(r1, g1, b1, r2, g2, b2, n int) [][3]int {
	gradient := make([][3]int, n)
	for i := 0; i < n; i++ {
		gradient[i][0] = r1 + (r2-r1)*i/n
		gradient[i][1] = g1 + (g2-g1)*i/n
		gradient[i][2] = b1 + (b2-b1)*i/n
	}
	return gradient
}

func Output(msgType, msg string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Format("15:04:05")
	timeColor := "\033[90m"
	reset := "\033[0m"

	var typeColor string
	var startRGB, endRGB [3]int

	switch msgType {
	case "INFO":
		typeColor = "\033[38;2;0;123;255m"
		startRGB = [3]int{0, 123, 255}
		endRGB = [3]int{0, 200, 255}
	case "CAPTCHA":
		typeColor = "\033[38;2;255;193;7m"
		startRGB = [3]int{255, 193, 7}
		endRGB = [3]int{255, 230, 100}
	case "SUCCESS":
		typeColor = "\033[38;2;0;200;83m"
		startRGB = [3]int{0, 200, 83}
		endRGB = [3]int{100, 255, 150}
	case "FAILED":
		typeColor = "\033[38;2;255;82;82m"
		startRGB = [3]int{255, 82, 82}
		endRGB = [3]int{255, 150, 150}
	default:
		typeColor = "\033[37m"
		startRGB = [3]int{255, 255, 255}
		endRGB = [3]int{200, 200, 200}
	}

	fmt.Printf("%s[%s]%s %s[%s]%s ", timeColor, now, reset, typeColor, msgType, reset)

	gradient := rgbGradient(startRGB[0], startRGB[1], startRGB[2], endRGB[0], endRGB[1], endRGB[2], len(msg))
	for i, c := range msg {
		r, g, b := gradient[i][0], gradient[i][1], gradient[i][2]
		fmt.Printf("\033[38;2;%d;%d;%dm%c", r, g, b, c)
	}
	fmt.Println(reset)
}
