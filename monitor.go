package main

import (
	"fmt"
	"log"
	"os"
)

type Monitor struct {
	h C_http
	d C_db
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

func (m *Monitor) Init(_s_ini_file_path, _s_section string) error {
	err := m.Create_log()
	if err != nil {
		return fmt.Errorf("log")
	}

	err = m.d.Load_db_config(_s_ini_file_path, _s_section)
	if err != nil {
		return fmt.Errorf("load")
	}

	err = m.d.SQL_connection()
	if err != nil {
		return fmt.Errorf("sql")
	}

	return nil
}

func (m *Monitor) Run() error {
	err := m.d.Select_target()
	if err != nil {
		return err
	}
	fmt.Println(m.d.s_id, m.d.s_taget, m.d.s_email)

	for _, url := range m.d.s_taget {
		err := m.h.Get_http_status(url)
		if err != nil {
			return err
		}

		err = m.d.Insert_status(m.d.s_taget, m.h.s_url, m.h.s_status, m.h.s_time)
		if err != nil {
			return err
		}
	}

	err = m.d.SQL_disconnect()

	return nil
}
