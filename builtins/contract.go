package builtins

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	ABI     *abi.ABI
	Address common.Address
}

func (c *Contract) Load(thor *thorgo.Thor) *accounts.Contract {
	return thor.Account(c.Address).Contract(c.ABI)
}

var (
	VTHO = &Contract{
		ABI:     mustParseABI(compiledEnergyAbi, "VTHO"),
		Address: common.HexToAddress("0x0000000000000000000000000000456e65726779"),
	}
	Authority = &Contract{
		ABI:     mustParseABI(compiledAuthorityAbi, "Authority"),
		Address: common.HexToAddress("0x0000000000000000000000417574686f72697479"),
	}
	Executor = &Contract{
		ABI:     mustParseABI(compiledExecutorAbi, "Executor"),
		Address: common.HexToAddress("0x0000000000000000000000004578656375746f72"),
	}
	Extension = &Contract{
		ABI:     mustParseABI(compiledExtensionv2Abi, "Extension"),
		Address: common.HexToAddress("0x0000000000000000000000457874656e73696f6e"),
	}
	Prototype = &Contract{
		ABI:     mustParseABI(compiledPrototypeAbi, "Prototype"),
		Address: common.HexToAddress("0x000000000000000000000050726f746f74797065"),
	}
	Params = &Contract{
		ABI:     mustParseABI(compiledParamsAbi, "Params"),
		Address: common.HexToAddress("0x0000000000000000000000000000506172616d73"),
	}
)

func mustParseABI(data []byte, name string) *abi.ABI {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("read %q: %v", name, err))
	}

	contractABI, err := abi.JSON(gz)
	if err != nil {
		panic(fmt.Errorf("parse %q: %v", name, err))
	}

	return &contractABI
}

var compiledEnergyAbi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\xcf\xaf\xd3\x30\x0c\xfe\x5f\x7c\xce\x09\x09\x84\x7a\x83\x03\x37\xc4\x01\x6e\x4f\x3b\xb8\xad\x8b\x22\x25\x76\x94\x38\x1b\xd5\xd3\xfb\xdf\xd1\xb6\xac\xab\xa0\xbf\x98\x86\xb6\x53\x2b\xd9\x8e\xbf\xcf\xf6\x97\xf8\xe5\x15\x1a\xe1\xa4\xc8\x0a\x95\xc6\x4c\x06\x2c\x87\xac\x09\xaa\x97\x9d\x01\x46\x4f\x50\x9d\x3f\x06\x24\x6b\x31\xbd\x5e\x2c\x60\x40\xfb\x70\xfc\x4b\x1a\x2d\xff\x84\xb7\x9d\x81\x80\x3d\xd6\x8e\xa0\xea\xd0\x25\x32\x90\x14\x95\xbe\x66\xc5\xda\x3a\xab\x3d\x54\x10\x72\xa4\x6b\x68\x97\xb9\x51\x2b\x0c\x6f\x66\x0c\xa7\x44\x0f\x78\x86\xa4\x29\x10\xb7\x14\xaf\x07\x60\xdb\x46\x4a\xe9\x14\x5f\x7c\xf6\xe8\xf2\x28\x45\xb6\xac\xef\xde\x7f\x38\xc1\x2b\x1e\x18\x42\x94\xfd\x0c\xaf\x94\x9b\xe6\x78\xe2\x70\x40\x2d\xe2\x36\x92\x63\xe1\x8b\xd3\x1a\xc5\xd9\x8a\xab\x28\xba\xef\x39\x04\xd7\xaf\x15\x7e\x4c\x6d\x1d\xdc\xde\xd2\xe1\xf6\xca\x77\x51\xfc\x62\xd9\x55\x16\xcd\xe8\x25\xb3\x2e\xb6\x45\x23\x72\xea\x28\x7e\x39\xa7\x7a\xc2\xde\xb4\xd4\x58\x8f\x2e\x6d\x69\xcc\xc7\x7b\x0a\xe2\x0f\x44\x43\x4a\x39\xf0\xa4\x1a\xae\x90\x6b\x74\xc8\x0d\x7d\xeb\xa6\x31\x17\xf3\x7f\x9d\xa9\xd9\x72\xa6\xde\xd7\xe2\x9e\xe8\x7a\xb9\xdf\x0c\x3f\x74\x7e\x1f\x29\x62\xff\xdc\x17\xeb\xe7\x1c\x99\xda\x07\x5c\xac\xff\xaa\x60\xb3\xe1\xcd\x1b\xbd\x68\xce\xc9\xa1\x08\x79\x82\x59\x24\x8f\x96\x8f\x2a\xba\x3f\x45\x64\xe1\xde\x4b\x4e\x53\x73\x67\xb9\xa5\x5f\xd4\x5e\xe8\xaf\x8e\xe1\xb4\xff\xdc\x54\x0e\xde\x25\xf1\xe6\x05\xe0\xc7\x55\xa5\xc5\x89\xf6\xc4\x7a\x2b\x9f\x85\x1e\x4e\x07\x2c\xae\x31\x37\xb3\xfa\x74\x5a\x6b\xd0\xfd\xc5\x6a\xf7\x3b\x00\x00\xff\xff\x3d\x94\x5b\x7e\xec\x09\x00\x00")

