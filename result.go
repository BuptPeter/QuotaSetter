package main

type Result struct {
	Code         int          `json:"code"`
	IsSuccess    bool          `json:"success"`
	Description  string       `json:"description"`

}

type Results []Result
