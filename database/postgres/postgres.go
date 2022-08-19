package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/richardbizik/db-schema-web-view/database"
)

var (
	ctx context.Context
	db  *sql.DB
)

const (
	ReadTableInfo        = "SELECT table_name  FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = $1;"
	ReadColumnInfo       = "SELECT column_name, ordinal_position, cast(column_default as text), is_nullable, udt_name, character_maximum_length, numeric_precision, numeric_precision_radix, numeric_scale, datetime_precision FROM information_schema.columns WHERE table_name = $1 AND table_schema = $2;"
	ReadTablesConstraint = "SELECT constraint_name, table_name, constraint_type FROM information_schema.table_constraints WHERE constraint_catalog = $1 and constraint_schema = $2"
	ReadFkUqPk           = "SELECT kcu.constraint_name, kcu.table_name, kcu.column_name, kcu.ordinal_position, rcu.table_name, rcu.column_name, tc.constraint_type  FROM information_schema.table_constraints tc FULL JOIN information_schema.key_column_usage kcu ON tc.constraint_schema = kcu.constraint_schema AND tc.constraint_name = kcu.constraint_name FULL JOIN information_schema.referential_constraints rc ON tc.constraint_schema = rc.constraint_schema AND tc.constraint_name = rc.constraint_name FULL JOIN information_schema.key_column_usage rcu ON rc.unique_constraint_schema = rcu.constraint_schema AND rc.unique_constraint_name = rcu.constraint_name AND kcu.ordinal_position = rcu.ordinal_position WHERE constraint_type in ('PRIMARY KEY', 'UNIQUE', 'FOREIGN KEY') AND (rcu.table_schema = $1 or tc.table_schema = $1) ORDER BY kcu.table_schema, kcu.table_name, kcu.ordinal_position;"
)

func GetDBSchema(dbName string, host string, port string, user string, pass string, schema string) (database.Model, error) {
	model := database.Model{
		Tables: []database.Table{},
		Name:   dbName,
		Schema: schema,
	}

	db, err := sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s", user, pass, host, port, dbName))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return database.Model{}, err
	}
	defer db.Close()

	//read table info
	rows, err := db.Query(
		ReadTableInfo, schema)
	defer rows.Close()
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		table := database.Table{
			Columns:          []database.Column{},
			Name:             tableName,
			TableConstraints: []database.TableConstraint{},
			//FKConstraints: []database.FKConstraint{},
		}
		if err != nil {
			log.Println(err)
			return database.Model{}, err
		}
		model.Tables = append(model.Tables, table)
	}
	rows.Close()

	//read column info
	for i := range model.Tables {
		rows, err := db.Query(
			ReadColumnInfo, model.Tables[i].Name, schema)
		if err != nil {
			log.Println(rows.Err())
			return database.Model{}, err
		}
		for rows.Next() {
			var (
				name            sql.NullString
				colType         sql.NullString
				typeDetail      sql.NullString
				ordinalPosition sql.NullInt32
				columnDefault   sql.NullString
				nullable        sql.NullString

				charMaxLen        sql.NullInt32
				numPrecision      sql.NullInt32
				numPrecisionRadix sql.NullInt32
				numScale          sql.NullInt32
				dateTimePrec      sql.NullInt32
			)
			err := rows.Scan(&name, &ordinalPosition, &columnDefault, &nullable, &colType, &charMaxLen, &numPrecision, &numPrecisionRadix, &numScale, &dateTimePrec)
			var nullableBool bool
			if nullable.Valid && nullable.String == "YES" {
				nullableBool = true
			} else {
				nullableBool = false
			}
			column := database.Column{
				Name:            name.String,
				ColType:         colType.String,
				TypeDetail:      typeDetail.String,
				OrdinalPosition: int(ordinalPosition.Int32),
				ColumnDefault:   columnDefault.String,
				Nullable:        nullableBool,
				Primary:         false,
				Unique:          []string{},
				FKConstraints:   []database.FKConstraint{},
			}
			if err != nil {
				log.Println(rows.Err())
				return database.Model{}, err
			}
			model.Tables[i].Columns = append(model.Tables[i].Columns, column)
		}
		rows.Close()
	}

	//read foreign keys
	rows, err = db.Query(
		ReadFkUqPk, schema)
	if err != nil {
		log.Printf("Unable to retrieve table constraints: %v\n", err)
		return database.Model{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var kcuConstraintName string
		var kcuTableName string
		var kcuColumnName string
		var kcuOrdinalPosition int
		var rcuTableName sql.NullString
		var rcuColumnName sql.NullString
		var constraintType string

		err := rows.Scan(&kcuConstraintName, &kcuTableName, &kcuColumnName, &kcuOrdinalPosition, &rcuTableName, &rcuColumnName, &constraintType)

		for i := range model.Tables {
			if model.Tables[i].Name == kcuTableName {
				for c := range model.Tables[i].Columns {
					if model.Tables[i].Columns[c].Name == kcuColumnName {
						if constraintType == "FOREIGN KEY" {
							model.Tables[i].Columns[c].FKConstraints = append(model.Tables[i].Columns[c].FKConstraints, database.FKConstraint{
								Name: kcuConstraintName,
								//KTable:    kcuTableName,
								//KColumn:   kcuColumnName,
								KPosition: kcuOrdinalPosition,
								RTable:    rcuTableName.String,
								RColumn:   rcuColumnName.String,
							})
							break
						} else if constraintType == "UNIQUE" {
							model.Tables[i].Columns[c].Unique = append(model.Tables[i].Columns[c].Unique, kcuConstraintName)
						} else if constraintType == "PRIMARY KEY" {
							model.Tables[i].Columns[c].Primary = true
						}
					}
				}
				break
			}
		}
		if err != nil {
			log.Println(err)
			return database.Model{}, err
		}
	}
	return model, nil
}
