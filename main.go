package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/gzip"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func fQueryQuestion(key string, db *sql.DB) string {
	var questionText string
	row := db.QueryRow("select question_text from userkey NATURAL JOIN question where user_key = $1", key)
	switch err := row.Scan(&questionText); err {
	case sql.ErrNoRows:
		questionText = "Poll not found"
	case nil:
	}
	fmt.Println(questionText)
	return questionText
}

//Option najlepsia vec
type Option struct {
	OptionID   string
	OptionText string
}

//OptionSubmit super vec
type OptionSubmit struct {
	OptionText string
	Votes      int
}

type pollNumbers struct {
	optionID    int
	optionValue int
}

var (
	ctx context.Context

	db *sql.DB
)

/*func createOptions(optionArray []option) string {
	var options string
	for i := range optionArray {

		options += "<input type=\"radio\" id=\"" + strconv.Itoa(optionArray[i].optionID) + "\" name=\"option\" value=\"" + optionArray[i].optionText + "\">"
	}
	fmt.Print(options)
	return options
}*/

func main() {
	//Postgres
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	var questionID int
	userKey := "0dfa16810a3852b0b397e40e834ac9a9"
	err = db.QueryRow("SELECT question_id FROM userkey natural join question WHERE user_key = $1", userKey).Scan(&questionID)
	switch {

	case err == sql.ErrNoRows:

		fmt.Println("no user with id ")

	case err != nil:

		fmt.Println(err)

	default:

		fmt.Println(questionID)

	}
	//fmt.Println(questionID)
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
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.HTML(http.StatusOK, "index.tmpl.html", nil)
		} else {
			var questionSystem int
			err = db.QueryRow("select question_system from userkey natural join question WHERE user_key = $1 AND used = FALSE AND question_date > now()", key).Scan(&questionSystem)
			switch {
			case err == sql.ErrNoRows:
				var date string
				fmt.Println(155)

				err = db.QueryRow("SELECT question_date FROM userkey natural join question WHERE user_key = $1 AND question_date > now()", key).Scan(&date)
				switch {

				case err == sql.ErrNoRows:
					err = db.QueryRow("SELECT question_id FROM userkey  WHERE user_key = $1", key).Scan()
					switch {

					case err == sql.ErrNoRows:

						c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
							"Error": "Enter valid key",
						})

					case err != nil:

						fmt.Println(err)

					default:

						fmt.Println(164)
						date = fmt.Sprintf("Poll has expired on %s", date)
						c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
							"OptionSubmited": strings.TrimSuffix(date, "T00:00:00Z"),
							"Key":            key,
						})

					}

				case err != nil:

					fmt.Println(170)
					c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
						"OptionSubmited": "Key already submitted",
						"Key":            key,
					})

				default:

					fmt.Println(170)
					c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
						"OptionSubmited": "Key already submitted",
						"Key":            key,
					})

				}
			case err != nil:
				fmt.Println("Nejde to :-(")
				fmt.Println(err)

			default:
				fmt.Println(questionSystem)
				fmt.Println("GET / Case nil")
				/*
					queryOption := "select option_id, option_text from userkey natural join option where user_key = '{"
					queryOption += key
					queryOption += "}';"*/
				Options := []Option{}

				var (
					optionID   int
					optionText string
				)
				rows, err := db.Query("select option_id, option_text from userkey natural join option where user_key = $1", key)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					err := rows.Scan(&optionID, &optionText)
					if err != nil {
						log.Fatal(err)
					}
					Options = append(Options, Option{strconv.Itoa(optionID), optionText})
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
				if questionSystem == 0 {
					c.HTML(http.StatusOK, "showPoll.tmpl.html", gin.H{
						"key":          key,
						"questionText": fQueryQuestion(key, db),
						"Options":      Options,
					})
				} else {
					c.HTML(http.StatusOK, "showPollNumbers.tmpl.html", gin.H{
						"key":          key,
						"questionText": fQueryQuestion(key, db),
						"Options":      Options,
					})

				}
			}

		}
	})

	router.POST("/", func(c *gin.Context) {
		//sqlStatement := "select user_key from userkey where user_key = '{"
		key := c.PostForm("key")
		//sqlStatement += Key
		//sqlStatement += "}' AND used = FALSE;"
		//fmt.Println(sqlStatement)
		row := db.QueryRow("select user_key from userkey natural join question where user_key = '$1' AND used = FALSE AND question_date > now()", key)
		var userKey string
		switch err := row.Scan(&userKey); err {
		case sql.ErrNoRows:
			fmt.Println("zapinam world")
			fmt.Println(key)
			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Issue with key try again",
				"Key":            key,
			})
		case nil:
			fmt.Println("Fsetko vporiadku")
			//UPDATE option SET votes = votes + 1 WHERE option_id = 1;
			//sqlSubmit := "UPDATE option SET votes = votes + 1 WHERE option_id ="
			//sqlSubmit += c.PostForm("option")
			//sqlSubmit += ";"
			_, err = db.Exec("UPDATE option SET votes = votes + 1 WHERE option_id =$1", c.PostForm("option"))
			if err != nil {
				panic(err)
			}
			// UPDATE userkey SET used = TRUE WHERE user_key = '{RS8e7vC5lfwV4u3YfkPKu3hZxU76UE45}'
			//userkeyFalse := "UPDATE userkey SET used = TRUE WHERE user_key = '{"
			//userkeyFalse += Key
			//userkeyFalse += "}'"
			_, err = db.Exec("UPDATE userkey SET used = TRUE WHERE user_key = $1", key)
			if err != nil {
				fmt.Println(err)
			}
			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Vote has been counted",
				"Key":            key,
			})
		}

	})

	router.POST("/submitPollNumbers", func(c *gin.Context) {
		c.MultipartForm()
		user := c.PostForm("key")
		fmt.Println(user)
		err = db.QueryRow("select user_key from userkey natural join question where user_key = $1 AND used = FALSE AND question_date > now()", user).Scan(&userKey)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("zapinam world")
			fmt.Println(user)
			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Issue with key try again",
				"Key":            user,
			})
		case err != nil:
			log.Fatalf("query error: %v\n", err)
		default:
			for key, value := range c.Request.PostForm {
				if key == "key" {
				} else {
					valueInt, _ := strconv.Atoi(strings.Join(value, ""))
					fmt.Printf("%v = %v", key, valueInt)
					_, err = db.Exec("UPDATE option SET votes = votes + $1 WHERE option_id =$2", valueInt, key)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf(key)
					_, err = db.Exec("UPDATE userkey SET used = TRUE WHERE user_key = $1", user)
					if err != nil {
						fmt.Println(err)
					}
				}
				log.Printf(user)
			}

			c.HTML(http.StatusOK, "optionSubmited.tmpl.html", gin.H{
				"OptionSubmited": "Vote has been counted",
				"Key":            user,
			})

		}

	})

	router.POST("/createPoll", func(c *gin.Context) {
		question := c.PostForm("question")
		participants, err := strconv.Atoi(c.PostForm("participants"))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Participants: %d", participants)
		system := c.PostForm("system")
		end := c.PostForm("end")
		option := c.PostFormArray("option")
		fmt.Println(question, participants, system, end, option)
		var questionID int
		var generatedKeys []string
		err = db.QueryRow("INSERT INTO question (question_text, question_system, question_date) VALUES ($1, $2, $3) RETURNING question_id", question, system, end).Scan(&questionID)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("no user with id \n")
		case err != nil:
			log.Fatalf("query error: %v\n", err)
		default:
			fmt.Printf("%d", questionID)
			for k, v := range option {
				fmt.Println(v)
				fmt.Println(k)
				result, err := db.Exec("INSERT INTO option (option_text, question_id) VALUES ($1, $2)", v, questionID)
				if err != nil {
					fmt.Println(err)
				}
				rows, err := result.RowsAffected()
				if err != nil {
					fmt.Println(err)
				}
				if rows != 1 {
					fmt.Printf("expected to affect 1 row, affected %d", rows)
				}
			}
			fmt.Printf("Participants %d", participants)
			for i := 0; i < participants; i++ {

				err = db.QueryRow("INSERT INTO userkey (question_id) VALUES ($1) RETURNING user_key", questionID).Scan(&userKey)
				switch {
				case err != nil:
					log.Fatalf("query error: %v\n", err)
				default:
					fmt.Println(userKey)
					generatedKeys = append(generatedKeys, userKey)
				}
			}

		}
		c.HTML(http.StatusOK, "showKeys.tmpl.html", gin.H{
			"Keys": generatedKeys,
		})

	})
	router.GET("/createPoll", func(c *gin.Context) {
		c.HTML(http.StatusOK, "createPoll.tmpl.html", nil)
	})

	router.GET("/showResults", func(c *gin.Context) {
		//select optionText, votes from userkey  NATURAL JOIN option WHERE user_key = '{AS8e7vC5lfwV4u3YfkPKu3hZxU76UE45}';
		fmt.Println("zapinam world")
		key := c.Query("key")
		OptionSubmits := []OptionSubmit{}
		//sqlDisplayVotes := "select option_text, votes from userkey  NATURAL JOIN option WHERE user_key = '{"
		//sqlDisplayVotes += key
		//sqlDisplayVotes += "}';"
		//fmt.Println(sqlDisplayVotes)
		var (
			optionText string
			votes      int
		)
		rows, err := db.Query("select option_text, votes from userkey  NATURAL JOIN option WHERE user_key = $1", key)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&optionText, &votes)
			if err != nil {
				log.Fatal(err)
			}
			OptionSubmits = append(OptionSubmits, OptionSubmit{optionText, votes})

		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "showResults.tmpl.html", gin.H{
			"questionText":   fQueryQuestion(key, db),
			"OptionsSubmits": OptionSubmits,
		})

	})

	router.Run(":" + port)

}
