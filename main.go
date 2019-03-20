package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

// stringService is a concrete implementation of StringService
type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

type BookService interface {
	Books() map[string]Book
	GetBook(int)  (Book, error)
	SetBook(Book) error
}
// bookService is a concrete implementation of BookService
type bookService struct{}

type Book struct {
	ID int
	Title string
	Author string
	Date string
	Publisher string
}
var bookmap map[string]Book  = make(map[string]Book)

func (bookService) Books() map[string]Book {
	if len(bookmap) == 0 {
		log.Println("Init Books")
		bookmap["1"]=Book{
			ID:1,
			Title:"A Spell for Chameleon",
			Author:"Piers Anthony",
			Date: "1977",
			Publisher:" Del Rey ",
		}
		bookmap["2"]=Book{
			ID:        2,
			Title:     "The Source of Magic",
			Author:    "Piers Anthony",
			Date:      "1979",
			Publisher: " Del Rey ",
		}
	}
	return bookmap
}

func (bookService) SetBook( book Book) error {
	id  := strconv.Itoa(book.ID)
	bookmap[id] = book
	return nil
}

func (bookService) GetBook(id int) (Book, error) {
	strid := strconv.Itoa(id)
	//Book = bookmap[id]
	book, _ := bookmap[strid]
	return book, nil
}


// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// For each method, we define request and response structs
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

type getBookRequest struct {
	ID int `json:"id"`
}

type getBookResponse struct {
	Book Book `json:"book"`
}

type setBookRequest struct {
	Book Book `json:"book"`
}

type setBookResponse struct {
	ok bool `json:"ok"`
}


type booksRequest struct {

}

type booksResponse struct {
	Books 	map[string]Book `json:"books"`
}


// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func makeBooksEndpoint(svc BookService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
//		?\request.(booksRequest)
		b := svc.Books()
		//if err != nil {
		//	return booksResponse{Books:b}, nil
		//}
		return booksResponse{Books:b}, nil
	}
}

func makeBookEndpoint(svc BookService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getBookRequest)
		//idString := req.ID
		//id, _  := strconv.Atoi(idString)
		//book, err := svc.GetBook(id)
		book, err := svc.GetBook(req.ID)
		if err != nil {
			return getBookResponse{Book:book}, nil
		}
		return getBookResponse{Book:book}, nil
	}
}


func makeSetBookEndpoint(svc BookService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(setBookRequest)
		//idString := req.ID
		//id, _  := strconv.Atoi(idString)
		//book, err := svc.GetBook(id)
		err := svc.SetBook(req.Book)
		if err != nil {
			return setBookResponse{ok:true}, nil
		}
		return setBookResponse{ok:true}, nil
	}
}
// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
func main() {
	svc := stringService{}
	bookSvc := bookService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	booksHandler := httptransport.NewServer(
		makeBooksEndpoint(bookSvc),
		decodeBooksRequest,
		encodeResponse,
	)

	bookHandler := httptransport.NewServer(
		makeBookEndpoint(bookSvc),
		decodeBookRequest,
		encodeResponse,
	)

	setbookHandler := httptransport.NewServer(
		makeSetBookEndpoint(bookSvc),
		decodeSetBookRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/books", booksHandler)
	http.Handle("/book", bookHandler)
	http.Handle("/setbook", setbookHandler)

	log.Fatal(http.ListenAndServe(":7070", nil))
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeBooksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request booksRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSetBookRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request setBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
