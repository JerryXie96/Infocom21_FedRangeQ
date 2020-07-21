// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	PubKey     []byte
	Tag        [32]uint8
	Ciphertext [32][]byte
}

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	Label      *big.Int
	PubKey     [32][]byte
	Tag        [32][1][1]uint8
	Ciphertext [32][1][]byte
}

// FedRangeQABIABI is the input ABI used to generate the binding from.
const FedRangeQABIABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"clearResult\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"pubKey\",\"type\":\"bytes\"},{\"name\":\"tag\",\"type\":\"uint8[32]\"},{\"name\":\"ciphertext\",\"type\":\"bytes[32]\"}],\"name\":\"query\",\"type\":\"tuple\"}],\"name\":\"search\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"label\",\"type\":\"uint256\"},{\"name\":\"pubKey\",\"type\":\"bytes[32]\"},{\"name\":\"tag\",\"type\":\"uint8[1][1][32]\"},{\"name\":\"ciphertext\",\"type\":\"bytes[1][32]\"}],\"name\":\"indexToBeAdded\",\"type\":\"tuple[]\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getResult\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// FedRangeQABI is an auto generated Go binding around an Ethereum contract.
type FedRangeQABI struct {
	FedRangeQABICaller     // Read-only binding to the contract
	FedRangeQABITransactor // Write-only binding to the contract
	FedRangeQABIFilterer   // Log filterer for contract events
}

// FedRangeQABICaller is an auto generated read-only Go binding around an Ethereum contract.
type FedRangeQABICaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FedRangeQABITransactor is an auto generated write-only Go binding around an Ethereum contract.
type FedRangeQABITransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FedRangeQABIFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FedRangeQABIFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FedRangeQABISession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FedRangeQABISession struct {
	Contract     *FedRangeQABI     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FedRangeQABICallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FedRangeQABICallerSession struct {
	Contract *FedRangeQABICaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// FedRangeQABITransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FedRangeQABITransactorSession struct {
	Contract     *FedRangeQABITransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// FedRangeQABIRaw is an auto generated low-level Go binding around an Ethereum contract.
type FedRangeQABIRaw struct {
	Contract *FedRangeQABI // Generic contract binding to access the raw methods on
}

// FedRangeQABICallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FedRangeQABICallerRaw struct {
	Contract *FedRangeQABICaller // Generic read-only contract binding to access the raw methods on
}

// FedRangeQABITransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FedRangeQABITransactorRaw struct {
	Contract *FedRangeQABITransactor // Generic write-only contract binding to access the raw methods on
}

