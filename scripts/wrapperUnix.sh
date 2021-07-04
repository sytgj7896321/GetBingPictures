#!/bin/bash
exiftool -overwrite_original -artist="$1" "$3"
exiftool -overwrite_original -usercomment="$2" "$3"