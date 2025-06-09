package model

type Bank struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
	Founded int    `json:"founded"`
	Rank    int    `json:"rank"`
	Profile string `json:"profile"`
}

type BankListData struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type BankProfile struct {
	Name      string         `json:"name"`
	CEO       string         `json:"ceo,omitempty"`
	Employees int            `json:"employees,omitempty"`
	Lists     []BankListData `json:"lists,omitempty"`
}
