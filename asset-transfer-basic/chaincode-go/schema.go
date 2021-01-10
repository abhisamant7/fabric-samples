package main

const (
	ORDER_ASSET_RECORD_TYPE string = "OrderAsset"
	FARMERTOCAS             uint64 = 1
	FARMERTOPC              uint64 = 2
	CASTOPC                 uint64 = 3
	FARMERTODC              uint64 = 4
	CASTODC                 uint64 = 5
	PCTODC                  uint64 = 6
	FARMERTOBUYER           uint64 = 7
	CASTOBUYER              uint64 = 8
	PCTOBUYER               uint64 = 9
	DCTOBUYER               uint64 = 10
	FARMERTOCC              uint64 = 11
	CCTOPH                  uint64 = 12
	CCTOPC                  uint64 = 13
	CCTOCAS                 uint64 = 14
)

type AssetData struct {
	OrderID string `json:"ORDERID,omitempty"`
}
type AssetDataChild struct {
	ChildOrderID string `json:"ChildOrderID,omitempty"`
}
type Trucker struct {
	FullName      string `json:"fullName,omitempty"`
	LicensePlate  string `json:"licensePlate,omitempty"`
	DriverLicense string `json:"driverLicense,omitempty"`
	Type          string `json:"type,omitempty"`
	Size          uint64 `json:"size,omitempty"`
}

type ReceiveData struct {
	TruckID          uint64   `json:"truckId,omitempty"`
	ActualQty        uint64   `json:"actualQty,omitempty"`
	AgreeToActualQty bool     `json:"agreeToActualQty,omitempty"`
	CSQty            uint64   `json:"csQty,omitempty"`
	Room             uint64   `json:"room,omitempty"`
	Bins             []string `json:"bins,omitempty"`
}

type TransportaionoFCAS struct {
	TruckID         uint64  `json:"truckId,omitempty"`
	OrderBy         string  `json:"orderBy,omitempty"`
	PickUpDate      string  `json:"pickUpDate,omitempty"`
	PickUpLocation  string  `json:"pickUpLocation,omitempty"`
	DropOffLocation string  `json:"dropOffLocation,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Insurance       string  `json:"insurance,omitempty"`
	Truck           Trucker `json:"trucker,omitempty"`
	Status          string  `json:"STATUS,omitempty"`
	ActualQty       uint64  `json:"actualQty,omitempty"`
	PickUpTime      string  `json:"pickUpTime,omitempty"`
}

type Unit struct {
	Name   string  `json:"NAME,omitempty"`
	Weight float64 `json:"WEIGHT,omitempty"`
}
type DistanceUnit struct {
	Name  string `json:"NAME,omitempty"`
	Value uint64 `json:"VALUE,omitempty"`
}

//Need to Discuss about Flow
// Coldstorage
type CSOrderAsset struct {
	DocType                  string               `json:"docType,omitempty"`
	OrderType                uint32               `json:"ORDERTYPE,omitempty"`
	OrderID                  string               `json:"ORDERID,omitempty"`
	OrderUnixTime            int64                `json:"orderUnixTime,omitempty"`
	SourceID                 string               `json:"SOURCEID,omitempty"`
	ProduceID                string               `json:"PRID,omitempty"`
	Variety                  string               `json:"variety,omitempty"`
	HarvestDate              string               `json:"HARVESTDATE,omitempty"`
	Qty                      uint64               `json:"QTY,omitempty"`
	DestinationID            string               `json:"DESTINATIONID,omitempty"`
	StorageFacility          string               `json:"STORAGE_FACILITY,omitempty"`
	QtyforSale               uint64               `json:"Quantity_of_sale,omitempty"`
	Status                   string               `json:"STATUS,omitempty"` //Order Requested from Farmer,Order approved from CAS,Picked up by trucker,in Transit,Delivered to CAS
	Transports               []TransportaionoFCAS `json:"Transportaion,omitempty"`
	Receives                 []ReceiveData        `json:"receive,omitempty"`
	ActualQty                uint64               `json:"TotalActualQuantity,omitempty"`
	TotaltransportationPrice float64              `json:"TotalTransportationCost,omitempty"`
	BaseUnit                 string               `json:"BASE_UNIT,omitempty"`
	SenderUnit               Unit                 `json:"SENDER_UNIT,omitempty"`
	Country                  string               `json:"COUNTRY,omitempty"`
}

//DC

type Transportaion struct {
	TruckID         uint64  `json:"truckId,omitempty"`
	OrderBy         string  `json:"orderBy,omitempty"`
	PickUpDate      string  `json:"pickUpDate,omitempty"`
	PickUpLocation  string  `json:"pickUpLocation,omitempty"`
	DropOffLocation string  `json:"dropOffLocation,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Insurance       string  `json:"insurance,omitempty"`
	Truck           Trucker `json:"trucker,omitempty"`
	Status          string  `json:"STATUS,omitempty"`
	PickUpTime      string  `json:"pickUpTime,omitempty"`
	ActualQty       uint64  `json:"actualQty,omitempty"`
	ReceiveDate     string  `json:"receiveDate,omitempty"`
}

