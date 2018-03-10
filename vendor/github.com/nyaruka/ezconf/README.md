# EZConf [![Build Status](https://travis-ci.org/nyaruka/ezconf.svg?branch=master)](https://travis-ci.org/nyaruka/ezconf) [![codecov](https://codecov.io/gh/nyaruka/ezconf/branch/master/graph/badge.svg)](https://codecov.io/gh/nyaruka/ezconf) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/ezconf)](https://goreportcard.com/report/github.com/nyaruka/ezconf)

EZConf provides a simple way of reading configuration settings from four sources, in order of priority (each level is higher priority than the previous ones):
 
 1. The default settings for your app
 2. A TOML file with settings
 3. Environment variables mapping to your top level settings
 4. Command line parameters mapping to your top level settings

To use it, you only need to create a struct representing the desired configuration and create an instance
with the defaults for your app. You can then pass that struct to EZConf and read the settings from all 
the sources above.

EZConf will automatically parse command line parameters and environment variables for all top level fields
in your struct of the following types:

 * int, int8, int16, int32, int64
 * uint, uint8, uint16, uint32, uint64
 * float32, float64
 * bool
 * string

It converts all CamelCase fields to snake_case in a manner that is compatible with the acronyms we work with
everyday. Some examples of how a struct name is converted to a TOML field, environment variable and command
line parameter can be found below. 

Environment variables are prefixed with your app name, in this case `courier`:

| Struct Field  | TOML Field       | Environment Variable         | Command line Parameter |
|---------------|------------------|------------------------------|------------------------|
| AWSRegion     | aws_region       | COURIER_AWS_REGION           | aws-region             |
| EC2InstanceID | ec2_instance_id  | COURIER_EC2_INSTANCE_ID      | ec2-instance-id        |
| DB            | db               | COURIER_DB                   | db                     |
| NumWorkers    | num_workers      | COURIER_NUM_WORKERS          | num-workers            |

EZConf will also automatically create the appropriate flags and help based on your struct definition, for example:

```
% courier -help
Courier - a fast message broker for IP and SMS messages

Usage of courier:
  -aws-region string
    	the aws region that S3 buckets are in (default "us-west-2")
  -db string
    	the url describing how to connect to the database (default "postgres://user@secret:rds-internal.foo.bar/courier")
  -debug-conf
    	print where config values are coming from
  -ec2-instance-id string
    	the id of the ec2 instance we are running on (default "i-12355111134a")
  -help
    	print usage information
  -num-workers int
    	the number of workers to start (default 32)

Environment variables:
          COURIER_AWS_REGION - string
                  COURIER_DB - string
     COURIER_EC2_INSTANCE_ID - string
         COURIER_NUM_WORKERS - int
```

## Example



```golang
package main

import (
	"fmt"

	"github.com/nyaruka/ezconf"
)

// Config is our apps configuration struct, we can use the `help` tag to add usage information
type Config struct {
	AWSRegion     string `help:"the aws region that S3 buckets are in"`
	DB            string `help:"the url describing how to connect to the database"`
	EC2InstanceID string `help:"the id of the ec2 instance we are running on"`
	NumWorkers    int    `help:"the number of workers to start"`
}

func main() {
	// instantiate our config with our defaults
	config := &Config{
		AWSRegion:     "us-west-2",
		DB:            "postgres://user@secret:rds-internal.foo.bar/courier",
		EC2InstanceID: "i-12355111134a",
		NumWorkers:    32,
	}

	// create our loader object, configured with configuration struct (must be a pointer), our name
	// and description, as well as any files we want to search for
	loader := ezconf.NewLoader(
		config,
		"courier", "Courier - a fast message broker for IP and SMS messages",
		[]string{"courier.toml"},
	)

	// load our configuration, exiting if we encounter any errors
	loader.MustLoad()

	// our settings have now been loaded into our config struct
	fmt.Printf("Final Settings:\n%+v\n", *config)

	// if we wish we can also validate our config using our favorite validation library
}
```
