package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	logger "log"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/snail007/proxy/core/lib/kcpcfg"
	encryptconn "github.com/snail007/proxy/core/lib/transport/encrypt"
	services "github.com/snail007/proxy/services"
	mux "github.com/snail007/proxy/services/mux"
	kcp "github.com/xtaci/kcp-go"

	"golang.org/x/crypto/pbkdf2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	log     = logger.New(os.Stderr, "", logger.Ldate|logger.Ltime)
	app     *kingpin.Application
	service *services.ServiceItem
	cmd     *exec.Cmd
	cpuProfilingFile, memProfilingFile, blockProfilingFile,
	goroutineProfilingFile, threadcreateProfilingFile *os.File
	isDebug bool
)

const APP_VERSION = "8.2"

func initConfig() (err error) {
	//define  args
	muxClientArgs := mux.MuxClientArgs{}
	kcpArgs := kcpcfg.KCPConfigArgs{}
	//build srvice args
	app = kingpin.New("gpclient", "happy with proxy nat client")
	app.Author("snail").Version(APP_VERSION)
	isDebug := app.Flag("debug", "debug log output").Default("false").Bool()
	daemon := app.Flag("daemon", "run proxy in background").Default("false").Bool()
	forever := app.Flag("forever", "run proxy in forever,fail and retry").Default("false").Bool()
	logfile := app.Flag("log", "log file path").Default("").String()
	nolog := app.Flag("nolog", "turn off logging").Default("false").Bool()
	kcpArgs.Key = app.Flag("kcp-key", "pre-shared secret between client and server").Default("secrect").String()
	kcpArgs.Crypt = app.Flag("kcp-method", "encrypt/decrypt method, can be: aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none").Default("aes").Enum("aes", "aes-128", "aes-192", "salsa20", "blowfish", "twofish", "cast5", "3des", "tea", "xtea", "xor", "sm4", "none")
	kcpArgs.Mode = app.Flag("kcp-mode", "profiles: fast3, fast2, fast, normal, manual").Default("fast").Enum("fast3", "fast2", "fast", "normal", "manual")
	kcpArgs.MTU = app.Flag("kcp-mtu", "set maximum transmission unit for UDP packets").Default("450").Int()
	kcpArgs.SndWnd = app.Flag("kcp-sndwnd", "set send window size(num of packets)").Default("1024").Int()
	kcpArgs.RcvWnd = app.Flag("kcp-rcvwnd", "set receive window size(num of packets)").Default("1024").Int()
	kcpArgs.DataShard = app.Flag("kcp-ds", "set reed-solomon erasure coding - datashard").Default("10").Int()
	kcpArgs.ParityShard = app.Flag("kcp-ps", "set reed-solomon erasure coding - parityshard").Default("3").Int()
	kcpArgs.DSCP = app.Flag("kcp-dscp", "set DSCP(6bit)").Default("0").Int()
	kcpArgs.NoComp = app.Flag("kcp-nocomp", "disable compression").Default("false").Bool()
	kcpArgs.AckNodelay = app.Flag("kcp-acknodelay", "be carefull! flush ack immediately when a packet is received").Default("true").Bool()
	kcpArgs.NoDelay = app.Flag("kcp-nodelay", "be carefull!").Default("0").Int()
	kcpArgs.Interval = app.Flag("kcp-interval", "be carefull!").Default("50").Int()
	kcpArgs.Resend = app.Flag("kcp-resend", "be carefull!").Default("0").Int()
	kcpArgs.NoCongestion = app.Flag("kcp-nc", "be carefull! no congestion").Default("0").Int()
	kcpArgs.SockBuf = app.Flag("kcp-sockbuf", "be carefull!").Default("4194304").Int()
	kcpArgs.KeepAlive = app.Flag("kcp-keepalive", "be carefull!").Default("10").Int()

	//########mux-client#########
	muxClientArgs.Parent = app.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	muxClientArgs.ParentType = app.Flag("parent-type", "parent protocol type <tls|tcp|tcps|kcp|tou|ws|wss>").Default("tls").Short('T').Enum("tls", "tcp", "tcps", "kcp", "tou", "ws", "wss")
	muxClientArgs.CertFile = app.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	muxClientArgs.KeyFile = app.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	muxClientArgs.Timeout = app.Flag("timeout", "tcp timeout with milliseconds").Short('i').Default("2000").Int()
	muxClientArgs.Key = app.Flag("k", "key same with server").Default("default").String()
	muxClientArgs.IsCompress = app.Flag("c", "compress data when tcp|tls mode").Default("false").Bool()
	muxClientArgs.SessionCount = app.Flag("session-count", "session count which connect to bridge").Short('n').Default("10").Int()
	muxClientArgs.Jumper = app.Flag("jumper", "http(s) or socks5 or ss proxies used when connecting to parent, only worked of -T is tls or tcp, format is  http://username:password@host:port http://host:port https://username:password@host:port https://host:port or socks5://username:password@host:port socks5://host:port socks5s://username:password@host:port socks5s://host:port").Short('J').Default("").String()
	muxClientArgs.TCPSMethod = app.Flag("tcps-method", "method of parent tcps's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxClientArgs.TCPSPassword = app.Flag("tcps-password", "password of parent tcps's encrpyt/decrypt").Default("snail007's_goproxy").String()
	muxClientArgs.TOUMethod = app.Flag("tou-method", "method of parent tou's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxClientArgs.TOUPassword = app.Flag("tou-password", "password of parent tou's encrpyt/decrypt").Default("snail007's_goproxy").String()
	muxClientArgs.WSMethod = app.Flag("ws-method", "method of local ws's encrpyt/decrypt, these below are supported :\n"+strings.Join(encryptconn.GetCipherMethods(), ",")).Default("aes-192-cfb").String()
	muxClientArgs.WSPassword = app.Flag("ws-password", "password of parent ws's encrpyt/decrypt").Default("snail007/goproxy").String()
	muxClientArgs.IsP2P = app.Flag("p2p", "using p2p when server connect to client").Default("false").Bool()
	muxClientArgs.P2PPort = app.Flag("p2p-port", "bridge udp port of p2p punnching, if leave empty as same as tcp port of bridge").Default("").String()
	muxClientArgs.KCP = kcpArgs

	//parse args
	kingpin.MustParse(app.Parse(os.Args[1:]))

	//set kcp config

	switch *kcpArgs.Mode {
	case "normal":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 0, 40, 2, 1
	case "fast":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 0, 30, 2, 1
	case "fast2":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 1, 20, 2, 1
	case "fast3":
		*kcpArgs.NoDelay, *kcpArgs.Interval, *kcpArgs.Resend, *kcpArgs.NoCongestion = 1, 10, 2, 1
	}
	pass := pbkdf2.Key([]byte(*kcpArgs.Key), []byte("snail007-goproxy"), 4096, 32, sha1.New)

	switch *kcpArgs.Crypt {
	case "sm4":
		kcpArgs.Block, _ = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		kcpArgs.Block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		kcpArgs.Block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		kcpArgs.Block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		kcpArgs.Block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		kcpArgs.Block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		kcpArgs.Block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		kcpArgs.Block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		kcpArgs.Block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		kcpArgs.Block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		*kcpArgs.Crypt = "aes"
		kcpArgs.Block, _ = kcp.NewAESBlockCrypt(pass)
	}
	//attach kcp config
	muxClientArgs.KCP = kcpArgs

	flags := logger.Ldate
	if *isDebug {
		flags |= logger.Lshortfile | logger.Lmicroseconds
		cpuProfilingFile, _ = os.Create("cpu.prof")
		memProfilingFile, _ = os.Create("memory.prof")
		blockProfilingFile, _ = os.Create("block.prof")
		goroutineProfilingFile, _ = os.Create("goroutine.prof")
		threadcreateProfilingFile, _ = os.Create("threadcreate.prof")
		pprof.StartCPUProfile(cpuProfilingFile)
	} else {
		flags |= logger.Ltime
	}
	log.SetFlags(flags)
	if *nolog {
		log.SetOutput(ioutil.Discard)
	} else if *logfile != "" {
		f, e := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			log.Fatal(e)
		}
		log.SetOutput(f)
	}
	if *daemon {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--daemon" {
				args = append(args, arg)
			}
		}
		cmd = exec.Command(os.Args[0], args...)
		cmd.Start()
		f := ""
		if *forever {
			f = "forever "
		}
		log.Printf("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if *forever {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--forever" {
				args = append(args, arg)
			}
		}
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
				}
			}()
			for {
				if cmd != nil {
					cmd.Process.Kill()
					time.Sleep(time.Second * 5)
				}
				cmd = exec.Command(os.Args[0], args...)
				cmdReaderStderr, err := cmd.StderrPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				scanner := bufio.NewScanner(cmdReader)
				scannerStdErr := bufio.NewScanner(cmdReaderStderr)
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
						}
					}()
					for scanner.Scan() {
						fmt.Println(scanner.Text())
					}
				}()
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("crashed, err: %s\nstack:%s", e, string(debug.Stack()))
						}
					}()
					for scannerStdErr.Scan() {
						fmt.Println(scannerStdErr.Text())
					}
				}()
				if err := cmd.Start(); err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				pid := cmd.Process.Pid
				log.Printf("worker %s [PID] %d running...\n", os.Args[0], pid)
				if err := cmd.Wait(); err != nil {
					log.Printf("ERR:%s,restarting...", err)
					continue
				}
				log.Printf("worker %s [PID] %d unexpected exited, restarting...\n", os.Args[0], pid)
			}
		}()
		return
	}
	if *logfile == "" {
		poster()
		if *isDebug {
			log.Println("[profiling] cpu profiling save to file : cpu.prof")
			log.Println("[profiling] memory profiling save to file : memory.prof")
			log.Println("[profiling] block profiling save to file : block.prof")
			log.Println("[profiling] goroutine profiling save to file : goroutine.prof")
			log.Println("[profiling] threadcreate profiling save to file : threadcreate.prof")
		}
	}

	//regist services and run service
	serviceName := "client"
	services.Regist(serviceName, "client", mux.NewMuxClient(serviceName), muxClientArgs, log)
	service, err = services.Run(serviceName, nil)
	if err != nil {
		log.Fatalf("run service [%s] fail, ERR:%s", serviceName, err)
	}
	return
}

func poster() {
	fmt.Printf(`Proxy Client v%s`+" by snail , blog : http://www.host900.com/\n\n", APP_VERSION)
}
func saveProfiling() {
	goroutine := pprof.Lookup("goroutine")
	goroutine.WriteTo(goroutineProfilingFile, 1)
	heap := pprof.Lookup("heap")
	heap.WriteTo(memProfilingFile, 1)
	block := pprof.Lookup("block")
	block.WriteTo(blockProfilingFile, 1)
	threadcreate := pprof.Lookup("threadcreate")
	threadcreate.WriteTo(threadcreateProfilingFile, 1)
	pprof.StopCPUProfile()
}

func main() {
	err := initConfig()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	if service != nil && service.S != nil {
		Clean(&service.S)
	} else {
		Clean(nil)
	}
}
func Clean(s *services.Service) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("crashed, err: %s\nstack:", e, string(debug.Stack()))
			}
		}()
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping services...")
			if s != nil && *s != nil {
				(*s).Clean()
			}
			if cmd != nil {
				log.Printf("clean process %d", cmd.Process.Pid)
				cmd.Process.Kill()
			}
			if isDebug {
				saveProfiling()
			}
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
