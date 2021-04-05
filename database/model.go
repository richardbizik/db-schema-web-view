package database

type ConstraintType int

const (
	PK ConstraintType = iota
	FK
	UQ
	CH
)

type Column struct {
	Name            string         `json:"name"`
	ColType         string         `json:"type"`
	TypeDetail      string         `json:"typeDetail"`
	OrdinalPosition int            `json:"position"`
	ColumnDefault   string         `json:"default"`
	Nullable        bool           `json:"nullable"`
	Primary         bool           `json:"primary"`
	Unique          []string       `json:"uniqueConstraints"`
	FKConstraints   []FKConstraint `json:"fkConstraints"`
}

type UniqueConstraint struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type FKConstraint struct {
	//KTable    string `json:"kTable"`
	//KColumn   string `json:"kColumn"`
	Name      string `json:"name"`
	KPosition int    `json:"kPosition"`
	RTable    string `json:"rTable"`
	RColumn   string `json:"rColumn"`
}

type TableConstraint struct {
	Name       string `json:"name"`
	Constraint string `json:"constraint"` //check constraints for postgres
}

type Table struct {
	Name             string            `json:"name"`
	Columns          []Column          `json:"columns"`
	TableConstraints []TableConstraint `json:"tableConstraints"`
	//FKConstraints    []FKConstraint    `json:"fkConstraints"`
}

type Model struct {
	Name   string  `json:"name"`
	Schema string  `json:"schema"`
	Tables []Table `json:"tables"`
}
