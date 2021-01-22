package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func queryToHtml(dobre string) string {
	dobre = strings.TrimPrefix(dobre, "{\"")
	dobre = strings.TrimSuffix(dobre, "\"}")
	return dobre
}

func fQueryQuestion(key string, db *sql.DB) string {
	queryQuestion := "select question_text from userkey NATURAL JOIN question where user_key = '{"
	queryQuestion += key
	queryQuestion += "}';"
	var question_text string
	row := db.QueryRow(queryQuestion)
	switch err := row.Scan(&question_text); err {
	case sql.ErrNoRows:
		question_text = "Poll not found"
	case nil:
		question_text = queryToHtml(question_text)

	}
	fmt.Println(question_text)
	return question_text
}

type Option struct {
	OptionID   string
	OptionText string
}

type OptionSubmit struct {
	OptionText string
	Votes      string
}

/*func createOptions(optionArray []option) string {
	var options string
	for i := range optionArray {

		options += "<input type=\"radio\" id=\"" + strconv.Itoa(optionArray[i].optionID) + "\" name=\"option\" value=\"" + optionArray[i].optionText + "\">"
	}
	fmt.Print(options)
	return options
}*/

func renderShowPoll(key string) {

}
func main() {
	//Postgres
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hello world")
	fmt.Println(db.Query("SELECT * FROM question"))

	/*
		var (
			id   int
			name string
		)

		rows, err := db.Query("SELECT * FROM question")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	*/

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.HTML(http.StatusOK, "index.tmpl.html", nil)
		} else {
			masina := "select question_system from userkey natural join question WHERE user_key = '{"
			masina += key
			masina += "}' AND used = FALSE;"
			fmt.Println(masina)
			row := db.QueryRow(masina)
			var question_system int
			switch err := row.Scan(&question_system); err {
			case sql.ErrNoRows:
				c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
					"OptionSubmited": "Key already submitted",
					"Key":            key,
				})
			case nil:
				fmt.Println(question_system)
				fmt.Println("GET / Case nil")

				queryOption := "select option_id, option_text from userkey natural join option where user_key = '{"
				queryOption += key
				queryOption += "}';"
				Options := []Option{}

				var (
					option_id   int
					option_text string
				)
				rows, err := db.Query(queryOption)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&option_id, &option_text)
					if err != nil {
						log.Fatal(err)
					}
					Options = append(Options, Option{strconv.Itoa(option_id), queryToHtml(option_text)})
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				if question_system == 0 {
					c.HTML(http.StatusOK, "showPoll.tmpl.html", gin.H{
						"key":           key,
						"question_text": fQueryQuestion(key, db),
						"Options":       Options,
					})
				} else {
					c.HTML(http.StatusOK, "showPollNumbers.tmpl.html", gin.H{
						"key":           key,
						"question_text": fQueryQuestion(key, db),
						"Options":       Options,
					})
				}

			default:
				fmt.Println("Nejde to :-(")
				fmt.Println(err)
			}

		}
	})

	router.POST("/", func(c *gin.Context) {
		sqlStatement := "select user_key from userkey where user_key = '{"
		Key := c.PostForm("key")
		sqlStatement += Key
		sqlStatement += "}' AND used = FALSE;"
		fmt.Println(sqlStatement)
		row := db.QueryRow(sqlStatement)
		var user_key string
		switch err := row.Scan(&user_key); err {
		case sql.ErrNoRows:
			fmt.Println("zapinam world")
			fmt.Println(Key)
			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Key already submitted",
				"Key":            Key,
			})
		case nil:
			//UPDATE option SET votes = votes + 1 WHERE option_id = 1;
			sqlSubmit := "UPDATE option SET votes = votes + 1 WHERE option_id ="
			sqlSubmit += c.PostForm("option")
			sqlSubmit += ";"
			_, err = db.Exec(sqlSubmit)
			if err != nil {
				panic(err)
			}
			// UPDATE userkey SET used = TRUE WHERE user_key = '{RS8e7vC5lfwV4u3YfkPKu3hZxU76UE45}'
			userkeyFalse := "UPDATE userkey SET used = TRUE WHERE user_key = '{"
			userkeyFalse += Key
			userkeyFalse += "}'"
			_, err = db.Exec(userkeyFalse)
			if err != nil {
				fmt.Println(err)
			}
			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Vote has been counted",
				"Key":            Key,
			})
		}

	})
	router.GET("/createPoll", func(c *gin.Context) {
		c.HTML(http.StatusOK, "createPoll.tmpl.html", nil)
	})

	router.GET("/showResults", func(c *gin.Context) {
		//select option_text, votes from userkey  NATURAL JOIN option WHERE user_key = '{AS8e7vC5lfwV4u3YfkPKu3hZxU76UE45}';
		fmt.Println("zapinam world")
		key := c.Query("key")
		OptionSubmits := []OptionSubmit{}
		sqlDisplayVotes := "select option_text, votes from userkey  NATURAL JOIN option WHERE user_key = '{"
		sqlDisplayVotes += key
		sqlDisplayVotes += "}';"
		fmt.Println(sqlDisplayVotes)
		var (
			option_text string
			votes       string
		)
		rows, err := db.Query(sqlDisplayVotes)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&option_text, &votes)
			if err != nil {
				log.Fatal(err)
			}
			OptionSubmits = append(OptionSubmits, OptionSubmit{queryToHtml(option_text), votes})
			fmt.Println("spustila sa vec world")
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "showResults.tmpl.html", gin.H{
			"Question_text":  fQueryQuestion(key, db),
			"OptionsSubmits": OptionSubmits,
		})

	})

	router.Run(":" + port)

}
