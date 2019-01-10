#!/bin/sh

# vars
new=$1
dir=$(echo $path | cut -d'/' -f1);
suffix='lambda'

declare -a lambdas
lambdas=(
	'connectionmanager' 
	'channelmanager' 
	'messagemanager'
)



buildanddep() {
	path=$1/main.go
	# build
	GOOS=linux go build -o $1-$suffix $path
	zip -j $1.zip $1-$suffix
	chmod 777 $1.zip

	# publish
	if [[ $new = 'new' ]]; then
		echo "creating"
		aws lambda create-function --profile jds --function-name $1 --runtime go1.x --role arn:aws:iam::671958020402:role/iam_for_lambda --handler $1-$suffix --zip-file fileb://./$1.zip
	else
		echo "updating"
		aws lambda update-function-code --profile jds --function-name $1 --zip-file fileb://./$1.zip 
	fi;

	rm $1.zip $1-$suffix
}

for i in "${lambdas[@]}"
do
	echo "lambda: " $i
	buildanddep $i
done;

# update config
# aws lambda --profile jds update-function-configuration --function-name ssmssh --handler ssmssh-lambda