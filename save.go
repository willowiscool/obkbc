package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/clbanning/mxj"
)

type Keybinding struct {
	key, command string
}

// readCombs reads the key combinations openbox has set up
// After using readRC to get the raw data in mxj.Map format, it will translate it into an array of Keybindings. It will ignore the keybindings that don't execute a command.
func readCombs() ([]Keybinding, error) {
	data, err := readRC()
	if err != nil {
		return nil, err
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
	return commandKeybinds, nil
}

// addComb adds a combination to the config and reloads openbox to use it
// After using readRC to get the raw data, it translates the keybinding given to it into a map[string]interface{} which it inserts into the data and then writes back using writeRC.
func addComb(comb Keybinding) error {
	data, err := readRC()
	if err != nil {
		return err
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
	return writeRC(data)
}

// deleteComb deletes a combination and reloads openbox to reflect that
// After reading the data with readRC, it finds the combinatoin it needs, deletes it, and uses writeRC to write the data back into the file.
//TODO: FIX THIS MESS
func deleteComb(comb Keybinding) error {
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
		//TODO: Fix this messs
		switch inspect["action"].(type) {
			case map[string]interface{}:
				action := inspect["action"].(map[string]interface{})
				if inspect["-key"].(string) == comb.key && action["command"].(string) == comb.command {
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
//TODO: FIX THIS MESS
func editComb(oldComb, comb Keybinding) error {
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
		switch inspect["action"].(type) {
			case map[string]interface{}:
				action := inspect["action"].(map[string]interface{})
				if inspect["-key"].(string) == oldComb.key && action["command"].(string) == oldComb.command {
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
