package test

import (
	"encoding/json"

	"os"
	"strconv"
	"testing"

	"github.com/blewjy/fire-gb/gb"
)

type CPUData struct {
	A  string `json:"a"`
	B  string `json:"b"`
	C  string `json:"c"`
	D  string `json:"d"`
	E  string `json:"e"`
	F  string `json:"f"`
	H  string `json:"h"`
	L  string `json:"l"`
	PC string `json:"pc"`
	SP string `json:"sp"`
}

type InitialData struct {
	CPU CPUData     `json:"cpu"`
	RAM [][2]string `json:"ram"`
}

type FinalData struct {
	CPU CPUData     `json:"cpu"`
	RAM [][2]string `json:"ram"`
}

type JSONData struct {
	Name    string      `json:"name"`
	Initial InitialData `json:"initial"`
	Final   FinalData   `json:"final"`
	Cycles  [][3]string `json:"cycles"`
}

func parseHexStringToUint8(s string) uint8 {
	hexValue, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		panic(err)
	}
	return uint8(hexValue)
}

func parseHexStringToUint16(s string) uint16 {
	hexValue, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		panic(err)
	}
	return uint16(hexValue)
}

func TestAdtennantCpu(t *testing.T) {

	testCaseFile, err := os.ReadFile("adtennant/v1/1a.json")
	if err != nil {
		panic(err)
	}

	var allData []JSONData
	if err := json.Unmarshal(testCaseFile, &allData); err != nil {
		panic(err)
	}

	for _, data := range allData {
		t.Run(data.Name, func(t *testing.T) {
			tgb := gb.InitTestGameboy(gb.TestGameboyState{
				CPU: gb.TestGameboyStateCPU{
					A:  parseHexStringToUint8(data.Initial.CPU.A),
					B:  parseHexStringToUint8(data.Initial.CPU.B),
					C:  parseHexStringToUint8(data.Initial.CPU.C),
					D:  parseHexStringToUint8(data.Initial.CPU.D),
					E:  parseHexStringToUint8(data.Initial.CPU.E),
					F:  parseHexStringToUint8(data.Initial.CPU.F),
					H:  parseHexStringToUint8(data.Initial.CPU.H),
					L:  parseHexStringToUint8(data.Initial.CPU.L),
					PC: parseHexStringToUint16(data.Initial.CPU.PC),
					SP: parseHexStringToUint16(data.Initial.CPU.SP),
				},
				RAM: func() map[uint16]uint8 {
					m := map[uint16]uint8{}
					for _, ramData := range data.Initial.RAM {
						m[parseHexStringToUint16(ramData[0])] = parseHexStringToUint8(ramData[1])
					}
					return m
				}(),
				ROM: make([]uint8, 0xFFFF),
			})

			tgb.StepCPU()

			finalState := tgb.ExportState()
			wantFinalState := gb.TestGameboyState{
				CPU: gb.TestGameboyStateCPU{
					A:  parseHexStringToUint8(data.Final.CPU.A),
					B:  parseHexStringToUint8(data.Final.CPU.B),
					C:  parseHexStringToUint8(data.Final.CPU.C),
					D:  parseHexStringToUint8(data.Final.CPU.D),
					E:  parseHexStringToUint8(data.Final.CPU.E),
					F:  parseHexStringToUint8(data.Final.CPU.F),
					H:  parseHexStringToUint8(data.Final.CPU.H),
					L:  parseHexStringToUint8(data.Final.CPU.L),
					PC: parseHexStringToUint16(data.Final.CPU.PC),
					SP: parseHexStringToUint16(data.Final.CPU.SP),
				},
				RAM: map[uint16]uint8{
					parseHexStringToUint16(data.Final.RAM[0][0]): parseHexStringToUint8(data.Final.RAM[0][1]),
				},
			}

			if finalState.CPU != wantFinalState.CPU {
				t.Errorf("%v: got = %+v, want = %+v", data.Name, finalState.CPU, wantFinalState.CPU)
			}
		})
	}
}
