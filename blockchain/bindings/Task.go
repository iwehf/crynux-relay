// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TaskTaskInfo is an auto generated low-level Go binding around an user-defined struct.
type TaskTaskInfo struct {
	Id                     *big.Int
	TaskType               *big.Int
	Creator                common.Address
	TaskHash               [32]byte
	DataHash               [32]byte
	VramLimit              *big.Int
	IsSuccess              bool
	SelectedNodes          []common.Address
	Commitments            [][32]byte
	Nonces                 [][32]byte
	CommitmentSubmitRounds []*big.Int
	Results                [][]byte
	ResultDisclosedRounds  []*big.Int
	ResultNode             common.Address
	Aborted                bool
	Timeout                *big.Int
	Balance                *big.Int
	TotalBalance           *big.Int
	GpuName                string
	GpuVram                *big.Int
}

// TaskMetaData contains all meta data concerning the Task contract.
var TaskMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractNode\",\"name\":\"nodeInstance\",\"type\":\"address\"},{\"internalType\":\"contractQOS\",\"name\":\"qosInstance\",\"type\":\"address\"},{\"internalType\":\"contractTaskQueue\",\"name\":\"taskQueueInstance\",\"type\":\"address\"},{\"internalType\":\"contractNetworkStats\",\"name\":\"netStatsInstance\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"TaskAborted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"}],\"name\":\"TaskNodeCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"}],\"name\":\"TaskNodeSlashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"TaskNodeSuccess\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"taskType\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"taskHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"TaskPending\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"TaskResultCommitmentsReady\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"TaskResultUploaded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"taskType\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"selectedNode\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"taskHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"TaskStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"resultNode\",\"type\":\"address\"}],\"name\":\"TaskSuccess\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"cancelTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskType\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"taskHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"vramLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cap\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gpuName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"gpuVram\",\"type\":\"uint256\"}],\"name\":\"createTask\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"name\":\"discloseTaskResult\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"}],\"name\":\"getNodeTask\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"getTask\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"taskType\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"taskHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"vramLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isSuccess\",\"type\":\"bool\"},{\"internalType\":\"address[]\",\"name\":\"selectedNodes\",\"type\":\"address[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"commitments\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"nonces\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"commitmentSubmitRounds\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256[]\",\"name\":\"resultDisclosedRounds\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"resultNode\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"aborted\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBalance\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gpuName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"gpuVram\",\"type\":\"uint256\"}],\"internalType\":\"structTask.TaskInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"root\",\"type\":\"address\"}],\"name\":\"nodeAvailableCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"reportResultsUploaded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"reportTaskError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"commitment\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"}],\"name\":\"submitTaskResultCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAbortedTasks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSuccessTasks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalTasks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"updateDistanceThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"t\",\"type\":\"uint256\"}],\"name\":\"updateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TaskABI is the input ABI used to generate the binding from.
// Deprecated: Use TaskMetaData.ABI instead.
var TaskABI = TaskMetaData.ABI

// Task is an auto generated Go binding around an Ethereum contract.
type Task struct {
	TaskCaller     // Read-only binding to the contract
	TaskTransactor // Write-only binding to the contract
	TaskFilterer   // Log filterer for contract events
}

