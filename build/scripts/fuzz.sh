#!/bin/bash

set -e

fuzzTime=${1:-10s}

files=$(grep -r --include='**_test.go' --files-with-matches 'func Fuzz' .)

for file in ${files}
do
	funcs=$(grep -o 'func Fuzz\w*' $file | sed 's/func //')
	for func in ${funcs}
	do
		echo "Fuzzing $func in $file"
		parentDir=$(dirname $file)
        # Add $ to the end of the func name, since go fuzzing gets stressed out by ambiguiity
        # ie "testing: will not fuzz, -fuzz matches more than one fuzz test: [FuzzCompareControls FuzzCompareControlsInt]" 
		go test $parentDir -run=$func -fuzz=$func\$ -fuzztime=${fuzzTime}
	done
done