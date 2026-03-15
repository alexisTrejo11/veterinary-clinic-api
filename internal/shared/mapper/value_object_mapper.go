package mapper

import (
	"clinic-vet-api/internal/core/customers"
	"clinic-vet-api/internal/core/employees"
	"clinic-vet-api/internal/shared"
)

type ValueObjcetMapper struct {
}

func PtrToEmployeeIDPtr(id *uint) *employees.EmployeeID {
	if id == nil {
		return nil
	}
	employeeID := employees.NewEmployeeID(*id)
	return &employeeID
}

func PtrToCustomerIDPtr(id *uint) *customers.CustomerID {
	if id == nil {
		return nil
	}
	customerID := customers.NewCustomerID(*id)
	return &customerID
}

func PaidByCustomerPtrToPtr(id *customers.CustomerID) *uint {
	if id == nil {
		return nil
	}
	customerID := id.Value()
	return &customerID
}

func PaidByEmployeePtrToPtr(id *employees.EmployeeID) *uint {
	if id == nil {
		return nil
	}
	employeeID := id.Value()
	return &employeeID
}

func MoneyPtrToFloat64Ptr(money *shared.Money) *float64 {
	if money == nil {
		return nil
	}
	amount := money.Amount().Float64()
	return &amount
}

func PaidToEmployeePtrToPtr(id *employees.EmployeeID) *uint {
	if id == nil {
		return nil
	}
	employeeID := id.Value()
	return &employeeID
}
