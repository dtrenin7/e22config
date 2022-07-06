package main

import (
	"fmt"
	"time"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"encoding/hex"
	"e22config/LoRa/E22"
	"go.bug.st/serial.v1"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
)

type Device struct {
	IP 						string
	SerialNumber 	string
	Checksum			string
}

type BTLP struct {
	ADDH         			int `json:"address-num-editable"`
	ADDL    					int `json:"address-num-editable"`
	NETID    					int `json:"address-num-editable"`
	UARTRate					string `json:"uart-string-select"`
	UARTParityBit			string `json:"uart-string-select"`
	Device						string `json:"main-string-select"`
	WirelessRate			string `json:"wireless-string-select"`
	SubPacketLength		string `json:"wireless-string-select"`
	AmbientNoise 			bool `json:"wireless-bool-check"`
	Power							string `json:"wireless-string-select"`
	Channel						string `json:"wireless-string-select"`
	RSSI 							bool `json:"wireless-bool-check"`
	TransmissionMode	string `json:"wireless-string-select"`
	Repeater 					bool `json:"wireless-bool-check"`
	LBT 							bool `json:"wireless-bool-check"`
	WOR								string `json:"wor-string-select"`
	WORCycle					string `json:"wor-string-select"`
	CryptH 						int `json:"cryptography-num-editable"`
	CryptL 						int `json:"cryptography-num-editable"`
	PID								string `json:"product-string-label"`

	names 	map[string]string
	labels  map[string]*widget.Label
	entries  map[string]*widget.Entry
	defaults map[string]string
	selects map[string]*widget.SelectEntry
	selectOptions map[string][]string
	checks map[string]*widget.Check
	devices []Device
	table *widget.Table

	filePath string
	Buttons map[string]*widget.Button
	Progress *widget.ProgressBarInfinite
}

var (
	lora E22.Device
	boot *BTLP
	noResponse string = "ERROR: No response"
	noBootloader string = "ERROR: Bootloader not found"
	width float32 = 680
)

func (x *BTLP) newLabel(name string) *widget.Label {
	w := widget.NewLabel("")
	x.labels[name] = w
	return w
}

