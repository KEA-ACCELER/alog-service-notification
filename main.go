package main

import (
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddMessageTmp struct {
	UserPk     int64
	MsgContent string
}

type SendMessageTmp struct {
	UserPk     int64
	MsgContent string
	Datetime   string
	IsChecked  bool
}

func main() {
	app := fiber.New()

	// 실라 디비와 연결하기
	cluster := gocql.NewCluster(os.Getenv("DB_HOST"))
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// alog라는 keyspace가 있는지 확인하고 없으면 만듭니다.
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS alog WITH REPLICATION={'class' : 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		log.Fatal(err)
	}

	// noti라는 table이 있는지 확인하고 없으면 만듭니다.
	err = session.Query("CREATE TABLE IF NOT EXISTS alog.noti (id uuid PRIMARY KEY, UserPk bigint, MsgContent text, Datetime date, IsChecked boolean);").Exec()
	if err != nil {
		log.Fatal(err)
	}

	/*
	* testing
	 */
	app.Get("/api/noti/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	/*
	* 프론트에서 호출하여 UserPk에 해당하는 전체 메시지 반환
	 */
	app.Get("/api/noti", func(c *fiber.Ctx) error {
		m := c.Queries()
		UserPk := m["userPk"]

		scanner := session.Query("SELECT UserPk, MsgContent, Datetime, IsChecked FROM alog.noti WHERE UserPk=? ALLOW FILTERING;", UserPk).Iter().Scanner()

		returnlist := []*SendMessageTmp{}
		for scanner.Next() {
			msg := &SendMessageTmp{}
			// date := date

			err = scanner.Scan(&msg.UserPk, &msg.MsgContent, &msg.Datetime, &msg.IsChecked)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("msg.Datetime : ", msg.Datetime)
			log.Println("msg.ischecked : ", msg.IsChecked)
			returnlist = append(returnlist, msg)
		}
		// scanner.Err() closes the iterator, so scanner nor iter should be used afterwards.
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return c.JSON(fiber.Map{"data": returnlist})

	})

	/*
	* 서비스에서 호출 하여 메시지를 받아서 db에 넣는것 ->  UserPk, message, time, IsChecked
	 */
	app.Post("/api/noti", func(c *fiber.Ctx) error {

		p := new(AddMessageTmp)
		if err := c.BodyParser(p); err != nil {
			return err
		}

		log.Println(p.UserPk)     // 1 -> 0
		log.Println(p.MsgContent) // hello ->
		now := time.Now().UTC()
		u, err := uuid.NewRandom()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		// TODO insert into noti (_pk , UserPk, message, time, IsChecked) values (p.userPk, p.MsgContent, now, false)
		if err := session.Query("INSERT INTO alog.noti (id, UserPk, MsgContent, Datetime, IsChecked) VALUES (?, ?, ?, ?, ?);", u.String(), p.UserPk, p.MsgContent, now, false).Exec(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Status(fiber.StatusOK).SendString("ok")
	})

	app.Listen(":" + os.Getenv("HOST_PORT"))
}