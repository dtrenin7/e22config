package E32

import (

)

var (
  FACTORY_DEFAULTS                    = [...]byte{0x00, 0x00, 0x1A, 0x17, 0x44}

  COMMAND_SET_PARAMETERS              = [...]byte{0xC0}
  COMMAND_GET_PARAMETERS              = [...]byte{0xC1, 0xC1, 0xC1}
  COMMAND_GET_MODULE_VERSION          = [...]byte{0xC4, 0xC4, 0xC4}
  COMMAND_SET_PARAMETERS_TEMPORARY    = [...]byte{0xC2}

  RESPONSE_WRONG_FORMAT               = [...]byte{0xFF, 0xFF, 0xFF}

  REGISTER_ADDH                       = [...]byte{0x01}
  REGISTER_ADDL                       = [...]byte{0x02}
  REGISTER_SPED                       = [...]byte{0x03}
  REGISTER_CHAN                       = [...]byte{0x04}
  REGISTER_OPTION                     = [...]byte{0x05}

  BROADCAST                           = [...]byte{0xFF} // ADDH, ADDL

  UART_BAUD_1200 byte                 = 0x00
  UART_BAUD_2400 byte                 = 0x08
  UART_BAUD_4800 byte                 = 0x10
  UART_BAUD_9600 byte                 = 0x18
  UART_BAUD_19200 byte                = 0x20
  UART_BAUD_38400 byte                = 0x28
  UART_BAUD_57600 byte                = 0x30
  UART_BAUD_115200 byte               = 0x38

  UART_8N1 byte                       = 0x00
  UART_8O1 byte                       = 0x40
  UART_8E1 byte                       = 0x80
  UART_8N1_2 byte                     = 0xC0

  AIR_BAUD_300 byte                   = 0x00
  AIR_BAUD_1200 byte                  = 0x01
  AIR_BAUD_2400 byte                  = 0x02
  AIR_BAUD_4800 byte                  = 0x03
  AIR_BAUD_9600 byte                  = 0x04
  AIR_BAUD_19200 byte                 = 0x05
  AIR_BAUD_19200_2 byte               = 0x06
  AIR_BAUD_19200_3 byte               = 0x07

  POWER_DBM_30 byte                   = 0x00
  POWER_DBM_27 byte                   = 0x01
  POWER_DBM_24 byte                   = 0x02
  POWER_DBM_21 byte                   = 0x03

  FEC_ENABLE byte                     = 0x02
  FEC_DISABLE byte                    = 0x00

  TRANSMISSION_MODE_FIXED byte        = 0x80
  TRANSMISSION_MODE_TRANSPARENT byte  = 0x00

  DRIVE_MODE_PUSH_PULL byte           = 0x40
  DRIVE_MODE_OPEN byte                = 0x00

  WAKE_UP_MS_250 byte                 = 0x00
  WAKE_UP_MS_500 byte                 = 0x08
  WAKE_UP_MS_750 byte                 = 0x10
  WAKE_UP_MS_1000 byte                = 0x18
  WAKE_UP_MS_1250 byte                = 0x20
  WAKE_UP_MS_1500 byte                = 0x28
  WAKE_UP_MS_1750 byte                = 0x30
  WAKE_UP_MS_2000 byte                = 0x38

  MASK_AIR_BAUD byte                  = 0x07
  MASK_UART_BAUD byte                 = 0x38
  MASK_UART_PARITY byte               = 0xC0
  MASK_CHANNEL byte                   = 0x1F
  MASK_POWER byte                     = 0x03
  MASK_FEC byte                       = 0x02
  MASK_WAKE_UP byte                   = 0x38
  MASK_DRIVE_MODE byte                = 0x40
  MASK_TRANSMISSION_MODE byte         = 0x80
)
