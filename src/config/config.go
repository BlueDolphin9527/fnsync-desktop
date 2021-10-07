package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cxfksword/fnsync-desktop/entity"
	"github.com/cxfksword/fnsync-desktop/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var App *AppConfig = newAppConfig()
var saveChan chan bool = make(chan bool, 10)
var saveDeviceChan chan entity.Device = make(chan entity.Device, 10)
var deleteDeviceChan chan entity.Device = make(chan entity.Device, 10)

type AppConfig struct {
	MachineId string

	ConnectOnStartup   bool
	HideOnStartup      bool
	DontToastConnected bool
	ClipboardSync      bool
	TextCastAutoCopy   bool
	Devices            map[string]entity.Device

	ListenPort       string
	randomListenPort int

	Log string
}

func newAppConfig() *AppConfig {
	defaultConf := AppConfig{
		MachineId:          uuid.NewString(),
		ConnectOnStartup:   true,
		HideOnStartup:      true,
		DontToastConnected: false,
		ClipboardSync:      true,
		TextCastAutoCopy:   true,
		Devices:            make(map[string]entity.Device),
		ListenPort:         "",
		randomListenPort:   9888, //utils.RandInt(10000, 60000)
	}
	conf := defaultConf

	viper.SetConfigName("config")     // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configPath()) // call multiple times to add many search paths
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		log.Error().Err(err).Msg("Load config file error")
		conf.Save()
	} else {
		err := viper.Unmarshal(&conf)
		if err != nil {
			log.Error().Err(err).Msg("unable to decode into struct")
			conf = defaultConf
		}
		if conf.MachineId == "" {
			conf.MachineId = uuid.NewString()
			conf.Save()
		}
	}

	log.Info().Msgf("Load config: \n%s", string(utils.ToJSON(conf)))

	go conf.startSaveRunner()
	return &conf
}

func configPath() string {
	dirname, _ := os.UserHomeDir()
	path := filepath.Join(dirname, ".fnsync")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(dirname, ".fnsync"), os.ModePerm); err != nil {
			log.Error().Err(err).Msg("")
		}
	}

	return path
}

func (a *AppConfig) Save() {
	saveChan <- true
}

func (a *AppConfig) SaveDevice(device entity.Device) {
	saveDeviceChan <- device
}

func (a *AppConfig) DeleteDevice(device entity.Device) {
	deleteDeviceChan <- device
}

func (a *AppConfig) startSaveRunner() {
	defer func() { log.Info().Msgf("Quit conifig SaveRunner.") }()
	log.Info().Msgf("Start conifig SaveRunner...")

	for {
		select {
		case <-saveChan:
			var mapConf map[string]interface{}
			jsonConf, _ := json.Marshal(*a)
			log.Debug().Msgf("Save config: %s", string(jsonConf))
			if err := json.Unmarshal(jsonConf, &mapConf); err != nil {
				log.Error().Err(err).Msg("unable to decode into struct")
			}
			for k, v := range mapConf {
				viper.Set(k, v)
			}
			if err := viper.SafeWriteConfig(); err != nil {
				err := viper.WriteConfig() // this force-saves if file is present
				if err != nil {
					log.Error().Err(err).Msg("save config error")
				}
			}
		case device := <-saveDeviceChan:
			a.Devices[device.Id] = device

			a.Save()

		}
	}
}

func (l *AppConfig) GetListenPort() int {
	port, err := strconv.Atoi(l.ListenPort)
	if err != nil && port > 0 {
		return port
	} else {
		return l.randomListenPort
	}
}
