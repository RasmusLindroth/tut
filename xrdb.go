package main

import (
	"os/exec"
	"strconv"
	"strings"
)

type xrdbField struct {
	Name  string
	Value string
}

type XrdbColors map[string]map[string]string

func getXrdb() ([]xrdbField, error) {
	var xf []xrdbField

	out, err := exec.Command("xrdb", "-query").CombinedOutput()
	if err != nil {
		return xf, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		xf = append(xf,
			xrdbField{
				Name: strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(
					strings.Join(parts[1:], ":"),
				),
			},
		)
	}

	return xf, nil
}

func GetXrdbColors() (XrdbColors, error) {
	colors := make(XrdbColors)

	xf, err := getXrdb()
	if err != nil {
		return colors, err
	}

	for _, field := range xf {
		parts := strings.Split(field.Name, ".")
		if len(parts) == 1 && strings.HasPrefix(parts[0], "*") {
			parts = append(parts, parts[0][1:])
			parts[0] = "*"
		}
		if len(parts) != 2 {
			continue
		}

		c := parts[1]
		if c != "background" && c != "foreground" && !strings.HasPrefix(c, "color") {
			continue
		}

		if strings.HasPrefix(c, "color") {
			num := strings.TrimPrefix(c, "color")
			i, err := strconv.Atoi(num)

			if err != nil || i < 0 || i > 15 {
				continue
			}
		}

		if _, ok := colors[parts[0]]; !ok {
			colors[parts[0]] = make(map[string]string)
		}
		colors[parts[0]][c] = field.Value
	}

	return colors, nil
}
