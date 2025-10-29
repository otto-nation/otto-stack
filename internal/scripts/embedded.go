package scripts

import _ "embed"

//go:embed localstack-init.sh
var LocalstackInitScript string

//go:embed kafka-topics-init.sh
var KafkaTopicsInitScript string
