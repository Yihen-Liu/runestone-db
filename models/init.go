package models

func InitTables() {
	if db.Migrator().HasTable(&RuneBalance{}) == false {
		if err := db.Migrator().CreateTable(&RuneBalance{}); err != nil {
			panic("create runebalance table err:" + err.Error())
		}
	}

	if db.Migrator().HasTable(&Etching{}) == false {
		if err := db.Migrator().CreateTable(&Etching{}); err != nil {
			panic("create etching table err:" + err.Error())
		}
	}

	if db.Migrator().HasTable(&BlockInfo{}) == false {
		if err := db.Migrator().CreateTable(&BlockInfo{}); err != nil {
			panic("create blockinfo table err:" + err.Error())
		}
	}
}
