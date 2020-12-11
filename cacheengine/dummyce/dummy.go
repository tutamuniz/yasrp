package dummyce

import (
	"fmt"
	"log"
)

type DummyCE struct{}

func NewDummyCE() (*DummyCE, error) {
	return &DummyCE{}, nil
}

func (_ DummyCE) InCache(_ string) bool {
	return false
}

func (_ DummyCE) Get(_ string) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (_ DummyCE) Put(string, []byte) error {
	return fmt.Errorf("Not implemented")

}
func (_ DummyCE) PutTTL(string, []byte) error {
	return fmt.Errorf("Not implemented")

}
func (_ DummyCE) StartEngine() {
	log.Printf("Starting Dummy Cache Engine\n")
}
