package main

import (
	"database/sql"
	"fmt"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
	_ "github.com/go-sql-driver/mysql"
)

type MyStruct struct {
    ID   int64
    post_content string
}

var myCntlimit = 10

func main() {
	db1, err := sql.Open("mysql", "root:welcome@tcp(192.168.1.222:3306)/mydbgo")
	checkErr(err)
	rows, err := db1.Query("SELECT ID, post_content FROM dbedu_posts where post_content_filtered = '-' limit " + strconv.Itoa(myCntlimit))
	checkErr(err)
	thisCollection := make([]map[string]interface{},myCntlimit)
	converter := md.NewConverter("", true, nil)
	i := 0
	for rows.Next() {
		thisCollection[i] = make(map[string]interface{}, myCntlimit)
        var thisRow MyStruct
        // for each row, scan the result into our tag composite object
        err = rows.Scan(&thisRow.ID, &thisRow.post_content)
		checkErr(err)
		markdown, err := converter.ConvertString(thisRow.post_content)
		checkErr(err)
		thisCollection[i]["ID"] = thisRow.ID      // int64
		thisCollection[i]["post_markdown"] = markdown  // string
		i++ 
	}

	for i := 0; i < myCntlimit; i++ {
		stmt, err := db1.Prepare("update dbedu_posts set post_content_filtered=? where ID=?")
		checkErr(err)
		res, err := stmt.Exec(thisCollection[i]["post_markdown"], thisCollection[i]["ID"])
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(affect)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
