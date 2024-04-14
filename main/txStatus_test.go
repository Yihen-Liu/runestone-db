package main

import (
	"galactic-monitor/models"
	"testing"
)

func TestTxStatus(t *testing.T) {
	txhash := "0x5171910dcbc50f7299f15a20d8474a8123b23fdd2e8dc851deb97ceaa0ed9b2b"
	chain := models.Blockchain{Label: "Bsc"}
	queryBscTxStatus(txhash, chain)
}
