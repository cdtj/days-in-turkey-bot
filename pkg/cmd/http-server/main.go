package main

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/cmd"
	"cdtj.io/days-in-turkey-bot/db"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"

	cep "cdtj.io/days-in-turkey-bot/entity/country/endpoint/http/echo"
	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	uep "cdtj.io/days-in-turkey-bot/entity/user/endpoint/http/echo"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"
)

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("CUSTOM", "📝", "", 60, 90, 180, true)
)

func main() {
	i18n, err := i18n.NewI18n("i18n", defaultLang)
	if err != nil {
		panic(err)
	}
	tgFmt := formatter.NewTelegramFormatter(i18n, false)
	countryDB := db.NewMapDB()
	countryRepo := cr.NewCountryRepo(countryDB)
	countrySvc := cs.NewCountryService(tgFmt, defaultCountry)
	countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)

	userDB := db.NewMapDB()
	userRepo := ur.NewUserRepo(userDB)
	userSvc := us.NewUserService(tgFmt, i18n, countrySvc)
	userUC := uuc.NewUserUsecase(userRepo, userSvc, countryUC)

	echoRouter := httpserver.NewEchoRouter()
	uep.RegisterHTTPEndpointsEcho(echoRouter, userUC)
	cep.RegisterHTTPEndpointsEcho(echoRouter, countryUC)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: echoRouter,
	})
	cmd.Serve(srv)
}