// TaskCaller is an auto generated read-only Go binding around an Ethereum contract.
type TaskCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TaskTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TaskFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TaskSession struct {
	Contract     *Task             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaskCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TaskCallerSession struct {
	Contract *TaskCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TaskTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TaskTransactorSession struct {
	Contract     *TaskTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaskRaw is an auto generated low-level Go binding around an Ethereum contract.
type TaskRaw struct {
	Contract *Task // Generic contract binding to access the raw methods on
}

// TaskCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TaskCallerRaw struct {
	Contract *TaskCaller // Generic read-only contract binding to access the raw methods on
}

// TaskTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TaskTransactorRaw struct {
	Contract *TaskTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTask creates a new instance of Task, bound to a specific deployed contract.
func NewTask(address common.Address, backend bind.ContractBackend) (*Task, error) {
	contract, err := bindTask(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Task{TaskCaller: TaskCaller{contract: contract}, TaskTransactor: TaskTransactor{contract: contract}, TaskFilterer: TaskFilterer{contract: contract}}, nil
}

// NewTaskCaller creates a new read-only instance of Task, bound to a specific deployed contract.
func NewTaskCaller(address common.Address, caller bind.ContractCaller) (*TaskCaller, error) {
	contract, err := bindTask(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TaskCaller{contract: contract}, nil
}

// NewTaskTransactor creates a new write-only instance of Task, bound to a specific deployed contract.
func NewTaskTransactor(address common.Address, transactor bind.ContractTransactor) (*TaskTransactor, error) {
	contract, err := bindTask(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TaskTransactor{contract: contract}, nil
}

// NewTaskFilterer creates a new log filterer instance of Task, bound to a specific deployed contract.
func NewTaskFilterer(address common.Address, filterer bind.ContractFilterer) (*TaskFilterer, error) {
	contract, err := bindTask(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TaskFilterer{contract: contract}, nil
}

// bindTask binds a generic wrapper to an already deployed contract.
func bindTask(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TaskMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Task *TaskRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Task.Contract.TaskCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Task *TaskRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Task.Contract.TaskTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Task *TaskRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Task.Contract.TaskTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Task *TaskCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Task.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Task *TaskTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Task.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Task *TaskTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Task.Contract.contract.Transact(opts, method, params...)
}

// GetNodeTask is a free data retrieval call binding the contract method 0x9877ad45.
//
// Solidity: function getNodeTask(address nodeAddress) view returns(uint256)
func (_Task *TaskCaller) GetNodeTask(opts *bind.CallOpts, nodeAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "getNodeTask", nodeAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeTask is a free data retrieval call binding the contract method 0x9877ad45.
//
// Solidity: function getNodeTask(address nodeAddress) view returns(uint256)
func (_Task *TaskSession) GetNodeTask(nodeAddress common.Address) (*big.Int, error) {
	return _Task.Contract.GetNodeTask(&_Task.CallOpts, nodeAddress)
}

// GetNodeTask is a free data retrieval call binding the contract method 0x9877ad45.
//
// Solidity: function getNodeTask(address nodeAddress) view returns(uint256)
func (_Task *TaskCallerSession) GetNodeTask(nodeAddress common.Address) (*big.Int, error) {
	return _Task.Contract.GetNodeTask(&_Task.CallOpts, nodeAddress)
}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns((uint256,uint256,address,bytes32,bytes32,uint256,bool,address[],bytes32[],bytes32[],uint256[],bytes[],uint256[],address,bool,uint256,uint256,uint256,string,uint256))
func (_Task *TaskCaller) GetTask(opts *bind.CallOpts, taskId *big.Int) (TaskTaskInfo, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "getTask", taskId)

	if err != nil {
		return *new(TaskTaskInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(TaskTaskInfo)).(*TaskTaskInfo)

	return out0, err

}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns((uint256,uint256,address,bytes32,bytes32,uint256,bool,address[],bytes32[],bytes32[],uint256[],bytes[],uint256[],address,bool,uint256,uint256,uint256,string,uint256))
func (_Task *TaskSession) GetTask(taskId *big.Int) (TaskTaskInfo, error) {
	return _Task.Contract.GetTask(&_Task.CallOpts, taskId)
}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns((uint256,uint256,address,bytes32,bytes32,uint256,bool,address[],bytes32[],bytes32[],uint256[],bytes[],uint256[],address,bool,uint256,uint256,uint256,string,uint256))
func (_Task *TaskCallerSession) GetTask(taskId *big.Int) (TaskTaskInfo, error) {
	return _Task.Contract.GetTask(&_Task.CallOpts, taskId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Task *TaskCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Task *TaskSession) Owner() (common.Address, error) {
	return _Task.Contract.Owner(&_Task.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Task *TaskCallerSession) Owner() (common.Address, error) {
	return _Task.Contract.Owner(&_Task.CallOpts)
}

// TotalAbortedTasks is a free data retrieval call binding the contract method 0x2ff6d2ca.
//
// Solidity: function totalAbortedTasks() view returns(uint256)
func (_Task *TaskCaller) TotalAbortedTasks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "totalAbortedTasks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAbortedTasks is a free data retrieval call binding the contract method 0x2ff6d2ca.
//
// Solidity: function totalAbortedTasks() view returns(uint256)
func (_Task *TaskSession) TotalAbortedTasks() (*big.Int, error) {
	return _Task.Contract.TotalAbortedTasks(&_Task.CallOpts)
}

// TotalAbortedTasks is a free data retrieval call binding the contract method 0x2ff6d2ca.
//
// Solidity: function totalAbortedTasks() view returns(uint256)
func (_Task *TaskCallerSession) TotalAbortedTasks() (*big.Int, error) {
	return _Task.Contract.TotalAbortedTasks(&_Task.CallOpts)
}

// TotalSuccessTasks is a free data retrieval call binding the contract method 0x775820f8.
//
// Solidity: function totalSuccessTasks() view returns(uint256)
func (_Task *TaskCaller) TotalSuccessTasks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "totalSuccessTasks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSuccessTasks is a free data retrieval call binding the contract method 0x775820f8.
//
// Solidity: function totalSuccessTasks() view returns(uint256)
func (_Task *TaskSession) TotalSuccessTasks() (*big.Int, error) {
	return _Task.Contract.TotalSuccessTasks(&_Task.CallOpts)
}

// TotalSuccessTasks is a free data retrieval call binding the contract method 0x775820f8.
//
// Solidity: function totalSuccessTasks() view returns(uint256)
func (_Task *TaskCallerSession) TotalSuccessTasks() (*big.Int, error) {
	return _Task.Contract.TotalSuccessTasks(&_Task.CallOpts)
}

// TotalTasks is a free data retrieval call binding the contract method 0xd22c81e5.
//
// Solidity: function totalTasks() view returns(uint256)
func (_Task *TaskCaller) TotalTasks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Task.contract.Call(opts, &out, "totalTasks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalTasks is a free data retrieval call binding the contract method 0xd22c81e5.
//
// Solidity: function totalTasks() view returns(uint256)
func (_Task *TaskSession) TotalTasks() (*big.Int, error) {
	return _Task.Contract.TotalTasks(&_Task.CallOpts)
}

// TotalTasks is a free data retrieval call binding the contract method 0xd22c81e5.
//
// Solidity: function totalTasks() view returns(uint256)
func (_Task *TaskCallerSession) TotalTasks() (*big.Int, error) {
	return _Task.Contract.TotalTasks(&_Task.CallOpts)
}

// CancelTask is a paid mutator transaction binding the contract method 0x7eec20a8.
//
// Solidity: function cancelTask(uint256 taskId) returns()
func (_Task *TaskTransactor) CancelTask(opts *bind.TransactOpts, taskId *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "cancelTask", taskId)
}

// CancelTask is a paid mutator transaction binding the contract method 0x7eec20a8.
//
// Solidity: function cancelTask(uint256 taskId) returns()
func (_Task *TaskSession) CancelTask(taskId *big.Int) (*types.Transaction, error) {
	return _Task.Contract.CancelTask(&_Task.TransactOpts, taskId)
}

// CancelTask is a paid mutator transaction binding the contract method 0x7eec20a8.
//
// Solidity: function cancelTask(uint256 taskId) returns()
func (_Task *TaskTransactorSession) CancelTask(taskId *big.Int) (*types.Transaction, error) {
	return _Task.Contract.CancelTask(&_Task.TransactOpts, taskId)
}

// CreateTask is a paid mutator transaction binding the contract method 0xb8300eaa.
//
// Solidity: function createTask(uint256 taskType, bytes32 taskHash, bytes32 dataHash, uint256 vramLimit, uint256 cap, string gpuName, uint256 gpuVram) payable returns()
func (_Task *TaskTransactor) CreateTask(opts *bind.TransactOpts, taskType *big.Int, taskHash [32]byte, dataHash [32]byte, vramLimit *big.Int, cap *big.Int, gpuName string, gpuVram *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "createTask", taskType, taskHash, dataHash, vramLimit, cap, gpuName, gpuVram)
}

// CreateTask is a paid mutator transaction binding the contract method 0xb8300eaa.
//
// Solidity: function createTask(uint256 taskType, bytes32 taskHash, bytes32 dataHash, uint256 vramLimit, uint256 cap, string gpuName, uint256 gpuVram) payable returns()
func (_Task *TaskSession) CreateTask(taskType *big.Int, taskHash [32]byte, dataHash [32]byte, vramLimit *big.Int, cap *big.Int, gpuName string, gpuVram *big.Int) (*types.Transaction, error) {
	return _Task.Contract.CreateTask(&_Task.TransactOpts, taskType, taskHash, dataHash, vramLimit, cap, gpuName, gpuVram)
}

// CreateTask is a paid mutator transaction binding the contract method 0xb8300eaa.
//
// Solidity: function createTask(uint256 taskType, bytes32 taskHash, bytes32 dataHash, uint256 vramLimit, uint256 cap, string gpuName, uint256 gpuVram) payable returns()
func (_Task *TaskTransactorSession) CreateTask(taskType *big.Int, taskHash [32]byte, dataHash [32]byte, vramLimit *big.Int, cap *big.Int, gpuName string, gpuVram *big.Int) (*types.Transaction, error) {
	return _Task.Contract.CreateTask(&_Task.TransactOpts, taskType, taskHash, dataHash, vramLimit, cap, gpuName, gpuVram)
}

// DiscloseTaskResult is a paid mutator transaction binding the contract method 0x63be8c33.
//
// Solidity: function discloseTaskResult(uint256 taskId, uint256 round, bytes result) returns()
func (_Task *TaskTransactor) DiscloseTaskResult(opts *bind.TransactOpts, taskId *big.Int, round *big.Int, result []byte) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "discloseTaskResult", taskId, round, result)
}

// DiscloseTaskResult is a paid mutator transaction binding the contract method 0x63be8c33.
//
// Solidity: function discloseTaskResult(uint256 taskId, uint256 round, bytes result) returns()
func (_Task *TaskSession) DiscloseTaskResult(taskId *big.Int, round *big.Int, result []byte) (*types.Transaction, error) {
	return _Task.Contract.DiscloseTaskResult(&_Task.TransactOpts, taskId, round, result)
}

// DiscloseTaskResult is a paid mutator transaction binding the contract method 0x63be8c33.
//
// Solidity: function discloseTaskResult(uint256 taskId, uint256 round, bytes result) returns()
func (_Task *TaskTransactorSession) DiscloseTaskResult(taskId *big.Int, round *big.Int, result []byte) (*types.Transaction, error) {
	return _Task.Contract.DiscloseTaskResult(&_Task.TransactOpts, taskId, round, result)
}

// NodeAvailableCallback is a paid mutator transaction binding the contract method 0xd06ac297.
//
// Solidity: function nodeAvailableCallback(address root) returns()
func (_Task *TaskTransactor) NodeAvailableCallback(opts *bind.TransactOpts, root common.Address) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "nodeAvailableCallback", root)
}

// NodeAvailableCallback is a paid mutator transaction binding the contract method 0xd06ac297.
//
// Solidity: function nodeAvailableCallback(address root) returns()
func (_Task *TaskSession) NodeAvailableCallback(root common.Address) (*types.Transaction, error) {
	return _Task.Contract.NodeAvailableCallback(&_Task.TransactOpts, root)
}

// NodeAvailableCallback is a paid mutator transaction binding the contract method 0xd06ac297.
//
// Solidity: function nodeAvailableCallback(address root) returns()
func (_Task *TaskTransactorSession) NodeAvailableCallback(root common.Address) (*types.Transaction, error) {
	return _Task.Contract.NodeAvailableCallback(&_Task.TransactOpts, root)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Task *TaskTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Task *TaskSession) RenounceOwnership() (*types.Transaction, error) {
	return _Task.Contract.RenounceOwnership(&_Task.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Task *TaskTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Task.Contract.RenounceOwnership(&_Task.TransactOpts)
}

// ReportResultsUploaded is a paid mutator transaction binding the contract method 0x95b1d89b.
//
// Solidity: function reportResultsUploaded(uint256 taskId, uint256 round) returns()
func (_Task *TaskTransactor) ReportResultsUploaded(opts *bind.TransactOpts, taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "reportResultsUploaded", taskId, round)
}

// ReportResultsUploaded is a paid mutator transaction binding the contract method 0x95b1d89b.
//
// Solidity: function reportResultsUploaded(uint256 taskId, uint256 round) returns()
func (_Task *TaskSession) ReportResultsUploaded(taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.Contract.ReportResultsUploaded(&_Task.TransactOpts, taskId, round)
}

// ReportResultsUploaded is a paid mutator transaction binding the contract method 0x95b1d89b.
//
// Solidity: function reportResultsUploaded(uint256 taskId, uint256 round) returns()
func (_Task *TaskTransactorSession) ReportResultsUploaded(taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.Contract.ReportResultsUploaded(&_Task.TransactOpts, taskId, round)
}

// ReportTaskError is a paid mutator transaction binding the contract method 0x695b1c8f.
//
// Solidity: function reportTaskError(uint256 taskId, uint256 round) returns()
func (_Task *TaskTransactor) ReportTaskError(opts *bind.TransactOpts, taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "reportTaskError", taskId, round)
}

// ReportTaskError is a paid mutator transaction binding the contract method 0x695b1c8f.
//
// Solidity: function reportTaskError(uint256 taskId, uint256 round) returns()
func (_Task *TaskSession) ReportTaskError(taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.Contract.ReportTaskError(&_Task.TransactOpts, taskId, round)
}

// ReportTaskError is a paid mutator transaction binding the contract method 0x695b1c8f.
//
// Solidity: function reportTaskError(uint256 taskId, uint256 round) returns()
func (_Task *TaskTransactorSession) ReportTaskError(taskId *big.Int, round *big.Int) (*types.Transaction, error) {
	return _Task.Contract.ReportTaskError(&_Task.TransactOpts, taskId, round)
}

// SubmitTaskResultCommitment is a paid mutator transaction binding the contract method 0x47f90738.
//
// Solidity: function submitTaskResultCommitment(uint256 taskId, uint256 round, bytes32 commitment, bytes32 nonce) returns()
func (_Task *TaskTransactor) SubmitTaskResultCommitment(opts *bind.TransactOpts, taskId *big.Int, round *big.Int, commitment [32]byte, nonce [32]byte) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "submitTaskResultCommitment", taskId, round, commitment, nonce)
}

// SubmitTaskResultCommitment is a paid mutator transaction binding the contract method 0x47f90738.
//
// Solidity: function submitTaskResultCommitment(uint256 taskId, uint256 round, bytes32 commitment, bytes32 nonce) returns()
func (_Task *TaskSession) SubmitTaskResultCommitment(taskId *big.Int, round *big.Int, commitment [32]byte, nonce [32]byte) (*types.Transaction, error) {
	return _Task.Contract.SubmitTaskResultCommitment(&_Task.TransactOpts, taskId, round, commitment, nonce)
}

// SubmitTaskResultCommitment is a paid mutator transaction binding the contract method 0x47f90738.
//
// Solidity: function submitTaskResultCommitment(uint256 taskId, uint256 round, bytes32 commitment, bytes32 nonce) returns()
func (_Task *TaskTransactorSession) SubmitTaskResultCommitment(taskId *big.Int, round *big.Int, commitment [32]byte, nonce [32]byte) (*types.Transaction, error) {
	return _Task.Contract.SubmitTaskResultCommitment(&_Task.TransactOpts, taskId, round, commitment, nonce)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Task *TaskTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Task *TaskSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Task.Contract.TransferOwnership(&_Task.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Task *TaskTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Task.Contract.TransferOwnership(&_Task.TransactOpts, newOwner)
}

// UpdateDistanceThreshold is a paid mutator transaction binding the contract method 0x9244e462.
//
// Solidity: function updateDistanceThreshold(uint256 threshold) returns()
func (_Task *TaskTransactor) UpdateDistanceThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "updateDistanceThreshold", threshold)
}