// NewBTLP returns a new BTLP app
func NewBTLP() *BTLP {
	rand.Seed(time.Now().UnixNano())
	b := &BTLP{
		names: make(map[string]string),
		labels: make(map[string]*widget.Label),
		entries: make(map[string]*widget.Entry),
		selects: make(map[string]*widget.SelectEntry),
		selectOptions: make(map[string][]string),
		checks: make(map[string]*widget.Check),
		defaults: make(map[string]string),
		Buttons: make(map[string]*widget.Button),
	}
	b.names["Device"] = "Device"
	b.names["ADDH"] = "High byte (ADDH)"
	b.names["ADDL"] = "Low byte (ADDL)"
	b.names["UARTRate"] = "Data rate (bps)"
	b.names["UARTParityBit"] = "Parity bit"
	b.names["NETID"] = "Network (NETID)"
	b.names["WirelessRate"] = "Data rate (bps)"
	b.names["SubPacketLength"] = "Sub packet length (bytes)"
	b.names["AmbientNoise"] = "Enable ambient noise (RSSI)"
	b.names["Power"] = "Transmitting power (dbm)"
	b.names["Channel"] = "Frequency (MHz)"
	b.names["RSSI"] = "Enable RSSI"
	b.names["TransmissionMode"] = "Transmission mode"
	b.names["Repeater"] = "Enable repeater"
	b.names["LBT"] = "Enable monitor before transmission (LBT)"
	b.names["WOR"] = "Transceiver mode"
	b.names["WORCycle"] = "Monitoring interval period (ms)"
	b.names["CryptH"] = "Key high byte (CRYPT_H)"
	b.names["CryptL"] = "Key low byte (CRYPT_L)"
	b.names["PID"] = "PID"

	b.defaults["ADDH"] = "0"
	b.defaults["ADDL"] = "0"
	b.defaults["UARTRate"] = "9600"
	b.defaults["UARTParityBit"] = "8N1"
	b.defaults["NETID"] = "0"
	b.defaults["WirelessRate"] = "2400"
	b.defaults["SubPacketLength"] = "240"
	b.defaults["Power"] = "30"
	b.defaults["Channel"] = "433"
	b.defaults["TransmissionMode"] = "Transparent"
	b.defaults["WOR"] = "Receiver"
	b.defaults["WORCycle"] = "500"
	b.defaults["CryptH"] = "0"
	b.defaults["CryptL"] = "0"
	b.defaults["PID"] = ""

	b.selectOptions["UARTRate"] = []string{"1200", "2400", "4800", "9600", "19200", "38400", "57600", "115200"}
	b.selectOptions["UARTParityBit"] = []string{"8N1", "8O1", "8E1"}
	b.selectOptions["WirelessRate"] = []string{"300", "1200", "2400", "4800", "9600", "19200", "38400", "62500"}
	b.selectOptions["SubPacketLength"] = []string{"240", "128", "64", "32"}
	b.selectOptions["Power"] = []string{"30", "27", "24", "21"}
	ports, err := serial.GetPortsList()
	if err == nil && len(ports) > 0 {
		b.selectOptions["Device"] = ports
		b.defaults["Device"] = ports[0]
	} else {
		b.defaults["Device"] = "/dev/tty.usbserial-0001"
		b.selectOptions["Device"] = []string{"/dev/tty.usbserial-0001"}
	}
	data := []string{}
	for i := 0; i < 84; i++ {
		data = append(data, fmt.Sprintf("%d", i + 410))
	}
	b.selectOptions["Channel"] = data
	b.selectOptions["TransmissionMode"] = []string{"Fixed point", "Transparent"}
	b.selectOptions["WOR"] = []string{"Transmitter", "Receiver"}
	data = []string{}
	for i := 1; i < 9; i++ {
		data = append(data, fmt.Sprintf("%d", i * 500))
	}
	b.selectOptions["WORCycle"] = data

	b.Buttons["read"] = widget.NewButton("Read", func() {
		var err error
		dev := boot.selects["Device"].Text
		log.Println("Trying to open " + dev)
		Serial = NewSerialPort(dev)
		if err = Serial.Open(); err == nil {
			TryCatchBlock {
		    Try: func() {
					boot.DisableButtons()
					boot.SetState("Port " + dev + " opened")
					data, err := Serial.Command(E22.COMMAND_GET_REGISTER[:], E22.GET_CONFIG[:])
					if len(data) == 10 {
						args := data[3:]
						log.Printf("ARGS: %s", hex.Dump(args))
						boot.entries["ADDH"].SetText(fmt.Sprintf("%d", args[E22.REGISTER_ADDH[0]]))
						boot.entries["ADDL"].SetText(fmt.Sprintf("%d", args[E22.REGISTER_ADDL[0]]))
						boot.entries["NETID"].SetText(fmt.Sprintf("%d", args[E22.REGISTER_NETID[0]]))
						switch args[E22.REGISTER_REG0[0]] & E22.MASK_UART_BAUD {
						case E22.UART_BAUD_1200:
							boot.selects["UARTRate"].SetText("1200")
						case E22.UART_BAUD_2400:
							boot.selects["UARTRate"].SetText("2400")
						case E22.UART_BAUD_4800:
							boot.selects["UARTRate"].SetText("4800")
						case E22.UART_BAUD_9600:
							boot.selects["UARTRate"].SetText("9600")
						case E22.UART_BAUD_19200:
							boot.selects["UARTRate"].SetText("19200")
						case E22.UART_BAUD_38400:
							boot.selects["UARTRate"].SetText("38400")
						case E22.UART_BAUD_57600:
							boot.selects["UARTRate"].SetText("57600")
						case E22.UART_BAUD_115200:
							boot.selects["UARTRate"].SetText("115200")
						default:
							Throw("Unknown UARTRate value!")
						}
						switch args[E22.REGISTER_REG0[0]] & E22.MASK_UART_PARITY {
						case E22.UART_8N1:
							boot.selects["UARTParityBit"].SetText("8N1")
						case E22.UART_8O1:
							boot.selects["UARTParityBit"].SetText("8O1")
						case E22.UART_8E1:
							boot.selects["UARTParityBit"].SetText("8E1")
						default:
							Throw("Unknown UARTParityBit value!")
						}
						switch args[E22.REGISTER_REG0[0]] & E22.MASK_AIR_BAUD {
						case E22.AIR_BAUD_300:
							boot.selects["WirelessRate"].SetText("300")
						case E22.AIR_BAUD_1200:
							boot.selects["WirelessRate"].SetText("1200")
						case E22.AIR_BAUD_2400:
							boot.selects["WirelessRate"].SetText("2400")
						case E22.AIR_BAUD_4800:
							boot.selects["WirelessRate"].SetText("4800")
						case E22.AIR_BAUD_9600:
							boot.selects["WirelessRate"].SetText("9600")
						case E22.AIR_BAUD_19200:
							boot.selects["WirelessRate"].SetText("19200")
						case E22.AIR_BAUD_38400:
							boot.selects["WirelessRate"].SetText("38400")
						case E22.AIR_BAUD_62500:
							boot.selects["WirelessRate"].SetText("62500")
						default:
							Throw("Unknown WirelessRate value!")
						}
						switch args[E22.REGISTER_REG1[0]] & E22.MASK_SUB_PACKET {
						case E22.SUB_PACKET_BYTES_240:
							boot.selects["SubPacketLength"].SetText("240")
						case E22.SUB_PACKET_BYTES_128:
							boot.selects["SubPacketLength"].SetText("128")
						case E22.SUB_PACKET_BYTES_64:
							boot.selects["SubPacketLength"].SetText("64")
						case E22.SUB_PACKET_BYTES_32:
							boot.selects["SubPacketLength"].SetText("32")
						default:
							Throw("Unknown SubPacketLength value!")
						}
						switch args[E22.REGISTER_REG1[0]] & E22.MASK_AMBIENT_NOISE {
						case E22.AMBIENT_NOISE_ENABLE:
							boot.checks["AmbientNoise"].SetChecked(true)
						case E22.AMBIENT_NOISE_DISABLE:
							boot.checks["AmbientNoise"].SetChecked(false)
						default:
							Throw("Unknown AmbientNoise value!")
						}
						switch args[E22.REGISTER_REG1[0]] & E22.MASK_POWER {
						case E22.POWER_DBM_30:
							boot.selects["Power"].SetText("30")
						case E22.POWER_DBM_27:
							boot.selects["Power"].SetText("27")
						case E22.POWER_DBM_24:
							boot.selects["Power"].SetText("24")
						case E22.POWER_DBM_21:
							boot.selects["Power"].SetText("21")
						default:
							Throw("Unknown Power value!")
						}
						boot.selects["Channel"].SetText(fmt.Sprintf("%d", int(args[E22.REGISTER_REG2[0]]) + 410 ))
						switch args[E22.REGISTER_REG3[0]] & E22.MASK_RSSI {
						case E22.RSSI_ENABLE:
							boot.checks["RSSI"].SetChecked(true)
						case E22.RSSI_DISABLE:
							boot.checks["RSSI"].SetChecked(false)
						default:
							Throw("Unknown RSSI value!")
						}
						switch args[E22.REGISTER_REG3[0]] & E22.MASK_TRANSMISSION_MODE {
						case E22.TRANSMISSION_MODE_FIXED:
							boot.selects["TransmissionMode"].SetText("Fixed point")
						case E22.TRANSMISSION_MODE_TRANSPARENT:
							boot.selects["TransmissionMode"].SetText("Transparent")
						default:
							Throw("Unknown TransmissionMode value!")
						}
						switch args[E22.REGISTER_REG3[0]] & E22.MASK_REPEATER {
						case E22.REPEATER_ENABLE:
							boot.checks["Repeater"].SetChecked(true)
						case E22.REPEATER_DISABLE:
							boot.checks["Repeater"].SetChecked(false)
						default:
							Throw("Unknown Repeater value!")
						}
						switch args[E22.REGISTER_REG3[0]] & E22.MASK_LBT {
						case E22.LBT_ENABLE:
							boot.checks["LBT"].SetChecked(true)
						case E22.LBT_DISABLE:
							boot.checks["LBT"].SetChecked(false)
						default:
							Throw("Unknown LBT value!")
						}
						switch args[E22.REGISTER_REG3[0]] & E22.MASK_WOR_CONTROL {
						case E22.WOR_TRANSMITTER:
							boot.selects["WOR"].SetText("Transmitter")
						case E22.WOR_RECEIVER:
							boot.selects["WOR"].SetText("Receiver")
						default:
							Throw("Unknown WOR value!")
						}
						boot.selects["WORCycle"].SetText(fmt.Sprintf("%d", (int(args[E22.REGISTER_REG3[0]] & E22.MASK_WOR_CYCLE) + 1) * 500 ))
						boot.SetState("Reading DONE")
					} else {
						logError("Reading", err)
					}
					data, err = Serial.Command(E22.COMMAND_GET_REGISTER[:], E22.GET_PRODUCT_INFO[:])
					if err == nil {
						boot.labels["PID"].SetText(hex.Dump(data[3:]))
					} else {
						logError("Reading", err)
					}
				},
				Catch: func(e Exception) {
					log.Printf("%v\n", e)
					logError("Writing", fmt.Errorf("%v", e))
				},
				Finally: func() {
					Serial.Close()
					boot.EnableButtons()
				},
			}.Do()
		} else {
			logError("Port opening", err)
		}
	})
	b.Buttons["write"] = widget.NewButton("Write", func() {
		var err error
		dev := boot.selects["Device"].Text
		log.Println("Trying to open " + dev)
		Serial = NewSerialPort(dev)
		if err = Serial.Open(); err == nil {
		  TryCatchBlock {
		    Try: func() {
					boot.DisableButtons()
					boot.SetState("Port " + dev + " opened")
					args := [9]byte{}
					args[E22.REGISTER_ADDH[0]] = s2b(boot.entries["ADDH"].Text)
					args[E22.REGISTER_ADDL[0]] = s2b(boot.entries["ADDL"].Text)
					args[E22.REGISTER_NETID[0]] = s2b(boot.entries["NETID"].Text)
					switch boot.selects["UARTRate"].Text {
					case "1200":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_1200
					case "2400":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_2400
					case "4800":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_4800
					case "9600":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_9600
					case "19200":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_19200
					case "38400":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_38400
					case "57600":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_57600
					case "115200":
						args[E22.REGISTER_REG0[0]] |= E22.UART_BAUD_115200
					default:
						Throw("Unknown UARTRate value!")
					}
					switch boot.selects["UARTParityBit"].Text {
					case "8N1":
						args[E22.REGISTER_REG0[0]] |= E22.UART_8N1
					case "8O1":
						args[E22.REGISTER_REG0[0]] |= E22.UART_8O1
					case "8E1":
						args[E22.REGISTER_REG0[0]] |= E22.UART_8E1
					default:
						Throw("Unknown UARTParityBit value!")
					}
					switch boot.selects["WirelessRate"].Text {
					case "300":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_300
					case "1200":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_1200
					case "2400":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_2400
					case "4800":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_4800
					case "9600":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_9600
					case "19200":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_19200
					case "38400":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_38400
					case "62500":
						args[E22.REGISTER_REG0[0]] |= E22.AIR_BAUD_62500
					default:
						Throw("Unknown WirelessRate value!")
					}
					switch boot.selects["SubPacketLength"].Text {
					case "240":
						args[E22.REGISTER_REG1[0]] |= E22.SUB_PACKET_BYTES_240
					case "128":
						args[E22.REGISTER_REG1[0]] |= E22.SUB_PACKET_BYTES_128
					case "64":
						args[E22.REGISTER_REG1[0]] |= E22.SUB_PACKET_BYTES_64
					case "32":
						args[E22.REGISTER_REG1[0]] |= E22.SUB_PACKET_BYTES_32
					default:
						Throw("Unknown SubPacketLength value!")
					}
					if boot.checks["AmbientNoise"].Checked {
						args[E22.REGISTER_REG1[0]] |= E22.AMBIENT_NOISE_ENABLE
					}
					switch boot.selects["Power"].Text {
					case "30":
						args[E22.REGISTER_REG1[0]] |= E22.POWER_DBM_30
					case "27":
						args[E22.REGISTER_REG1[0]] |= E22.POWER_DBM_27
					case "24":
						args[E22.REGISTER_REG1[0]] |= E22.POWER_DBM_24
					case "21":
						args[E22.REGISTER_REG1[0]] |= E22.POWER_DBM_21
					default:
						Throw("Unknown Power value!")
					}
					channel := getInt(boot.selects["Channel"].Text) - 410
					if channel < 0 || channel > 83 {
						Throw("Unknown Channel value!")
					}
					args[E22.REGISTER_REG2[0]] = byte(channel)
					if boot.checks["RSSI"].Checked {
						args[E22.REGISTER_REG3[0]] |= E22.RSSI_ENABLE
					}
					switch boot.selects["TransmissionMode"].Text {
					case "Fixed point":
						args[E22.REGISTER_REG3[0]] |= E22.TRANSMISSION_MODE_FIXED
					case "Transparent":
						args[E22.REGISTER_REG3[0]] |= E22.TRANSMISSION_MODE_TRANSPARENT
					default:
						Throw("Unknown TransmissionMode value!")
					}
					if boot.checks["Repeater"].Checked {
						args[E22.REGISTER_REG3[0]] |= E22.REPEATER_ENABLE
					}
					if boot.checks["LBT"].Checked {
						args[E22.REGISTER_REG3[0]] |= E22.LBT_ENABLE
					}
					switch boot.selects["WOR"].Text {
					case "Transmitter":
						args[E22.REGISTER_REG3[0]] |= E22.WOR_TRANSMITTER
					case "Receiver":
						args[E22.REGISTER_REG3[0]] |= E22.WOR_RECEIVER
					default:
						Throw("Unknown WOR value!")
					}
					cycle := byte((getInt(boot.selects["WORCycle"].Text) / 500) - 1)
					if cycle > E22.WOR_CYCLE_MS_4000 {
						Throw("Unknown WORCycle value!")
					}
					args[E22.REGISTER_REG3[0]] |= cycle
					args[E22.REGISTER_CRYPT_H[0]] = s2b(boot.entries["CryptH"].Text)
					args[E22.REGISTER_CRYPT_L[0]] = s2b(boot.entries["CryptL"].Text)
					data, err := Serial.Command(E22.COMMAND_SET_REGISTER[:], E22.SET_CONFIG[:], args[:])
					if err != nil {
						logError("Writing", err)
					} else {
						log.Printf("WRITE: %s", hex.Dump(data[3:]))
						boot.SetState("Writing DONE")
					}
		    },
		    Catch: func(e Exception) {
		      log.Printf("%v\n", e)
		      logError("Writing", fmt.Errorf("%v", e))
		    },
				Finally: func() {
					Serial.Close()
					boot.EnableButtons()
				},
		  }.Do()
		}
	})
	return b
}

