#!/bin/sh -e

OPENAPI_VALIDATOR="swagger-cli"
OPENAPI_VALIDATOR_INSTALL="npm -g install @apidevtools/swagger-cli"
OPENAPI_VALIDATOR_UPGRADE="npm -g install @apidevtools/swagger-cli@latest"

OPENAPI_BUNDLER="swagger-cli"
OPENAPI_BUNDLER_INSTALL="npm -g install @apidevtools/swagger-cli"
OPENAPI_BUNDLER_UPGRADE="npm -g install @apidevtools/swagger-cli@latest"

OPENAPI_GENERATOR="${HOME}/bin/oapi-codegen"
OPENAPI_GENERATOR_REPO="github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
OPENAPI_GENERATOR_INSTALL="go get $OPENAPI_GENERATOR_REPO"
OPENAPI_GENERATOR_UPGRADE="go get -u $OPENAPI_GENERATOR_REPO"

openapi_usage()
{
	echo "Usage: $0 ( validate INDEXFILE | bundle INDEXFILE SPECFILE | generate SPECFILE OUTDIR PACKAGE | upgrade )"
	exit 1
}

openapi_check()
{
	which $1 >/dev/null 2>&1
}

openapi_validate()
{
	# check argument
	[ $# -eq 1 ] || openapi_usage

	OPENAPI_INDEXFILE=$1

	# install if not installed
	openapi_check $OPENAPI_VALIDATOR || $OPENAPI_VALIDATOR_INSTALL

	# validate
	$OPENAPI_VALIDATOR validate $OPENAPI_INDEXFILE
}

openapi_bundle()
{
	# check argument
	[ $# -eq 2 ] || openapi_usage

	OPENAPI_INDEXFILE=$1
	OPENAPI_SPECFILE=$2

	# install if not installed
	openapi_check $OPENAPI_BUNDLER || $OPENAPI_BUNDLER_INSTALL

	# bundle
	$OPENAPI_BUNDLER bundle -t yaml $OPENAPI_INDEXFILE -o $OPENAPI_SPECFILE
}

openapi_generate()
{
	# check argument
	[ $# -eq 2 ] || openapi_usage

	OPENAPI_SPECFILE=$1
	OPENAPI_OUTDIR=$2

	# install if not installed
	openapi_check $OPENAPI_GENERATOR || $OPENAPI_GENERATOR_INSTALL

	# generate
	mkdir -p $OPENAPI_OUTDIR
	for i in fiber-server types spec; do
	  echo $OPENAPI_GENERATOR --package api -generate $i -o $OPENAPI_OUTDIR/$i.go $OPENAPI_SPECFILE
		$OPENAPI_GENERATOR --package api -generate $i -o $OPENAPI_OUTDIR/$i.go $OPENAPI_SPECFILE
		echo $OPENAPI_OUTDIR/$i.go generated
	done
}

openapi_upgrade()
{
	$OPENAPI_VALIDATOR_UPGRADE
	$OPENAPI_BUNDLER_UPGRADE
	$OPENAPI_GENERATOR_UPGRADE
}

case $1 in
validate|bundle|generate|upgrade)
	FUNC=$1
	shift
	openapi_$FUNC $@
	;;
*)
	openapi_usage
	;;
esac