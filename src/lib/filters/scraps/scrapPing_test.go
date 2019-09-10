package scraps

import (
	"testing"
)

var input = []string{`
PING google.co.in (172.217.26.227) 56(84) bytes of data.
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=1 ttl=57 time=52.6 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=2 ttl=57 time=54.4 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=3 ttl=57 time=52.5 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=4 ttl=57 time=53.0 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=5 ttl=57 time=52.9 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=6 ttl=57 time=53.3 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=7 ttl=57 time=52.6 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=8 ttl=57 time=52.3 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=9 ttl=57 time=52.1 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=10 ttl=57 time=53.3 ms

--- google.co.in ping statistics ---
10 packets transmitted, 10 received, 0% packet loss, time 20ms
rtt min/avg/max/mdev = 52.065/52.894/54.366/0.676 ms
`,
`
PING google.co.in (172.217.26.227) 56(84) bytes of data.
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=1 ttl=57 time=52.4 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=2 ttl=57 time=52.10 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=3 ttl=57 time=52.3 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=4 ttl=57 time=51.8 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=5 ttl=57 time=54.3 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=6 ttl=57 time=51.4 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=7 ttl=57 time=51.10 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=8 ttl=57 time=51.5 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=9 ttl=57 time=52.2 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=10 ttl=57 time=50.10 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=11 ttl=57 time=51.10 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=12 ttl=57 time=51.9 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=13 ttl=57 time=52.2 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=14 ttl=57 time=51.2 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=15 ttl=57 time=51.5 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=16 ttl=57 time=51.7 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=17 ttl=57 time=51.10 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=18 ttl=57 time=66.3 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=19 ttl=57 time=52.1 ms
64 bytes from bom05s09-in-f3.1e100.net (172.217.26.227): icmp_seq=20 ttl=57 time=51.9 ms

--- google.co.in ping statistics ---
20 packets transmitted, 20 received, 0% packet loss, time 49ms
rtt min/avg/max/mdev = 50.967/52.729/66.345/3.211 ms
`,
};

func TestCLIPingScrap(t *testing.T) {
	for _, samples := range input {
		a := CLIPingScrap(&samples)
		if *a == (TypePingScrap{}) {
			t.Errorf("invalid scrapping of targets for ping module")
		}
	}
}