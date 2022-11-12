package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/scottcagno/angular-refresher/cmd/testing/main/example1/users"
)

func main() {

	//                     DEFAULT API ACTIONS
	// +--------+---------------------+------------------+--------------+
	// | METHOD |         HOST        |        PATH      |    HANDLER   |
	// +--------+---------------------+------------------+--------------+
	// |   GET  | http://0.0.0.0:3000 | /items           |  ReturnAll   |
	// |  POST  | http://0.0.0.0:3000 | /items           |  Create      |
	// |   GET  | http://0.0.0.0:3000 | /items/{id}      |  Return      |
	// |   PUT  | http://0.0.0.0:3000 | /items/{id}      |  Update      |
	// | DELETE | http://0.0.0.0:3000 | /items/{id}      |  Delete      |
	// | *ANY	| http://0.0.0.0:3000 | /items/{custom}  |  Custom      |
	// | *ANY	| http://0.0.0.0:3000 | /items/{custom}/{id} |  CustomOne   |
	// +--------+---------------------+------------------+--------------+

	fp, err := os.OpenFile(
		"cmd/testing/main/example1/main.go", os.O_RDONLY,
		0666,
	)
	if err != nil {
		panic(err)
	}

	// init app muxer
	app := http.NewServeMux()

	// init user service
	userService := users.NewUserService(users.NewUserRepoistory())

	// register user service
	app.Handle("/api/users/", userService)
	app.Handle("/file", ServeAFileHandler(fp))
	log.Panicln(http.ListenAndServe(":3000", app))
}

func ServeAFileHandler(file *os.File) http.HandlerFunc {
	// init buffer
	buf := make([]byte, 512)
	// read the fist bit
	_, err := file.Read(buf)
	if err != nil {
		panic(err)
	}
	// detect content type
	var contentType string
	contentType = http.DetectContentType(buf)
	// rewind reader back to beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	// now we can actually handle our web request
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, err := io.CopyBuffer(w, file, buf)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusExpectationFailed), http.StatusExpectationFailed)
			return
		}
	}
}
