package unused

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type CustomDB interface {
	Save(item interface{})
	Load() interface{}
}

type CustomDBImpl struct {
	hhhhh int
}

func (cdbi CustomDBImpl) Save(item interface{}) {
	fmt.Printf("%T:%v\n", item, item)

}

func (cdbi CustomDBImpl) Load() interface{} {
	return nil
}

type FFF struct {
	gg int
	pp int
}

func (cdbi * CustomDBImpl) Init(){
	var customDB CustomDB
	customDB = &CustomDBImpl{33221}

	// unsafe {{{
	db, err := sql.Open("mysql", "root:qwer1234@/")
	db.Exec("CREATE DATABASE IF NOT EXISTS gocrawlerdemo")
	db.Exec("USE gocrawlerdemo")
	db.Exec("CREATE TABLE IF NOT EXISTS `users` (`name` VARCHAR(128) NOT NULL,PRIMARY KEY (`name`))")
	db.Exec("CREATE TABLE IF NOT EXISTS `rels` (  `idrels` INT NOT NULL,  `followee` VARCHAR(128) NOT NULL,  `follower` VARCHAR(128) NOT NULL,  PRIMARY KEY (`idrels`),  INDEX `fk_rels_1_idx` (`followee` ASC),  INDEX `fk_rels_2_idx` (`follower` ASC),  CONSTRAINT `fk_rels_1`    FOREIGN KEY (`followee`)    REFERENCES `users` (`name`)    ON DELETE NO ACTION    ON UPDATE NO ACTION,  CONSTRAINT `fk_rels_2`    FOREIGN KEY (`follower`)    REFERENCES `users` (`name`)    ON DELETE NO ACTION    ON UPDATE NO ACTION)")
	// }}}unsafe

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	customDB.Save(FFF{999, 234})
}

