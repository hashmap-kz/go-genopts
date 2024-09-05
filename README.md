# Generate bash getopt boilerplate.

This tool is designed for fast setup bash scripts boilerplates.

It generates the main function, usage, argument parsing routine and checks.

Mostly all parameters are configured.

### Input config example

```
opts:
  # required args
  - name: dbname
    desc: "database to dump"
  - name: host
    desc: "database server host or socket directory"
  - name: port
    desc: "database server port number"
  - name: username
    short: U
    desc: "connect as specified database user"
  - name: output
    short: O
    desc: "output path"
  # optional args
  - name: schema
    type: list
    short: "n"
    desc: "dump only schemas matching pattern"
    optional: true
  - name: exclude-schema
    type: list
    short: "N"
    desc: "do not dump any schemas matching pattern"
    optional: true
  - name: verbose
    defaultValue: "true"
    type: bool
    optional: true

```

### Generated code

```
#!/bin/bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $(basename "$0") [OPTION]

Options:
  -d, --dbname           database to dump
  -h, --host             database server host or socket directory
  -p, --port             database server port number
  -U, --username         connect as specified database user
  -O, --output           output path
  -n, --schema           dump only schemas matching pattern
  -N, --exclude-schema   do not dump any schemas matching pattern
  -v, --verbose
EOF
}

main() {
  local dbname=''
  local host=''
  local port=''
  local username=''
  local output=''
  local schema=()
  local exclude_schema=()
  local verbose="true"

  VALID_ARGS=$(getopt -o d:h:p:U:O:n:N:v --long dbname:,host:,port:,username:,output:,schema:,exclude-schema:,verbose -- "$@")

  # shellcheck disable=SC2181
  if [ $? != 0 ]; then
    printf "error parsing options"
    usage
    exit 1
  fi

  eval set -- "$VALID_ARGS"
  while true; do
    case "$1" in

    -d | --dbname)
      dbname="${2}"
      shift 2
      ;;
    -h | --host)
      host="${2}"
      shift 2
      ;;
    -p | --port)
      port="${2}"
      shift 2
      ;;
    -U | --username)
      username="${2}"
      shift 2
      ;;
    -O | --output)
      output="${2}"
      shift 2
      ;;
    -n | --schema)
      schema+=("${2}")
      shift 2
      ;;
    -N | --exclude-schema)
      exclude_schema+=("${2}")
      shift 2
      ;;
    -v | --verbose)
      verbose=true
      shift
      ;;

    --)
      shift
      break
      ;;
    *)
      printf "unexpected argument ${1}"
      usage
      exit 1
      ;;
    esac
  done

  # check remaining
  shift $((OPTIND - 1))
  remaining_args="${*}"
  if [ -n "${remaining_args}" ]; then
    echo "remaining args are not allowed: ${remaining_args[*]}"
    usage
    exit 1
  fi

  # set checks
  if [ -z "${dbname}" ]; then
    printf "\n[error] required arg is empty: dbname\n\n"
    usage
    exit 1
  fi
  if [ -z "${host}" ]; then
    printf "\n[error] required arg is empty: host\n\n"
    usage
    exit 1
  fi
  if [ -z "${port}" ]; then
    printf "\n[error] required arg is empty: port\n\n"
    usage
    exit 1
  fi
  if [ -z "${username}" ]; then
    printf "\n[error] required arg is empty: username\n\n"
    usage
    exit 1
  fi
  if [ -z "${output}" ]; then
    printf "\n[error] required arg is empty: output\n\n"
    usage
    exit 1
  fi

  # debug variables
  echo "dbname=${dbname}"
  echo "host=${host}"
  echo "port=${port}"
  echo "username=${username}"
  echo "output=${output}"
  echo "schema=${schema[*]}"
  echo "exclude_schema=${exclude_schema[*]}"
  echo "verbose=${verbose}"
}

main "${@}"

```

### Generate script:

```
go run main.go -config=config.yml > test.sh
```

### Example usage of final script:

```
bash test.sh -d keycloak_base -h 10.40.240.30 -p 5432 -U postgres --verbose -O "/mnt/backup" -n "public|data_audit"
bash test.sh --dbname=keycloak_base --host=10.40.240.30 --port=5432 --username=postgres --verbose --output="/mnt/backup" --schema="public|data_audit"
```
