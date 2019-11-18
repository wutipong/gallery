package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

var fileItemTemplace *template.Template

func init() {
	var err error
	fileItemTemplace, err = template.New("fileitem.gohtml").ParseFiles("template/fileitem.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

type fileItem struct {
	ImageURL string
	Name     string
}

//WriteBreadcrumb write breadcrumb component.
func WriteBreadcrumb(writer io.Writer, path string) {
	io.WriteString(writer, "<nav>")
	io.WriteString(writer, `<ol class="breadcrumb">`)
	if path == "" {
		io.WriteString(writer, `<li class="breadcrumb-item active" aria-current="page">Home</li>`)
	} else {
		io.WriteString(writer, `<li class="breadcrumb-item active" aria-current="page"><a href="/browse">Home</a></li>`)

		parts := strings.Split(path, "/")
		for i, part := range parts {
			if i == len(parts)-1 {
				io.WriteString(writer, fmt.Sprintf(`<li class="breadcrumb-item active" aria-current="page">%s</li>`, part))
			} else {
				url := strings.Join(parts[0:i+1], "/")
				url = "/browse/" + url
				io.WriteString(writer, fmt.Sprintf(`<li class="breadcrumb-item active" aria-current="page"><a href="%s">%s</a></li>`, url, part))
			}
		}
	}
	io.WriteString(writer, "</ol>")
	io.WriteString(writer, "</nav>")
}

// WriteDirectories write directory entries.
func WriteDirectories(writer io.Writer, path string, dirs []FileEntry) {
	io.WriteString(writer, `<div class="container">`)
	length := len(dirs)
	for i := 0; i < length; i++ {
		dir := dirs[i]
		if i%3 == 0 {
			io.WriteString(writer, `<div class="row">`)
		}
		var url string

		if path == "" {
			url = "/browse/" + dir.Filename
		} else {
			url = "/browse/" + path + "/" + dir.Filename
		}

		io.WriteString(writer, fmt.Sprintf(`<div class="col"><a href="%s">%s</a></div>`, url, dir.Filename))

		if i%3 == 2 || i == length-1 {
			io.WriteString(writer, `</div>`)
		}
	}
	io.WriteString(writer, `</div>`)
}

// WriteFiles write file entries.
func WriteFiles(writer io.Writer, path string, files []FileEntry) {
	io.WriteString(writer, `<div class="container">`)
	length := len(files)
	for i := 0; i < length; i++ {
		file := files[i]
		if i%3 == 0 {
			io.WriteString(writer, `<div class="row">`)
		}
		var url string

		if path == "" {
			url = "/get_image/" + file.Filename
		} else {
			url = "/get_image/" + path + "/" + file.Filename
		}

		fileItemTemplace.Execute(writer, fileItem{ImageURL: url, Name: file.Filename})

		if i%3 == 2 || i == length-1 {
			io.WriteString(writer, `</div>`)
		}
	}
	io.WriteString(writer, `</div>`)
}

// Handler
func browse(c echo.Context) error {
	builder := strings.Builder{}

	WriteHeader(&builder, Header{Title: "Hello"})

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	dirs, files, err := ListDir(p)
	if err != nil {
		return err
	}

	WriteBreadcrumb(&builder, p)
	WriteDirectories(&builder, p, dirs)
	WriteFiles(&builder, p, files)

	builder.WriteString(`
	<script>
	$(function() {
		$('.thumb').Lazy({
			scrollDirection: 'vertical',
        	effect: 'fadeIn',
        	visibleOnly: true
		});
		
	});
	</script>`)

	return c.HTML(http.StatusOK, builder.String())
}
