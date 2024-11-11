module ciba-example

go 1.23.1

require (
	github.com/daedaluz/goauth2 v0.0.0-20210819154807-3b3b3b3b3b3b
)

replace (
	github.com/daedaluz/goauth2 latest => ../../ latest
)