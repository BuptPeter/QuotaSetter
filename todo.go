package main

// import "time"

type Todo struct {
	Path       string    `json:"target"`
	Max_bytes  int       `json:"max_bytes"`
	Max_files  int       `json:"max_files"`
}

type Todos []Todo
