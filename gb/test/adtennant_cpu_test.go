package test

import (
	"encoding/json"
	"fmt"
	"github.com/blewjy/fire-gb/gb"
	"os"
	"strconv"
	"testing"
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
	dir := "./adtennant/v1" // Specify the directory you want to read

	// Open the directory
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read the directory content
	files, err := f.Readdir(-1) // Use -1 to read all entries, or a positive number to limit the number of entries read
	if err != nil {
		panic(err)
	}

	skip := []string{"10.json"}

	// Loop through each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		shouldSkip := false
		for _, s := range skip {
			if file.Name() == s {
				shouldSkip = true
			}
		}
		if shouldSkip {
			continue
		}

		fileName := fmt.Sprintf("adtennant/v1/%s", file.Name())
		testCaseFile, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}

		var allData []JSONData
		if err := json.Unmarshal(testCaseFile, &allData); err != nil {
			panic(err)
		}

		for _, data := range allData {
			skipThis := false
			for _, cycleData := range data.Cycles {
				if cycleData[0] == "0xff04" {
					skipThis = true
				}
			}

			if skipThis {
				continue
			}

			t.Run(data.Name, func(t *testing.T) {
				tgb := gb.InitWithoutDisplay(make([]uint8, 0xFFFF))
				tgb.SetState(
					gb.State{
						CPU: gb.CPUState{
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
					})

				tgb.StepCPU()

				var targetRamAddresses []uint16
				for _, ramData := range data.Final.RAM {
					targetRamAddresses = append(targetRamAddresses, parseHexStringToUint16(ramData[0]))
				}

				finalState := tgb.ExportStateWithAddresses(targetRamAddresses)
				wantFinalState := gb.State{
					CPU: gb.CPUState{
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
					RAM: func() map[uint16]uint8 {
						m := map[uint16]uint8{}
						for _, ramData := range data.Final.RAM {
							m[parseHexStringToUint16(ramData[0])] = parseHexStringToUint8(ramData[1])
						}
						return m
					}(),
				}

				if finalState.CPU != wantFinalState.CPU {
					t.Errorf("%v: got = %+v, want = %+v", data.Name, finalState.CPU, wantFinalState.CPU)
				}

				for ramAddr, ramData := range wantFinalState.RAM {
					if finalState.RAM[ramAddr] != ramData {
						t.Errorf("%v: RAM addr %04x, got = %+v, want = %+v", data.Name, ramAddr, finalState.RAM[ramAddr], ramData)
					}
				}
			})
		}
	}
}
