package main

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionContext is a mock implementation of TransactionContextInterface
type MockTransactionContext struct {
	mock.Mock
	contractapi.TransactionContextInterface
}

// MockChaincodeStub is a mock implementation of ChaincodeStubInterface
type MockChaincodeStub struct {
	mock.Mock
	shim.ChaincodeStubInterface
}

func (m *MockChaincodeStub) PutState(key string, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockChaincodeStub) GetState(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(0)
}

func TestCreateRecord(t *testing.T) {
	scc := new(SimpleChaincode)
	ctx := new(MockTransactionContext)
	stub := new(MockChaincodeStub)

	ctx.On("GetStub").Return(stub)

	stub.On("PutState", "drone1", mock.Anything).Return(nil)

	err := scc.CreateRecord(ctx, "drone1", "10001", 100, "record1", "reserved1")
	assert.NoError(t, err)

	// stub.AssertCalled(t, "PutState", "drone1", mock.Anything)
}

func TestRecordExists(t *testing.T) {
	scc := new(SimpleChaincode)
	ctx := new(MockTransactionContext)
	stub := new(MockChaincodeStub)

	ctx.On("GetStub").Return(stub)

	stub.On("GetState", "drone1_100").Return([]byte("record exists"), nil)
	stub.On("GetState", "nonexistent").Return(nil, nil)

	exists, err := scc.RecordExists(ctx, "drone1_100")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = scc.RecordExists(ctx, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)
}
