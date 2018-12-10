package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/clbanning/mxj"
)

// readCombs reads the key combinations openbox has set up
// After using readRC to get the raw data in mxj.Map format, it will translate it into an array of Kcommands. It will ignore the keybindings that don't execute a command.
func readCombs() ([]Kcommand, []Kmrt, error) {
	data, err := readRC()
	if err != nil {
		return nil, nil, err
	}
	keybindings := data[
		"openbox_config"].(map[string]interface{})[
		"keyboard"].(map[string]interface{})[
		"keybind"].([]interface{})
	var commandKeybinds []Kcommand
	var mrtKeybinds []Kmrt
	for _, value := range keybindings {
		keybind := value.(map[string]interface{})
		switch keybind["action"].(type) {
			case map[string]interface{}:
				if keybind["action"].(map[string]interface{})["-name"] == "Execute" {
					commandKeybinds = append(commandKeybinds, Kcommand{
						key: keybind["-key"].(string),
						command: keybind["action"].(map[string]interface{})["command"].(string),
					})
				}
			case []interface{}:
				actions := keybind["action"].([]interface{})
				for _, action := range actions {
					if action.(map[string]interface{})["-name"].(string) == "MoveResizeTo" {
						mrtAction := action.(map[string]interface{})
						translated := Kmrt{key: keybind["-key"].(string)}
						if mrtAction["height"].(string) != "" {
							translated.height = mrtAction["height"].(string)
						}
						if mrtAction["width"].(string) != "" {
							translated.width = mrtAction["width"].(string)
						}
						if mrtAction["x"].(string) != "" {
							translated.x = mrtAction["x"].(string)
						}
						if mrtAction["y"].(string) != "" {
							translated.y = mrtAction["y"].(string)
						}
						mrtKeybinds = append(mrtKeybinds, translated)
					}
				}
			default:
				continue
		}
	}
	return commandKeybinds, mrtKeybinds, nil
}

// addComb adds a combination to the config and reloads openbox to use it
// After using readRC to get the raw data, it translates the keybinding given to it into a map[string]interface{} which it inserts into the data and then writes back using writeRC.
func addComb(comb Command) error {
	data, err := readRC()
	if err != nil {
		return err
	}
	var (
		config   = data["openbox_config"].(map[string]interface{})
		keyboard = config["keyboard"].(map[string]interface{})
		keybind  = keyboard["keybind"].([]interface{})
	)
	translated := comb.Translate()
	keyboard["keybind"] = append(keybind, translated)
	return writeRC(data)
}

// deleteComb deletes a combination and reloads openbox to reflect that
// After reading the data with readRC, it finds the combinatoin it needs, deletes it, and uses writeRC to write the data back into the file.
func deleteComb(key string) error {
	data, err := readRC()
	if err != nil {
		return err
	}
	var (
		config   = data["openbox_config"].(map[string]interface{})
		keyboard = config["keyboard"].(map[string]interface{})
		keybind  = keyboard["keybind"].([]interface{})
	)
	for i, inspectRAW := range keybind {
		inspect := inspectRAW.(map[string]interface{})
		switch inspect["action"].(type) {
			case map[string]interface{}:
				if inspect["-key"].(string) == key {
					keyboard["keybind"] = append(keybind[:i], keybind[i+1:]...)
				}
			default:
				continue
		}
	}
	return writeRC(data)
}

// editComb edits a keyboard combination and reloads openbox to relfect that
// After reading the data with readRC, it changes the keyboard combination and uses writeRC to insert the edit into the file and reload openbox.
func editComb(oldComb, comb Command) error {
	data, err := readRC()
	if err != nil {
		return err
	}
	keybinds := data[
		"openbox_config"].(map[string]interface{})[
		"keyboard"].(map[string]interface{})[
		"keybind"].([]interface{})
	for i, inspectRAW := range keybinds {
		inspect := inspectRAW.(map[string]interface{})
		if inspect["-key"].(string) == oldComb.Key() {
			keybinds[i] = comb.Translate()
		}
	}
	return writeRC(data)
}

// readRC returns a mxj.Map containing the data found in $HOME/.config/openbox/rc.xml
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

// writeRC takes a mxj.Map, converts it to XML, and then writes it into $HOME/.config/openbox/rc.xml
func writeRC(data mxj.Map) error {
	result, err := data.XmlIndent("", "	")
	if err != nil {
		return err
	}
	result = []byte(strings.Replace(string(result), "&", "&amp;", -1))
	err = ioutil.WriteFile(os.Getenv("HOME") + "/.config/openbox/rc.xml", result, 0644)
	if err != nil {
		return err
	}
	err = exec.Command("openbox", "--reconfigure").Run()
	if err != nil {
		return err
	}
	return nil
}
