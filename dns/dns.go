package dns

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

const (
	DefaultDNS     = "udp:[8.8.8.8]:53"
	DefaultDNSPort = 53
)

type Opt struct {
	Domain string
	Type   QueryType
	Server string
}

func Run(o Opt) (*Response, error) {
	if o.Domain == "" {
		return nil, errors.New("valid query domain name must be entered")
	}
	if o.Type == UnknownQueryType {
		o.Type = A
	}
	if o.Server == "" {
		o.Server = GetDNSFromSystem()
	}

	// 定义正则表达式
	re := regexp.MustCompile(`^(\w+):\[((?:\d{1,3}\.){3}\d{1,3}|[0-9A-Fa-f:]+)]:(\d+)$`)
	// 进行匹配
	matches := re.FindStringSubmatch(o.Server)
	if len(matches) != 4 {
		return nil, errors.New("the server format is incorrect and should be: protocol:[address]:port")
	}

	protocol, address, port := matches[1], matches[2], cast.ToInt(matches[3])
	switch protocol {
	case "udp", "tcp":
	default:
		return nil, errors.New("unsupported protocol type: " + protocol)
	}

	dial, err := net.Dial(protocol, fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	_, err = dial.Write(MarshalRequest(NewRequest(o.Domain, o.Type)))
	if err != nil {
		return nil, err
	}

	buff := make([]byte, BufferSize)
	n, err := dial.Read(buff)
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	_, err = UnmarshalResponse(buff[:n], resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetDNSFromSystem
// Try to get dns from the system, when there are multiple dns server in the system, return the first one.
// If the acquisition fails, return the default dns.
// The current default dns is: udp:[8.8.8.8]:53
func GetDNSFromSystem() (dns string) {
	dns = DefaultDNS

	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return
	}
	defer file.Close()

	// Get the first nameserver in the system.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return
		}

		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 1 {
			continue
		}

		switch f[0] {
		case "nameserver": // add one name server
			if len(f) > 1 {
				// One more check: make sure server name is
				// just an IP address. Otherwise, we need DNS
				// to look it up.
				dns = fmt.Sprintf("udp:[%s]:%d", f[1], DefaultDNSPort)
				return
			}
		}
	}
	return
}
