package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"fmt"
)

type R map[string]interface{}

const (
	ResSuccess                  = iota
	ResErrUnknown
	ResErrIdLenMismatch
	ResErrAccountAlreadyExists
	ResErrIncorrectPassword
	ResErrInvalidToken
	ResErrNoToken
	ResErrNoGold
)

var tokens = TokenManager {[]Token{}}

func RegisterHandler(e echo.Context) error {
	id := toId(e.FormValue("grade"), e.FormValue("class"), e.FormValue("id"))
	if len(id) != 5 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrIdLenMismatch,
		})
	}
	password := e.FormValue("password")

	err := mongo.AddAccount(id, password, 0)
	if err == ErrAccountExist {
		return e.JSON(http.StatusTeapot, R {
			"res": ResErrAccountAlreadyExists,
		})
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
	})
}

func VerifyAccount(e echo.Context) error {
	id := toId(e.FormValue("grade"), e.FormValue("class"), e.FormValue("id"))
	if len(id) != 5 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrIdLenMismatch, // 유저이름의 문자열 길이가 틀리다
		})
	}
	password := e.FormValue("password")

	err := mongo.VerifyAccount(id, password)
	if err != nil {
		if err == ErrNoAccount {
			err := mongo.AddAccount(id, password, 0)
			if err == ErrAccountExist {
				return e.JSON(http.StatusTeapot, R {
					"res": ResErrAccountAlreadyExists, // ???????!?????!?!?!?!?
				})
			}
		}else if err == ErrIncorrectPassword {
			return e.JSON(http.StatusUnauthorized, R {
				"res": ResErrIncorrectPassword,
			})
		}else{
			return e.JSON(http.StatusInternalServerError, R {
				"res": ResErrUnknown,
			})
		}
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess, // 성공
		"token": tokens.New(id).Token,
	})
}

func GetAccount(e echo.Context) error {
	token := e.FormValue("token")
	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	account, err := mongo.GetAccount(t.User)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	return e.JSON(http.StatusOK, R {
		"id": account.Id,
		"gold": account.Gold,
	})
}

func RenewToken(e echo.Context) error {
	token := e.FormValue("token")

	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	t.Renew()
	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
	})
}

func SubtractGold(e echo.Context) error {
	token := e.FormValue("token")
	if len(token) != 32 {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrInvalidToken,
		})
	}

	g := e.FormValue("gold")
	gold, err := strconv.Atoi(g)
	if err != nil {
		return e.JSON(http.StatusBadRequest, R {
			"res": ResErrUnknown,
		})
	}

	t, err := tokens.Get(token)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, R {
			"res": ResErrNoToken,
		})
	}

	account, err := mongo.GetAccount(t.User)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	if gold < 0 {
		gold = -gold
	}

	if account.Gold < gold {
		return e.JSON(http.StatusTeapot, R {
			"res": ResErrNoGold,
		})
	}

	err = mongo.SubtractGold(account.Id, gold)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, R {
			"res": ResErrUnknown,
		})
	}

	return e.JSON(http.StatusOK, R {
		"res": ResSuccess,
	})
}

func toId(grade, class, no string) string {
	g, err := strconv.Atoi(grade)
	if err != nil {
		return ""
	}
	c, err := strconv.Atoi(class)
	if err != nil {
		return ""
	}
	n, err := strconv.Atoi(no)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%01d%02d%02d", g, c, n)
}