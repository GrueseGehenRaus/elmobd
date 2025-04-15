package elmobd

import (
	"fmt"
	"strings"
	"time"
)

var lehm int32 = -1

/*==============================================================================
 * External
 */

// MockResult represents the raw text output of running a raw command,
// including information used in debugging to show what input caused what
// error, how long the command took, etc.
type MockResult struct {
	input     string
	outputs   []string
	error     error
	writeTime time.Duration
	readTime  time.Duration
	totalTime time.Duration
}

// Failed checks if the result is successful or not
func (res *MockResult) Failed() bool {
	return res.error != nil
}

// GetError returns the results current error
func (res *MockResult) GetError() error {
	return res.error
}

// GetOutputs returns the outputs of the result
func (res *MockResult) GetOutputs() []string {
	return res.outputs
}

// FormatOverview formats a result as an overview of what command was run and
// how long it took.
func (res *MockResult) FormatOverview() string {
	lines := []string{
		"=======================================",
		" Mocked command \"%s\"",
		"=======================================",
	}

	return fmt.Sprintf(
		strings.Join(lines, "\n"),
		res.input,
	)
}

// MockDevice represent a mocked serial connection
type MockDevice struct {
}

// RunCommand mocks the given AT/OBD command by just returning a result for the
// mocked outputs set earlier.
func (dev *MockDevice) RunCommand(command string) RawResult {
	return &MockResult{
		input:     command,
		outputs:   mockOutputs(command),
		writeTime: 0,
		readTime:  0,
		totalTime: 0,
	}
}

/*==============================================================================
 * Internal
 */

func mockMode1Outputs(subcmd string) []string {
	if strings.HasPrefix(subcmd, "00") {
		// PIDs supported part 1
		return []string{
			"41 00 0C 10 00 00", // Means PIDs supported: 05, 06, 0C
		}
	} else if strings.HasPrefix(subcmd, "20") {
		// PIDs supported part 2
		return []string{
			"41 20 00 00 00 00",
		}
	} else if strings.HasPrefix(subcmd, "40") {
		// PIDs supported part 3
		return []string{
			"41 40 00 00 00 00",
		}
	} else if strings.HasPrefix(subcmd, "60") {
		// PIDs supported part 4
		return []string{
			"41 60 00 00 00 00",
		}
	} else if strings.HasPrefix(subcmd, "80") {
		// PIDs supported part 5
		return []string{
			"41 80 00 00 00 00",
		}
	} else if strings.HasPrefix(subcmd, "01") {
		return []string{
			"41 01 FF 00 00 00",
		}
	} else if strings.HasPrefix(subcmd, "05") { // Engine coolant temperature
		return []string{
			"41 05 4F", // 39 C
		}
	} else if strings.HasPrefix(subcmd, "06") { // Short term fuel trim - Bank 1
		return []string{
			"41 06 02", // -98.4375%
		}
	} else if strings.HasPrefix(subcmd, "0C") { // Engine speed
		if lehm > 26 {
			lehm = -1
		}
		lehm += 1
		lehms := []string{
			"41 0C 13 40",
			"41 0C 15 28",
			"41 0C 17 10",
			"41 0C 19 F8",
			"41 0C 1B E0",
			"41 0C 1D C8",
			"41 0C 1F B0",
			"41 0C 21 98",
			"41 0C 23 80",
			"41 0C 25 68",
			"41 0C 27 50",
			"41 0C 29 38",
			"41 0C 2A 20",
			"41 0C 2C 08",
			"41 0C 2E F0",
			"41 0C 31 D8",
			"41 0C 2E F0",
			"41 0C 2C 08",
			"41 0C 2A 20",
			"41 0C 29 38",
			"41 0C 27 50",
			"41 0C 25 68",
			"41 0C 23 80",
			"41 0C 21 98",
			"41 0C 1F B0",
			"41 0C 1D C8",
			"41 0C 1B E0",
			"41 0C 19 F8",
			"41 0C 17 10",
			"41 0C 15 28",
			"41 0C 13 40",
		}
		return []string{
			lehms[lehm],
		}
	} else if strings.HasPrefix(subcmd, "2F") { // Fuel tank level input
		return []string{
			"41 2F 6B", // 41.96%
		}
	} else if strings.HasPrefix(subcmd, "0D") { // Vehicle speed
		lehms := []string{
			"41 0D 00",
			"41 0D 03",
			"41 0D 05",
			"41 0D 0A",
			"41 0D 0F",
			"41 0D 14",
			"41 0D 19",
			"41 0D 1E",
			"41 0D 23",
			"41 0D 28",
			"41 0D 2D",
			"41 0D 32",
			"41 0D 37",
			"41 0D 3C",
			"41 0D 41",
			"41 0D 46",
			"41 0D 4B",
			"41 0D 50",
			"41 0D 55",
			"41 0D 5A",
			"41 0D 5F",
			"41 0D 64",
			"41 0D 69",
			"41 0D 6B",
			"41 0D 71",
			"41 0D 76",
			"41 0D 7A",
			"41 0D 7F",
			"41 0D 80",
			"41 0D 85",
			"41 0D 89",
		}
		return []string{
			lehms[lehm],
		}
	} else if strings.HasPrefix(subcmd, "31") { // Distance traveled since codes cleared
		return []string{
			"41 31 02 0C", // 524 km
		}
	} else if strings.HasPrefix(subcmd, "42") { // Control Module Voltage
		return []string{
			"41 42 33 90", // 13.2 volts
		}
	} else if strings.HasPrefix(subcmd, "A6") { // Odometer
		return []string{
			"41 A6 00 06 68 a0", // 42,000.00 km
		}
	} else if strings.HasPrefix(subcmd, "46") { // Ambient Air Temperature
		return []string{
			"41 46 3A", // 18.0 C
		}
	} else if strings.HasPrefix(subcmd, "A4") { // Transmission Actual Gear
		return []string{
			"41 A4 27 10 00 00", // 10.0:1
		}
	} else if strings.HasPrefix(subcmd, "11") { // Throttle Position
		return []string{
			"41 11 50", // 50% throttle
		}
	}

	return []string{"NOT SUPPORTED"}
}

func mockOutputs(cmd string) []string {
	if cmd == "ATSP0" {
		return []string{"OK"}
	} else if cmd == "AT@1" {
		return []string{"OBDII by elm329@gmail.com"}
	} else if cmd == "AT RV" {
		return []string{"12.1234"}
	} else if strings.HasPrefix(cmd, "01") {
		return mockMode1Outputs(cmd[2:])
	}

	return []string{"NOT SUPPORTED"}
}
