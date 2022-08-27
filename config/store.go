package config

import (
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	fileMode os.FileMode = 0666
	dirMode  os.FileMode = 0755
)

var Store = Storage{
	initLoad: true,
}

type Storage struct {
	lock        sync.Mutex
	initLoad    bool
	config      *RawConfig
	uiAbsFolder string
	uiFile      string
	uiAbsPath   string
}

func createName(name, ext string) string {
	return name + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ext
}

func splitNameExt(fName string) (string, string) {
	ext := path.Ext(fName)
	name := strings.TrimSuffix(fName, ext)
	return name, ext
}

func (s *Storage) Initial(dir, configFile string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Errorln("can't create UI storage directory %s: %s", dir, err.Error())
		}
	}
	currentDir, _ := os.Getwd()
	s.uiAbsFolder = filepath.Join(currentDir, dir)

	filename := filepath.Base(configFile)
	name, ext := splitNameExt(filename)
	s.uiFile = createName(name, ext)
	s.uiAbsPath = filepath.Join(s.uiAbsFolder, s.uiFile)
}

func (s *Storage) GetStorePath() string {
	return s.uiAbsPath
}

func (s *Storage) IsFirstLoad() bool {
	return s.initLoad
}

func (s *Storage) SetFirstLoad(load bool) {
	s.initLoad = load
}

func (s *Storage) GetConfig() *RawConfig {
	//buffer, err := yaml.Marshal(s.config)
	//if err != nil {
	//	return fmt.Sprintf("fail to yaml marshal config: %s", err.Error())
	//}
	//str, err := json.Marshal(string(buffer))
	//if err != nil {
	//	return fmt.Sprintf("fail to json marshal config: %s", err.Error())
	//}
	//return str
	return s.config
}

func (s *Storage) SetConfig(c *RawConfig) {
	s.config = c
}

func (s *Storage) WriteStore() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	buf, err := yaml.Marshal(s.config)
	if err != nil {
		log.Errorln("error on yaml marshal raw config %s", err.Error())
		return err
	}
	return ioutil.WriteFile(s.uiAbsPath, buf, fileMode)
}
