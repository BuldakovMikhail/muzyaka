package testhelpers

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	tc "github.com/romnn/testcontainers"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// Options ...
type Options struct {
	tc.ContainerOptions
	ImageTag     string
	RootUser     string
	RootPassword string
}

// Container ...
type Container struct {
	Container testcontainers.Container
	tc.ContainerConfig
	Host         string
	Port         uint
	RootUser     string
	RootPassword string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),

		//src/muzyaka/internal/domain/album/repository/postgres/postgres_test.go
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "..", "..", "..", "database", "docker-entrypoint-initdb.d", "01-init.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Minute)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

// Terminate ...
func (c *Container) Terminate(ctx context.Context) {
	if c.Container != nil {
		c.Container.Terminate(ctx)
	}
}

// ConnectionURI ...
func (c *Container) ConnectionURI() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Start ...
func Start(ctx context.Context, options Options) (Container, error) {
	var container Container
	container.RootUser = options.RootUser
	container.RootPassword = options.RootPassword

	port, err := nat.NewPort("", "9000")
	if err != nil {
		return container, fmt.Errorf("failed to build port: %v", err)
	}

	env := make(map[string]string)
	if options.RootUser != "" && options.RootPassword != "" {
		env["MINIO_ROOT_USER"] = options.RootUser
		env["MINIO_ROOT_PASSWORD"] = options.RootPassword
	}

	tag := "latest"
	if options.ImageTag != "" {
		tag = options.ImageTag
	}

	timeout := options.ContainerOptions.StartupTimeout
	if int64(timeout) < 1 {
		timeout = 5 * time.Minute // Default timeout
	}

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("minio/minio:%s", tag),
		Env:          env,
		Cmd:          []string{"server", "/data"},
		ExposedPorts: []string{string(port)},
		WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(timeout),
	}

	tc.MergeRequest(&req, &options.ContainerOptions.ContainerRequest)

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return container, fmt.Errorf("failed to start container: %v", err)
	}
	container.Container = minioContainer

	host, err := minioContainer.Host(ctx)
	if err != nil {
		return container, fmt.Errorf("failed to get container host: %v", err)
	}
	container.Host = host

	realPort, err := minioContainer.MappedPort(ctx, port)
	if err != nil {
		return container, fmt.Errorf("failed to get exposed container port: %v", err)
	}
	container.Port = uint(realPort.Int())

	return container, nil
}
