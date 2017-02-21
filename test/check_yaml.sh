#! /bin/bash

for yml in ../*.yaml; do 
  if ! ruby -e "require 'yaml';puts YAML.load_file('./$yml')" > /dev/null  2>&1; then
     echo -e "Error ian $yml";
     exit 1
  fi
done
echo "Check successfull!!"
