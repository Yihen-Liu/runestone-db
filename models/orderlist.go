package models

import "time"

type Etching struct {
	ID                int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Owner             string    `gorm:"column:owner;" json:"owner"`
	RuneId            string    `gorm:"column:runeid;not null" json:"runeid"`
	RuneName          string    `gorm:"column:rune_name;not null; unique" json:"rune_name"`
	BlockHeight       int32     `gorm:"column:block_height;not null" json:"block_height"`
	TxIndex           int32     `gorm:"column:tx_index;not null" json:"tx_index"`
	BlockHash         string    `gorm:"column:block_hash;not null" json:"block_hash"`
	PreviousBlockHash string    `gorm:"column:previous_block_hash;not null" json:"previous_block_hash"`
	Txid              string    `gorm:"column:txid;not null" json:"txid"`
	Valid             bool      `gorm:"column:valid;not null; default:true" json:"valid"`
	Divisibility      int32     `gorm:"column:divisibility;not null; default:0" json:"divisibility"`
	Premine           string    `gorm:"column:premine;not null; default:''" json:"premine"`
	Spacers           string    `gorm:"column:spacers;not null; default:''" json:"spacers"`
	Symbol            string    `gorm:"column:symbol;not null; default:0" json:"symbol"`
	Capacity          string    `gorm:"column:capacity;not null; default:0" json:"capacity"`
	MintAmount        int64     `gorm:"column:mint_amount;not null; default:0" json:"mint_amount"`
	Offset            int32     `gorm:"column:offset;not null; default:0" json:"offset"`
	StartHeight       int32     `gorm:"column:start_height;not null; default:0" json:"start_height"`
	EndHeight         int32     `gorm:"column:end_height;not null; default:0" json:"end_height"`
	UpdateAt          time.Time `gorm:"column:update_at; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP; autoUpdateTime" json:"update_at,omitempty"`
	CreateAt          time.Time `gorm:"column:create_at; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;<-:create" json:"create_at,omitempty"`
}

// TableName Orderlist's table name
func (*Etching) TableName() string {
	return "etching"
}

type RuneBalance struct {
	ID                int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Address           string    `gorm:"column:address;not null" json:"address"`
	RuneId            string    `gorm:"column:runeid;not null" json:"runeid"`
	RuneName          string    `gorm:"column:rune_name;not null" json:"rune_name"`
	Balance           int64     `gorm:"column:balance;not null" json:"balance"`
	Txid              string    `gorm:"column:txid;not null" json:"txid"`
	Vout              int32     `gorm:"column:vout;not null" json:"vout"`
	BlockHeight       int32     `gorm:"column:block_height;not null" json:"block_height"`
	BlockHash         string    `gorm:"column:block_hash;not null" json:"block_hash"`
	PreviousBlockHash string    `gorm:"column:previous_block_hash;not null" json:"previous_block_hash"`
	LastUpdated       time.Time `gorm:"column:last_updated; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP; autoUpdateTime" json:"last_updated,omitempty"`
	CreateAt          time.Time `gorm:"column:create_at; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;<-:create" json:"create_at,omitempty"`
}

// TableName Orderlist's table name
func (*RuneBalance) TableName() string {
	return "runebalance"
}

type Blockchain struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Label       string    `gorm:"column:label;not null" json:"label"`
	Address     string    `gorm:"column:address;not null" json:"address"`
	ChainId     int32     `gorm:"column:chain_id;comment:状态" json:"chain_id"`
	LastUpdated time.Time `gorm:"column:last_updated; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP; autoUpdateTime" json:"last_updated,omitempty"`
}

// TableName Orderlist's table name
func (*Blockchain) TableName() string {
	return "blockchain"
}

// Orderlist mapped from table <orderlist>
type Order struct {
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	InitiatorAddress string    `gorm:"column:initiator_address;not null" json:"initiator_address"`
	ReceiverAdddress string    `gorm:"column:receiver_address;not null" json:"receiver_address"`
	Type             string    `gorm:"column:type;comment:类型"  json:"type"`
	Amount           float64   `gorm:"column:amount;not null; comment: 数额" json:"amount"`
	Status           string    `gorm:"column:status;comment:状态" json:"status"`
	InsertTime       time.Time `gorm:"column:insert_time; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;<-:create" json:"insert_time,omitempty"`
	TxHash           string    `gorm:"column:txhash;not null;unique" json:"txhash"`
	ChainId          int32     `gorm:"column:chain_id;comment:链Id" json:"chain_id"`
}

// TableName Orderlist's table name
func (*Order) TableName() string {
	return "orders"
}

type Registration struct {
	ID                     int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserAddress            string    `gorm:"column:user_address;not null" json:"user_address"`
	IsActive               int32     `gorm:"column:is_active" json:"is_active"`
	PersonalScore          int32     `gorm:"column:personal_score" json:"personal_score"`
	InvitedIds             int32     `gorm:"column:invited_ids" json:"invited_ids"`
	DepositAmount          float64   `gorm:"column:deposit_amount" json:"deposit_amount"`
	PersonalInvitationCode string    `gorm:"column:personal_invitation_code" json:"personal_invitation_code"`
	ParentUserAddress      string    `gorm:"column:parent_user_address" json:"parent_user_address"`
	CreateAt               time.Time `gorm:"column:create_at; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;<-:create" json:"create_at,omitempty"`
	UpdateAt               time.Time `gorm:"column:update_at; type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP; autoUpdateTime" json:"update_at,omitempty"`
	TwitterAddress         string    `gorm:"column:twitter_address" json:"twitter_address"`
	ParentUserCode         string    `gorm:"column:parent_user_code" json:"parent_user_code"`
	TwitterScreenName      string    `gorm:"column:twitter_screen_name" json:"twitter_screen_name"`
	USDCCount              float64   `gorm:"column:usdccount" json:"usdccount"`
	USDTCount              float64   `gorm:"column:usdtcount" json:"usdtcount"`
	ETHCount               float64   `gorm:"column:ethcount" json:"ethcount"`
	TRIASCount             float64   `gorm:"column:triascount" json:"triascount"`
}

func (*Registration) TableName() string {
	return "registration"
}