package main

import (
  "fmt"
	"time"
  "log"
  "context"
  "errors"
  "encoding/hex"
  "go.bug.st/serial.v1"
)

type SerialPort struct {
  Device string
  Config serial.Mode
  port serial.Port
}

var (
	Serial *SerialPort
  NoResponse = errors.New("No response from device")
)

func (s *SerialPort) Open() error {
  var err error
  s.port, err = serial.Open(s.Device, &s.Config)
  return err
}

func (s *SerialPort) Write(data []byte) error {
  if s.port == nil {
    return makeError(fmt.Errorf("port is nil"), FileLine())
  }
  _, err := s.port.Write(data)
  return err
}

type SerialResponse struct {
  Data []byte
  Err error
}

func (s *SerialPort) Command(data ...[]byte) ([]byte, error) {
  cmd := []byte{}
  for _, item := range data {
    cmd = append(cmd, item...)
  }

  channel := make(chan SerialResponse)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

  go func() {
    if err := s.Write(cmd); err != nil {
      channel <- SerialResponse{[]byte{}, NoResponse}
      return
    }
    log.Printf("<--- (%d bytes) %s", len(cmd), hex.Dump(cmd))
    time.Sleep(500 * time.Millisecond)
    res := SerialResponse{}
    res.Data, res.Err = s.Read()
    channel <- res
  }()

  select {
  case result := <-channel:
    return result.Data, result.Err
	case <-ctx.Done():
    return []byte{}, NoResponse
	}
  return []byte{}, NoResponse
}

func (s *SerialPort) Read() ([]byte, error) {
  if s.port == nil {
    return []byte{}, makeError(fmt.Errorf("port is nil"), FileLine())
  }
  buffer := make([]byte, 1024)
  n, err := s.port.Read(buffer)
  if err == nil && n > 0 {
    log.Printf("===> (%d bytes) %s", n, hex.Dump(buffer[:n]))
    return buffer[:n], nil
  }
  return []byte{}, err
}

func (s *SerialPort) Close() error {
  if s.port != nil {
    return s.port.Close()
  }
  return nil
}

func NewSerialPort(device string) *SerialPort {
	return &SerialPort{
    Device: device,
    Config: serial.Mode{
       BaudRate: 9600,
       DataBits: 8,
       StopBits: serial.OneStopBit,
       Parity: serial.NoParity,
     },
	}
}
