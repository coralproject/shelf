// Package cfg provides a package level configuration loader which loaders a
// given set of configuration options using a given namespace and a map as the
// storage endpoint.
//		os.Setenv("MYAPP_PROC_ID", "322")
//		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
//		os.Setenv("MYAPP_PORT", "4034")
//		os.Setenv("MYAPP_Stamp", "2013-10-03 10:43:32.200")
// To load the set of configuration values from your environment, simple do:
//
//      cfg.Init("myapp")
//
// To Retrieve keys:
//
//  porc := cfg.String("proc_id")
//  port := cfg.Int("port")
//  ms := cfg.Time("stamp")
package cfg
