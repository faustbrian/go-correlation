package correlation_test

import (
	"fmt"

	correlation "github.com/faustbrian/go-correlation"
)

func Example() {
	factory, err := correlation.NewFactory(correlation.FactoryOptions{})
	if err != nil {
		fmt.Println("factory error")
		return
	}
	root, err := factory.Start()
	if err != nil {
		fmt.Println("root error")
		return
	}
	child, err := factory.Next(root)
	if err != nil {
		fmt.Println("child error")
		return
	}
	fmt.Println(
		root.CorrelationID != "",
		root.RequestID != "",
		child.CorrelationID == root.CorrelationID,
		child.CausationID.String() == root.RequestID.String(),
		child.RequestID != root.RequestID,
	)
	// Output: true true true true true
}

func ExampleFactory_Create() {
	generator := &sequenceGenerator{values: []string{"workflow", "root-request"}}
	factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	if err != nil {
		fmt.Println("factory error")
		return
	}
	root, err := factory.Create()
	if err != nil {
		fmt.Println("create error")
		return
	}
	fmt.Println(root.CorrelationID, root.RequestID, root.CausationID)
	// Output: workflow root-request
}

func ExampleFactory_Next() {
	generator := &sequenceGenerator{values: []string{"child-request"}}
	factory, err := correlation.NewFactory(correlation.FactoryOptions{Generator: generator})
	if err != nil {
		fmt.Println("factory error")
		return
	}
	parent := correlation.Values{
		CorrelationID: correlation.MustCorrelationID("workflow", correlation.Policy{}),
		RequestID:     correlation.MustRequestID("parent-request", correlation.Policy{}),
	}
	child, err := factory.Next(parent)
	if err != nil {
		fmt.Println("next error")
		return
	}
	fmt.Println(child.CorrelationID, child.RequestID, child.CausationID)
	// Output: workflow child-request parent-request
}
