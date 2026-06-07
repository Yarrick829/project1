package domain

import "time"

type Device struct {
	Id              uint64
	OrganizationId  uint64
	RoomId          *uint64
	GUID            string
	InventoryNumber string
	SerialNumber    string
	Characteristic  string
	Category        DeviceCathegory
	Units           *string
	PowerConsumtion *float64
	CreatedDate     time.Time
	UpdatedDate     time.Time
	DeletedDate     *time.Time
	Measurements    []Measurement
}

type DeviceCathegory string

const (
	Sensor   DeviceCathegory = "SENSOR"
	Actuator DeviceCathegory = "ACTUATOR"
)
