package client

type ListResponse struct {
	Data []Account `json:"data,omitempty"`
}

type FetchResponse struct {
	Data DeprecatedAccount `json:"data,omitempty"`
}

type CreateRequest struct {
	Data Account `json:"data,omitempty"`
}

type CreateResponse struct {
	Data Account `json:"data,omitempty"`
}

type Account struct {
	Attributes Attributes `json:"attributes,omitempty"`
	BaseAccount
}

type BaseAccount struct {
	Type           string `json:"type,omitempty"`
	Id             string `json:"id,omitempty"`
	OrganisationId string `json:"organisation_id,omitempty"`
	Version        int    `json:"version,omitempty"`
}

type Attributes struct {
	BaseAttributes
	PrivateIdentification      PrivateIdentification      `json:"private_identification,omitempty"`
	OrganisationIdentification OrganisationIdentification `json:"organisation_identification,omitempty"`
}

type BaseAttributes struct {
	Country                 string   `json:"country,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	BankId                  string   `json:"bank_id,omitempty"`
	BankIdCode              string   `json:"bank_id_code,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	CustomerId              string   `json:"customer_id,omitempty"`
	Name                    []string `json:"name,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
}

type PrivateIdentification struct {
	BirthDate      string   `json:"birth_date,omitempty"`
	BirthCountry   string   `json:"birth_country,omitempty"`
	Identification string   `json:"identification,omitempty"`
	Address        []string `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
}

type OrganisationIdentification struct {
	Identification string   `json:"identification,omitempty"`
	Address        []string `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
	Actors         []Actor  `json:"actors,omitempty"`
}

type Actor struct {
	Name      []string `json:"name,omitempty"`
	BirthDate string   `json:"birth_date,omitempty"`
	Residency string   `json:"residency,omitempty"`
}

type DeprecatedAccount struct {
	Title                       string               `json:"title,omitempty"`
	FirstName                   string               `json:"first_name,omitempty"`
	BankAccountName             string               `json:"bank_account_name,omitempty"`
	AlternativeBankAccountNames []string             `json:"alternative_bank_account_names,omitempty"`
	Attributes                  DeprecatedAttributes `json:"attributes,omitempty"`
	BaseAccount
}

type DeprecatedAttributes struct {
	BaseAttributes
	PrivateIdentification      DeprecatedPrivateIdentification      `json:"private_identification,omitempty"`
	OrganisationIdentification DeprecatedOrganisationIdentification `json:"organisation_identification,omitempty"`
}

type DeprecatedPrivateIdentification struct {
	Title          string `json:"title,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
	PrivateIdentification
}

type DeprecatedOrganisationIdentification struct {
	Name               string         `json:"name,omitempty"`
	RegistrationNumber string         `json:"registration_number,omitempty"`
	Representative     Representative `json:"representative,omitempty"`
	OrganisationIdentification
}

type Representative struct {
	Name      string `json:"name,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	Residency string `json:"residency,omitempty"`
}

/*

relationships.account_eventsarray // not in list
relationships.master_accountarray

*/
