package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"log"
	"strings"
	"github.com/clbanning/mxj"
)

type Keybinding struct {
	key, command string
}

func readCombs() []Keybinding {
	data, err := readRC()
	if err != nil {
		log.Fatal("Problem reading rc.xml: " + err.Error())
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

func addComb(comb Keybinding) {
	data, err := readRC()
	if err != nil {
		log.Fatal("Problem reading rc.xml: " + err.Error())
	}
	translated := map[string]interface{}{
		"action": map[string]string{
			"-name": "Execute",
			"command": comb.command,
		},
		"-key": comb.key,
	}
	//Thanks to Kale Blankenship from the golang slack for this solution!!!!!!!!
	var (
		config   = data["openbox_config"].(map[string]interface{})
		keyboard = config["keyboard"].(map[string]interface{})
		keybind  = keyboard["keybind"].([]interface{})
	)
	keyboard["keybind"] = append(keybind, translated)
	writeRC(data)
}
//TODO: FIX THIS MESS
func deleteComb(comb Keybinding) {
	data, err := readRC()
	if err != nil {
		log.Fatal("Problem reading rc.xml: " + err.Error())
	}
	var (
		config   = data["openbox_config"].(map[string]interface{})
		keyboard = config["keyboard"].(map[string]interface{})
		keybind  = keyboard["keybind"].([]interface{})
	)
	for i, inspect := range keybind {
		//TODO: Fix this messs
		switch inspect.(map[string]interface{})["action"].(type) {
			case map[string]interface{}:
				if inspect.(map[string]interface{})["-key"].(string) == comb.key && inspect.(map[string]interface{})["action"].(map[string]interface{})["command"].(string) == comb.command {
					keyboard["keybind"] = append(keybind[:i], keybind[i+1:]...)
				}
			default:
				continue
		}
	}
	writeRC(data)
}
//TODO: FIX THIS MESS
func editComb(oldComb, comb Keybinding) {
	data, err := readRC()
	if err != nil {
		log.Fatal("Problem reading rc.xml: " + err.Error())
	}
	keybinds := data[
		"openbox_config"].(map[string]interface{})[
		"keyboard"].(map[string]interface{})[
		"keybind"].([]interface{})
	for i, inspect := range keybinds {
		//TODO: Fix this messs
		switch inspect.(map[string]interface{})["action"].(type) {
			case map[string]interface{}:
				if inspect.(map[string]interface{})["-key"].(string) == oldComb.key && inspect.(map[string]interface{})["action"].(map[string]interface{})["command"].(string) == oldComb.command {
					keybinds[i] = map[string]interface{}{
						"action": map[string]string{
							"-name": "Execute",
							"command": comb.command,
						},
						"-key": comb.key,
					}
				}
			default:
				continue
		}
	}
	writeRC(data)
}

func readRC() (mxj.Map, error) {
	rc, err := ioutil.ReadFile(os.Getenv("HOME") + "/.config/openbox/rc.xml")
	if err != nil {
		return nil, err
	}
	data, err := mxj.NewMapXml(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeRC(data mxj.Map) {
	result, err := data.XmlIndent("", "	")
	if err != nil {
		log.Fatal("Problem re-writing XML: " + err.Error())
	}
	result = []byte(strings.Replace(string(result), "&", "&amp;", -1))
	err = ioutil.WriteFile(os.Getenv("HOME") + "/.config/openbox/rc.xml", result, 0644)
	if err != nil {
		log.Fatal("Problem re-writing XML: " + err.Error())
	}
	err = exec.Command("openbox", "--reconfigure").Run()
	if err != nil {
		log.Fatal("Problem reconfiguring OpenBox: " + err.Error())
	}
}
