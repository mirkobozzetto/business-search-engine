package models

type CompanySearchCriteria struct {
	NaceCode     string `json:"nace_code,omitempty"`
	ZipCode      string `json:"zipcode,omitempty"`
	Status       string `json:"status,omitempty"`
	Denomination string `json:"denomination,omitempty"`
}

type CompanyResult struct {
	EntityNumber   string           `json:"entitynumber"`
	Denominations  []map[string]any `json:"denominations,omitempty"`
	JuridicalForm  string           `json:"juridical_form,omitempty"`
	StartDate      string           `json:"start_date,omitempty"`
	Status         string           `json:"status,omitempty"`
	Addresses      []map[string]any `json:"addresses,omitempty"`
	Contacts       []map[string]any `json:"contacts,omitempty"`
	Activities     []map[string]any `json:"activities,omitempty"`
	Establishments []map[string]any `json:"establishments,omitempty"`
	Enterprise     map[string]any   `json:"enterprise,omitempty"`

	// Legacy fields for compatibility
	Denomination    string `json:"denomination,omitempty"`
	ZipCode         string `json:"zipcode,omitempty"`
	City            string `json:"city,omitempty"`
	Street          string `json:"street,omitempty"`
	HouseNumber     string `json:"house_number,omitempty"`
	Email           string `json:"email,omitempty"`
	Website         string `json:"web,omitempty"`
	Phone           string `json:"tel,omitempty"`
	Fax             string `json:"fax,omitempty"`
	NaceCode        string `json:"nace_code,omitempty"`
	NaceDescription string `json:"nace_description,omitempty"`
}

type CompanySearchResult struct {
	Criteria CompanySearchCriteria `json:"criteria"`
	Results  []CompanyResult       `json:"results"`
	Meta     Meta                  `json:"meta"`
}
