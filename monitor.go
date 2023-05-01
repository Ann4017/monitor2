package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Monitor struct {
	h C_http
	d C_db
	s C_ses
}

func (m *Monitor) Init(_s_ini_file_path, _s_db_section, _s_aws_section string) error {
	err := m.Create_log()
	if err != nil {
		return err
	}

	err = m.d.Load_db_config(_s_ini_file_path, _s_db_section)
	if err != nil {
		return err
	}

	err = m.d.SQL_connection()
	if err != nil {
		return err
	}

	err = m.s.Load_aws_config(_s_ini_file_path, _s_aws_section)
	if err != nil {
		return err
	}

	err = m.s.Set_config()
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Run(_time_interval time.Duration) error {
	m.d.s_crated_tables = make(map[string]bool)

	err := m.d.Select_target()
	if err != nil {
		return err
	}
	fmt.Println(m.d.s_id, m.d.s_taget, m.d.s_email)

	ticker := time.NewTicker(_time_interval)
	for range ticker.C {

		for _, url := range m.d.s_taget {
			err := m.h.Get_http_status(url)
			if err != nil {
				return err
			}

			// err = m.d.Create_table(url)
			// if err != nil {
			// 	return err
			// }

			err = m.d.Insert_status(url, m.h.s_url, m.h.s_status, m.h.s_time)
			if err != nil {
				return err
			}

			err = m.d.Select_err_status(url)
			if err != nil {
				return err
			}
		}

		for _, email := range m.d.s_email {
			if m.d.s_err_rows != "" {
				m.s.Send_email(m.s.pc_client, email, m.d.s_email, "server_err", m.d.s_err_rows)
			}
		}

	}

	err = m.d.SQL_disconnect()
	if err != nil {
		return err
	}

	err = m.d.Close_rows()
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Create_log() error {
	log_file, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer log_file.Close()

	log.SetOutput(log_file)

	return nil
}

func (m *Monitor) Insert_server(_url, _manager_email string) error {
	_, err := m.d.pc_sql_db.Exec("insert into server (url, manager_email) values (?, ?)", "`"+_url+"`", _manager_email)
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Delete_serer(_id int) error {
	_, err := m.d.pc_sql_db.Exec("delete from server where id = ?", _id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Update_server(_id int, _url, _manager_email string) error {
	_, err := m.d.pc_sql_db.Exec("update server set url = ?, manager_email = ? where = ?", "`"+_url+"`", _manager_email, _id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Monitor) Select_server() error {
	rows, err := m.d.pc_sql_db.Query("select * from server")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var manager_email string
		if err := rows.Scan(&id, &url, &manager_email); err != nil {
			return err
		}
		fmt.Printf("id:%d, url:%s, email:%s\n", id, url, manager_email)
	}

	return nil
}
