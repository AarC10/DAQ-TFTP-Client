/**
TFTP client that reads and writes config files.

@author Aaron Chan, RIT Launch Initiative's Future Star Programmer 😎
*/

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
)

// Struct representing config file entries
type config struct {
	srcIP    *widget.Entry
	dstIP    *widget.Entry
	gwIP     *widget.Entry
	subnetIP *widget.Entry

	srcUDP  *widget.Entry
	adc0UDP *widget.Entry
	adc1UDP *widget.Entry
	tcUDP   *widget.Entry

	adc0Rate *widget.Select
	adc1Rate *widget.Select

	resetCheck *widget.Check
}

// Struct representing user settings and extra widgets.
type extras struct { // TODO: Brainstorm a better struct name
	broadcastAddr *widget.Entry
	guiResponses  *widget.Entry
	loadingBar    *widget.ProgressBar
	inputResponse *widget.Label
}

const TFTP_PORT = ":69"

/**
Error Handling
*/
func check(e error, messageLabel **widget.Label) {
	if e != nil {
		(*messageLabel).SetText(e.Error())
		log.Println(e)
	}
}

/**
Checks if selection box has a selection and return a string representation
*/
func adcRateSelected(adcRate *widget.Select) string {
	out := ""

	switch adcRate.Selected {
	case "Slow (8 kHz)":
		out = "slow"

	case "Fast (43 kHz)":
		out = "fast"
	}

	return out
}

// Main method
func main() {

	// Initialize GUI
	gui := app.New()

	window := gui.NewWindow("RIT Launch Initiative TFTP Client")
	window.Resize(fyne.NewSize(1920, 1080))

	extras := extras{
		makeEntryField("Current DAQ IP", "ip"),
		widget.NewEntry(),
		widget.NewProgressBar(),
		widget.NewLabel(""),
	}

	// Extras modifiers
	extras.broadcastAddr.Validator = daqIPAddrValidator
	extras.inputResponse.Alignment = fyne.TextAlignCenter

	// Entries for config file
	configData := config{
		srcIP:    makeEntryField("Source IP", "ip"),
		dstIP:    makeEntryField("Destination IP", "ip"),
		gwIP:     makeEntryField("Gateway IP", "ip"),
		subnetIP: makeEntryField("Subnet IP", "ip"),

		srcUDP:  makeEntryField("Source UDP Port", "port"),
		adc0UDP: makeEntryField("ADC0 UDP Port", "port"),
		adc1UDP: makeEntryField("ADC1 UDP Port", "port"),
		tcUDP:   makeEntryField("TC Port", "port"),

		adc0Rate: makeNewSelection("ADC0 Rate", []string{"Slow (8 kHz)", "Fast (43 kHz)"}),
		adc1Rate: makeNewSelection("ADC1 Rate", []string{"Slow (8 kHz)", "Fast (43 kHz)"}),

		resetCheck: widget.NewCheck("Reset", func(b bool) {}),
	}

	// Instructions message
	instructionsOne := makeNewText("The first entry should be a valid IP address, otherwise the program will not attempt to send/receive files")
	instructionsTwo := makeNewText("All other entries, if invalid will be marked red and will not be written to the config file")
	instructionsThree := makeNewText("The write file button will write all valid entries as commands into a config file")
	instructionsFour := makeNewText("The upload button will send the config file that is currently in the application directory to the DAQ")
	instructionsFive := makeNewText("The receive button will overwrite the config file in the application directory and input its fields into the application")

	// Set and Display Content
	window.SetContent(
		container.NewVBox(
			// Basic Settings
			extras.broadcastAddr,
			widget.NewSeparator(),

			configData.srcIP,
			configData.dstIP,

			container.NewHBox(
				widget.NewLabel("ADC0 Rate: "),
				configData.adc0Rate,

				widget.NewLabel("ADC1 Rate: "),
				configData.adc1Rate,

				configData.resetCheck,
			),

			// Advanced Options dropdown
			widget.NewAccordion(
				widget.NewAccordionItem(
					"Advanced Options",

					container.NewVBox(
						makeNewText("Do not modify unless you know what you are doing. Changing these options risks breaking the DAQ."),
						configData.gwIP,
						configData.subnetIP,
						configData.srcUDP,
						configData.adc0UDP,
						configData.adc1UDP,
						configData.tcUDP,
					)),
			),

			// Instructions dropdown
			widget.NewAccordion(
				widget.NewAccordionItem("Instructions",
					container.NewVBox(
						instructionsOne,
						instructionsTwo,
						instructionsThree,
						instructionsFour,
						instructionsFive,
					)),
			),

			layout.NewSpacer(),
			widget.NewSeparator(),

			// Buttons
			widget.NewButton("Create Config", func() {
				createConfig(&configData, &extras)
			}),

			widget.NewButton("Receive File", func() {
				receiveFile(&configData, &extras)
			}),

			widget.NewButton("Upload File", func() {
				uploadFile(&extras)
			}),

			// extras.loadingBar,

			// Response to user interactions
			widget.NewSeparator(),
			extras.inputResponse,
		),
	)

	window.ShowAndRun()
}