func (x *BTLP) SetState(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	x.labels["State"].SetText(text)
}

// NewForm generates a new BTLP form
func (x *BTLP) NewForm(w fyne.Window, tab string) *widget.Form {
	form := &widget.Form{}
	tt := reflect.TypeOf(x).Elem()
	for i := 0; i < tt.NumField(); i++ {
		fld := tt.Field(i)
		tag := fld.Tag.Get("json")
		if !strings.HasPrefix(tag, tab) {
			continue
		}
		if strings.HasSuffix(tag, "editable") {
			entry := widget.NewEntry()
			entry.SetText(x.defaults[fld.Name])
			form.Append(x.names[fld.Name], entry)
			x.entries[fld.Name] = entry
			continue
		}
		if strings.HasSuffix(tag, "select") {
			entry := widget.NewSelectEntry(x.selectOptions[fld.Name])
			entry.SetText(x.defaults[fld.Name])
			form.Append(x.names[fld.Name], entry)
			x.selects[fld.Name] = entry
			continue
		}
		if strings.HasSuffix(tag, "check") {
			entry := widget.NewCheck("", func(bool) {})
			form.Append(x.names[fld.Name], entry)
			x.checks[fld.Name] = entry
			continue
		}
		entry := x.newLabel(tag)
		entry.SetText(x.defaults[fld.Name])
		form.Append(x.names[fld.Name], entry)
		x.labels[fld.Name] = entry
	}
	return form
}

