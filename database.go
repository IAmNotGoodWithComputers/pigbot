package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Database struct {

}

type CooldownPolicy struct {
	Command string `db:"command"`
	Cooldown int `db:"cooldown"`
}

type Blacklist struct {
	UserId string `db:"user_id"`
}

type CommandPolicy struct {
	GuildId string `db:"guild_id"`
	Allowed int `db:"allowed"`
	Command string `db:"command"`
}

type FourchanPost struct {
	DateTime string `db:"date_time"`
	PostId int `db:"post_id"`
	Content string `db:"content"`
	FileUrl string `db:"file_url"`
}

type AdminUser struct {
	UserId string `db:"user_id"`
}

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Open("sqlite3", "pigbot.db")
	if err != nil {
		log.Fatal("can't open db", err.Error())
	}

	mustExec(`CREATE TABLE IF NOT EXISTS cooldown_policy (
  		 command TEXT NOT NULL,
  		 cooldown INT NOT NULL
	);`)

	mustExec(`CREATE TABLE IF NOT EXISTS blacklist (
  		user_id TEXT NOT NULL
	)`)

	mustExec(`CREATE TABLE IF NOT EXISTS command_policy (
  		guild_id TEXT NOT NULL,
		command TEXT NOT NULL,
		allowed INT NOT NULL,

		CONSTRAINT unique_policy UNIQUE(guild_id, command)
	)`)

	mustExec(`CREATE TABLE IF NOT EXISTS fourchan_post (
		post_id INT PRIMARY KEY NOT NULL,
  		date_time TEXT NOT NULL,
		content TEXT NOT NULL,
		file_url TEXT NOT NULL
	)`)

	mustExec(`CREATE TABLE IF NOT EXISTS admin_user (
		user_id TEXT NOT NULL,

		CONSTRAINT unique_admin UNIQUE(user_id)
	)`)

	mustExec("REPLACE INTO admin_user (user_id) VALUES ('157636475117240320')",)
	mustExec("REPLACE INTO admin_user (user_id) VALUES ('79921519756574720')",)
	mustExec("REPLACE INTO admin_user (user_id) VALUES ('168425277821681664')",)
	mustExec("REPLACE INTO admin_user (user_id) VALUES ('168473459838550016')",)
	mustExec("REPLACE INTO admin_user (user_id) VALUES ('193731247091089408')",)
}

func mustExec(query string) {
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("can not exec: ", query, ": ", err.Error())
		os.Exit(-1)
	}
}

func CommandIsAllowed(command, guild string) bool {
	entity := CommandPolicy{}
	err := db.Get(&entity, "SELECT * FROM command_policy WHERE command = $1 AND guild_id = $2", command, guild)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			fmt.Println("CommandIsAllowed error: ", err.Error())
		}
		return true
	}
	if entity.Allowed == 1 {
		return true
	}
	return false
}

func SetCommandPolicy(command, guild string, allowed int) {
	db.Exec("REPLACE INTO command_policy (guild_id, command, allowed) VALUES ($1, $2, $3)", guild, command, allowed)
}

func UserIsBlacklisted(userid string) bool {
	entity := Blacklist{}
	err := db.Get(&entity, "SELECT * FROM blacklist WHERE user_id = $1", userid)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			fmt.Println("UserIsBlacklisted error: ", err.Error())
		}
		return false
	}
	return true
}

func ToggleUserBlacklist(userid string) bool {
	if UserIsBlacklisted(userid) {
		_, err := db.Exec("DELETE FROM blacklist WHERE user_id = $1", userid)
		return err == nil
	} else {
		_, err := db.Exec("INSERT INTO blacklist (user_id) VALUES ($1)", userid)
		return err == nil
	}
}

func UserIsAdmin(userid string) bool {
	entity := AdminUser{}
	err := db.Get(&entity, "SELECT * FROM admin_user WHERE user_id = $1", userid)
	if err != nil {
		fmt.Println("UserIsAdmin error: ", err.Error())
		return false
	}
	return true
}