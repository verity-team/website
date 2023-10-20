package server

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/verity-team/dws/api"
	"github.com/verity-team/dws/internal/delphi/db"
)

type DelphiServer struct {
	db *sqlx.DB
}

func NewDelphiServer(db *sqlx.DB) *DelphiServer {
	return &DelphiServer{
		db: db,
	}
}

func (s *DelphiServer) ConnectWallet(ctx echo.Context) error {
	var cr api.ConnectionRequest
	err := ctx.Bind(&cr)
	if err != nil {
		log.Error("failed to bind POST param (AffiliateRequest)")
		return err
	}
	return db.ConnectWalletToAffiliate(s.db, cr)
}

func (s *DelphiServer) DonationData(ctx echo.Context) error {
	dd, err := db.GetDonationData(s.db)
	if err != nil {
		return err
	}
	ra, present := os.LookupEnv("DWS_DONATION_ADDRESS")
	if !present {
		err = errors.New("DWS_DONATION_ADDRESS environment variable not set")
		log.Error(err)
		return err
	}
	dd.ReceivingAddress = ra
	// if we failed to fetch an ETH price, the status should be set to "paused"
	if dd.Prices[0].Price == "0.00" {
		dd.Status = api.Paused
	}
	return ctx.JSON(http.StatusOK, *dd)
}

func (s *DelphiServer) UserData(ctx echo.Context, address string) error {
	address = strings.ToLower(address)
	dd, err := db.GetUserDonationData(s.db, address)
	if err != nil {
		return err
	}
	if dd == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	us, err := db.GetUserData(s.db, address)
	if err != nil {
		return err
	}
	var result api.UserDataResult = api.UserDataResult{
		Donations: dd,
		UserData:  *us,
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *DelphiServer) Alive(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "{}\n")
}

func (s *DelphiServer) Ready(ctx echo.Context) error {
	err := s.db.Ping()
	if err != nil {
		return ctx.String(http.StatusServiceUnavailable, "{}\n")
	}
	return ctx.String(http.StatusOK, "{}\n")

}

func (s *DelphiServer) GenerateCode(ctx echo.Context, params api.GenerateCodeParams) error {
	return nil
}