func (x *BTLP) DisableButtons() {
	for _, item := range x.entries {
		item.Disable()
	}
	for _, item := range x.selects {
		item.Disable()
	}
	for _, item := range x.checks {
		item.Disable()
	}
	x.Buttons["read"].Disable()
	x.Buttons["write"].Disable()
	x.Progress.Show()
}

func (x *BTLP) EnableButtons() {
	for _, item := range x.entries {
		item.Enable()
	}
	for _, item := range x.selects {
		item.Enable()
	}
	for _, item := range x.checks {
		item.Enable()
	}
	x.Buttons["read"].Enable()
	x.Buttons["write"].Enable()
	x.Progress.Hide()
}

func b2s(data byte) string {
	return fmt.Sprintf("%v", data)
}

func s2b(str string) byte {
	if str == "" {
		Throw(makeError(fmt.Errorf("uint not parsed"), FileLine()).Error())
	}
	return getByte(str)
}

func logError(op string, err error) {
	errtext := fmt.Sprintf("%s FAILED: %v", op, err)
	log.Println(errtext)
	boot.SetState(errtext)
}

func addTabItem(win fyne.Window, name, heading string, read, write bool) *container.TabItem {
	form := boot.NewForm(win, name)
	form.Append("", layout.NewSpacer())
	readButton := layout.NewSpacer()
	if read {
		readButton = boot.Buttons["read"]
	}
	writeButton := layout.NewSpacer()
	if write {
		writeButton = boot.Buttons["write"]
	}
	buttons := container.NewHBox(
		layout.NewSpacer(),
		readButton,
		writeButton,
		layout.NewSpacer(),
	)
	return container.NewTabItem(heading,
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(form, buttons, nil, nil),
			form, buttons),
	)
}

