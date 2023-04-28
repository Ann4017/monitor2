package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type C_db struct {
	s_user     string
	s_pwd      string
	s_host     string
	s_port     string
	s_database string
	s_engine   string
	pc_sql_db  *sql.DB
	s_taget    []string
	s_email    []string
	s_id       []int
}

func (c *C_db) Load_db_config(_s_ini_file_path, _s_section string) error {
	if _, err := os.Stat(_s_ini_file_path); os.IsNotExist(err) {
		cfg := ini.Empty()

		section, _ := cfg.NewSection("section")
		section.NewKey("user", "value")
		section.NewKey("pwd", "value")
		section.NewKey("host", "value")
		section.NewKey("port", "value")
		section.NewKey("database", "value")
		section.NewKey("engine", "value")

		cfg.SaveTo(_s_ini_file_path)

		return fmt.Errorf("Since the %s does not exist, a new %s has been created", _s_ini_file_path, _s_ini_file_path)
	}

	file, err := ini.Load(_s_ini_file_path)
	if err != nil {
		return err
	}

	section, err := file.GetSection(_s_section)
	if err != nil {
		return fmt.Errorf("Failed to get %s section from %s configuration file", _s_section, _s_ini_file_path)
	}

	c.Set_db_config(section)

	return nil
}

func (c *C_db) Set_db_config(section *ini.Section) {
	c.s_user = section.Key("user").String()
	c.s_pwd = section.Key("pwd").String()
	c.s_host = section.Key("host").String()
	c.s_port = section.Key("port").String()
	c.s_database = section.Key("database").String()
	c.s_engine = section.Key("engine").String()
}

func (c *C_db) SQL_connection() error {
	source := c.s_user + ":" + c.s_pwd + "@tcp(" + c.s_host + ":" + c.s_port + ")/" + c.s_database
	sql_db, err := sql.Open(c.s_engine, source)
	if err != nil {
		return err
	}
	c.pc_sql_db = sql_db

	return nil
}

func (c *C_db) SQL_disconnect() error {
	if c.pc_sql_db != nil {
		return c.pc_sql_db.Close()
	}

	return nil
}

func (c *C_db) Select_target() error {
	query := fmt.Sprintf("select * from server")

	rows, err := c.pc_sql_db.Query(query)
	if err != nil {
		return fmt.Errorf("query")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var manager_email string

		if err := rows.Scan(&id, &url, &manager_email); err != nil {
			return fmt.Errorf("scan")
		}

		c.s_id = append(c.s_id, id)
		c.s_taget = append(c.s_taget, url)
		c.s_email = append(c.s_email, manager_email)
	}

	return nil
}

func (c *C_db) Create_table(_s_table string) error {
	query := fmt.Sprintf(`create table if not exists %s (
		id int primary key auto_increment,
		url varchar(255),
		status varchar(50),
		time varchar(50)
	)`, _s_table)

	_, err := c.pc_sql_db.Exec(query)
	if err != nil {
		return fmt.Errorf("create table err")
	}

	return nil
}

func (c *C_db) Insert_status(_s_taget_tables []string, _s_url string, _s_status string, _s_time string) error {
	for _, table := range _s_taget_tables {
		add_row_query := fmt.Sprintf(`insert into %s (url, status, time) 
		values (?, ?, ?)`, table)

		_, err := c.pc_sql_db.Exec(add_row_query, _s_url, _s_status, _s_time)
		if err != nil {
			if strings.Contains(err.Error(), "doesn't exist") {
				if err := c.Create_table(table); err != nil {
					return fmt.Errorf("failed to create table")
				}
				_, err := c.pc_sql_db.Exec(add_row_query, _s_url, _s_status, _s_time)
				if err != nil {
					return fmt.Errorf("failed to insert row new")
				}
				continue
			}
			return fmt.Errorf("failed to insert row")
		}
	}

	return nil
}
