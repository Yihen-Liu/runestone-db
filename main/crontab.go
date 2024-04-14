package main

import (
	"galactic-monitor/log"
	"galactic-monitor/models"

	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

/*
# 每天定时触发,支持linux crontab语法 @daily is 0 0 * * *
BscSync = @daily
TMSync = @daily
UpdateBscBalance = @every 10s
HandleFailBlock = @every 1m
HandleFailTX = @every 1m
HandleFailTXCount = @every 2m
HandleFailCalTransactionName = @every 2m
HandleTxnumstatisticsAllAddress = 00 23 * * *
HandleShardingInfoTsZero = @every 2m
HandleEsFailBlock = @every 2m
HandleFailEsTX = @every 2m
HandleTxTokenName = @every 2m
HandleToConfirmEsBlock = @every 1m
ReportGorountineNum = @every 2m
HandlePaymentFailBlock = @every 2m
HandleFailMqTX = @every 3m
HandleFailedMqBatchTX = @every 2m
HandleCountTableEqualFalse = @every 1h
GetPaymentCountLastIntervalTime = @every 10m
*/
func cronTask() {
	c := cron.New()

	updateUserBalance := "@every 10s"

	c.AddFunc(updateUserBalance, func() {
		users, err := models.QueryAllUsers()
		if err != nil {
			return
		}

		for _, user := range users {
			if err, amount := models.SumTokenAmount("ETH", user.UserAddress); err == nil {
				models.UpdateUserDepositAmount("ETH", user.UserAddress, amount)
			}

			if err, amount := models.SumTokenAmount("USDC", user.UserAddress); err == nil {
				models.UpdateUserDepositAmount("USDC", user.UserAddress, amount)
			}

			if err, amount := models.SumTokenAmount("USDT", user.UserAddress); err == nil {
				models.UpdateUserDepositAmount("USDT", user.UserAddress, amount)
			}

			if err, amount := models.SumTokenAmount("Trias", user.UserAddress); err == nil {
				models.UpdateUserDepositAmount("Trias", user.UserAddress, amount)
			}
			activeUserStatus(user.UserAddress)
		}
		log.Info("crontab task")
	})

	c.Start()

}

func activeUserStatus(user string) {
	registration, err := models.QueryInactiveUser(user)
	if err != nil {
		return
	}

	triasAmount := decimal.NewFromFloat(registration.TRIASCount).Mul(decimal.NewFromInt32(15))
	ethAmount := decimal.NewFromFloat(registration.ETHCount).Mul(decimal.NewFromInt32(4000))
	usdcAmount := decimal.NewFromFloat(registration.USDCCount)
	usdtAmount := decimal.NewFromFloat(registration.USDTCount)
	total := triasAmount.Add(ethAmount).Add(usdcAmount).Add(usdtAmount)

	if total.Cmp(decimal.NewFromInt32(100)) > 0 {
		models.UpdateUserStatus(user)
	}
}
