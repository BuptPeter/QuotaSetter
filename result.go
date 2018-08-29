package main

type SetResult struct {
	Code         int          `json:"code"`
	IsSuccess    bool          `json:"success"`
	Description  string       `json:"description"`

}

type SetResults []SetResult

type GetResult struct {
	Code         int          `json:"code"`
	IsSuccess    bool          `json:"success"`
	Description  string       `json:"description"`
	Max_bytes    string       `json:"max_bytes"`
	Max_files    string       `json:"max_files"`

}

type GetResults []SetResult