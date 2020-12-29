package cliutils

import (
	"net"
	"strings"
	"regexp"
	"strconv"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

const (
	DNSName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
)

var	rxDNSName = regexp.MustCompile(DNSName)

// IsCIDR check if the string is an valid CIDR notiation (IPV4 & IPV6)
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

// IsIP checks if a string is either IP version 4 or 6. Alias for `net.ParseIP`
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

// IsDNSName will validate the given string as a DNS name
func IsDNS(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !IsIP(str) && rxDNSName.MatchString(str)
}

func ParseDNS(domain string) (string, error) {
	addr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		return "", err
	}

	return addr.String(), nil
}


// cidr解析
func ParseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func FilterPort(str string) (int, error) {
	port, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("port out of range")
	}
	return port, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

//数组去重
func ArrayUnique(arr []int) []int {
	var newArr []int
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func copyFile(source, dist string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
	   return err
	}

	err = ioutil.WriteFile(dist, input, 0644)
	if err != nil {
	   return err
	}

	return nil
}

// Dir copies a whole directory recursively
func copyDir(src string, dst string) error {
    var err error
    var fds []os.FileInfo
    var srcinfo os.FileInfo

    if srcinfo, err = os.Stat(src); err != nil {
        return err
    }

    if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
        return err
    }

    if fds, err = ioutil.ReadDir(src); err != nil {
        return err
    }
    for _, fd := range fds {
        srcfp := path.Join(src, fd.Name())
        dstfp := path.Join(dst, fd.Name())

        if fd.IsDir() {
            if err = copyDir(srcfp, dstfp); err != nil {
                return err
            }
        } else {
            if err = copyFile(srcfp, dstfp); err != nil {
                return err
            }
        }
    }
    return nil
}
