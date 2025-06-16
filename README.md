# Project Status : As of now, Seacrate does not offer any kind of authentication, and as such, is not recommended for any use. It's still a work in progress project built for fun.

# Seacrate

Seacrate is a simple key-value store for all your infrastructure secrets.

```bash
Seacrate is an easy secret management application

Usage:
  seacrate [flags]
  seacrate [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Init the application database and create an encryption key
  run         Start the server
  validate    Validate the configuration file is valid
  version     Obtain information about version and compilation

Flags:
  -h, --help   help for seacrate

Use "seacrate [command] --help" for more information about a command.
```

## Configuration

The first step is to configure Seacrate through a configuration file. Here is an example of a configuration file.

```yaml
encryption:
  algorithm: aes

database:
  host: 127.0.0.1
  port: 5432
  database: seacrate
  username: seacrate
  password: seacrate
```

## Initialize the application

Then once the configuration is ready, you can initialize the application through the use of the `init` command.

```bash
./seacrate init
```

You'll be asked of few questions.

```
2
How many key part are required to unseal the instance ?
2
Create file `results.json`
```

At the end it should have created a `results.json` file.

```
{
  "keys": [
    "CWg7Lx4olu4NCsEmQdY858M+bLKya8zILTgcfzpwLRPB",
    "1x/QyGufBzEIjXUmo0PdIcxHi/tRl5fod0LdAZcXX96e"
  ]
}
```

## Unsealing Seacrate

By default Seacrate is in a sealed status, meaning that it does not hold the key that will allow him to decrypt and encrypt seacrate. To unseal it you have to give him the keys shard that have been created.

```bash
curl -XPOST -H 'Content-Type: application/json' http://127.0.0.1:3000/api/v1/system/seal -d '{"part": "GiQXzjIUyDZhrcZtXQwjvUagalpanHP/r3/w0IzaP861"}'
curl -XPOST -H 'Content-Type: application/json' http://127.0.0.1:3000/api/v1/system/seal -d '{"part": "JNEP+/6SJkamMaPdGHBbaGHFeZBcCpq8H40AOBsFz8yl"}'
```

The application will let you know once it's unsealed.

## Using The K/V Store

You can then use the K/V store to store any secrets.

```bash
curl -XPOST -H 'Content-Type: application/json' http://127.0.0.1:3000/api/v1/secrets/facebook -d '{"value": "MyFacebookPassword92120!"}'
```
## Building It

You can simply build the application by running the following command

```bash
go build -o seacrate main.go
```

# TODO

In order for the application to be considered usable in any way it should : 

- Have the ability to take in a SSL/TLS certificate for the HTTP listener in order to prevent secrets to transit clearly at any point;
- Have a proper authentication solution, with IAM capabilities (through the use of an ACL for example);
- Make sure that the application can be deployed as a cluster and work flawlessly.
