package mailroom

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/nyaruka/goflow/mailroom/config"
	"github.com/nyaruka/goflow/mailroom/store"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Mailroom is our interface for the main Mailroom app
type Mailroom interface {
	Start() error
	Stop()
}

type mailroom struct {
	config *config.Mailroom

	db        *sqlx.DB
	redisPool *redis.Pool
	s3Client  *s3.S3

	httpServer *http.Server
	router     *mux.Router

	awsCreds *credentials.Credentials

	waitGroup  *sync.WaitGroup
	stopChan   chan bool
	workerChan chan *worker
	stopped    bool
}

// New creates and returns a new Mailroom instance given the passed in config
func New(config *config.Mailroom) Mailroom {
	return &mailroom{
		config: config,

		stopChan:   make(chan bool),
		workerChan: make(chan *worker, config.Workers),
		waitGroup:  &sync.WaitGroup{},
	}
}

func (m *mailroom) Start() error {
	// parse and test our db config
	dbURL, err := url.Parse(m.config.DB)
	if err != nil {
		return fmt.Errorf("unable to parse DB URL '%s': %s", m.config.DB, err)
	}

	if dbURL.Scheme != "postgres" {
		return fmt.Errorf("invalid DB URL: '%s', only postgres is supported", m.config.DB)
	}

	// test our db connection
	db, err := sqlx.Connect("postgres", m.config.DB)
	if err != nil {
		log.Printf("[ ] DB: error connecting: %s\n", err)
	} else {
		log.Println("[X] DB: connection ok")
	}
	m.db = db

	// parse and test our redis config
	redisURL, err := url.Parse(m.config.Redis)
	if err != nil {
		return fmt.Errorf("unable to parse Redis URL '%s': %s", m.config.Redis, err)
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
	m.redisPool = redisPool

	// test our redis connection
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		log.Printf("[ ] Redis: error connecting: %s\n", err)
	} else {
		log.Println("[X] Redis: connection ok")
	}

	// create our s3 client
	s3Session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(m.config.AWS_Access_Key_ID, m.config.AWS_Secret_Access_Key, ""),
		Region:      aws.String(m.config.S3_Region),
	})
	if err != nil {
		return err
	}
	m.s3Client = s3.New(s3Session)

	// test out our S3 credentials
	err = testS3(m)
	if err != nil {
		log.Printf("[ ] S3: bucket inaccessible, media may not be accessible: %s\n", err)
	} else {
		log.Println("[X] S3: bucket accessible")
	}

	// stand up our http server
	m.router = mux.NewRouter()
	m.router.HandleFunc("/", m.handleIndex)

	m.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", m.config.Port),
		Handler:      m.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// start serving HTTP
	go func() {
		m.waitGroup.Add(1)
		defer m.waitGroup.Done()
		err := m.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("ERROR: %s", err)
		}
	}()

	// start our fanout
	startFanout(m)

	return nil
}

func (m *mailroom) Stop() {
	log.Println("Stopping mailroom processes")

	if m.db != nil {
		m.db.Close()
	}

	if m.redisPool != nil {
		m.redisPool.Close()
	}

	m.stopped = true
	close(m.stopChan)

	// shut down our HTTP server
	if err := m.httpServer.Shutdown(nil); err != nil {
		log.Printf("ERROR gracefully shutting down server: %s\n", err)
	}

	m.waitGroup.Wait()

	log.Printf("[X] Mailroom: stopped\n")
}

func (m *mailroom) popNextMsg() (*store.Msg, error) {
	return nil, nil
}

func (m *mailroom) handleIndex(w http.ResponseWriter, r *http.Request) {
	// test redis
	rc := m.redisPool.Get()
	_, redisErr := rc.Do("PING")
	defer rc.Close()

	// test our db
	_, dbErr := m.db.Exec("SELECT 1")

	if redisErr == nil && dbErr == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	var buf bytes.Buffer
	buf.WriteString("<title>mailroom</title><body><pre>\n")
	buf.WriteString(splash)

	if redisErr != nil {
		buf.WriteString(fmt.Sprintf("\n% 16s: %v", "redis err", redisErr))
	}
	if dbErr != nil {
		buf.WriteString(fmt.Sprintf("\n% 16s: %v", "db err", dbErr))
	}

	buf.WriteString("\n\n")
	buf.WriteString("</pre></body>")
	w.Write(buf.Bytes())
}

var splash = `
         __  ___        _  __                             
        /  |/  /____ _ (_)/ /_____ ____   ____   ____ ___ 
       / /|_/ // __ '// // // ___// __ \ / __ \ / __ '__ \
      / /  / // /_/ // // // /   / /_/ // /_/ // / / / / /
     /_/  /_/ \__,_//_//_//_/    \____/ \____//_/ /_/ /_/ 
`