var compiledAuthorityAbi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x93\xcd\x6e\xc2\x40\x0c\x84\xdf\xc5\xe7\x3d\xb5\xb7\x5c\x7b\xe6\x09\x10\x07\x13\x0f\xd5\xaa\xc1\x46\xb1\x37\x25\x42\xbc\x7b\x05\x0a\x24\x45\xb4\xa5\x6a\x51\xd5\x5b\x24\xff\xcc\x37\x93\xf5\x7c\x47\xb5\xa9\x07\x6b\x50\x15\x6d\x41\xa2\xac\x9b\x12\x4e\xd5\x7c\x91\x48\x79\x0d\xaa\x68\x95\x5b\x0f\x4a\x64\x25\x86\xda\xee\x54\xa2\x44\xd1\x6f\x0e\x5f\x2c\xd2\xc2\x9d\xf6\x8b\x44\x1b\xee\x79\xd9\x80\xaa\x15\x37\x8e\x44\x1e\x1c\x98\x95\xe0\x65\x6e\x72\xf4\x54\x51\x97\xf1\x3a\xce\xae\x8a\xd6\x91\x4d\x69\x9f\xa6\x40\xc3\xf4\x99\xe8\xac\xaa\x26\x98\xb1\x07\xda\xeb\xfa\x43\x5b\x8b\xce\x5e\xf0\x0e\xfc\x26\x38\x35\x3d\x35\x7d\x85\x78\x91\xd9\x77\x09\x15\xdb\xbf\x08\xf6\x87\xd4\xcf\xf8\x00\xba\xc9\x1e\x90\x71\x76\x69\xd6\x1c\xa5\x87\x3a\x54\xac\x75\xbb\xb6\x7d\x6c\xca\x02\x8d\x83\x99\x71\x4d\x1f\xf0\xc7\x87\x69\x13\xd7\x91\x3b\x5c\x2a\xdd\x2d\xa1\xd1\x3b\xb6\xa8\x4b\x1c\x2d\xfc\x8f\x73\xf8\xbd\xf4\xc7\x10\x58\xe4\x0e\x57\xc5\x6a\xda\xaf\xad\xf8\x35\xab\x59\x05\x5b\xc8\xe9\xcf\xdc\xe8\xfc\x3c\x35\x2c\x9c\xbc\x1e\xd3\x4f\x2d\x3e\xb1\x4a\x16\x8e\x09\x2c\x3a\x68\xd0\x7e\xf1\x16\x00\x00\xff\xff\x5e\x1e\xdc\x95\x35\x05\x00\x00")

