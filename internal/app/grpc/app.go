package grpcapp

import (
	authrpc "GoAuthGRPC/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log *slog.Logger

	gRPCServer  *grpc.Server
	port        int
	authService authrpc.Auth
}

func New(
	log *slog.Logger,
	authService authrpc.Auth,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	authrpc.Register(gRPCServer, authService)

	return &App{
		log:         log,
		gRPCServer:  gRPCServer,
		port:        port,
		authService: authService,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpc.Stop"

	a.log.With(slog.String("op", op)).Info("stopping grpc server", slog.Int("addr", a.port))

	a.gRPCServer.GracefulStop()
}
