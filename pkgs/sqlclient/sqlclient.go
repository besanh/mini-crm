package sqlclient

import (
	"errors"
	"fmt"
	"net/url"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/besanh/mini-crm/common/log"
	"github.com/besanh/mini-crm/ent"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	MYSQL      = "mysql"
	POSTGRESQL = "postgresql"
)

type (
	ISqlClientConn interface {
		GetDB() *ent.Client
		GetDriver() string
		Connect() (err error)
	}

	SqlConfig struct {
		Driver       string
		Host         string
		Port         int
		Database     string
		Username     string
		Password     string
		Timeout      int
		DialTimeout  int
		ReadTimeout  int
		WriteTimeout int
		PoolSize     int
		MaxIdleConns int
		MaxOpenConns int
	}

	SqlClientConn struct {
		SqlConfig
		DB *ent.Client
	}
)

func NewSqlClient(config SqlConfig) ISqlClientConn {
	client := &SqlClientConn{
		SqlConfig: config,
	}

	if err := client.Connect(); err != nil {
		log.Fatal("Connect error:", err)
		return nil
	}

	log.Info("Connected to database")

	return client
}

func (c *SqlClientConn) Connect() (err error) {
	switch c.Driver {
	case MYSQL:
		//username:password@protocol(address)/dbname?param=value
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?readTimeout=%ds&writeTimeout=%ds", c.Username, c.Password, c.Host, c.Port, c.Database, c.ReadTimeout, c.WriteTimeout)
		driver, err := sql.Open(dialect.MySQL, dsn)
		if err != nil {
			return err
		}
		driver.DB().SetMaxIdleConns(c.MaxIdleConns)
		driver.DB().SetMaxOpenConns(c.MaxOpenConns)

		if err := driver.DB().Ping(); err != nil {
			log.Error(err)
		}

		client := ent.NewClient(ent.Driver(driver))
		c.DB = client

		return nil
	case POSTGRESQL:
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
			c.Username,
			url.QueryEscape(c.Password),
			c.Host,
			c.Port,
			c.Database,
			c.Timeout,
		)

		driver, err := sql.Open("pgx", dsn)
		if err != nil {
			return err
		}
		driver.DB().SetMaxIdleConns(c.MaxIdleConns)
		driver.DB().SetMaxOpenConns(c.MaxOpenConns)

		if err := driver.DB().Ping(); err != nil {
			log.Error(err)
		}

		client := ent.NewClient(ent.Driver(driver))
		c.DB = client

		return nil
	default:
		log.Fatal("driver is missing")
		return errors.New("driver is missing")
	}
}

func (c *SqlClientConn) GetDB() *ent.Client {
	return c.DB
}

func (c *SqlClientConn) GetDriver() string {
	return c.Driver
}
