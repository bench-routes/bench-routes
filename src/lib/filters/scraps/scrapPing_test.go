package scraps

import (
	"testing"

	"github.com/zairza-cetb/bench-routes/src/lib/utils"
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
	`PING yahoo.com (98.138.219.232) 56(84) bytes of data.
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=1 ttl=45 time=457 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=2 ttl=45 time=482 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=3 ttl=45 time=405 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=4 ttl=45 time=427 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=5 ttl=45 time=451 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=6 ttl=45 time=474 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=7 ttl=45 time=397 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=8 ttl=45 time=420 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=9 ttl=45 time=444 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=10 ttl=45 time=468 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=11 ttl=45 time=389 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=12 ttl=45 time=413 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=13 ttl=45 time=438 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=14 ttl=45 time=461 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=15 ttl=45 time=487 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=16 ttl=45 time=416 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=17 ttl=45 time=432 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=18 ttl=45 time=456 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=19 ttl=45 time=480 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=20 ttl=45 time=400 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=21 ttl=45 time=426 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=22 ttl=45 time=448 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=23 ttl=45 time=474 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=24 ttl=45 time=396 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=25 ttl=45 time=418 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=26 ttl=45 time=442 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=27 ttl=45 time=475 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=28 ttl=45 time=490 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=29 ttl=45 time=409 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=30 ttl=45 time=434 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=31 ttl=45 time=458 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=32 ttl=45 time=482 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=33 ttl=45 time=407 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=34 ttl=45 time=431 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=35 ttl=45 time=453 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=36 ttl=45 time=474 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=37 ttl=45 time=398 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=38 ttl=45 time=421 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=39 ttl=45 time=372 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=40 ttl=45 time=367 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=41 ttl=45 time=389 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=42 ttl=45 time=424 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=43 ttl=45 time=440 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=44 ttl=45 time=562 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=45 ttl=45 time=485 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=46 ttl=45 time=405 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=47 ttl=45 time=439 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=48 ttl=45 time=452 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=49 ttl=45 time=478 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=50 ttl=45 time=397 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=51 ttl=45 time=424 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=52 ttl=45 time=447 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=53 ttl=45 time=469 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=54 ttl=45 time=394 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=55 ttl=45 time=416 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=56 ttl=45 time=506 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=57 ttl=45 time=464 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=58 ttl=45 time=488 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=59 ttl=45 time=408 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=60 ttl=45 time=432 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=61 ttl=45 time=457 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=62 ttl=45 time=480 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=63 ttl=45 time=402 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=64 ttl=45 time=429 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=65 ttl=45 time=453 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=66 ttl=45 time=474 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=67 ttl=45 time=399 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=68 ttl=45 time=423 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=69 ttl=45 time=445 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=70 ttl=45 time=468 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=71 ttl=45 time=407 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=72 ttl=45 time=415 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=73 ttl=45 time=438 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=74 ttl=45 time=461 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=75 ttl=45 time=486 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=76 ttl=45 time=408 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=77 ttl=45 time=433 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=78 ttl=45 time=458 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=79 ttl=45 time=479 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=80 ttl=45 time=403 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=81 ttl=45 time=427 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=82 ttl=45 time=454 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=83 ttl=45 time=476 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=84 ttl=45 time=393 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=85 ttl=45 time=421 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=86 ttl=45 time=367 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=87 ttl=45 time=465 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=88 ttl=45 time=488 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=89 ttl=45 time=412 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=90 ttl=45 time=432 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=91 ttl=45 time=456 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=92 ttl=45 time=481 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=93 ttl=45 time=405 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=94 ttl=45 time=427 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=95 ttl=45 time=451 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=96 ttl=45 time=479 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=97 ttl=45 time=398 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=98 ttl=45 time=420 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=99 ttl=45 time=442 ms
64 bytes from media-router-fp2.prod1.media.vip.ne1.yahoo.com (98.138.219.232): icmp_seq=100 ttl=45 time=466 ms

--- yahoo.com ping statistics ---
100 packets transmitted, 100 received, 0% packet loss, time 111ms
rtt min/avg/max/mdev = 366.832/439.717/562.375/34.218 ms
`,
}

func TestCLIPingScrap(t *testing.T) {
	for _, samples := range input {
		a := CLIPingScrap(&samples)
		if *a == (utils.TypePingScrap{}) {
			t.Errorf("invalid scrapping of targets for ping module")
		}
	}
}
