//import "github.com/solomonooo/mercury"
//author : solomonooo
//create time : 2016-09-08

package mercury

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	CONFIG_FILE         = "mercury.conf"
	DEFAULT_LISTEN_IP   = "0.0.0.0"
	DEFAULT_LISTEN_PORT = 7531

	DEFAULT_RECV_BUFF_SIZE = 32 * 1024
	DEFAULT_RECV_TIMEOUT   = 1000
	DEFAULT_SEND_TIMEOUT   = 1000

	DEFAULT_STDERR_2_FILE = false
)

type MercuryConfig struct {
	Ip       string
	Port     uint32
	LogDir   string
	LogLevel string

	RecvBuffSize uint32
	RecvTimeout  uint32
	SendTimeout  uint32

	StdErr2File bool

	StatusCycle uint32

	reserved map[string]string
}

func NewConfig() *MercuryConfig {
	return &MercuryConfig{
		Ip:           DEFAULT_LISTEN_IP,
		Port:         DEFAULT_LISTEN_PORT,
		LogDir:       DEFAULT_LOG_DIR,
		LogLevel:     DEFAULT_LOG_LEVEL,
		RecvBuffSize: DEFAULT_RECV_BUFF_SIZE,
		RecvTimeout:  DEFAULT_RECV_TIMEOUT,
		SendTimeout:  DEFAULT_SEND_TIMEOUT,
		StatusCycle:  DEFAULT_STATUS_CYCLE,
		StdErr2File:  DEFAULT_STDERR_2_FILE,
	}
}

func (config *MercuryConfig) Init(configPath string) error {
	err := config.ParseConfig(configPath)
	if err != nil {
		return err
	}

	//load server config
	if config.Find("server", "ip") {
		config.Ip, _ = config.Get("server", "ip")
	}
	if config.Find("server", "port") {
		config.Port, _ = config.GetUInt32("server", "port")
	}
	if config.Find("server", "log_dir") {
		config.LogDir, _ = config.Get("server", "log_dir")
	}
	if config.Find("server", "log_level") {
		config.LogLevel, _ = config.Get("server", "log_level")
	}
	if config.Find("server", "buffer_size") {
		config.RecvBuffSize, _ = config.GetUInt32("server", "buffer_size")
	}
	if config.Find("server", "recv_timeout") {
		config.RecvTimeout, _ = config.GetUInt32("server", "recv_timeout")
	}
	if config.Find("server", "send_timeout") {
		config.SendTimeout, _ = config.GetUInt32("server", "send_timeout")
	}
	if config.Find("server", "status_cycle") {
		config.StatusCycle, _ = config.GetUInt32("server", "status_cycle")
	}
	if config.Find("server", "stderr2file") {
		config.StdErr2File, _ = config.GetBool("server", "stderr2file")
	}

	return nil
}

func (config *MercuryConfig) ParseConfig(configPath string) error {
	handler, err := os.Open(configPath)
	if err != nil {
		Fatal("open mercury.conf failed, conf=%s", configPath)
		return err
	}
	defer handler.Close()

	config.reserved = make(map[string]string)
	buf := bufio.NewReader(handler)
	section := ""
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			Fatal("read conf from mercury.conf failed, err=%s", err.Error())
			return err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
		} else {
			idx := strings.Index(line, "=")
			if idx == -1 {
				config.reserved[section+"."+line] = ""
			} else {
				config.reserved[section+"."+strings.TrimSpace(line[0:idx])] = strings.TrimSpace(line[idx+1:])
			}
		}
	}

	return nil
}

func (config MercuryConfig) Find(section string, name string) bool {
	if config.reserved == nil {
		return false
	}

	realname := section + "." + name
	_, ok := config.reserved[realname]
	return ok
}

func (config MercuryConfig) Get(section string, name string) (string, error) {
	if config.reserved == nil {
		return "", fmt.Errorf("%s.%s not exist", section, name)
	}

	realname := section + "." + name
	value, ok := config.reserved[realname]
	if false == ok {
		return "", fmt.Errorf("%s.%s not exist", section, name)
	}
	return value, nil
}

func (config MercuryConfig) GetInt64(section string, name string) (int64, error) {
	str, err := config.Get(section, name)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(str, 10, 64)
}

func (config MercuryConfig) GetUInt64(section string, name string) (uint64, error) {
	str, err := config.Get(section, name)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(str, 10, 64)
}

func (config MercuryConfig) GetUInt32(section string, name string) (uint32, error) {
	v, err := config.GetUInt64(section, name)
	return uint32(v), err
}

func (config MercuryConfig) GetInt32(section string, name string) (int32, error) {
	v, err := config.GetInt64(section, name)
	return int32(v), err
}

func (config MercuryConfig) GetInt(section string, name string) (int, error) {
	v, err := config.GetInt64(section, name)
	return int(v), err
}

func (config MercuryConfig) GetBool(section string, name string) (bool, error) {
	str, err := config.Get(section, name)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(str)
}

//////////////////////
func SetConfig(configPath string) error {
	err := mercury.config.Init(configPath)
	if err != nil {
		panic(err)
	}
	return err
}

func GetConfig(section string, name string) (string, error) {
	return mercury.config.Get(section, name)
}

func GetConfigInt32(section string, name string) (int32, error) {
	return mercury.config.GetInt32(section, name)
}

func GetConfigUInt32(section string, name string) (uint32, error) {
	return mercury.config.GetUInt32(section, name)
}

func GetConfigInt(section string, name string) (int, error) {
	return mercury.config.GetInt(section, name)
}

func GetConfigBool(section string, name string) (bool, error) {
	return mercury.config.GetBool(section, name)
}