// NewFedRangeQABI creates a new instance of FedRangeQABI, bound to a specific deployed contract.
func NewFedRangeQABI(address common.Address, backend bind.ContractBackend) (*FedRangeQABI, error) {
	contract, err := bindFedRangeQABI(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FedRangeQABI{FedRangeQABICaller: FedRangeQABICaller{contract: contract}, FedRangeQABITransactor: FedRangeQABITransactor{contract: contract}, FedRangeQABIFilterer: FedRangeQABIFilterer{contract: contract}}, nil
}

// NewFedRangeQABICaller creates a new read-only instance of FedRangeQABI, bound to a specific deployed contract.
func NewFedRangeQABICaller(address common.Address, caller bind.ContractCaller) (*FedRangeQABICaller, error) {
	contract, err := bindFedRangeQABI(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FedRangeQABICaller{contract: contract}, nil
}

// NewFedRangeQABITransactor creates a new write-only instance of FedRangeQABI, bound to a specific deployed contract.
func NewFedRangeQABITransactor(address common.Address, transactor bind.ContractTransactor) (*FedRangeQABITransactor, error) {
	contract, err := bindFedRangeQABI(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FedRangeQABITransactor{contract: contract}, nil
}

// NewFedRangeQABIFilterer creates a new log filterer instance of FedRangeQABI, bound to a specific deployed contract.
func NewFedRangeQABIFilterer(address common.Address, filterer bind.ContractFilterer) (*FedRangeQABIFilterer, error) {
	contract, err := bindFedRangeQABI(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FedRangeQABIFilterer{contract: contract}, nil
}

// bindFedRangeQABI binds a generic wrapper to an already deployed contract.
func bindFedRangeQABI(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FedRangeQABIABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FedRangeQABI *FedRangeQABIRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FedRangeQABI.Contract.FedRangeQABICaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FedRangeQABI *FedRangeQABIRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.FedRangeQABITransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FedRangeQABI *FedRangeQABIRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.FedRangeQABITransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FedRangeQABI *FedRangeQABICallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FedRangeQABI.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FedRangeQABI *FedRangeQABITransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FedRangeQABI *FedRangeQABITransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.contract.Transact(opts, method, params...)
}

// GetResult is a free data retrieval call binding the contract method 0xde292789.
//
// Solidity: function getResult() view returns(uint256[])
func (_FedRangeQABI *FedRangeQABICaller) GetResult(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _FedRangeQABI.contract.Call(opts, out, "getResult")
	return *ret0, err
}

// GetResult is a free data retrieval call binding the contract method 0xde292789.
//
// Solidity: function getResult() view returns(uint256[])
func (_FedRangeQABI *FedRangeQABISession) GetResult() ([]*big.Int, error) {
	return _FedRangeQABI.Contract.GetResult(&_FedRangeQABI.CallOpts)
}

// GetResult is a free data retrieval call binding the contract method 0xde292789.
//
// Solidity: function getResult() view returns(uint256[])
func (_FedRangeQABI *FedRangeQABICallerSession) GetResult() ([]*big.Int, error) {
	return _FedRangeQABI.Contract.GetResult(&_FedRangeQABI.CallOpts)
}

// ClearResult is a paid mutator transaction binding the contract method 0x6765350a.
//
// Solidity: function clearResult() returns()
func (_FedRangeQABI *FedRangeQABITransactor) ClearResult(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FedRangeQABI.contract.Transact(opts, "clearResult")
}

// ClearResult is a paid mutator transaction binding the contract method 0x6765350a.
//
// Solidity: function clearResult() returns()
func (_FedRangeQABI *FedRangeQABISession) ClearResult() (*types.Transaction, error) {
	return _FedRangeQABI.Contract.ClearResult(&_FedRangeQABI.TransactOpts)
}

// ClearResult is a paid mutator transaction binding the contract method 0x6765350a.
//
// Solidity: function clearResult() returns()
func (_FedRangeQABI *FedRangeQABITransactorSession) ClearResult() (*types.Transaction, error) {
	return _FedRangeQABI.Contract.ClearResult(&_FedRangeQABI.TransactOpts)
}

// Search is a paid mutator transaction binding the contract method 0xbab85532.
//
// Solidity: function search((bytes,uint8[32],bytes[32]) query) returns()
func (_FedRangeQABI *FedRangeQABITransactor) Search(opts *bind.TransactOpts, query Struct0) (*types.Transaction, error) {
	return _FedRangeQABI.contract.Transact(opts, "search", query)
}

// Search is a paid mutator transaction binding the contract method 0xbab85532.
//
// Solidity: function search((bytes,uint8[32],bytes[32]) query) returns()
func (_FedRangeQABI *FedRangeQABISession) Search(query Struct0) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.Search(&_FedRangeQABI.TransactOpts, query)
}

// Search is a paid mutator transaction binding the contract method 0xbab85532.
//
// Solidity: function search((bytes,uint8[32],bytes[32]) query) returns()
func (_FedRangeQABI *FedRangeQABITransactorSession) Search(query Struct0) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.Search(&_FedRangeQABI.TransactOpts, query)
}

// Store is a paid mutator transaction binding the contract method 0xc20c7978.
//
// Solidity: function store((uint256,bytes[32],uint8[1][1][32],bytes[1][32])[] indexToBeAdded) returns()
func (_FedRangeQABI *FedRangeQABITransactor) Store(opts *bind.TransactOpts, indexToBeAdded []Struct1) (*types.Transaction, error) {
	return _FedRangeQABI.contract.Transact(opts, "store", indexToBeAdded)
}

// Store is a paid mutator transaction binding the contract method 0xc20c7978.
//
// Solidity: function store((uint256,bytes[32],uint8[1][1][32],bytes[1][32])[] indexToBeAdded) returns()
func (_FedRangeQABI *FedRangeQABISession) Store(indexToBeAdded []Struct1) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.Store(&_FedRangeQABI.TransactOpts, indexToBeAdded)
}

// Store is a paid mutator transaction binding the contract method 0xc20c7978.
//
// Solidity: function store((uint256,bytes[32],uint8[1][1][32],bytes[1][32])[] indexToBeAdded) returns()
func (_FedRangeQABI *FedRangeQABITransactorSession) Store(indexToBeAdded []Struct1) (*types.Transaction, error) {
	return _FedRangeQABI.Contract.Store(&_FedRangeQABI.TransactOpts, indexToBeAdded)
}
