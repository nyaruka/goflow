package handler

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/nyaruka/goflow/handler/config"
)

type server struct {
	config *config.Handler

	db        *sqlx.DB
	redisPool *redis.Pool
}

func NewServer(config *config.Handler) *server {
	return &server{config: config}
}

func (s *server) Start() error {
	// parse and test our db config
	dbURL, err := url.Parse(s.config.DB)
	if err != nil {
		return fmt.Errorf("unable to parse DB URL '%s': %s", s.config.DB, err)
	}

	if dbURL.Scheme != "postgres" {
		return fmt.Errorf("invalid DB URL: '%s', only postgres is supported", s.config.DB)
	}

	// test our db connection
	db, err := sqlx.Connect("postgres", s.config.DB)
	if err != nil {
		log.Printf("[ ] DB: error connecting: %s\n", err)
	} else {
		log.Println("[X] DB: connection ok")
	}
	s.db = db

	// parse and test our redis config
	redisURL, err := url.Parse(s.config.Redis)
	if err != nil {
		return fmt.Errorf("unable to parse Redis URL '%s': %s", s.config.Redis, err)
	}

	// create our pool
	redisPool := &redis.Pool{
		Wait:        true,              // makes callers wait for a connection
		MaxActive:   5,                 // only open this many concurrent connections at once
		MaxIdle:     2,                 // only keep up to 2 idle
		IdleTimeout: 240 * time.Second, // how long to wait before reaping a connection
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s", redisURL.Host))
			if err != nil {
				return nil, err
			}

			// switch to the right DB
			_, err = conn.Do("SELECT", strings.TrimLeft(redisURL.Path, "/"))
			return conn, err
		},
	}
	s.redisPool = redisPool

	// test our redis connection
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		log.Printf("[ ] Redis: error connecting: %s\n", err)
	} else {
		log.Println("[X] Redis: connection ok")
	}

	/*
		// create our s3 client
		s3Session, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(s.config.AWS_Access_Key_ID, s.config.AWS_Secret_Access_Key, ""),
			Region:      aws.String(s.config.S3_Region),
		})
		if err != nil {
			return err
		}
		s.s3Client = s3.New(s3Session)

		// test out our S3 credentials
		err = testS3(s)
		if err != nil {
			log.Printf("[ ] S3: bucket inaccessible, media may not save: %s\n", err)
		} else {
			log.Println("[X] S3: bucket accessible")
		}
	*/
	return nil
}
