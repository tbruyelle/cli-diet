module github.com/tendermint/starport

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Microsoft/hcsshim v0.8.17 // indirect
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/charmbracelet/glow v1.4.0
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/dariubs/percent v0.0.0-20200128140941-b7801cf1c7e2
	github.com/docker/docker v20.10.7+incompatible
	github.com/emicklei/proto v1.9.0
	github.com/fatih/color v1.10.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/logger v1.0.3
	github.com/gobuffalo/packd v1.0.0
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.8.0
	github.com/gookit/color v1.2.7
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/rpc v1.2.0
	github.com/iancoleman/strcase v0.1.3
	github.com/imdario/mergo v0.3.11
	github.com/jpillora/chisel v1.7.3
	github.com/kr/pretty v0.2.1
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-zglob v0.0.3
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/otiai10/copy v1.4.2
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/spm v0.0.0-20210524110815-6d7452d2dc4a
	github.com/tendermint/spn v0.0.0-20210406123257-decaff8dcaf9
	github.com/tendermint/tendermint v0.34.9
	github.com/tendermint/vue v0.1.49
	go.opencensus.io v0.22.6 // indirect
	golang.org/x/mod v0.4.2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
