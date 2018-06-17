package main

import "flag"

type flags struct {
	redisAddr string
	redisPwd  string
	redisDB   int
}

func getFlags() (flags, error) {

	var cmdFlags flags

	const (
		defaultAddr = "127.0.0.1:6379"
		addrUsage   = "The address of the redis instance"
		defaultPwd  = ""
		pwdUsage    = "The password of the redis instance"
		defaultDB   = 0
		dbUsage     = "The database to be selected after connecting to redis instance"
	)

	flag.StringVar(&(cmdFlags.redisAddr), "redisAddr", defaultAddr, addrUsage)
	flag.StringVar(&(cmdFlags.redisAddr), "a", defaultAddr, addrUsage+"(shorthand)")
	flag.StringVar(&(cmdFlags.redisPwd), "redisPwd", defaultPwd, pwdUsage)
	flag.StringVar(&(cmdFlags.redisPwd), "p", defaultPwd, pwdUsage+"(shorthand)")
	flag.IntVar(&(cmdFlags.redisDB), "redisDB", defaultDB, dbUsage)
	flag.IntVar(&(cmdFlags.redisDB), "d", defaultDB, dbUsage+"(shorthand)")

	flag.Parse()

	return cmdFlags, nil
}