var compiledExecutorAbi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x55\x3d\x6f\xdb\x40\x0c\xfd\x2f\x9c\x35\xb5\x45\x51\x68\x33\xd2\xa5\x43\x01\x4f\x5d\x82\xc0\xa0\x75\x4c\x7a\xa8\x4c\xaa\x77\x3c\x25\x42\x90\xff\x5e\xc4\xd2\xf5\x24\xeb\xa3\x6e\x6a\x58\xe8\x66\xc0\x8f\x4f\x8f\x8f\x8f\xbc\xdb\x67\x28\x84\xbd\x22\x2b\xe4\xea\x02\x65\x60\xb9\x0a\xea\x21\xbf\x7d\x06\xc6\x03\x41\x0e\x90\x81\x36\xd5\xeb\x2f\x34\xc6\x91\xf7\xf0\x72\x97\xc5\x3f\xb1\xaa\x9c\xd4\xe4\x3c\x64\x20\x41\x4f\x6b\xad\x21\x56\xab\x4d\xe2\xd8\x37\x4a\xfe\xfd\x3b\x78\xc9\x12\x88\xb7\xf2\x48\xae\x87\x11\x29\x8f\x1f\xa9\xb0\xc1\x7d\x49\x90\xdf\x63\xe9\x29\x03\xaf\xa8\xf4\x35\x28\xee\x6d\xf9\xca\x9a\x43\x6d\xe9\x31\x15\xde\x07\x2e\xd4\x0a\x1f\xd9\x67\x3b\x1b\xab\xbf\x91\xc0\x3a\xdd\x41\x22\x0f\x96\xf5\xd3\x25\x65\x75\xd5\x63\xc7\x77\x51\xd7\xa2\xf5\x8e\x6a\xf9\x41\x9b\x04\x4d\xea\xcf\xd2\xc8\xc2\x11\xf4\x97\x06\x4e\x98\xf3\x7b\xac\x49\x5f\xe5\xa4\x12\x8f\xe5\x4c\x34\xd4\x1e\x68\x7b\x84\x90\x19\x9a\xfc\xf1\x43\x3f\x1d\x2d\xcd\xa4\x17\x09\xf4\x33\x88\x0b\x87\xd1\xac\x12\xa0\x75\x14\xcb\x38\xe9\x39\x1c\x3d\x51\x11\xb4\xaf\xa8\x0d\x63\x42\x28\xba\x07\xd2\x45\x35\x06\x15\x4f\xac\x59\x3d\x37\x49\xdd\x6e\x61\x29\x7b\xab\x61\xcc\x55\xa2\x35\xdf\xcc\x19\x46\xef\xe6\x9c\x1e\x66\xe7\x4f\x9b\xdd\x6f\xff\x4a\xcd\x15\xc2\xea\xb0\x98\x6a\xaf\x37\x03\x55\x2c\xbe\x7f\x13\xb5\xfc\x70\x93\x0a\xd6\x18\x46\x5c\xe6\x2f\x9f\x97\x53\xd3\x46\x66\x25\x91\x67\x99\x6a\xe8\x7f\x33\xb5\xbb\x49\x2b\x5f\xf8\x09\x2b\xeb\x81\x89\x33\x77\xfe\xd2\xcf\x3a\xb2\x70\x73\x90\xe0\xa7\x9c\xb5\x6c\xe8\x89\x4c\xec\xe7\xe4\x25\x9a\xf4\x39\xeb\x55\x75\x84\x31\xcc\xed\x47\x97\x46\xb3\xed\x98\x13\x88\x6a\x62\x7d\xab\xd0\xe5\xfb\xfd\x76\x99\x9b\x11\xef\x3f\xc9\x8c\x7b\xb6\x31\xe6\xe2\x52\x47\x8b\x39\x14\x7c\xf7\x2b\x00\x00\xff\xff\xa9\x52\x06\x45\xb3\x0a\x00\x00")

