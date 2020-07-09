package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yuanyu90221/cyoa"
)

func main() {
	port := flag.Int("port", 9090, "the port start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", story)
	tpl := template.Must(template.New("").Parse(storyTmpl))
	h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPathFunc(pathFn))
	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting the server at: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "/story" || path == "/story/" { // default intro
		path = "/story/intro"
	}
	// /intro -> intro
	return path[len("/story/"):]
}

var storyTmpl = `
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Choose Your Own Adventure</title>
	</head>
	
	<body>
		<section class="page">
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
			<p>{{.}}</p>
			{{end}}
			<ul>
				{{range .Options}}
					<li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</section>
		<style>
			body {
				font-family: helvetica, arial;
			}
			h1 {
				text-align: center;
				position:relative;
			}
			.page {
				width: 80%;
				max-width: 500px;
				margin: auto;
				margin-top: 40px;
				margin-bottom: 40px;
				padding: 80px;
				background: #FFFCF6;
				border: 1px solid #eee;
				box-shadow: 0 10px 6px -6px #777;
			}
			ul {
				border-top: 1px dotted #ccc;
				padding: 10px 0 0 0;
				-webkit-padding-start: 0;
			}
			li {
				padding-top: 10px;
			}
			a,
			a: visited {
				text-decoration: none;
				color: #6295b5;
			}
			a: active,
			a: hover {
				color: #7792a2;
			}
			p {
				text-indent: 1em;
			}
		</style>
    </body>
</html>`
