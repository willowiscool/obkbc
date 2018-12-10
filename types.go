package main

type Command interface {
	Translate() map[string]interface{}
	Key() string
}

type Kcommand struct {
	key, command string
}
func (comb Kcommand) Translate() map[string]interface{} {
	return map[string]interface{}{
		"action": map[string]string{
			"-name": "Execute",
			"command": comb.command,
		},
		"-key": comb.key,
	}
}
func (comb Kcommand) Key() string {
	return comb.key
}

type Kmrt struct {
	key, height, width, x, y string
}
func (comb Kmrt) Translate() map[string]interface{} {
	return map[string]interface{}{
		"action": []map[string]interface{}{
			{
				"-name": "Unmaximize",
			},
			{
				"-name": "MoveResizeTo",
				"height": comb.height,
				"width": comb.width,
				"x": comb.x,
				"y": comb.y,
			},
		},
		"-key": comb.key,
	}
}
func (comb Kmrt) Key() string {
	return comb.key
}
