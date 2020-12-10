package dockertest
// as the package suggests, this is not meant for use in an actual scenario.
// it is meant to set up a postgresql instance using dockertest.
import (
	"errors"
	"fmt"
	_ "github.com/gobuffalo/pop/v5"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const hostname = "localhost" // make this configurable?

var resources = []*dockertest.Resource{}
var pool *dockertest.Pool


func getRunOpts(containerExposedPort, containerName, pgUsername, pgPassword, pgDbName string) dockertest.RunOptions {
	opts := dockertest.RunOptions{
		Repository: "postgresql",
		Tag:        "12.5-alpine",
		Env: []string{
			"POSTGRES_USER=" + pgUsername,
			"POSTGRES_PASSWORD=" + pgPassword,
			"POSTGRES_DB=" + pgDbName,
		},
		ExposedPorts: []string{"5432"},
		Name:         containerName,
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{
					HostIP:   "0.0.0.0",
					HostPort: containerExposedPort,
				},
			},
		},
	}
	return opts
}
o
// runs postgresql based on the variables passed into it.
func RunTestPostgreSQL(t *testing.T, containerName, containerExposedPort, pgUsername, pgPassword, pgDbName string) {
	opts := getRunOpts(containerExposedPort,containerName,pgUsername,pgPassword,pgDbName)
	_, err := initalizePostgresDb(opts)
	require.NoError(t,err)
	bootstrap(t,containerExposedPort,pgUsername,pgPassword, pgDbName)
}

func initalizePostgresDb(opts dockertest.RunOptions)(*dockertest.Resource, error){
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to docker")
	}
	resource, err := pool.RunWithOptions(opts)

	if err == nil {
		resources = append(resources, resource)
	}
	return resource, err
}


func bootstrap(t *testing.T,containerExposedPort, pgUsername, pgPassword, pgDbName string) (db *sqlx.DB){
	// uses sqlx to test for an alive connection,
	if err := Retry(time.Second*5, time.Minute*5, func() error {
		databaseConnStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			hostname,
			containerExposedPort,
			pgUsername,
			pgDbName,
			pgPassword)
		var err error
		db, err = sqlx.Connect("postgres", databaseConnStr)
		require.NoError(t, err)

		return db.Ping()
	}); err != nil {
		if pErr := pool.Purge(resource); pErr != nil {
			log.Fatalf("Could not connect to docker and unable to remove image: %s - %s", err, pErr)
			require.NoError(t, pErr)
		}
		log.Fatalf("Could not connect to docker: %s", err)
		require.NoError(t, err)
	}
	return
}
