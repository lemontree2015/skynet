package config

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// CONF解析模块

// 例子:
// $skynet.conf
// [DEFAULT]
//
// host = 127.0.0.1
// region = development
//
// client.conn.max = 5
// client.conn.idle = 2
//
// service.port.min = 9000
// service.port.max = 9999
//
// # Override values at the service level
// [TestService-1.0.0]
// service.port.min = 8000
// service.port.max = 8999
//
// 代码:
// cfg, err := config.New("/etc/skynet/skynet.conf")
// if err != nil {
//    ...处理错误...
// }
//
// if vHost, err := cfg.DefaultString("host"); err == nil {
//	  ...配置存在...
// } else {
//	  ...配置不存在...
// }
//
// if vServicePortMin, err := cfg.DefaultInt("service.port.min"); err == nil {
//	  ...配置存在...
// } else {
//	  ...配置不存在...
// }
//
//
// if vServicePortMin, err := cfg.String("TestService", "1.0.0", "service.port.min"); err == nil {
//	  ...配置存在...
// } else {
//	  ...配置不存在...
// }
//

// 每个实例对应一个*.conf文件
type Configuration struct {
	filePath string
	sections map[string]*Section
	lock     *sync.RWMutex
}

// 每个实例对应一个Section
//
// 例子:
// [TestService-1.0.0]
// service.port.min = 8000
// service.port.max = 8999
type Section struct {
	name    string
	options map[string]string
	lock    *sync.RWMutex
}

// 解析*.conf文件, 构造一个Configuration Instance
func Parse(filePath string) (*Configuration, error) {
	filePath = path.Clean(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := newConfiguration(filePath)
	var activeSection *Section

	scanner := bufio.NewScanner(bufio.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		if !(strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";")) && len(line) > 0 {
			if strings.HasPrefix(line, "[") {
				// 遇到一个新的Section
				name := strings.Trim(line, " []")
				activeSection = getSection(config, name)
				continue
			} else {
				// 处理当前的Section
				addOption(activeSection, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// Configuration

func (c *Configuration) FilePath() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.filePath
}

func (c *Configuration) Section(name string) (*Section, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if s, ok := c.sections[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Unable to find config section: %v", name)
}

func (c *Configuration) DefaultString(option string) (value string, err error) {
	return c.String("DEFAULT", "", option)
}

func (c *Configuration) DefaultInt(option string) (value int, err error) {
	return c.Int("DEFAULT", "", option)
}

func (c *Configuration) DefaultBool(option string) (value bool, err error) {
	return c.Bool("DEFAULT", "", option)
}

func (c *Configuration) String(serviceName, serviceVersion, option string) (value string, err error) {
	if serviceVersion == "" {
		return c.StringValue(serviceName, option)
	} else {
		return c.StringValue(serviceName+"-"+serviceVersion, option)
	}
}

func (c *Configuration) Int(serviceName, serviceVersion, option string) (value int, err error) {
	if serviceVersion == "" {
		return c.IntValue(serviceName, option)
	} else {
		return c.IntValue(serviceName+"-"+serviceVersion, option)
	}
}

func (c *Configuration) Bool(serviceName, serviceVersion, option string) (value bool, err error) {
	if serviceVersion == "" {
		return c.BoolValue(serviceName, option)
	} else {
		return c.BoolValue(serviceName+"-"+serviceVersion, option)
	}
}

func (c *Configuration) StringValue(section, option string) (value string, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if s, ok := c.sections[section]; ok {
		// Section存在
		return s.StringValue(option)
	} else {
		// Section不存在
		return "", fmt.Errorf("Unable to find config section: %v", section)
	}
}

func (c *Configuration) IntValue(section, option string) (value int, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if s, ok := c.sections[section]; ok {
		// Section存在
		return s.IntValue(option)
	} else {
		// Section不存在
		return 0, fmt.Errorf("Unable to find config section: %v", section)
	}
}

func (c *Configuration) BoolValue(section, option string) (value bool, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if s, ok := c.sections[section]; ok {
		// Section存在
		return s.BoolValue(option)
	} else {
		// Section不存在
		return false, fmt.Errorf("Unable to find config section: %v", section)
	}
}

// Section

func (s *Section) Exists(option string) (ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok = s.options[option]
	return
}

func (s *Section) StringValue(option string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if value, ok := s.options[option]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("Option not exist: %v", option)
	}
}

func (s *Section) IntValue(option string) (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if value, ok := s.options[option]; ok {
		// 解析成Int
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int(i), nil
		} else {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("Option not exist: %v", option)
	}
}

func (s *Section) BoolValue(option string) (bool, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if value, ok := s.options[option]; ok {
		// 解析成Bool
		if b, err := strconv.ParseBool(value); err == nil {
			return b, nil
		} else {
			return false, err
		}
	} else {
		return false, fmt.Errorf("Option not exist: %v", option)
	}
}

func (s *Section) Set(option string, value string) (oldValue string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	oldValue = s.options[option]
	s.options[option] = value

	return oldValue
}

func (s *Section) Options() map[string]string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.options
}

// Private

func newConfiguration(filePath string) *Configuration {
	return &Configuration{
		filePath: filePath,
		sections: make(map[string]*Section),
		lock:     new(sync.RWMutex),
	}
}

func newSection(name string) *Section {
	return &Section{
		name:    name,
		options: make(map[string]string),
		lock:    new(sync.RWMutex),
	}
}

// 解析一行数据, 分隔符可以是=或者:
//
// 例子:
// parseOption("key1 = val1") => opt="key1", value="val1"
// parseOption("key1 = val1=val=val") => opt="key1", value="val1=val=val"
// parseOption("key1:val1") => opt="key1", value="val1"
// parseOption(" xxxxx ") => opt="xxxxx", value=""
func parseOption(option string) (opt, value string) {

	splitFun := func(i int, delim string) (opt, value string) {
		opt = strings.Trim(option[:i], " ")
		value = strings.Trim(option[i+1:], " ")
		return
	}

	if i := strings.Index(option, "="); i != -1 {
		opt, value = splitFun(i, "=")
	} else if i := strings.Index(option, ":"); i != -1 {
		opt, value = splitFun(i, ":")
	} else {
		opt = option
	}
	return
}

func addOption(s *Section, option string) {
	var opt, value string
	if opt, value = parseOption(option); value != "" {
		s.options[opt] = value
	} else {
		s.options[opt] = ""
	}
}

func getSection(c *Configuration, name string) *Section {
	if s, ok := c.sections[name]; ok {
		return s
	}

	section := newSection(name)
	c.sections[name] = section
	return section
}
