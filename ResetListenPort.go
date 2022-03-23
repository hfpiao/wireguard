package main

import (
	"flag"
	util "github.com/hfpiao/wireguard/util"
)

func main() {
	amqp_url := flag.String("url", "", "The broker URI (must). ex: amqps://user:password@url/vhost")
	device_name := flag.String("i", "wg0", "The wireguard interface (optional)")
	exchange := flag.String("e", "wireguard", "The rabbitmq exchange (optional)")
	routing_key := flag.String("rk", "endpoint", "The rabbitmq exchange type (optional) ")
	host := flag.String("h", "160.119.69.126", "The wireguard device host ip address. (optional) ")
	flag.Parse()

	util.SetDevicePort(*amqp_url, *device_name, *exchange, *routing_key, *host)
}
