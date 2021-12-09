module github.com/philippgille/gokv

go 1.17

replace github.com/philippgille/gokv => ./
replace github.com/philippgille/gokv/encoding => ./encoding
replace github.com/philippgille/gokv/bbolt => ./bbolt

require (
	github.com/mattn/go-colorable v0.1.12
	github.com/mitchellh/cli v1.1.2
	github.com/philippgille/gokv/bbolt v0.6.0
	github.com/posener/complete v1.1.1
)

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/philippgille/gokv/encoding v0.0.0-20191011213304-eb77f15b9c61 // indirect
	github.com/philippgille/gokv/util v0.0.0-20191011213304-eb77f15b9c61 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)
