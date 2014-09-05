package workflow

import (
	"testing"
)

const (
	PLACE_ORDER     = Step("PlaceOrder")
	CAPTURE_PAYMENT = Step("CapturePayment")
	COMPLETE        = Step("Complete")
	FAIL            = Step("Fail")

	ORDERED        = Event("ordered")
	PAID           = Event("paid")
	PAYMENT_FAILED = Event("payment_failed")
)

type OrderWorkflow interface {
	NewOrderHandler() Event
	PlaceOrderHandler() Event
	CapturePaymentHandler() Event
}

type SuccessDocument struct {
	Name   string
	Status string
}

func (t *SuccessDocument) String() string {
	return "Name [" + t.Name + "] Status [" + t.Status + "]"
}

func (doc *SuccessDocument) NewOrderHandler() Event {
	doc.Status = "Started"
	return START
}

func (doc *SuccessDocument) PlaceOrderHandler() Event {
	doc.Status = "Ordered"
	return ORDERED
}

func (doc *SuccessDocument) CapturePaymentHandler() Event {
	doc.Status = "Paid"
	return PAID
}

type FailDocument struct {
	Name   string
	Status string
}

func (t *FailDocument) String() string {
	return "Name [" + t.Name + "] Status [" + t.Status + "]"
}

func (doc *FailDocument) NewOrderHandler() Event {
	doc.Status = "Started"
	return START
}

func (doc *FailDocument) PlaceOrderHandler() Event {
	doc.Status = "Ordered"
	return ORDERED
}

func (doc *FailDocument) CapturePaymentHandler() Event {
	doc.Status = "Payemtn not authorized"
	return PAYMENT_FAILED
}

func TestOrderWorkflow(t *testing.T) {
	// Setup the workflow
	wf := NewWorkflow()
	wf.AddHandler(NEW, Method("NewOrderHandler"))
	wf.AddRoute(START, PLACE_ORDER)

	wf.AddHandler(PLACE_ORDER, Method("PlaceOrderHandler"))
	wf.AddRoute(ORDERED, CAPTURE_PAYMENT)

	wf.AddHandler(CAPTURE_PAYMENT, Method("CapturePaymentHandler"))
	wf.AddRoute(PAID, COMPLETE)
	wf.AddRoute(PAYMENT_FAILED, FAIL)

	// test the success branch
	td1 := &SuccessDocument{Name: "Test Doc", Status: "Beginning"}
	wf.Run(td1, NEW)

	if td1.Status != "Paid" {
		t.Errorf("Expecting Status == Paid got %v", td1.Status)
	}

	if wf.Step() != COMPLETE {
		t.Errorf("Expecting Step == COMPLETE got %v", wf.Step())
	}

	// test eh failure branch
	td2 := &FailDocument{Name: "2nd Doc", Status: "Beginning again"}
	wf.Run(td2, NEW)

	if wf.Step() != FAIL {
		t.Errorf("Expecting Step == FAIL got %v", wf.Step())
	}

}
