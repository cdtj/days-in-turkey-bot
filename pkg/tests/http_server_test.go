package tests

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"testing"

	"context"

	"cdtj.io/days-in-turkey-bot/cmd"
	"cdtj.io/days-in-turkey-bot/db"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"github.com/go-chi/chi/v5/middleware"

	cep "cdtj.io/days-in-turkey-bot/entity/country/endpoint/http/chi"
	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"

	uep "cdtj.io/days-in-turkey-bot/entity/user/endpoint/http/chi"
	ur "cdtj.io/days-in-turkey-bot/entity/user/repo"
	us "cdtj.io/days-in-turkey-bot/entity/user/service"
	uuc "cdtj.io/days-in-turkey-bot/entity/user/usecase"

	_ "net/http/pprof"
)

var (
	hostAddr = "http://localhost:8080"
)

const RunMultiplier = 100

func BenchmarkHttpServer(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	incr := atomic.Int64{}

	srv := srvr()
	go func() {
		cmd.Serve(srv)
	}()

	cli := cli()
	for _, err := cli.Get(hostAddr + "/user/info/" + strconv.FormatInt(incr.Add(1), 10)); err != nil; {
		b.Error(err)
		return
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cli.Get(hostAddr + "/user/info/" + strconv.FormatInt(incr.Add(1), 10))
		}
	})
	b.StopTimer()
	srv.Shutdown(context.Background())
}

func BenchmarkHttpServerWaitProfile(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	incr := atomic.Int64{}

	srv := srvr()
	go cmd.Serve(srv)

	cli := cli()
	for _, err := cli.Get(hostAddr + "/user/info/" + strconv.FormatInt(incr.Add(1), 10)); err != nil; {
		b.Error(err)
		return
	}

	cpuFile, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	memFile, _ := os.Create("mem.pprof")
	pprof.WriteHeapProfile(memFile)
	defer memFile.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < RunMultiplier; i++ {
				body := bytes.NewBufferString(`{ "dates" : "11/06/2023 11/07/2023 11/08/2023 11/09/2023 11/10/2023 11/11/2023 11/12/2023 11/01/2024 11/02/2024" }`)
				resp, err := cli.Post(hostAddr+"/user/calc/"+strconv.FormatInt(incr.Add(1), 10), "application/json", body)
				if err != nil {
					b.Error(err)
					return
				}
				_, err = io.ReadAll(resp.Body)
				if err != nil {
					b.Error(err)
					return
				}
			}
		}
	})
	b.StopTimer()

	srv.Shutdown(context.Background())
}

func cli() *http.Client {
	return &http.Client{}
}

var (
	defaultLang    = "en"
	defaultCountry = model.NewCountry("RU", "RU", 60, 90, 180)
)

func srvr() *httpserver.HttpServer {
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

	router := httpserver.NewChiRouter()
	router.Mount("/debug", middleware.Profiler())

	uep.RegisterHTTPEndpointsChi(router, userUC)
	cep.RegisterHTTPEndpointsChi(router, countryUC)

	srv := httpserver.NewHttpServer(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
	return srv
}