func makeTableTab(win fyne.Window) *widget.Table {
	t := widget.NewTable(
		func() (int, int) { return len(boot.devices), 3 }, // number of cells (rows, cols)
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell 000, 000")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(boot.devices[id.Row].IP)
			case 1:
				label.SetText(boot.devices[id.Row].SerialNumber)
			case 2:
				label.SetText(boot.devices[id.Row].Checksum)
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	return t
}

func Show(win fyne.Window) fyne.CanvasObject {
	form := boot.NewForm(win, "main")
	state := boot.newLabel("State")
	state.SetText("Ready to work.")
	boot.Progress = widget.NewProgressBarInfinite()
	boot.Progress.Hide()
	states := container.NewVBox(
		form,
		state,
		boot.Progress,
		layout.NewSpacer(),
	)
	border := container.NewBorder(nil, nil, nil, nil,
		container.NewAppTabs(
			addTabItem(win, "address", "Address", true, true),
			addTabItem(win, "uart", "UART", true, true),
			addTabItem(win, "wireless", "Wireless", true, true),
			addTabItem(win, "wor", "WOR", true, true),
			addTabItem(win, "crypto", "Cryptography", false, true),
			addTabItem(win, "product", "Product Information", true, false),
		))
	box := container.NewVBox(
		states,
		layout.NewSpacer(),
		border,
		layout.NewSpacer(),
	)
	return box
}

func createGUI() fyne.Window {
	log.Println("Starting GUI...")
	a := app.New()
	w := a.NewWindow("E22-400T30D Module Configuration Utility")
	w.SetContent(Show(w))
	w.Resize(fyne.NewSize(width, 200))
	return w
}

func main() {
	rand.Seed(time.Now().Unix())
	boot = NewBTLP()
	window := createGUI()
	window.ShowAndRun()
}
