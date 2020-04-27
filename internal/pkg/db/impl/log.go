package impl

import (
	exception "clicktweak/internal/pkg/error"
	"clicktweak/internal/pkg/model"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	clickBasedQuery = `SELECT COUNT() FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id`

	userBasedQuery = `SELECT uniq(remote_address) FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id`

	clickBasedGroupByBrowser = `SELECT COUNT(), browser FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id
			GROUP BY browser`

	clickBasedGroupByDevice = `SELECT COUNT(), device FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id
			GROUP BY device`

	userBasedGroupByBrowser = `SELECT uniq(remote_address), browser FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id 
			GROUP BY browser`

	userBasedGroupByDevice = `SELECT uniq(remote_address), device FROM logs
			WHERE created_at >= toDateTime(:from) AND created_at <= toDateTime(:until) AND id = :id
			GROUP BY device`
)

var queies []string = []string{
	clickBasedQuery,
	clickBasedGroupByBrowser,
	clickBasedGroupByDevice,
	userBasedQuery,
	userBasedGroupByBrowser,
	userBasedGroupByDevice,
}

type Log struct {
	db *sqlx.DB
}

func NewLogDB(db *sqlx.DB) (*Log, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS logs (
										id 				String,
										created_at 		DateTime,
										device 			String,
										browser 		String,
										remote_address 	String
									) ENGINE = 	MergeTree()
									PRIMARY KEY created_at
									ORDER BY 	(created_at, id)
									SETTINGS 	index_granularity=8192`,
	)
	if err != nil {
		log.Error(err)
		return nil, exception.InternalServerError
	}

	return &Log{db}, nil
}

func (l *Log) GetStats(id, from, until string) (*model.Report, error) {
	queryParams := map[string]interface{}{
		"from":  from,
		"until": until,
		"id":    id,
	}

	var result = new(model.Report)
	result.Id = id

	for _, query := range queies {
		rows, err := l.db.NamedQuery(query, queryParams)
		if err != nil {
			log.Error(err)
			return nil, exception.InternalServerError
		}

		// single column query
		if query == clickBasedQuery || query == userBasedQuery {
			if !rows.Next() {
				log.Error("clickhouse count query error")
				return nil, exception.InternalServerError
			}

			var total int
			err := rows.Scan(&total)
			if err != nil {
				log.Error(err)
				return nil, exception.InternalServerError
			}
			switch query {
			case userBasedQuery:
				result.Visitors.Total = total
			case clickBasedQuery:
				result.Clicks.Total = total
			}
			continue
		}

		var val = make([]map[string]interface{}, 0)
		i := 0
		for rows.Next() {
			val = append(val, map[string]interface{}{})
			err = rows.MapScan(val[i])
			i++
		}

		var dst *map[string]int
		len := i
		switch query {
		case clickBasedGroupByBrowser:
			dst = &result.Clicks.PerBrowser
			*dst = make(map[string]int)
			for i := 0; i < len; i++ {
				(*dst)[val[i]["browser"].(string)] = int(val[i]["COUNT()"].(uint64))
			}
		case clickBasedGroupByDevice:
			dst = &result.Clicks.PerDevice
			*dst = make(map[string]int)
			for i := 0; i < len; i++ {
				(*dst)[val[i]["device"].(string)] = int(val[i]["COUNT()"].(uint64))
			}
		case userBasedGroupByBrowser:
			dst = &result.Visitors.PerBrowser
			*dst = make(map[string]int)
			for i := 0; i < len; i++ {
				(*dst)[val[i]["browser"].(string)] = int(val[i]["uniq(remote_address)"].(uint64))
			}
		case userBasedGroupByDevice:
			dst = &result.Visitors.PerDevice
			*dst = make(map[string]int)
			for i := 0; i < len; i++ {
				(*dst)[val[i]["device"].(string)] = int(val[i]["uniq(remote_address)"].(uint64))
			}
		}
	}

	return result, nil
}

func (l *Log) Save(elem []*model.Log, len int) error {
	// begin transaction for batch insert
	tx, err := l.db.Begin()
	if err != nil {
		log.Error(err)
		return exception.InternalServerError
	}

	// make prepared statement
	stmt, err := tx.Prepare("INSERT INTO logs (id, created_at, device, browser, remote_address) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return exception.InternalServerError
	}
	defer stmt.Close()

	for i := 0; i < len; i++ {
		createTime, _ := time.Parse(time.RFC3339, elem[i].CreatedAt)
		if _, err = stmt.Exec(
			elem[i].Id,
			createTime,
			elem[i].Device,
			elem[i].Browser,
			elem[i].RemoteAddr); err != nil {
			log.Error(err)
			tx.Rollback()
			return exception.InternalServerError
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return exception.InternalServerError
	}

	return nil
}
