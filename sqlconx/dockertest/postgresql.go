package dockertest
// as the package suggests, this is not meant for use in an actual scenario.
// it is meant to set up a postgresql instance using dockertest.
import (
	"fmt"
	_ "github.com/gobuffalo/pop/v5"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"github.com/pkg/errors"
)

const hostname = "localhost" // make this configurable?

var resources = []*dockertest.Resource{}
var pool *dockertest.Pool


// KillAllTestDatabases deletes all test databases.
func KillAllTestDatabases() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	for _, r := range resources {
		if err := pool.Purge(r); err != nil {
			panic(err)
		}
	}
}



// Retry executes a f until no error is returned or failAfter is reached.
func Retry(maxWait time.Duration, failAfter time.Duration, f func() error) (err error) {
	var lastStart time.Time
	err = errors.New("did not connect")
	loopWait := time.Millisecond * 100
	retryStart := time.Now().UTC()
	for retryStart.Add(failAfter).After(time.Now().UTC()) {
		lastStart = time.Now().UTC()
		if err = f(); err == nil {
			return nil
		}

		if lastStart.Add(maxWait * 2).Before(time.Now().UTC()) {
			retryStart = time.Now().UTC()
		}

		log.Errorf("Retrying in %f seconds...\n", loopWait.Seconds())
		time.Sleep(loopWait)
		loopWait = loopWait * time.Duration(int64(2))
		if loopWait > maxWait {
			loopWait = maxWait
		}
	}
	return err
}

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
		PortBindings: map[dc.Port][]dc.PortBinding{
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

// runs postgresql based on the variables passed into it.
func RunTestPostgreSQL(t *testing.T, containerName, containerExposedPort, pgUsername, pgPassword, pgDbName string) {
	opts := getRunOpts(containerExposedPort,containerName,pgUsername,pgPassword,pgDbName)
	resource, err := initalizePostgresDb(opts)
	require.NoError(t,err)
	bootstrap(t,containerExposedPort,pgUsername,pgPassword, pgDbName, resource)
}

func initalizePostgresDb(opts dockertest.RunOptions)(*dockertest.Resource, error){
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to docker")
	}
	resource, err := pool.RunWithOptions(&opts)

	if err == nil {
		resources = append(resources, resource)
	}
	return resource, err
}


func bootstrap(t *testing.T,containerExposedPort, pgUsername, pgPassword, pgDbName string, resource *dockertest.Resource) (db *sqlx.DB){
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
