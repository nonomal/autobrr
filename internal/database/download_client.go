package database

import (
	"database/sql"
	"encoding/json"

	"github.com/autobrr/autobrr/internal/domain"

	"github.com/rs/zerolog/log"
)

type DownloadClientRepo struct {
	db *sql.DB
}

func NewDownloadClientRepo(db *sql.DB) domain.DownloadClientRepo {
	return &DownloadClientRepo{db: db}
}

func (r *DownloadClientRepo) List() ([]domain.DownloadClient, error) {

	rows, err := r.db.Query("SELECT id, name, type, enabled, host, port, ssl, username, password, settings FROM client")
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not query download client rows")
		return nil, err
	}

	defer rows.Close()

	clients := make([]domain.DownloadClient, 0)

	for rows.Next() {
		var f domain.DownloadClient
		var settingsJsonStr string

		if err := rows.Scan(&f.ID, &f.Name, &f.Type, &f.Enabled, &f.Host, &f.Port, &f.SSL, &f.Username, &f.Password, &settingsJsonStr); err != nil {
			log.Error().Stack().Err(err).Msg("could not scan download client to struct")
			return nil, err
		}

		if settingsJsonStr != "" {
			if err := json.Unmarshal([]byte(settingsJsonStr), &f.Settings); err != nil {
				log.Error().Stack().Err(err).Msgf("could not marshal download client settings %v", settingsJsonStr)
				return nil, err
			}
		}

		clients = append(clients, f)
	}
	if err := rows.Err(); err != nil {
		log.Error().Stack().Err(err).Msg("could not query download client rows")
		return nil, err
	}

	return clients, nil
}

func (r *DownloadClientRepo) FindByID(id int32) (*domain.DownloadClient, error) {

	query := `
		SELECT id, name, type, enabled, host, port, ssl, username, password, settings FROM client WHERE id = ?
	`

	row := r.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		log.Error().Stack().Err(err).Msg("could not query download client rows")
		return nil, err
	}

	var client domain.DownloadClient
	var settingsJsonStr string

	if err := row.Scan(&client.ID, &client.Name, &client.Type, &client.Enabled, &client.Host, &client.Port, &client.SSL, &client.Username, &client.Password, &settingsJsonStr); err != nil {
		log.Error().Stack().Err(err).Msg("could not scan download client to struct")
		return nil, err
	}

	if settingsJsonStr != "" {
		if err := json.Unmarshal([]byte(settingsJsonStr), &client.Settings); err != nil {
			log.Error().Stack().Err(err).Msgf("could not marshal download client settings %v", settingsJsonStr)
			return nil, err
		}
	}

	return &client, nil
}

func (r *DownloadClientRepo) Store(client domain.DownloadClient) (*domain.DownloadClient, error) {
	var err error

	settings := domain.DownloadClientSettings{
		APIKey: client.Settings.APIKey,
		Basic:  client.Settings.Basic,
		Rules:  client.Settings.Rules,
	}

	settingsJson, err := json.Marshal(&settings)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("could not marshal download client settings %v", settings)
		return nil, err
	}

	if client.ID != 0 {
		_, err = r.db.Exec(`
			UPDATE 
    			client 
			SET 
			    name = ?, 
			    type = ?, 
			    enabled = ?, 
			    host = ?, 
			    port = ?, 
			    ssl = ?, 
			    username = ?, 
			    password = ?, 
			    settings = json_set(?) 
			WHERE
			    id = ?`,
			client.Name,
			client.Type,
			client.Enabled,
			client.Host,
			client.Port,
			client.SSL,
			client.Username,
			client.Password,
			settingsJson,
			client.ID,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("could not update download client: %v", client)
			return nil, err
		}
	} else {
		var res sql.Result

		res, err = r.db.Exec(`INSERT INTO 
    		client(
    		       name,
    		       type, 
    		       enabled,
    		       host,
    		       port,
    		       ssl,
    		       username,
    		       password,
    		       settings)
			VALUES (?, ?, ?, ?, ?, ? , ?, ?, json_set(?)) ON CONFLICT DO NOTHING`,
			client.Name,
			client.Type,
			client.Enabled,
			client.Host,
			client.Port,
			client.SSL,
			client.Username,
			client.Password,
			settingsJson,
		)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("could not store new download client: %v", client)
			return nil, err
		}

		resId, _ := res.LastInsertId()
		client.ID = int(resId)

		log.Trace().Msgf("download_client: store new record %d", client.ID)
	}

	log.Info().Msgf("store download client: %v", client.Name)
	log.Trace().Msgf("store download client: %+v", client)

	return &client, nil
}

func (r *DownloadClientRepo) Delete(clientID int) error {
	res, err := r.db.Exec(`DELETE FROM client WHERE client.id = ?`, clientID)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("could not delete download client: %d", clientID)
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return err
	}

	log.Info().Msgf("delete download client: %d", clientID)

	return nil
}
