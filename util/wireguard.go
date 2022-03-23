package util

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"math/rand"
	"net"
	"strconv"
)

type WireguardMessage struct {
	Peer string `json:peer`
	Host string `json:host`
	Port string `json:port`
}

func RandPort(seed int) int {
	var temp int64
	r := rand.New(rand.NewSource(int64(seed)))
	temp = int64(r.Intn(56000))
	if temp > 10000 {
		return int(temp)
	} else {
		return RandPort(int(temp))
	}
}

func SetDevicePort(amqp_url, deviceName, exchange, routing_key, host string) error {
	wgctrlClient, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer wgctrlClient.Close()

	d, _ := wgctrlClient.Device(deviceName)
	port := RandPort(d.ListenPort)
	log.Println(strconv.Itoa(port))
	config := wgtypes.Config{
		ReplacePeers: false,
		ListenPort:   &port,
	}
	err = wgctrlClient.ConfigureDevice(deviceName, config)
	if err != nil {
		log.Fatalf("failed to config device: %v", err)
	}
	message := WireguardMessage{
		Peer: d.PublicKey.String(),
		Host: host,
		Port: strconv.Itoa(port),
	}
	return Publish(amqp_url, exchange, routing_key, message)
}

func SetPeer(deviceName string, message WireguardMessage) error {
	wgctrlClient, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer wgctrlClient.Close()

	d, _ := wgctrlClient.Device(deviceName)
	var peerConfig []wgtypes.PeerConfig
	for i := range d.Peers {
		publicKey, err := wgtypes.ParseKey(message.Peer)
		if d.Peers[i].PublicKey == publicKey {
			endpoint := d.Peers[i].Endpoint
			endpoint.IP = net.ParseIP(message.Host)
			endpoint.Port, err = strconv.Atoi(message.Port)
			if err != nil {
				fmt.Println(err)
			}
			peerConfig = append(peerConfig, wgtypes.PeerConfig{
				PublicKey: d.Peers[i].PublicKey,
				Endpoint:  endpoint,
			})
		} else {
			peerConfig = append(peerConfig, wgtypes.PeerConfig{
				PublicKey: d.Peers[i].PublicKey,
				Endpoint:  d.Peers[i].Endpoint,
			})
		}
	}
	config := wgtypes.Config{
		Peers:        peerConfig,
		ReplacePeers: false,
	}
	return wgctrlClient.ConfigureDevice(deviceName, config)
}
