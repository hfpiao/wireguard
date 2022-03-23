package main

import (
	"flag"
	"github.com/hfpiao/wireguard/util"
)

/*
Options:
 [-help]                      Display help
 [-amqp_url <Url>]                 Action pub (publish) or sub (subscribe)
 [-device_name <DeviceName>]               Payload to send0
 [-exchange <Exchange>]                Number of messages to send or receive
 [-exchange_type direct|topic|fanout|headers]                   Quality of Service
*/
func main() {
	amqp_url := flag.String("url", "", "The broker URI (must). ex: amqps://user:password@url/vhost")
	device_name := flag.String("i", "wg0", "The wireguard interface (optional)")
	exchange := flag.String("e", "wireguard", "The rabbitmq exchange (optional)")
	exchange_type := flag.String("et", "direct", "The rabbitmq exchange type (optional) ")
	queue := flag.String("q", "", "The rabbitmq exchange type (must) ")
	routing_key := flag.String("rk", "endpoint", "The rabbitmq exchange type (optional) ")
	flag.Parse()
	util.Run(*amqp_url,
		*device_name,
		*exchange,
		*exchange_type,
		*routing_key,
		*queue)
}