// UpdateDistanceThreshold is a paid mutator transaction binding the contract method 0x9244e462.
//
// Solidity: function updateDistanceThreshold(uint256 threshold) returns()
func (_Task *TaskSession) UpdateDistanceThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Task.Contract.UpdateDistanceThreshold(&_Task.TransactOpts, threshold)
}

// UpdateDistanceThreshold is a paid mutator transaction binding the contract method 0x9244e462.
//
// Solidity: function updateDistanceThreshold(uint256 threshold) returns()
func (_Task *TaskTransactorSession) UpdateDistanceThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Task.Contract.UpdateDistanceThreshold(&_Task.TransactOpts, threshold)
}

// UpdateTimeout is a paid mutator transaction binding the contract method 0xa330214e.
//
// Solidity: function updateTimeout(uint256 t) returns()
func (_Task *TaskTransactor) UpdateTimeout(opts *bind.TransactOpts, t *big.Int) (*types.Transaction, error) {
	return _Task.contract.Transact(opts, "updateTimeout", t)
}

// UpdateTimeout is a paid mutator transaction binding the contract method 0xa330214e.
//
// Solidity: function updateTimeout(uint256 t) returns()
func (_Task *TaskSession) UpdateTimeout(t *big.Int) (*types.Transaction, error) {
	return _Task.Contract.UpdateTimeout(&_Task.TransactOpts, t)
}

