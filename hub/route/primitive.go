package route

import (
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func primitiveRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getPrimitive)
	r.Put("/", putPrimitive)
	return r
}

func jsonConverter(input interface{}) interface{} {
	switch it := input.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range it {
			m2[k.(string)] = jsonConverter(v)
		}
		return m2
	case []interface{}:
		for mk, mv := range it {
			it[mk] = jsonConverter(mv)
		}
	}
	return input
}

func getPrimitive(w http.ResponseWriter, r *http.Request) {
	if config.Store.GetConfig().Profile.UiStorage != "" {
		//render.JSON(w, r, config.Store.GetConfig())
		render.JSON(w, r, config.Store.GetConfig().BuildJson())
	} else {
		render.JSON(w, r, render.M{
			"profile": render.M{
				"ui-storage": "",
			},
		})
	}

}

func putPrimitive(w http.ResponseWriter, r *http.Request) {
	if config.Store.GetConfig().Profile.UiStorage != "" {
		req := updateConfigRequest{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, ErrBadRequest)
			return
		}

		var rawCfg *config.RawConfig
		var cfg *config.Config
		var err error

		rawCfg, cfg, err = ParsePayLoad([]byte(req.Payload))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError(err.Error()))
			return
		}

		config.Store.SetConfig(rawCfg)
		err = config.Store.WriteStore()
		if err != nil {
			log.Warnln("fail to write UI Store file %s with %s", config.Store.GetStorePath(), err.Error())
		}
		executor.ApplyConfig(cfg, true)
		log.Infoln("Update configuration and write to %s successfully", config.Store.GetStorePath())
	} else {
		render.JSON(w, r, render.M{
			"uiEnable": false,
		})
	}

}

func ParsePayLoad(buf []byte) (*config.RawConfig, *config.Config, error) {
	rawCfg, err := config.UnmarshalRawConfig(buf)
	if err != nil {
		return nil, nil, err
	}

	cfg, err := config.ParseRawConfig(rawCfg)
	if err != nil {
		return nil, nil, err
	}

	return rawCfg, cfg, nil
}
