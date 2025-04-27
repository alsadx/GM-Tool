module campaigntool

go 1.24.2

require (
	auth v0.0.0
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/fatih/color v1.18.0
	github.com/golang-migrate/migrate/v4 v4.18.2
	github.com/golang/mock v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v5 v5.7.4
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.72.0
)

require (
	github.com/alsadx/protos v0.0.19
	protos v0.0.0-00010101000000-000000000000
)

replace auth => ../auth

replace protos => ../protos

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