// UpdateTimeout is a paid mutator transaction binding the contract method 0xa330214e.
//
// Solidity: function updateTimeout(uint256 t) returns()
func (_Task *TaskTransactorSession) UpdateTimeout(t *big.Int) (*types.Transaction, error) {
	return _Task.Contract.UpdateTimeout(&_Task.TransactOpts, t)
}

// TaskOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Task contract.
type TaskOwnershipTransferredIterator struct {
	Event *TaskOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskOwnershipTransferred represents a OwnershipTransferred event raised by the Task contract.
type TaskOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Task *TaskFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaskOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaskOwnershipTransferredIterator{contract: _Task.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Task *TaskFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TaskOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskOwnershipTransferred)
				if err := _Task.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Task *TaskFilterer) ParseOwnershipTransferred(log types.Log) (*TaskOwnershipTransferred, error) {
	event := new(TaskOwnershipTransferred)
	if err := _Task.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskAbortedIterator is returned from FilterTaskAborted and is used to iterate over the raw logs and unpacked data for TaskAborted events raised by the Task contract.
type TaskTaskAbortedIterator struct {
	Event *TaskTaskAborted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskAbortedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskAborted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskAborted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskAbortedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskAbortedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskAborted represents a TaskAborted event raised by the Task contract.
type TaskTaskAborted struct {
	TaskId *big.Int
	Reason string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTaskAborted is a free log retrieval operation binding the contract event 0x346216bb718ede47b40f85e7caf58cd997eebda879d146a8b49091c106c851d9.
//
// Solidity: event TaskAborted(uint256 indexed taskId, string reason)
func (_Task *TaskFilterer) FilterTaskAborted(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskAbortedIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskAborted", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskAbortedIterator{contract: _Task.contract, event: "TaskAborted", logs: logs, sub: sub}, nil
}

// WatchTaskAborted is a free log subscription operation binding the contract event 0x346216bb718ede47b40f85e7caf58cd997eebda879d146a8b49091c106c851d9.
//
// Solidity: event TaskAborted(uint256 indexed taskId, string reason)
func (_Task *TaskFilterer) WatchTaskAborted(opts *bind.WatchOpts, sink chan<- *TaskTaskAborted, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskAborted", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskAborted)
				if err := _Task.contract.UnpackLog(event, "TaskAborted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskAborted is a log parse operation binding the contract event 0x346216bb718ede47b40f85e7caf58cd997eebda879d146a8b49091c106c851d9.
//
// Solidity: event TaskAborted(uint256 indexed taskId, string reason)
func (_Task *TaskFilterer) ParseTaskAborted(log types.Log) (*TaskTaskAborted, error) {
	event := new(TaskTaskAborted)
	if err := _Task.contract.UnpackLog(event, "TaskAborted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskNodeCancelledIterator is returned from FilterTaskNodeCancelled and is used to iterate over the raw logs and unpacked data for TaskNodeCancelled events raised by the Task contract.
type TaskTaskNodeCancelledIterator struct {
	Event *TaskTaskNodeCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskNodeCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskNodeCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskNodeCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskNodeCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskNodeCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskNodeCancelled represents a TaskNodeCancelled event raised by the Task contract.
type TaskTaskNodeCancelled struct {
	TaskId      *big.Int
	NodeAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTaskNodeCancelled is a free log retrieval operation binding the contract event 0x5451b4da8bffa59b36dc84ae091811086db04bcdaf13db59567eba01a3d7d51e.
//
// Solidity: event TaskNodeCancelled(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) FilterTaskNodeCancelled(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskNodeCancelledIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskNodeCancelled", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskNodeCancelledIterator{contract: _Task.contract, event: "TaskNodeCancelled", logs: logs, sub: sub}, nil
}

// WatchTaskNodeCancelled is a free log subscription operation binding the contract event 0x5451b4da8bffa59b36dc84ae091811086db04bcdaf13db59567eba01a3d7d51e.
//
// Solidity: event TaskNodeCancelled(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) WatchTaskNodeCancelled(opts *bind.WatchOpts, sink chan<- *TaskTaskNodeCancelled, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskNodeCancelled", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskNodeCancelled)
				if err := _Task.contract.UnpackLog(event, "TaskNodeCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskNodeCancelled is a log parse operation binding the contract event 0x5451b4da8bffa59b36dc84ae091811086db04bcdaf13db59567eba01a3d7d51e.
//
// Solidity: event TaskNodeCancelled(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) ParseTaskNodeCancelled(log types.Log) (*TaskTaskNodeCancelled, error) {
	event := new(TaskTaskNodeCancelled)
	if err := _Task.contract.UnpackLog(event, "TaskNodeCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskNodeSlashedIterator is returned from FilterTaskNodeSlashed and is used to iterate over the raw logs and unpacked data for TaskNodeSlashed events raised by the Task contract.
type TaskTaskNodeSlashedIterator struct {
	Event *TaskTaskNodeSlashed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskNodeSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskNodeSlashed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskNodeSlashed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskNodeSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskNodeSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskNodeSlashed represents a TaskNodeSlashed event raised by the Task contract.
type TaskTaskNodeSlashed struct {
	TaskId      *big.Int
	NodeAddress common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTaskNodeSlashed is a free log retrieval operation binding the contract event 0x0f297f1fd7c66c9836a5636092ca18255164bf531f1f9b8717e1214e86620136.
//
// Solidity: event TaskNodeSlashed(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) FilterTaskNodeSlashed(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskNodeSlashedIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskNodeSlashed", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskNodeSlashedIterator{contract: _Task.contract, event: "TaskNodeSlashed", logs: logs, sub: sub}, nil
}

// WatchTaskNodeSlashed is a free log subscription operation binding the contract event 0x0f297f1fd7c66c9836a5636092ca18255164bf531f1f9b8717e1214e86620136.
//
// Solidity: event TaskNodeSlashed(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) WatchTaskNodeSlashed(opts *bind.WatchOpts, sink chan<- *TaskTaskNodeSlashed, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskNodeSlashed", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskNodeSlashed)
				if err := _Task.contract.UnpackLog(event, "TaskNodeSlashed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskNodeSlashed is a log parse operation binding the contract event 0x0f297f1fd7c66c9836a5636092ca18255164bf531f1f9b8717e1214e86620136.
//
// Solidity: event TaskNodeSlashed(uint256 indexed taskId, address nodeAddress)
func (_Task *TaskFilterer) ParseTaskNodeSlashed(log types.Log) (*TaskTaskNodeSlashed, error) {
	event := new(TaskTaskNodeSlashed)
	if err := _Task.contract.UnpackLog(event, "TaskNodeSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskNodeSuccessIterator is returned from FilterTaskNodeSuccess and is used to iterate over the raw logs and unpacked data for TaskNodeSuccess events raised by the Task contract.
type TaskTaskNodeSuccessIterator struct {
	Event *TaskTaskNodeSuccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskNodeSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskNodeSuccess)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskNodeSuccess)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskNodeSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskNodeSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskNodeSuccess represents a TaskNodeSuccess event raised by the Task contract.
type TaskTaskNodeSuccess struct {
	TaskId      *big.Int
	NodeAddress common.Address
	Fee         *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTaskNodeSuccess is a free log retrieval operation binding the contract event 0x7c0de5a9de40ae20fc84128fa8611c1afa9ef6e2dad146d43fe44aa25081ccde.
//
// Solidity: event TaskNodeSuccess(uint256 indexed taskId, address nodeAddress, uint256 fee)
func (_Task *TaskFilterer) FilterTaskNodeSuccess(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskNodeSuccessIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskNodeSuccess", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskNodeSuccessIterator{contract: _Task.contract, event: "TaskNodeSuccess", logs: logs, sub: sub}, nil
}

// WatchTaskNodeSuccess is a free log subscription operation binding the contract event 0x7c0de5a9de40ae20fc84128fa8611c1afa9ef6e2dad146d43fe44aa25081ccde.
//
// Solidity: event TaskNodeSuccess(uint256 indexed taskId, address nodeAddress, uint256 fee)
func (_Task *TaskFilterer) WatchTaskNodeSuccess(opts *bind.WatchOpts, sink chan<- *TaskTaskNodeSuccess, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskNodeSuccess", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskNodeSuccess)
				if err := _Task.contract.UnpackLog(event, "TaskNodeSuccess", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskNodeSuccess is a log parse operation binding the contract event 0x7c0de5a9de40ae20fc84128fa8611c1afa9ef6e2dad146d43fe44aa25081ccde.
//
// Solidity: event TaskNodeSuccess(uint256 indexed taskId, address nodeAddress, uint256 fee)
func (_Task *TaskFilterer) ParseTaskNodeSuccess(log types.Log) (*TaskTaskNodeSuccess, error) {
	event := new(TaskTaskNodeSuccess)
	if err := _Task.contract.UnpackLog(event, "TaskNodeSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskPendingIterator is returned from FilterTaskPending and is used to iterate over the raw logs and unpacked data for TaskPending events raised by the Task contract.
type TaskTaskPendingIterator struct {
	Event *TaskTaskPending // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskPendingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskPending)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskPending)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskPendingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskPendingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskPending represents a TaskPending event raised by the Task contract.
type TaskTaskPending struct {
	TaskId   *big.Int
	TaskType *big.Int
	Creator  common.Address
	TaskHash [32]byte
	DataHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTaskPending is a free log retrieval operation binding the contract event 0xa8b23b175c46381c07d9001d248005cebb7fd414d3f64bf1cebfe320a23f6050.
//
// Solidity: event TaskPending(uint256 taskId, uint256 taskType, address indexed creator, bytes32 taskHash, bytes32 dataHash)
func (_Task *TaskFilterer) FilterTaskPending(opts *bind.FilterOpts, creator []common.Address) (*TaskTaskPendingIterator, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskPending", creatorRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskPendingIterator{contract: _Task.contract, event: "TaskPending", logs: logs, sub: sub}, nil
}

// WatchTaskPending is a free log subscription operation binding the contract event 0xa8b23b175c46381c07d9001d248005cebb7fd414d3f64bf1cebfe320a23f6050.
//
// Solidity: event TaskPending(uint256 taskId, uint256 taskType, address indexed creator, bytes32 taskHash, bytes32 dataHash)
func (_Task *TaskFilterer) WatchTaskPending(opts *bind.WatchOpts, sink chan<- *TaskTaskPending, creator []common.Address) (event.Subscription, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskPending", creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskPending)
				if err := _Task.contract.UnpackLog(event, "TaskPending", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskPending is a log parse operation binding the contract event 0xa8b23b175c46381c07d9001d248005cebb7fd414d3f64bf1cebfe320a23f6050.
//
// Solidity: event TaskPending(uint256 taskId, uint256 taskType, address indexed creator, bytes32 taskHash, bytes32 dataHash)
func (_Task *TaskFilterer) ParseTaskPending(log types.Log) (*TaskTaskPending, error) {
	event := new(TaskTaskPending)
	if err := _Task.contract.UnpackLog(event, "TaskPending", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskResultCommitmentsReadyIterator is returned from FilterTaskResultCommitmentsReady and is used to iterate over the raw logs and unpacked data for TaskResultCommitmentsReady events raised by the Task contract.
type TaskTaskResultCommitmentsReadyIterator struct {
	Event *TaskTaskResultCommitmentsReady // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskResultCommitmentsReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskResultCommitmentsReady)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskResultCommitmentsReady)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskResultCommitmentsReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskResultCommitmentsReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskResultCommitmentsReady represents a TaskResultCommitmentsReady event raised by the Task contract.
type TaskTaskResultCommitmentsReady struct {
	TaskId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTaskResultCommitmentsReady is a free log retrieval operation binding the contract event 0xb16812d8924e5125b1e331ca0097225a801aaa45b056a9ab12ab6ba658e6c9e5.
//
// Solidity: event TaskResultCommitmentsReady(uint256 indexed taskId)
func (_Task *TaskFilterer) FilterTaskResultCommitmentsReady(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskResultCommitmentsReadyIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskResultCommitmentsReady", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskResultCommitmentsReadyIterator{contract: _Task.contract, event: "TaskResultCommitmentsReady", logs: logs, sub: sub}, nil
}

// WatchTaskResultCommitmentsReady is a free log subscription operation binding the contract event 0xb16812d8924e5125b1e331ca0097225a801aaa45b056a9ab12ab6ba658e6c9e5.
//
// Solidity: event TaskResultCommitmentsReady(uint256 indexed taskId)
func (_Task *TaskFilterer) WatchTaskResultCommitmentsReady(opts *bind.WatchOpts, sink chan<- *TaskTaskResultCommitmentsReady, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskResultCommitmentsReady", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskResultCommitmentsReady)
				if err := _Task.contract.UnpackLog(event, "TaskResultCommitmentsReady", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskResultCommitmentsReady is a log parse operation binding the contract event 0xb16812d8924e5125b1e331ca0097225a801aaa45b056a9ab12ab6ba658e6c9e5.
//
// Solidity: event TaskResultCommitmentsReady(uint256 indexed taskId)
func (_Task *TaskFilterer) ParseTaskResultCommitmentsReady(log types.Log) (*TaskTaskResultCommitmentsReady, error) {
	event := new(TaskTaskResultCommitmentsReady)
	if err := _Task.contract.UnpackLog(event, "TaskResultCommitmentsReady", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskResultUploadedIterator is returned from FilterTaskResultUploaded and is used to iterate over the raw logs and unpacked data for TaskResultUploaded events raised by the Task contract.
type TaskTaskResultUploadedIterator struct {
	Event *TaskTaskResultUploaded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskResultUploadedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskResultUploaded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskResultUploaded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskResultUploadedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskResultUploadedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskResultUploaded represents a TaskResultUploaded event raised by the Task contract.
type TaskTaskResultUploaded struct {
	TaskId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTaskResultUploaded is a free log retrieval operation binding the contract event 0x11f72037be5686bd54d7042f5a9392a74a7c127754536da6a05710bfdf1f9a0e.
//
// Solidity: event TaskResultUploaded(uint256 indexed taskId)
func (_Task *TaskFilterer) FilterTaskResultUploaded(opts *bind.FilterOpts, taskId []*big.Int) (*TaskTaskResultUploadedIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskResultUploaded", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskResultUploadedIterator{contract: _Task.contract, event: "TaskResultUploaded", logs: logs, sub: sub}, nil
}

// WatchTaskResultUploaded is a free log subscription operation binding the contract event 0x11f72037be5686bd54d7042f5a9392a74a7c127754536da6a05710bfdf1f9a0e.
//
// Solidity: event TaskResultUploaded(uint256 indexed taskId)
func (_Task *TaskFilterer) WatchTaskResultUploaded(opts *bind.WatchOpts, sink chan<- *TaskTaskResultUploaded, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskResultUploaded", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskResultUploaded)
				if err := _Task.contract.UnpackLog(event, "TaskResultUploaded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskResultUploaded is a log parse operation binding the contract event 0x11f72037be5686bd54d7042f5a9392a74a7c127754536da6a05710bfdf1f9a0e.
//
// Solidity: event TaskResultUploaded(uint256 indexed taskId)
func (_Task *TaskFilterer) ParseTaskResultUploaded(log types.Log) (*TaskTaskResultUploaded, error) {
	event := new(TaskTaskResultUploaded)
	if err := _Task.contract.UnpackLog(event, "TaskResultUploaded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskStartedIterator is returned from FilterTaskStarted and is used to iterate over the raw logs and unpacked data for TaskStarted events raised by the Task contract.
type TaskTaskStartedIterator struct {
	Event *TaskTaskStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskStarted represents a TaskStarted event raised by the Task contract.
type TaskTaskStarted struct {
	TaskId       *big.Int
	TaskType     *big.Int
	Creator      common.Address
	SelectedNode common.Address
	TaskHash     [32]byte
	DataHash     [32]byte
	Round        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTaskStarted is a free log retrieval operation binding the contract event 0x1399dffd65fe4e2d81d94992e9227fc0a75a5289ba36619c5afe826c0e686440.
//
// Solidity: event TaskStarted(uint256 taskId, uint256 taskType, address indexed creator, address indexed selectedNode, bytes32 taskHash, bytes32 dataHash, uint256 round)
func (_Task *TaskFilterer) FilterTaskStarted(opts *bind.FilterOpts, creator []common.Address, selectedNode []common.Address) (*TaskTaskStartedIterator, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var selectedNodeRule []interface{}
	for _, selectedNodeItem := range selectedNode {
		selectedNodeRule = append(selectedNodeRule, selectedNodeItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskStarted", creatorRule, selectedNodeRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskStartedIterator{contract: _Task.contract, event: "TaskStarted", logs: logs, sub: sub}, nil
}

// WatchTaskStarted is a free log subscription operation binding the contract event 0x1399dffd65fe4e2d81d94992e9227fc0a75a5289ba36619c5afe826c0e686440.
//
// Solidity: event TaskStarted(uint256 taskId, uint256 taskType, address indexed creator, address indexed selectedNode, bytes32 taskHash, bytes32 dataHash, uint256 round)
func (_Task *TaskFilterer) WatchTaskStarted(opts *bind.WatchOpts, sink chan<- *TaskTaskStarted, creator []common.Address, selectedNode []common.Address) (event.Subscription, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var selectedNodeRule []interface{}
	for _, selectedNodeItem := range selectedNode {
		selectedNodeRule = append(selectedNodeRule, selectedNodeItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskStarted", creatorRule, selectedNodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskStarted)
				if err := _Task.contract.UnpackLog(event, "TaskStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskStarted is a log parse operation binding the contract event 0x1399dffd65fe4e2d81d94992e9227fc0a75a5289ba36619c5afe826c0e686440.
//
// Solidity: event TaskStarted(uint256 taskId, uint256 taskType, address indexed creator, address indexed selectedNode, bytes32 taskHash, bytes32 dataHash, uint256 round)
func (_Task *TaskFilterer) ParseTaskStarted(log types.Log) (*TaskTaskStarted, error) {
	event := new(TaskTaskStarted)
	if err := _Task.contract.UnpackLog(event, "TaskStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskTaskSuccessIterator is returned from FilterTaskSuccess and is used to iterate over the raw logs and unpacked data for TaskSuccess events raised by the Task contract.
type TaskTaskSuccessIterator struct {
	Event *TaskTaskSuccess // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskTaskSuccessIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskTaskSuccess)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskTaskSuccess)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskTaskSuccessIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskTaskSuccessIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskTaskSuccess represents a TaskSuccess event raised by the Task contract.
type TaskTaskSuccess struct {
	TaskId     *big.Int
	Result     []byte
	ResultNode common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTaskSuccess is a free log retrieval operation binding the contract event 0x3cd644a756687c7dd87bbe74e2fd65bde569a3335e26f68f8fa334b2849673f7.
//
// Solidity: event TaskSuccess(uint256 indexed taskId, bytes result, address indexed resultNode)
func (_Task *TaskFilterer) FilterTaskSuccess(opts *bind.FilterOpts, taskId []*big.Int, resultNode []common.Address) (*TaskTaskSuccessIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	var resultNodeRule []interface{}
	for _, resultNodeItem := range resultNode {
		resultNodeRule = append(resultNodeRule, resultNodeItem)
	}

	logs, sub, err := _Task.contract.FilterLogs(opts, "TaskSuccess", taskIdRule, resultNodeRule)
	if err != nil {
		return nil, err
	}
	return &TaskTaskSuccessIterator{contract: _Task.contract, event: "TaskSuccess", logs: logs, sub: sub}, nil
}

// WatchTaskSuccess is a free log subscription operation binding the contract event 0x3cd644a756687c7dd87bbe74e2fd65bde569a3335e26f68f8fa334b2849673f7.
//
// Solidity: event TaskSuccess(uint256 indexed taskId, bytes result, address indexed resultNode)
func (_Task *TaskFilterer) WatchTaskSuccess(opts *bind.WatchOpts, sink chan<- *TaskTaskSuccess, taskId []*big.Int, resultNode []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	var resultNodeRule []interface{}
	for _, resultNodeItem := range resultNode {
		resultNodeRule = append(resultNodeRule, resultNodeItem)
	}

	logs, sub, err := _Task.contract.WatchLogs(opts, "TaskSuccess", taskIdRule, resultNodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskTaskSuccess)
				if err := _Task.contract.UnpackLog(event, "TaskSuccess", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskSuccess is a log parse operation binding the contract event 0x3cd644a756687c7dd87bbe74e2fd65bde569a3335e26f68f8fa334b2849673f7.
//
// Solidity: event TaskSuccess(uint256 indexed taskId, bytes result, address indexed resultNode)
func (_Task *TaskFilterer) ParseTaskSuccess(log types.Log) (*TaskTaskSuccess, error) {
	event := new(TaskTaskSuccess)
	if err := _Task.contract.UnpackLog(event, "TaskSuccess", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
