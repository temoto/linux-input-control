package main

import (
	"encoding/binary"
	"log"
	"os"
)

var (
	senable  = []byte("a")
	sdisable = []byte("d")
)

func main() {
	// TODO: command line flag
	masterPath := "/dev/input/event0"
	slavePath := "/sys/bus/i2c/devices/0-0050/neocmd"

	var master, slave *os.File
	var err error
	if master, err = os.Open(masterPath); err != nil {
		log.Fatal("open ", masterPath, err)
	}

	slaveState := true
	evb := make([]byte, 16)
	for {
		n, err := master.Read(evb)
		if err != nil {
			log.Print("master.Read: ", err)
			continue
		}
		if n != 16 {
			log.Printf("master.Read invalid event size=%d, ignore", n)
			continue
		}
		evType := binary.LittleEndian.Uint16(evb[8:])
		evCode := binary.LittleEndian.Uint16(evb[10:])
		evValue := binary.LittleEndian.Uint32(evb[12:])
		log.Printf("event type=%d code=%d value=%d", evType, evCode, evValue)

		if evType == 1 && evCode == 116 && evValue == 1 {
			sb := senable
			slaveState = !slaveState
			if !slaveState {
				sb = sdisable
			}
			log.Printf("power key pressed, new slavestate=%v writing=%s", slaveState, string(sb))
			if slave, err = os.OpenFile(slavePath, os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
				log.Fatal("open ", slavePath, err)
			}
			if _, err = slave.WriteAt(sb, 0); err != nil {
				log.Fatal("slave.Write ", err)
			}
			if err = slave.Close(); err != nil {
				log.Fatal("slave.Close ", err)
			}
		}
	}
}
