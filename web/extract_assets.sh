#!/bin/bash

assets_source_folder='./src/ms2k/assets/files/'
assets_dest_folder='./web/assets/'

IFS=$'\n'
for filename in $(find -type f -path ${assets_source_folder}'**/*' | grep -v '.json'); do
    mkdir -p ${assets_dest_folder}$(dirname ${filename#$assets_source_folder})
    mv ${filename} ${assets_dest_folder}${filename#$assets_source_folder}
done