var compiledExtensionv2Abi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\xd3\x41\x4b\xc3\x30\x14\x07\xf0\xef\xf2\xce\x39\x55\x1d\xd2\xa3\x28\xe2\x41\x18\x6e\xe0\x61\xf4\xf0\xd2\xbe\x4a\x68\x9a\x84\xe4\x65\x2e\x8c\x7d\x77\x69\x99\x5b\x0f\xe2\x2a\xa2\xed\xad\xd0\xfc\x79\x3f\xfe\xbc\xb7\xd9\x43\x69\x4d\x60\x34\x0c\x39\xfb\x48\x02\x94\x71\x91\x03\xe4\x9b\x42\x80\xc1\x96\x20\x07\xb6\x8c\x7a\x15\x9d\xd3\x09\x04\xd8\xc8\xc7\x17\xfb\xcf\x07\x20\x80\x93\xeb\xbe\xa2\x32\x9c\xdd\x2c\xe0\x50\x08\x70\x98\x50\x6a\x82\xbc\x46\x1d\x48\x40\x60\x64\x7a\x8e\x8c\x52\x69\xc5\x09\x72\xd8\x2a\x7a\x3f\x67\xeb\x68\x4a\x56\xd6\xc0\x41\x7c\xc3\x3a\x0d\xad\x90\xf1\x1c\x96\x89\x29\xf4\x63\x8f\xbf\xa5\xc6\x86\x32\xd9\x61\x2e\x98\xfb\xe8\x55\xf6\x1f\x66\x13\xdb\xaf\xbb\x3a\xa1\x6d\xd9\xac\x55\x4b\x33\xea\x79\x94\x79\xa5\xde\x0c\xf9\x4b\x6a\xac\x2a\x4f\x21\xcc\x46\xbd\xee\x17\xbb\xb4\x7e\x54\xdf\x8b\xeb\x3f\x84\x0f\xae\x6d\xf7\x88\x61\x89\x69\x0e\x75\x0e\x55\x0f\x3b\xa7\x3c\xf6\x99\xc9\x97\x73\xe8\x7a\xba\x9f\xfe\xc0\x87\x9e\xa5\xb7\x5b\xaa\x5e\xad\x6f\xa6\xef\xe9\x67\xe7\x30\xb7\x26\xef\x3a\xd4\x0b\xd5\xa3\x54\xb7\xbf\x44\x15\x1f\x01\x00\x00\xff\xff\x30\xf3\x42\x68\x0c\x07\x00\x00")

var compiledParamsAbi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x92\xc1\x6a\xc3\x30\x0c\x86\xdf\xe5\x3f\xfb\xd4\xb1\x1d\xf2\x0e\x3b\xed\x58\xc2\x50\x12\x75\x98\xa5\x72\x88\xa4\xac\xa6\xe4\xdd\xc7\x4a\x12\xc3\xe8\x08\x65\x47\xe3\xcf\xbf\x3e\xa3\xff\x78\x45\x9b\x44\x8d\xc4\x50\x9d\xa8\x57\x0e\x88\x32\xb8\x29\xaa\xe3\x15\x42\x67\x46\x85\xf7\x4f\xce\x08\xb0\x3c\xfc\x9c\x9a\x6c\xac\x4f\x07\xcc\xa1\x00\x13\xf5\xce\x05\xf1\x28\x76\x78\x7e\xc1\x5c\x87\x15\x51\x36\x04\x24\xb7\x25\xbc\x0e\x18\x28\x53\xd3\xf3\x36\x58\x8d\x8c\x5f\xdd\xa8\x89\x7d\xb4\x8c\x0a\x92\x64\x85\xb6\xec\x93\x4b\x6b\x31\xc9\x6d\x7e\x91\xb7\xd1\x1f\x71\x2f\x62\x1f\xbf\xc4\xb6\x97\xf7\xbf\xb3\x2f\x3d\x45\xfe\x7a\x54\xb7\xe8\xf0\x85\x5b\xb7\x34\xee\x39\x51\xd7\x8d\xac\xfa\x7f\x27\x92\x24\xf9\x9c\x5c\xef\x15\x20\x4a\xc7\x17\xee\x56\xdf\xc5\xe2\xcf\x3a\x6c\xf8\x92\xb4\xf0\xfb\xed\x78\xbb\x2d\x61\xb9\xe7\x89\xc5\x30\xd7\xdf\x01\x00\x00\xff\xff\xfb\x8f\x43\xc8\x9d\x02\x00\x00")

