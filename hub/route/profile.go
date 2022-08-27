package route

import (
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func profileRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProfile)
	r.Put("/", putProfile)
	return r
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	if config.Store.IsFirstLoad() {
		render.JSON(w, r, render.M{
			"enable": false,
		})
	} else {
		//render.JSON(w, r, render.M{
		//	"payload": config.Store.GetConfig(),
		//})
		render.JSON(w, r, config.Store.GetConfig())
	}

}

func putProfile(w http.ResponseWriter, r *http.Request) {
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
	// need to write
	config.Store.SetConfig(rawCfg)
	err = config.Store.WriteStore()
	if err != nil {
		log.Warnln("fail to write UI Store file %s with %s", config.Store.GetStorePath(), err.Error())
	}
	executor.ApplyConfig(cfg, true)
	log.Infoln("apply updated configuration and write to %s successfully", config.Store.GetStorePath())
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
