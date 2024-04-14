package models

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const TX_STATUS_MINING = "mining"
const TX_STATUS_MINED = "mined"
const TX_STATUS_TODETERMINE = "todetermine"

const ORDER_STATUS_PENDING = "pending"
const ORDER_STATUS_DISTRIBUTE = "distribute"
const ORDER_STATUS_COMPLETED = "completed"
const ORDER_STATUS_UNISATVERIFY = "unisatVerify"
const ORDER_STATUS_FAILED = "failed"

const ORDER_VERIFIED_TOBE = "tobeVerified"
const ORDER_VERIFIED_FAILED = "failed"
const ORDER_VERIFIED_SUCCESS = "success"

const TOKEN_TYPE_FT = "FT"
const TOKEN_TYPE_NFT = "NFT"

const ORDER_STAGE_WHITELIST = "whitelist"
const ORDER_STAGE_PUBLIC = "public"

const PERIOD_PAYMENT = 0
const PERIOD_COMMIT = 1
const PERIOD_REVEL = 2
const PERIOD_TRANSFER = 3

const ADDRESS_TYPE_FROM = "from地址"
const ADDRESS_TYPE_FUND = "fund地址"
const ADDRESS_TYPE_TRANSMIT = "中转地址"
const ADDRESS_TYPE_RECEIVE = "用户接收地址"
const ADDRESS_TYPE_MINTER = "资产铸造发起地址"

func Db() *gorm.DB {
	return db
}

/*
	func QueryReceiveAddressFromOrderlist(orderId string, tokenType string) (error, Orderlist) {
		var order Orderlist

		err := db.Model(&Orderlist{}).Where("\"orderId\" = ? and \"type\"=? ", orderId, tokenType).First(&order).Error
		if err != nil {
			return err, Orderlist{}
		}

		return nil, order
	}
*/
func QueryChainByReceiver(address string, chainId int32) (Blockchain, error) {
	var blockchain Blockchain
	err := db.Model(&Blockchain{}).Where(" \"address\" = ? and \"chain_id\" = ?", address, chainId).Find(&blockchain).Error
	return blockchain, err
}

func QueryInactiveUser(userAddr string) (Registration, error) {
	var user Registration
	err := db.Model(&Registration{}).Where(" \"is_active\" = ? and \"user_address\"=? ", 2, userAddr).First(&user).Error
	if err != nil {
		return Registration{}, err
	}
	return user, nil
}

func UpdateUserStatus(user string) error {
	if err := db.Model(&Registration{}).Where("\"user_address\" = ?", user).UpdateColumn("is_active", 1).Error; err != nil {
		return err
	}
	return nil
}

func QueryPendingOrders() ([]*Order, error) {
	var orders []*Order
	err := db.Model(&Order{}).Where(" \"status\" = ?", "pending").Find(&orders).Error
	return orders, err
}

func QueryAllUsers() ([]*Registration, error) {
	var users []*Registration
	err := db.Model(&Registration{}).Find(&users).Error
	return users, err

}

func UpdateOrderStatus(txHash string, status string) error {
	if err := db.Model(&Order{}).Where("\"txhash\" = ?", txHash).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

func SumTokenAmount(token string, owner string) (error, float64) {
	var sum decimal.Decimal
	var result []float64
	err := db.Model(&Order{}).Where("\"type\"=? and  \"status\"= ? and \"initiator_address\" = ? ", token, "success", owner).Pluck("\"amount\"", &result).Error
	for _, v := range result {
		sum = sum.Add(decimal.NewFromFloat(v))
	}

	floatSum, _ := sum.Float64()
	return err, floatSum
}

func UpdateUserDepositAmount(token string, owner string, amount float64) error {
	column := ""
	switch token {
	case "USDC":
		column = "usdccount"
	case "ETH":
		column = "ethcount"
	case "Trias":
		column = "triascount"
	case "USDT":
		column = "usdtcount"
	default:
		return errors.New("token name error")
	}
	if err := db.Model(&Registration{}).Where("\"user_address\" = ?", owner).UpdateColumn(column, amount).Error; err != nil {
		return err
	}
	return nil
}

// ///////////////////////////////////////////////////////