type TableVariety struct {
	Name            string  `json:"name,omitempty"`
	Value           float64 `json:"quantity,omitempty"`
	SelecteUnitdata Unit    `json:"selected_unit,omitempty"`
}

type ParentIDInfo struct {
	ProduceID      string         `json:"PRID,omitempty"`
	TableVarieties []TableVariety `json:"QTY,omitempty"`
	ParentOrderId  string         `json:"ParentOrderId,omitempty"`
	ChildOrderID   string         `json:"ChildOrderID,omitempty"`
}

type DCOrder struct {
	SourceID                 string          `json:"SOURCEID,omitempty"`
	DestinationID            string          `json:"DESTINATIONID,omitempty"`
	OrderType                uint32          `json:"ORDERTYPE,omitempty"`
	Produce                  string          `json:"Produce,omitempty"`
	Variety                  string          `json:"Variety,omitempty"`
	BaseUnit                 string          `json:"BASE_UNIT,omitempty"`
	RequiredDate             string          `json:"requiredDate,omitempty"`
	Qty                      uint64          `json:"QTY,omitempty"`
	OrderID                  string          `json:"ORDERID,omitempty"`
	Status                   string          `json:"STATUS,omitempty"`
	Transports               []Transportaion `json:"Transportaion,omitempty"`
	ParentInfo               []ParentIDInfo  `json:"bcParentorderinfo,omitempty"`
	TotaltransportationPrice float64         `json:"totalTransportationcost,omitempty"`
}

type DCOrderAsset struct {
	DocType                  string          `json:"docType,omitempty"`
	OrderType                uint32          `json:"ORDERTYPE,omitempty"`
	SourceID                 string          `json:"SOURCEID,omitempty"`
	DestinationID            string          `json:"DESTINATIONID,omitempty"`
	OrderUnixTime            int64           `json:"orderUnixTime,omitempty"`
	Produce                  string          `json:"Produce,omitempty"`
	Variety                  string          `json:"variety,omitempty"`
	RequiredDate             string          `json:"requiredDate,omitempty"`
	Qty                      uint64          `json:"QTY,omitempty"`
	ProduceID                string          `json:"PRID,omitempty"`
	TableVarieties           []TableVariety  `json:"tableVarietyDC,omitempty"`
	ParentOrderID            string          `json:"ParentOrderId,omitempty"`
	OrderID                  string          `json:"ORDERID,omitempty"`
	ChildOrderID             string          `json:"ChildOrderID,omitempty"`
	Status                   string          `json:"STATUS,omitempty"`
	Transports               []Transportaion `json:"Transportaion,omitempty"`
	TotaltransportationPrice float64         `json:"totalTransportationCost,omitempty"`
}

