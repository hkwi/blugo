# args:
# 1 : base
# 2 : seed suffix
# 3 : out os
# 4 : out arch
set -e
# arm, arm64 char defaults to unsigned-char and the others use singed-char
# we force unsigned-char for string, []byte, []uint8 conversion
go tool cgo -godefs -- -funsigned-char $1_$2.go | gofmt > $1_$3_$4.go
find $1_$3_$4.go -size 0 -delete
