package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) handlerPing(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary	send to email one time password
// @Schemes
// @Description	send code to email
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			email	body	tRequestOTP	true	"User email"
// @Success		200
// @Failure		400
// @Failure		500
// @Router			/auth/otp [post]
func (s *Server) handlerAuthOTP(c *gin.Context) {
	unauthorize(c)

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tRequestOTP{}

	err := json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !jBody.Email.IsValid() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.sport.NewOTP(c.Request.Context(), jBody.Email.String())
	if err != nil {
		s.log.Error("failed send otp", zap.String("email", jBody.Email.String()), zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary	authorization
// @Schemes
// @Description	authorization
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			email	body	tAuthorization	true	"User email and password"
// @Success		200
// @Failure		400
// @Failure		500
// @Router			/auth/login [post]
func (s *Server) handlerLogin(c *gin.Context) {
	unauthorize(c)

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tAuthorization{}

	err := json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, message := s.login(c, jBody.Email, jBody.OTP)
	if message != "" {
		c.JSON(statusCode, gin.H{
			"message": message,
		})
		return
	}
	c.Writer.WriteHeader(statusCode)
}

// @Summary	logout
// @Schemes
// @Description	logout
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400
// @Failure		500
// @Router			/auth/logout [get]
func (s *Server) handlerLogout(c *gin.Context) {
	unauthorize(c)
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary	все турниры
// @Schemes
// @Description	все турниры
// @Tags			guest
// @Accept			json
// @Produce		json
// @Success		200 {object} tTournament
// @Failure		400
// @Failure		500
// @Router			/tournaments [get]
func (s *Server) handlerGetAllTournament(c *gin.Context) {
	tournaments, err := s.sport.GetAllTournaments(c.Request.Context())
	if err != nil {
		s.log.Error("failed get all tournaments", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := []tTournament{}
	for _, t := range *tournaments {
		res = append(res, tTournament{
			ID:    t.ID,
			Title: t.Title,
		})
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	user info
// @Schemes
// @Description	user info
// @Tags			user
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		401
// @Failure		500
// @Router			/user [get]
func (s *Server) handlerUser(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mail": user.Email,
	})
}

// @Summary	создать турнир
// @Schemes
// @Description	создать турнир
// @Tags			user tournament
// @Accept			json
// @Produce		json
// @Param			tournamet	body		tCreateTournament	true	"tournament"
// @Success		201			{object}	tTournament
// @Failure		400
// @Failure		500
// @Router			/user/tournaments [post]
func (s *Server) handlerUserNewTournament(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tCreateTournament{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.NewTournament(c.Request.Context(), &models.Tournament{
		UserID: user.ID,
		Title:  jBody.Title,
	})
	if err != nil {
		s.log.Error("filed create tournament", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tTournament{
		ID:    tournament.ID,
		Title: tournament.Title,
	})
}

// @Summary	турниры пользователя
// @Schemes
// @Description	создать турнир
// @Tags			user tournament
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]tTournament
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/tournaments [get]
func (s *Server) handlerUserTournaments(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	tournaments, err := s.sport.GetTournaments(c.Request.Context(), user)
	if err != nil {
		s.log.Error("failed get tournaments by user", zap.Uint("userID", user.ID), zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tournaments == nil {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	result := []tTournament{}
	for _, t := range *tournaments {
		result = append(result, tTournament{
			ID:    t.ID,
			Title: t.Title,
		})
	}

	c.JSON(http.StatusOK, result)
}

// @Summary	информация турнира пользователя
// @Schemes
// @Description	информация турнира пользователя
// @Tags			user tournament
// @Param			id	path	int	true	"tournament id"
// @Produce		json
// @Success		200	{object}	tTournament
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/tournaments/{id} [get]
func (s *Server) handlerUserTournament(c *gin.Context) {
	// _, statusCode, err := s.checkUser(c)
	// if err != nil {
	// 	c.Writer.WriteHeader(statusCode)
	// 	return
	// }

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tournament.ID == 0 {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, tTournament{
		ID:    tournament.ID,
		Title: tournament.Title,
	})
}

// @Summary	Обновить турнир
// @Schemes
// @Description	Обновить турнир
// @Tags			user tournament
// @Accept			json
// @Produce		json
// @Param			tournamet	body	tUpdateTournament	true	"tournament"
// @Success		200
// @Success		204
// @Failure		400
// @Failure		500
// @Router			/user/tournaments [put]
func (s *Server) handlerUserUpdTournament(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdateTournament{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.UpdTournament(c.Request.Context(), &models.Tournament{
		ID:     jBody.ID,
		Title:  jBody.Title,
		UserID: user.ID,
	})
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update tournament", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tTournament{
		ID:    tournament.ID,
		Title: tournament.Title,
	})
}

// @Summary	создать команду
// @Schemes
// @Description	создать команду
// @Tags			user team
// @Accept			json
// @Produce		json
// @Param			tournamet	body		tCreateTeam	true	"team"
// @Success		201			{object}	tTeam
// @Failure		400
// @Failure		500
// @Router			/user/teams [post]
func (s *Server) handlerUserNewTeam(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tCreateTeam{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.NewTeam(c.Request.Context(), &models.Team{
		Title:  jBody.Title,
		UserID: user.ID,
	})
	if err != nil {
		s.log.Error("failed create team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tTeam{
		ID:    team.ID,
		Title: team.Title,
	})
}

// @Summary	команды пользователя
// @Schemes
// @Description	команды пользователя
// @Tags			user team
// @Produce		json
// @Success		200	{object}	[]tTeam
// @Failure		400
// @Failure		500
// @Router			/user/teams [get]
func (s *Server) handlerUserTeams(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	teams, err := s.sport.GetTeams(c.Request.Context(), user)
	if err != nil {
		s.log.Error("failed get teams: %w", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := []tTeam{}
	if teams != nil {
		for _, t := range *teams {
			res = append(res, tTeam{
				ID:    t.ID,
				Title: t.Title,
			})
		}
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	информация команды пользователя
// @Schemes
// @Description	информация команды пользователя
// @Tags			user team
// @Param			id	path	int	true	"team id"
// @Produce		json
// @Success		200	{object}	tTeam
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/teams/{id} [get]
func (s *Server) handlerUserTeam(c *gin.Context) {
	// _, statusCode, err := s.checkUser(c)
	// if err != nil {
	// 	c.Writer.WriteHeader(statusCode)
	// 	return
	// }

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.GetTeamByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if team.ID == 0 {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, tTeam{
		ID:    team.ID,
		Title: team.Title,
	})
}

// @Summary	обновление команды пользователя
// @Schemes
// @Description	обновление команды пользователя
// @Tags			user team
// @Param			id	body	tUpdateTeam	true	"team"
// @Produce		json
// @Success		200	{object}	tTeam
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/teams [put]
func (s *Server) handlerUserUptTeam(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdateTeam{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.UpdTeam(c.Request.Context(), &models.Team{
		ID:     jBody.ID,
		Title:  jBody.Title,
		UserID: user.ID,
	})
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tTeam{
		ID:    team.ID,
		Title: jBody.Title,
	})
}

// @Summary	Добавить игрока
// @Schemes
// @Description	Добавить игрока
// @Tags			user players
// @Param			player	body	tNewPlayer	true	"player"
// @Produce		json
// @Success		201	{object}	tPlayer
// @Failure		400
// @Failure		500
// @Router			/user/players [post]
func (s *Server) handlerUserNewPlayer(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tNewPlayer{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	player, err := s.sport.NewPlayer(c.Request.Context(), &models.Player{
		FirstName:  jBody.FirstName,
		SecondName: jBody.SecondName,
		LastName:   jBody.LastName,
		UserID:     user.ID,
	})
	if err != nil {
		s.log.Error("failed create player", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tPlayer{
		ID:         player.ID,
		FirstName:  player.FirstName,
		SecondName: player.SecondName,
		LastName:   player.LastName,
	})
}

// @Summary	Все игроки
// @Schemes
// @Description	Все игроки
// @Tags			user players
// @Produce		json
// @Success		200	{object}	[]tPlayer
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/players [get]
func (s *Server) handlerUserPlayers(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	players, err := s.sport.GetPlayers(c.Request.Context(), userID)
	if err != nil {
		s.log.Error("failed get players", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res []tPlayer
	if players != nil {
		for _, p := range *players {
			res = append(res, tPlayer{
				ID:         p.ID,
				FirstName:  p.FirstName,
				SecondName: p.SecondName,
				LastName:   p.LastName,
			})
		}
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	обновить игрока
// @Schemes
// @Description	обновить игрока
// @Tags			user players
// @Param			id	body	tUpdatePlayer	true	"player"
// @Produce		json
// @Success		200	{object}	tPlayer
// @Failure		204
// @Failure		400
// @Failure		500
// @Router			/user/players [put]
func (s *Server) handlerUserUpdatePlayer(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdatePlayer{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	player, err := s.sport.UpdPlayer(c.Request.Context(), &models.Player{
		ID:         jBody.ID,
		FirstName:  jBody.FirstName,
		SecondName: jBody.SecondName,
		LastName:   jBody.LastName,
		UserID:     userID,
	})
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update player", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tPlayer{
		ID:         player.ID,
		FirstName:  player.FirstName,
		SecondName: player.SecondName,
		LastName:   player.LastName,
	})
}

// @Summary	добавить игрока в команду
// @Schemes
// @Description	пдобавить игрока в команду
// @Tags			user team
// @Param			id		path	int		true	"team id"
// @Param			player	body	[]int	true	"player id`s"
// @Produce		json
// @Success		200
// @Failure		204
// @Failure		400
// @Failure		409
// @Failure		500
// @Router			/user/teams/{id}/players [post]
func (s *Server) handlerUserTeamAddPlayer(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := []uint{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.sport.AddPlayersTeam(c.Request.Context(), &jBody, uint(teamID), userID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		if errors.Is(err, errstore.ErrConflictData) {
			c.Writer.WriteHeader(http.StatusConflict)
			return
		}
		s.log.Error("failed add player's to team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary	игроки в команде
// @Schemes
// @Description	игроки в команде
// @Tags			user team
// @Param			id	path	int	true	"team id"
// @Produce		json
// @Success		200	{object}	tPlayer
// @Failure		400
// @Failure		500
// @Router			/user/teams/{id}/players [get]
func (s *Server) handlerUserTeamPlayers(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	players, err := s.sport.GetPlayersTeam(c.Request.Context(), &models.Team{
		ID:     uint(teamID),
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed get player's by team", zap.Error(err), zap.Int("team_id", teamID))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := []tPlayer{}
	for _, player := range *players {
		res = append(res, tPlayer{
			ID:         player.ID,
			FirstName:  player.FirstName,
			SecondName: player.SecondName,
			LastName:   player.LastName,
		})
	}

	c.JSON(http.StatusOK, res)
}

// @Summary	удалить игроков из команду
// @Schemes
// @Description	удалить игроков из команду
// @Tags			user team
// @Param			id		path	int		true	"team id"
// @Param			player	body	[]int	true	"player id`s"
// @Produce		json
// @Success		200
// @Failure		204
// @Failure		400
// @Failure		409
// @Failure		500
// @Router			/user/teams/{id}/players [delete]
func (s *Server) handlerUserTeamRemovePlayers(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := []uint{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil || len(jBody) == 0 {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.sport.RemovePlayersTeam(c.Request.Context(), &jBody, uint(teamID), userID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		if errors.Is(err, errstore.ErrConflictData) {
			c.Writer.WriteHeader(http.StatusConflict)
			return
		}
		s.log.Error("failed remove player's from team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

//	@Summary	поданые заявки на турнир
//	@Schemes
//	@Description	поданые заявки на турнир
//	@Tags			user team
//	@Param			id	path	int	true	"tournament id"
//	@Produce		json
//	@Success		200
//	@Failure		204
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/user/tournaments/{id}/applications/ [get]
// func (s *Server) handlerGetTournamentApplications(c *gin.Context) {
// 	userID, err := s.checkAuth(c)
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	tournamentID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// }

//	@Summary	ответ на заявку
//	@Schemes
//	@Description	ответ на заявку
//	@Tags			user team
//	@Param			id		path	int		true	"team id"
//	@Param			player	body	[]int	true	"player id`s"
//	@Produce		json
//	@Success		200
//	@Failure		204
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/user/tournaments/{id}/applications/{aid} [put]
// func (s *Server) handlerUpdTournamentApplication(c *gin.Context) {
// 	userID, err := s.checkAuth(c)
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	teamID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	bBody, statusCode := s.readBody(c)
// 	if statusCode > 0 {
// 		c.Writer.WriteHeader(statusCode)
// 		return
// 	}

// 	jBody := []uint{}
// 	c.Writer.WriteHeader(http.StatusOK)
// }

//	@Summary	подать заявку
//	@Schemes
//	@Description	подать заявку
//	@Tags			user team application
//	@Param			id			path	int				true	"team id"
//	@Param			application	body	tNewApplication	true	"application"
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/user/team/{id}/applications/ [put]
// func (s *Server) handlerNewTeamApplication(c *gin.Context) {
// 	userID, err := s.checkAuth(c)
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	teamID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	bBody, statusCode := s.readBody(c)
// 	if statusCode > 0 {
// 		c.Writer.WriteHeader(statusCode)
// 		return
// 	}

// 	jBody := tNewApplication{}

// 	err = json.Unmarshal(bBody, &jBody)
// 	if err != nil {
// 		s.log.Debug("failed parse body", zap.Error(err))
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// }

//	@Summary	заявки команды
//	@Schemes
//	@Description	заявки команды
//	@Tags			user team
//	@Param			id	path	int	true	"team id"
//	@Produce		json
//	@Success		200
//	@Failure		204
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/user/team/{id}/applications/ [put]
// func (s *Server) handlerGetTeamApplications(c *gin.Context) {
// 	userID, err := s.checkAuth(c)
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	teamID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	bBody, statusCode := s.readBody(c)
// 	if statusCode > 0 {
// 		c.Writer.WriteHeader(statusCode)
// 		return
// 	}

// 	jBody := []uint{}
// }
