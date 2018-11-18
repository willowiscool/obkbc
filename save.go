package main

import (
	"io/ioutil"
	"os"
	"log"
	"github.com/clbanning/mxj"
)

type Keybinding struct {
	key, command string
}

func readCombs() []Keybinding {
	rc, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/openbox/rc.xml")
	if err != nil {
		log.Fatal("Problem reading rc.xml: " + err.Error())
	}
	data, err := mxj.NewMapXml(rc)
	if err != nil {
		log.Fatal("Problem parsing rc.xml: " + err.Error())
	}
	keybindings := data[
		"openbox_config"].(map[string]interface{})[
		"keyboard"].(map[string]interface{})[
		"keybind"].([]interface{})
	var commandKeybinds []Keybinding
	for _, value := range keybindings {
		keybind := value.(map[string]interface{})
		var action map[string]interface{}
		switch keybind["action"].(type) {
			case map[string]interface{}:
				action = keybind["action"].(map[string]interface{})
			default:
				continue
		}
		if action["-name"].(string) == "Execute" {
			commandKeybinds = append(commandKeybinds, Keybinding{
				key: keybind["-key"].(string),
				command: action["command"].(string),
			})
		}
	}
	return commandKeybinds
}
