package main

type Backpack struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
	Items     string `json:"backpack_json"`
}

type Backpacks []Backpack
