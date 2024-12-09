package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

func validateOctet(part int) error {
	if part >= 0 && part <= 255 {
		return nil
	}

	return fmt.Errorf("invalid octet %q", part)
}

type IpAddress struct {
	octets []int
	subnet int
}

func NewIP(ip string) (*IpAddress, error) {
	ipStruct := new(IpAddress)

	partStrs := strings.Split(ip, ".")

	if len(partStrs) != 4 {
		return nil, fmt.Errorf("invalid IP: %q", ip)
	}

	lastPart := strings.Split(partStrs[3], "/")

	if len(lastPart) != 2 {
		return nil, fmt.Errorf("no subnet found in given ip")
	}

	subnetStr := lastPart[1]

	subnet, err := strconv.Atoi(subnetStr)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet part: %q", subnetStr)
	}

	if subnet < 0 || subnet > 32 {
		return nil, fmt.Errorf("subnet part must be within 32 for five octets")
	}

	ipStruct.subnet = subnet

	for _, partStr := range append(append([]string{}, partStrs[:3]...), lastPart[0]) {
		var part int
		var err error

		part, err = strconv.Atoi(partStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ip part: %q", partStr)
		}

		if err := validateOctet(part); err != nil {
			return nil, fmt.Errorf("ip octet must be within 255 range, found: %q", partStr)
		}

		ipStruct.octets = append(ipStruct.octets, part)
	}

	return ipStruct, nil
}

func (ip *IpAddress) String() string {
	var strs [4]string
	for i, octet := range ip.octets {
		strs[i] = strconv.Itoa(octet)
	}

	freeStr := strings.Join(strs[:], ".")

	if ip.subnet == -1 {
		return freeStr
	}

	return fmt.Sprintf("%s/%d", freeStr, ip.subnet)
}

func (ip *IpAddress) HostMin() string {
	newIp := new(IpAddress)

	// skip processing of these octets
	skip := ip.subnet / 8

	for i := 0; i < skip; i++ {
		newIp.octets = append(newIp.octets, ip.octets[i])
	}

	if skip == 4 {
		return newIp.String()
	}

	bitsToProcess := 32 - ip.subnet // these change

	incompleteOctet := bitsToProcess % 8 // first octet that needs partial processing

	fullResetOctets := bitsToProcess / 8 // these get all reset to 255 or 0

	currIndex := skip

	if incompleteOctet != 0 {
		binary := toBinary(ip.octets[currIndex])

		skipBits := 8 - incompleteOctet

		// reset rest of the bits
		for i := 0; i < 8-skipBits; i++ { // resets in reverse order since we store in reverse for conversion between types
			binary[i] = 0
		}

		newIp.octets = append(newIp.octets, toDecimal(binary))
	}

	for i := 0; i < fullResetOctets; i++ {
		newIp.octets = append(newIp.octets, 0)
	}

	if newIp.octets[3] == 0 {
		newIp.octets[3] += 1
	}
	newIp.subnet = -1

	return newIp.String()
}

func toBinary(num int) [8]int {
	bits := [8]int{}

	remainder := 0

	// intentional reverse ordering, helps re-encode to decimal later
	for i := 0; i < 8; i++ {
		remainder = num % 2
		num /= 2

		bits[i] = remainder
	}

	return bits
}

func toDecimal(bits [8]int) int {
	var num int = bits[0] * 1

	for i := 1; i < 8; i++ {
		if bits[i] == 1 {
			num += 2 << (i - 1)
		}
	}

	return num
}

func (ip *IpAddress) HostMax() string {
	newIp := new(IpAddress)

	// skip processing of these octets
	skip := ip.subnet / 8

	for i := 0; i < skip; i++ {
		newIp.octets = append(newIp.octets, ip.octets[i])
	}

	if skip == 4 {
		return newIp.String()
	}

	bitsToProcess := 32 - ip.subnet // these change

	incompleteOctet := bitsToProcess % 8 // first octet that needs partial processing

	fullResetOctets := bitsToProcess / 8 // these get all reset to 255 or 0

	currIndex := skip

	if incompleteOctet != 0 {
		binary := toBinary(ip.octets[currIndex])

		skipBits := 8 - incompleteOctet

		// reset rest of the bits
		for i := 0; i < 8-skipBits; i++ { // resets in reverse order since we store in reverse for conversion between types
			binary[i] = 1
		}

		newIp.octets = append(newIp.octets, toDecimal(binary))
	}

	for i := 0; i < fullResetOctets; i++ {
		newIp.octets = append(newIp.octets, 1)
	}

	if newIp.octets[3] == 255 {
		newIp.octets[3] -= 1
	}

	newIp.subnet = -1

	return newIp.String()
}
