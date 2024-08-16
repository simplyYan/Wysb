package WysbSniffer

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	device     string
	handle     *pcap.Handle
	packetChan chan gopacket.Packet
}

func NewSniffer(device string) (*Sniffer, error) {
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	sniffer := &Sniffer{
		device:     device,
		handle:     handle,
		packetChan: make(chan gopacket.Packet),
	}

	go sniffer.startCapture()
	return sniffer, nil
}

func (s *Sniffer) startCapture() {
	defer s.handle.Close()
	packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	for packet := range packetSource.Packets() {
		s.packetChan <- packet
	}
}

func (s *Sniffer) GetPackets() <-chan gopacket.Packet {
	return s.packetChan
}

func (s *Sniffer) PrintPacketInfo(packet gopacket.Packet) {
	// Displays basic packet information
	fmt.Println("Captured packet:")
	fmt.Println(packet.String())
}

func (s *Sniffer) Start() {
	fmt.Printf("Starting packet capture on interface: %s\n", s.device)
	for packet := range s.GetPackets() {
		s.PrintPacketInfo(packet)
	}
}

func Disclaimer() {
	fmt.Println("This library is for educational and ethical purposes. The creator is not responsible for misuse.")
}