//PC
type TransportaionoFPC struct {
	TruckID    uint64 `json:"truckId,omitempty"`
	OrderBy    string `json:"orderBy,omitempty"`
	PickUpDate string `json:"pickUpDate,omitempty"`
	// PickUpLocation  string  `json:"pickUpLocation,omitempty"`
	DropOffLocation string  `json:"dropOffLocation,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Insurance       string  `json:"insurance,omitempty"`
	Truck           Trucker `json:"trucker,omitempty"`
	Status          string  `json:"STATUS,omitempty"`
	ActualQty       uint64  `json:"actualQty,omitempty"` // need to set data type
	PickUpTime      string  `json:"pickUpTime,omitempty"`
	ReceiveDate     string  `json:"receiveDate,omitempty"`
}

type Processe struct {
	Qty uint64 `json:"qty,omitempty"`
}

type Wastage struct {
	Qty uint64 `json:"qty,omitempty"`
}

//change in mongodb
//add  TableVariety to Order
type ProcessData struct {
	TableVarieties []TableVariety `json:"tableVariety,omitempty"`
	Processed      []TableVariety `json:"processVariety,omitempty"`
	Wastages       Wastage        `json:"wastages,omitempty"`
}

type PCOrderAsset struct {
	DocType                  string              `json:"docType,omitempty"`
	OrderType                uint32              `json:"ORDERTYPE,omitempty"`
	SourceID                 string              `json:"SOURCEID,omitempty"`
	DestinationID            string              `json:"DESTINATIONID,omitempty"`
	OrderUnixTime            int64               `json:"orderUnixTime,omitempty"`
	ProduceID                string              `json:"PRID,omitempty"`
	Variety                  string              `json:"variety,omitempty"`
	Qty                      uint64              `json:"QTY,omitempty"`
	OrderID                  string              `json:"ORDERID,omitempty"`
	Status                   string              `json:"STATUS,omitempty"`
	Transports               []TransportaionoFPC `json:"Transportaion,omitempty"`
	ProcessProduce           ProcessData         `json:"processData,omitempty"`
	ParentID                 string              `json:"parentorderId,omitempty"`
	TotaltransportationPrice float64             `json:"totalTransportationcost,omitempty"`
	TotalavailableQty        []TableVariety      `json:"TotalavailableQty,omitempty"` // added new
	Baseunit                 string              `json:"BASE_UNIT,omitempty"`
}

type BuyerOrderAsset struct {
	DocType          string            `json:"docType,omitempty"`
	OrderType        uint32            `json:"ORDERTYPE,omitempty"`
	DCID             string            `json:"SOURCEID,omitempty"`
	OrderUnixTime    int64             `json:"orderUnixTime,omitempty"`
	ProduceID        string            `json:"PRID,omitempty"`
	Variety          string            `json:"variety,omitempty"`
	Qty              uint64            `json:"totalQTY,omitempty"`
	BuyerID          string            `json:"DESTINATIONID,omitempty"`
	TotalPrice       float64           `json:"price,omitempty"` // total price
	QtyBreakdowns    []TableVariety    `json:"QTY,omitempty"`
	ParticipantInfos []ParticipantInfo `json:"PARTICIPANTINFOS,omitempty"`
	OrderID          string            `json:"ORDERID,omitempty"`
	Status           string            `json:"STATUS,omitempty"`
	ParentID         string            `json:"parentorderId,omitempty"` // NEED To discuss with prasanna
	//TransportRate    Rate              `json:"TRANSPORTRATE,omitempty"`
	PayableAssetID string `json:"payableAssetID,omitempty"`
}

//Rate for any kind of produce varity quality rate.

type QualityPriceRate struct {
	QualityName string `json:"QUALITYNAME,omitempty"`
	Rates       []Rate `json:"RATES,omitempty"`
}

// From here CC order Datastructure added
type ProducelistData struct {
	ItemId          int     `json:"itemId,omitempty"`
	ID              string  `json:"MAINID,omitempty"`
	OrderBy         string  `json:"orderBy,omitempty"`
	PickUpLocation  string  `json:"pickUpLocation,omitempty"`
	DropOffLocation string  `json:"dropOffLocation,omitempty"`
	Price           string  `json:"price,omitempty"`
	Insurance       string  `json:"insurance,omitempty"`
	BARCODEID       string  `json:"BARCODEID,omitempty"`
	ActualQty       string  `json:"actualQty,omitempty"`
	RejectedQty     string  `json:"rejectedQty,omitempty"`
	RecevieDate     string  `json:"RECEVIE_DATE,omitempty"`
	ReceviveTime    string  `json:"RECEVIVE_TIME,omitempty"`
	ReceiveStatus   string  `json:"RECEIVE_STATUS,omitempty"`
	Truck           Trucker `json:"trucker,omitempty"`
}

type CCOrderAsset struct {
	DocType             string            `json:"docType,omitempty"`
	OrderType           uint32            `json:"ORDERTYPE,omitempty"`
	OrderID             string            `json:"ORDERID,omitempty"`
	SourceID            string            `json:"SOURCEID,omitempty"`
	ProduceID           string            `json:"PRID,omitempty"`
	OrderUnixTime       int64             `json:"orderUnixTime,omitempty"`
	Produce             string            `json:"produce,omitempty"`
	Variety             string            `json:"variety,omitempty"`
	HarvestDate         string            `json:"HARVESTDATE,omitempty"`
	Qty                 uint64            `json:"QTY,omitempty"`
	DestinationID       string            `json:"DESTINATIONID,omitempty"`
	Status              string            `json:"STATUS,omitempty"` //Order Requested from Farmer,Order approved from CAS,Picked up by trucker,in Transit,Delivered to CAS
	ReceiveType         string            `json:"RECEIVE_TYPE,omitempty"`
	BaseUnit            string            `json:"BASE_UNIT,omitempty"`
	SenderUnit          Unit              `json:"SENDER_UNIT,omitempty"`
	Country             string            `json:"COUNTRY,omitempty"`
	State               string            `json:"SOURCE_STATE,omitempty"`
	TotalavailableQty   uint32            `json:"TotalavailableQty,omitempty"` // added new
	ActualHarvest       string            `json:"ACTUAL_HARVEST,omitempty"`
	TotalActualQuantity uint32            `json:"TotalActualQuantity,omitempty"` // added new
	ParticipantInfos    []ParticipantInfo `json:"PARTICIPANTINFOS,omitempty"`
	Producelist         []ProducelistData `json:"PRODUCELIST,omitempty"`
}

//PHOrderAsset Order info for pack house
type PHOrderAsset struct {
	DocType       string `json:"docType,omitempty"`
	OrderUnixTime int64  `json:"orderUnixTime,omitempty"`
	SourceID      string `json:"SOURCEID,omitempty"`
	DestinationID string `json:"DESTINATIONID,omitempty"`
	OrderID       string `json:"ORDERID,omitempty"`
	OrderType     uint64 `json:"ORDERTYPE,omitempty"`
	ProduceID     string `json:"PRID,omitempty"`
	Variety       string `json:"variety,omitempty"`
	Qty           uint64 `json:"QTY,omitempty"`
	BaseUnit      string `json:"BASE_UNIT,omitempty"`
	SenderUnit    Unit   `json:"SENDER_UNIT,omitempty"`
	Status        string `json:"STATUS,omitempty"`
}
