package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func main() {
	fmt.Printf(getDefaultConnectionSettingsValue())
}

func getDefaultConnectionSettingsValue() string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings\Connections`, registry.ALL_ACCESS)
	defer key.Close()
	if err != nil {
		panic(err)
	}

	s, _, _ := key.GetBinaryValue(`DefaultConnectionSettings`)
	d := ""
	for _, x := range s {
		d = d + fmt.Sprintf("%02x", x)
	}
	return d
}
