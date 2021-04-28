package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/rs/zerolog/log"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

func UnixMicro() int64 {
	return time.Now().UnixNano() / 1000000
}

func ToJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)

	return data
}

func FromJSON(data []byte, v interface{}) (err error) {
	err = json5.Unmarshal(data, v)
	return
}

func RandInt(min int, max int) int {

	return min + rand.Intn(max-min)
}

func RemoveSpace(data []byte) []byte {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAll(data, []byte(""))
}

func BuildQuery(params map[string]string) string {
	builder := url.Values{}
	for k, v := range params {
		builder.Add(k, v)
	}
	return builder.Encode()
}

func WorkDirFilePath(fileName string) string {
	ex, _ := os.Executable()

	return filepath.Join(filepath.Dir(ex), fileName)
}

func ToInt(str string) int {
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0
	} else {
		return v
	}
}

func MachineName() string {
	name, err := os.Hostname()
	if err != nil {
		return "(Unknown)"
	}

	return name
}

func GetAllInterface() string {
	ips := []string{}
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Error().Err(err).Msg("")
			continue
		}
		// handle err
		// fmt.Println(string(ToJSON(i)))
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// process IP address
			if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() {
				continue
			}

			if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	return strings.Join(ips, "|")
}

func GetLocalBroadcastInterface() []string {
	ips := []string{}
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {

		addrs, err := i.Addrs()
		if err != nil {
			log.Error().Err(err).Msg("")
			continue
		}

		for _, addr := range addrs {
			ipnet := addr.(*net.IPNet)
			if ipnet == nil {
				continue
			}
			ip := ipnet.IP
			// process IP address
			if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() {
				continue
			}

			if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() {
				continue
			}

			if ip.To4() == nil {
				continue
			}

			ipbroadcast := make(net.IP, len(ip.To4()))
			binary.BigEndian.PutUint32(ipbroadcast, binary.BigEndian.Uint32(ip.To4())|^binary.BigEndian.Uint32(net.IP(ipnet.Mask).To4()))
			ips = append(ips, ipbroadcast.String())

		}
	}

	return ips
}

func StringToUTF16Bytes(s string) []byte {
	runes := utf16.Encode([]rune(s))
	bytes := make([]byte, len(runes)*2)
	for i, r := range runes {
		binary.LittleEndian.PutUint16(bytes[i*2:], r)
	}
	return bytes
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	err := binary.Read(bytebuff, binary.BigEndian, &data)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return int(data)
}

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return bytesBuffer.Bytes()
}

func BytesToBinaryString(bs []byte) string {
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

func Base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func Base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}

func HashSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
