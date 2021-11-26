package E22

import (

)

var (
  FACTORY_DEFAULTS                    = [...]byte{0xC0, 0x00, 0x00, 0x62, 0x00, 0x00}
  GET_CONFIG                          = [...]byte{0x00, 0x07}
  SET_CONFIG                          = [...]byte{0x00, 0x09}
  GET_PRODUCT_INFO                        = [...]byte{0x80, 0x07}

  COMMAND_SET_REGISTER                = [...]byte{0xC0}
  COMMAND_GET_REGISTER                = [...]byte{0xC1}
  COMMAND_SET_TEMPORARY_REGISTER      = [...]byte{0xC2}
  COMMAND_WIRELESS_CONFIG             = [...]byte{0xCF, 0xCF}

  RESPONSE_WRONG_FORMAT               = [...]byte{0xFF, 0xFF, 0xFF}

  REGISTER_ADDH                       = [...]byte{0x00}
  REGISTER_ADDL                       = [...]byte{0x01}
  REGISTER_NETID                      = [...]byte{0x02}
  REGISTER_REG0                       = [...]byte{0x03}
  REGISTER_REG1                       = [...]byte{0x04}
  REGISTER_REG2                       = [...]byte{0x05}
  REGISTER_REG3                       = [...]byte{0x06}
  REGISTER_CRYPT_H                    = [...]byte{0x07}
  REGISTER_CRYPT_L                    = [...]byte{0x08}
  REGISTER_PID0                       = [...]byte{0x80}
  REGISTER_PID1                       = [...]byte{0x81}
  REGISTER_PID2                       = [...]byte{0x82}
  REGISTER_PID3                       = [...]byte{0x83}
  REGISTER_PID4                       = [...]byte{0x84}
  REGISTER_PID5                       = [...]byte{0x85}
  REGISTER_PID6                       = [...]byte{0x86}

  BROADCAST                           = [...]byte{0xFF} // ADDH, ADDL

  UART_BAUD_1200 byte                 = 0x00
  UART_BAUD_2400 byte                 = 0x20
  UART_BAUD_4800 byte                 = 0x40
  UART_BAUD_9600 byte                 = 0x60
  UART_BAUD_19200 byte                = 0x80
  UART_BAUD_38400 byte                = 0xA0
  UART_BAUD_57600 byte                = 0xC0
  UART_BAUD_115200 byte               = 0xE0

  UART_8N1 byte                       = 0x00
  UART_8O1 byte                       = 0x08
  UART_8E1 byte                       = 0x10

  AIR_BAUD_300 byte                   = 0x00
  AIR_BAUD_1200 byte                  = 0x01
  AIR_BAUD_2400 byte                  = 0x02
  AIR_BAUD_4800 byte                  = 0x03
  AIR_BAUD_9600 byte                  = 0x04
  AIR_BAUD_19200 byte                 = 0x05
  AIR_BAUD_38400 byte                 = 0x06
  AIR_BAUD_62500 byte                 = 0x07

  POWER_DBM_30 byte                   = 0x00
  POWER_DBM_27 byte                   = 0x01
  POWER_DBM_24 byte                   = 0x02
  POWER_DBM_21 byte                   = 0x03

  AMBIENT_NOISE_ENABLE byte           = 0x20
  AMBIENT_NOISE_DISABLE byte          = 0x00

  SUB_PACKET_BYTES_240 byte           = 0x00
  SUB_PACKET_BYTES_128 byte           = 0x40
  SUB_PACKET_BYTES_64 byte            = 0x80
  SUB_PACKET_BYTES_32 byte            = 0xC0

  RSSI_ENABLE byte                    = 0x80
  RSSI_DISABLE byte                   = 0x00

  TRANSMISSION_MODE_FIXED byte        = 0x40
  TRANSMISSION_MODE_TRANSPARENT byte  = 0x00

  REPEATER_ENABLE byte                = 0x20
  REPEATER_DISABLE byte               = 0x00

  LBT_ENABLE byte                     = 0x10
  LBT_DISABLE byte                    = 0x00

  WOR_TRANSMITTER byte                = 0x08
  WOR_RECEIVER byte                   = 0x00

  WOR_CYCLE_MS_500 byte               = 0x00
  WOR_CYCLE_MS_1000 byte              = 0x01
  WOR_CYCLE_MS_1500 byte              = 0x02
  WOR_CYCLE_MS_2000 byte              = 0x03
  WOR_CYCLE_MS_2500 byte              = 0x04
  WOR_CYCLE_MS_3000 byte              = 0x05
  WOR_CYCLE_MS_3500 byte              = 0x06
  WOR_CYCLE_MS_4000 byte              = 0x07

  MASK_AIR_BAUD byte                  = 0x07
  MASK_UART_BAUD byte                 = 0xE0
  MASK_UART_PARITY byte               = 0x18
  MASK_POWER byte                     = 0x03
  MASK_AMBIENT_NOISE byte             = 0x20
  MASK_SUB_PACKET byte                = 0xC0
  MASK_RSSI byte                      = 0x80
  MASK_TRANSMISSION_MODE byte         = 0x40
  MASK_REPEATER byte                  = 0x20
  MASK_LBT byte                       = 0x10
  MASK_WOR_CONTROL byte               = 0x08
  MASK_WOR_CYCLE byte                 = 0x07
)

type Device struct {
  IsOpen bool
}

func (d *Device) Open(name string) error {
  return nil
}

func (d *Device) Close() {

}
