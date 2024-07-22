package timezones

import (
	"os"
	"strings"
	"unicode"
)

var zoneDirs = []string{
	"/usr/share/zoneinfo/",
}

var timezones = []string{}

func walkTzDir(path string, zones []string) []string {
	dirInfo, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	isRuneNotLetter := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}
	for _, info := range dirInfo {
		if info.Name() != strings.ToUpper(info.Name()[:1])+info.Name()[1:] {
			continue
		}
		_, err := os.Readlink(info.Name())
		if err == nil {
			continue
		}
		if !isRuneNotLetter(info.Name()[:1]) {
			continue
		}
		newPath := path + "/" + info.Name()
		if info.IsDir() {
			zones = walkTzDir(newPath, zones)
		} else {
			zones = append(zones, newPath)
		}
	}
	return zones
}

func List() []string {
	if len(timezones) > 0 {
		return timezones
	}
	for _, zd := range zoneDirs {
		timezones = walkTzDir(zd, timezones)
		for idx, zone := range timezones {
			timezones[idx] = strings.ReplaceAll(zone, zd+"/", "")
		}
	}
	return timezones
}

func IsAvailable(name string) bool {
	l := List()
	for _, n := range l {
		if name == n {
			return true
		}
	}
	return false
}
