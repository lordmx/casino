package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type App struct {
	Config     *Config
	Location   *time.Location
	DB         *sql.DB

	configPath string
}

func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close()
	}
}

func NewApp(configPath string) (app *App, err error) {
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("Cannot find config file " + configPath)
	}

	app = &App{}

	app.configPath = configPath
	app.Config, err = app.getConfig()

	if err != nil {
		return nil, err
	}

	err = app.init()
	go app.watchConfig()

	return app, err
}

func (app *App) init() error {
	var err error

	app.DB, err = app.getPsql()

	if err != nil {
		return err
	}

	timezone := "UTC"

	if app.Config.Timezone != "" {
		timezone = app.Config.Timezone
	}

	location, err := time.LoadLocation(timezone)

	if err != nil {
		return err
	}

	app.Location = location

	return nil
}

func (app *App) getConfig() (*Config, error) {
	return NewConfig(app.configPath)
}

func (app *App) getPsql() (*sql.DB, error) {
	config := app.Config

	db, err := NewPsql(
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Database,
		config.DB.Encoding,
		config.Timezone,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *App) watchConfig() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write && event.Name == app.configPath {
					log.Printf("[%s] modified config: %s", time.Now().Format(time.RFC3339), event.Name)

					app.Config, err = app.getConfig()
					
					if err != nil {
						log.Fatal(err)
					}
					
					app.Close()
					err = app.init()

					if err != nil {
						log.Fatal(err)
					}
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(app.configPath)

	if err != nil {
		log.Fatal(err)
	}

	<-done
}