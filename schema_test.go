package sqlofi

import (
	"database/sql"
	"fmt"
	"testing"
)

var (
	aiconf struct {
		ID     int64         `sql:"primary key;autoincrement"`
		Apikey sql.NullInt64 `sql:"foreign key:apikey (id)"`
		Model  string        `sql:"NOT NULL"`
		Url    string        `sql:"NOT NULL"`
		Airole string        `sql:"NOT NULL"`
		ChatID sql.NullInt64 `sql:"foreign key:chat (id)"`
	}
)

type AiConfig struct {
	ID     int64         `sql:"primary key;autoincrement"`
	Apikey sql.NullInt64 `sql:"foreign key:apikey (id)"`
	Model  string        `sql:"NOT NULL"`
	Url    string        `sql:"NOT NULL"`
	Airole string        `sql:"NOT NULL"`
	ChatID sql.NullInt64 `sql:"foreign key:chat (id)"`
}

type Apikey struct {
	ID     int64  `sql:"primary key;autoincrement"`
	Label  string `sql:"NOT NULL"`
	Apikey string `sql:"NOT NULL;UNIQUE"`
}

type Chat struct {
	ID    int64  `sql:"primary key;autoincrement"`
	Label string `sql:"NOT NULL"`
}

type Message struct {
	ID       int64         `sql:"primary key;autoincrement"`
	ChatID   sql.NullInt64 `sql:"foreign key:chat (id)"`
	Role     string        `sql:"NOT NULL"`
	Message  string        `sql:"NOT NULL"`
	Included sql.NullBool  `sql:"DEFAULT:true"`
}

func testGeneraingSchema(t *testing.T) {
	ddl, err := CreateSchema(aiconf, Apikey{}, Chat{}, Message{})
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	for _, table := range ddl.Tables {
		fmt.Println(table.String())
	}
}
