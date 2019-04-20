package model

import "time"

type Example struct {
	ExampleID   int
	CustomerID  int
	ProductID   int
	ProductItem []string
	OrderDate   time.Time
}