var compiledPrototypeAbi = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x56\x4d\x8f\xda\x30\x10\xfd\x2f\x3e\xe7\x44\xd5\x1e\x72\x45\xea\x8d\xaa\x6a\xd5\x13\x42\xd5\xc4\x19\xa8\x85\x33\x8e\x3c\x63\x50\x84\xfa\xdf\x57\xa0\x25\x31\x4b\x84\x77\xf9\xd8\x84\x23\xe2\x4d\x3c\xcf\xf3\xde\xf3\xcc\x77\x4a\x3b\x62\x01\x12\x95\x2f\xc1\x32\x66\xca\x50\x1d\x84\x55\x3e\xdf\x29\x82\x0a\x55\xae\xfe\x32\xda\xa5\xca\x94\x34\xf5\xfe\x27\x94\xa5\x47\x66\xf5\x3f\xeb\x10\x84\xdb\x19\xb0\xa0\xef\x81\x2d\xb2\x23\x8c\x51\x5a\x94\x0b\xf2\x7a\xce\x22\x53\x35\x34\x50\x58\x6c\x7b\x60\x01\xc1\x59\x10\x28\x8c\x35\xd2\xa8\x5c\x91\xa3\x23\xa8\x3d\x61\x19\x48\x8b\x71\x74\xe8\xa4\xe3\x21\x3e\x5c\x4b\x23\x70\x82\x81\xe1\x3f\xfc\xa6\xfd\xb6\xba\x2b\x2c\x9c\xb3\x87\xaa\x34\xaf\x8d\xc1\xed\x23\x19\xad\xb1\x89\xfa\x6a\x04\xf9\xcb\xe4\x64\x24\xe2\x3c\xac\xf0\xbb\x4b\x93\x8a\x8a\x87\xe7\x55\x58\xa7\xd7\x3f\x42\x55\xc4\x03\x0b\x86\x64\xf2\xf5\x5b\xcc\x0f\x09\xfd\xaa\x49\x71\x8b\x0b\xef\xc4\xed\x16\x37\x25\x65\xe8\xb1\x72\x1b\x3c\x93\xe2\xd0\x4e\xea\x1a\xd4\xc1\x7b\x24\xf9\x5d\x3b\xe2\xb4\xb4\xe2\x0f\x8c\xe0\xfa\xb5\xc7\xd2\x48\x8f\x3a\x22\x8c\x47\xed\x36\xe8\x9b\x5f\x20\x78\x51\x80\x8c\x32\x3d\x7c\xef\xa7\x05\x7a\xe8\xb4\x6e\xa1\xcc\xed\xa0\x2e\xa5\xb7\x45\xdd\x3b\xd2\xa1\x75\x77\x45\x2e\x14\x60\x81\x34\x3e\x43\x30\x44\x13\x18\xe1\xdd\x47\x9e\xef\x97\xf9\xee\xf4\xff\x8b\xb6\x4a\xbb\x6a\x04\xf1\x90\x4c\x67\x28\xcb\x11\x47\xf3\x3f\xe0\xa9\x2b\x93\xca\x1f\x7a\x87\xe9\x1a\xae\xce\x57\xc6\xcf\x79\x43\x1e\xb9\x48\xee\x01\xd3\xa3\x21\x9e\x28\x82\x02\x8d\x31\x84\x3e\xf8\x92\x19\x7e\xe7\x62\x72\x07\x13\x2c\x5e\x02\x00\x00\xff\xff\xc1\xa7\x82\x58\x64\x0d\x00\x00")
