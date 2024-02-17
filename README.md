# XGO

## Usage

### As submodule

```bash
git submodule add git@github.com:anoaland/xgo.git example/xgo
```

### Project Integration

#### Workspace

Add `./example/xgo` to your `go.work` file, example:

```go
go 1.21.6

use (
	./example/api
	./example/xgo
)
```

#### Libs

Add required libraries to your main project:

```bash
go get github.com/gofiber/fiber/v2
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres       # or other db you want, see: https://gorm.io/docs/#Install
go get ariga.io/atlas-go-sdk/atlasexec
go get ariga.io/atlas-provider-gorm     # https://github.com/ariga/atlas-provider-gorm
go get github.com/Nerzal/gocloak/v13
go get github.com/joho/godotenv
```
