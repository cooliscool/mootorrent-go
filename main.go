package main

import (
	"errors"
	"fmt"
	s "strings"
)

// parse magnet
// http://www.bittorrent.org/beps/bep_0009.html

// for a real world implementation of magnet parser :
// https://github.com/webtorrent/magnet-uri/blob/master/index.js

func main() {
	magnet := `magnet:?xt=urn:btih:046B18E7D4AABB3A05A793C90D0D7A081AD29B4F&dn=Fear.City.New.York.vs.The.Mafia.S01.COMPLETE.1080p.WEB.H264-OATH%5BTGx%5D+%E2%AD%90&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce`

	magparsed, _ := parseMagnet(magnet)

	fmt.Printf("%v", magparsed)
}

type Magnet struct {
	values map[string][]string
	raw    string
}

func parseMagnet(magstr string) (Magnet, error) {
	// log := fmt.Println
	magstr = s.TrimSpace(magstr)
	// :::magnet URI format:::
	// all spl characters are expected to be urlencoded
	// v1: magnet:?xt=urn:btih:<info-hash>&dn=<name>&tr=<tracker-url>&x.pe=<peer-address>
	// v2: magnet:?xt=urn:btmh:<tagged-info-hash>&dn=<name>&tr=<tracker-url>&x.pe=<peer-address>

	magobj := Magnet{raw: magstr, values: make(map[string][]string)}

	magprefix := "magnet:?"
	if s.HasPrefix(magstr, magprefix) {
		magstr = s.TrimPrefix(magstr, magprefix)

		for _, item := range s.Split(magstr, "&") {
			key := s.Split(item, "=")[0]
			value := s.Split(item, "=")[1]
			magobj.values[key] = append(magobj.values[key], value)
		}
	} else {
		return Magnet{}, errors.New("Invalid magnet URI. No 'magnet:?' prefix")
	}

	// rule out invalid magnets
	// xt is the only mandatory parameter.
	if _, ok := magobj.values["xt"]; !ok {
		return Magnet{}, errors.New("Invalid magnet URI. No ''xt' param in URI")
	}

	// TODO: xt has either btih or btmh or both(hybrid)
	// TODO: if btih, <info-hash> is either 40 characters(hex-encoded) or 32 characters(base32 encoded)

	return magobj, nil
}
