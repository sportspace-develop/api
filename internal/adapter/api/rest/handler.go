package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sport-space/internal/adapter/errsport"
	"sport-space/internal/adapter/models"
	"sport-space/internal/adapter/storage/errstore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) handlerPing(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

//	@Summary	send to email one time password
//	@Schemes
//	@Description	send code to email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body	tRequestOTP	true	"User email"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/auth/otp [post]
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

//	@Summary	authorization
//	@Schemes
//	@Description	authorization
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			email	body	tAuthorization	true	"User email and password"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/auth/login [post]
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

//	@Summary	logout
//	@Schemes
//	@Description	logout
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/auth/logout [get]
func (s *Server) handlerLogout(c *gin.Context) {
	unauthorize(c)
	c.Writer.WriteHeader(http.StatusOK)
}

//	@Summary	все турниры
//	@Schemes
//	@Description	все турниры
//	@Tags			guest
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"page number"
//	@Param			limit	query		int	false	"limit size"
//	@Success		200		{object}	tGetTorunamentsResponse
//	@Failure		400
//	@Failure		500
//	@Router			/tournaments [get]
func (s *Server) handlerGetAllTournament(c *gin.Context) {
	tournaments, err := s.sport.GetAllTournaments(c.Request.Context())
	if err != nil {
		s.log.Error("failed get all tournaments", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	pg := s.getPagination(c, len(*tournaments))

	res := []tTournamentResponse{}
	for _, t := range (*tournaments)[pg.StartRow:pg.EndRow] {
		res = append(res, tTournamentResponse{
			ID:                t.ID,
			Title:             t.Title,
			StartDate:         formatDateTime(t.StartDate),
			EndDate:           formatDateTime(t.EndDate),
			RegisterStartDate: formatDateTime(t.RegisterStartDate),
			RegisterEndDate:   formatDateTime(t.RegisterEndDate),
			LogoURL:           t.LogoURL,
			LogoExternalURL:   s.getFullUploadURL(t.LogoURL),
		})
	}

	c.JSON(http.StatusOK, tGetTorunamentsResponse{
		Pagination: pg,
		Data:       res,
	})
}

//	@Summary	user info
//	@Schemes
//	@Description	user info
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401
//	@Failure		500
//	@Router			/user [get]
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

//	@Summary	создать турнир
//	@Schemes
//	@Description	создать турнир
//	@Tags			user tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournamet	body		tCreateTournamentRequest	true	"tournament"
//	@Success		201			{object}	tTournamentResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments [post]
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

	jBody := tCreateTournamentRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !jBody.IsValid() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	t := &models.Tournament{
		UserID:            user.ID,
		Title:             jBody.Title,
		StartDate:         jBody.StartDate.DateTime(),
		EndDate:           jBody.EndDate.DateTime(),
		RegisterStartDate: jBody.RegisterStartDate.DateTime(),
		RegisterEndDate:   jBody.RegisterEndDate.DateTime(),
		LogoURL:           jBody.LogoURL,
	}

	tournament, err := s.sport.NewTournament(c.Request.Context(), t)
	if err != nil {
		s.log.Error("filed create tournament", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tTournamentResponse{
		ID:                tournament.ID,
		Title:             tournament.Title,
		StartDate:         formatDateTime(tournament.StartDate),
		EndDate:           formatDateTime(tournament.EndDate),
		RegisterStartDate: formatDateTime(tournament.RegisterStartDate),
		RegisterEndDate:   formatDateTime(tournament.RegisterEndDate),
		LogoURL:           tournament.LogoURL,
		LogoExternalURL:   s.getFullUploadURL(tournament.LogoURL),
	})
}

//	@Summary	турниры пользователя
//	@Schemes
//	@Description	турниры пользователя
//	@Tags			user tournament
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"page number"
//	@Param			limit	query		int	false	"limit size"
//	@Success		200		{object}	tGetTorunamentsResponse
//	@Failure		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments [get]
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

	pg := s.getPagination(c, len(*tournaments))

	result := []tTournamentResponse{}
	for _, t := range (*tournaments)[pg.StartRow:pg.EndRow] {
		result = append(result, tTournamentResponse{
			ID:                t.ID,
			Title:             t.Title,
			StartDate:         formatDateTime(t.StartDate),
			EndDate:           formatDateTime(t.EndDate),
			RegisterStartDate: formatDateTime(t.RegisterStartDate),
			RegisterEndDate:   formatDateTime(t.RegisterEndDate),
			LogoURL:           t.LogoURL,
			LogoExternalURL:   s.getFullUploadURL(t.LogoURL),
		})
	}

	c.JSON(http.StatusOK, tGetTorunamentsResponse{
		Pagination: pg,
		Data:       result,
	})
}

//	@Summary	информация турнира пользователя
//	@Schemes
//	@Description	информация турнира пользователя
//	@Tags			user tournament
//	@Param			tournament_id	path	int	true	"tournament id"
//	@Produce		json
//	@Success		200	{object}	tTournamentResponse
//	@Failure		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments/{tournament_id} [get]
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

	c.JSON(http.StatusOK, tTournamentResponse{
		ID:                tournament.ID,
		Title:             tournament.Title,
		StartDate:         formatDateTime(tournament.StartDate),
		EndDate:           formatDateTime(tournament.EndDate),
		RegisterStartDate: formatDateTime(tournament.RegisterStartDate),
		RegisterEndDate:   formatDateTime(tournament.RegisterEndDate),
		LogoURL:           tournament.LogoURL,
		LogoExternalURL:   s.getFullUploadURL(tournament.LogoURL),
	})
}

//	@Summary	Обновить турнир
//	@Schemes
//	@Description	Обновить турнир
//	@Tags			user tournament
//	@Accept			json
//	@Produce		json
//	@Param			tournament_id	path	int						true	"tournament id"
//	@Param			tournamet		body	tUpdTournamentRequest	true	"tournament"
//	@Success		200
//	@Success		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments/{tournament_id} [put]
func (s *Server) handlerUserUpdTournament(c *gin.Context) {
	user, statusCode, err := s.checkUser(c)
	if err != nil {
		c.Writer.WriteHeader(statusCode)
		return
	}

	tournamentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdTournamentRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !jBody.IsValid() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.UpdTournament(c.Request.Context(), &models.Tournament{
		ID:                uint(tournamentID),
		Title:             jBody.Title,
		StartDate:         jBody.StartDate.DateTime(),
		EndDate:           jBody.EndDate.DateTime(),
		RegisterStartDate: jBody.RegisterStartDate.DateTime(),
		RegisterEndDate:   jBody.RegisterEndDate.DateTime(),
		LogoURL:           jBody.LogoURL,
		UserID:            user.ID,
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

	c.JSON(http.StatusOK, tTournamentResponse{
		ID:                tournament.ID,
		Title:             tournament.Title,
		StartDate:         formatDateTime(tournament.StartDate),
		EndDate:           formatDateTime(tournament.EndDate),
		RegisterStartDate: formatDateTime(tournament.RegisterStartDate),
		RegisterEndDate:   formatDateTime(tournament.RegisterEndDate),
		LogoURL:           tournament.LogoURL,
		LogoExternalURL:   s.getFullUploadURL(tournament.LogoURL),
	})
}

//	@Summary	создать команду
//	@Schemes
//	@Description	создать команду
//	@Tags			user team
//	@Accept			json
//	@Produce		json
//	@Param			tournamet	body		tCreateTeam	true	"team"
//	@Success		201			{object}	tTeam
//	@Failure		400
//	@Failure		500
//	@Router			/user/teams [post]
func (s *Server) handlerUserNewTeam(c *gin.Context) {
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

	jBody := tCreateTeam{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.NewTeam(c.Request.Context(), &models.Team{
		Title:    jBody.Title,
		UserID:   userID,
		PhotoURL: jBody.PhotoURL,
		LogoURL:  jBody.LogotURL,
	})
	if err != nil {
		s.log.Error("failed create team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tTeam{
		ID:               team.ID,
		Title:            team.Title,
		LogoURL:          team.LogoURL,
		LogoExternalURL:  s.getFullUploadURL(team.LogoURL),
		PhotoURL:         team.PhotoURL,
		PhotoExternalURL: s.getFullUploadURL(team.PhotoURL),
		CreatedAt:        formatDateTime(&team.CreatedAt),
	})
}

//	@Summary	команды пользователя
//	@Schemes
//	@Description	команды пользователя
//	@Tags			user team
//	@Produce		json
//	@Param			page	query		int	false	"page number"
//	@Param			limit	query		int	false	"limit size"
//	@Success		200		{object}	tGetTeamsResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/teams [get]
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

	pg := s.getPagination(c, len(*teams))

	res := []tTeam{}
	if teams != nil && pg.TotalRecords > 0 {
		for _, t := range (*teams)[pg.StartRow:pg.EndRow] {
			res = append(res, tTeam{
				ID:               t.ID,
				Title:            t.Title,
				LogoURL:          t.LogoURL,
				LogoExternalURL:  s.getFullUploadURL(t.LogoURL),
				PhotoURL:         t.PhotoURL,
				PhotoExternalURL: s.getFullUploadURL(t.PhotoURL),
				CreatedAt:        formatDateTime(&t.CreatedAt),
			})
		}
	}

	c.JSON(http.StatusOK, tGetTeamsResponse{
		Pagination: pg,
		Data:       res,
	})
}

//	@Summary	информация команды пользователя
//	@Schemes
//	@Description	информация команды пользователя
//	@Tags			user team
//	@Param			team_id	path	int	true	"team id"
//	@Produce		json
//	@Success		200	{object}	tGetTeamResponse
//	@Failure		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/teams/{team_id} [get]
func (s *Server) handlerUserTeam(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.GetTeamByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed get players from team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if team.UserID != userID {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	resPlayers := []tPlayerResponse{}
	for _, p := range team.Players {
		resPlayers = append(resPlayers, tPlayerResponse{
			ID:               p.ID,
			FirstName:        p.FirstName,
			SecondName:       p.SecondName,
			LastName:         p.LastName,
			PhotoURL:         p.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(p.PhotoURL),
		})
	}

	c.JSON(http.StatusOK, tGetTeamResponse{
		ID:               team.ID,
		Title:            team.Title,
		Players:          resPlayers,
		LogoURL:          team.LogoURL,
		LogoExternalURL:  s.getFullUploadURL(team.LogoURL),
		PhotoURL:         team.PhotoURL,
		PhotoExternalURL: s.getFullUploadURL(team.PhotoURL),
		CreatedAt:        formatDateTime(&team.CreatedAt),
	})
}

//	@Summary	обновление команды пользователя
//	@Schemes
//	@Description	обновление команды пользователя
//	@Tags			user team
//	@Param			team_id	path	int				true	"team id"
//	@Param			team	body	tUpdTeamRequest	true	"team"
//	@Produce		json
//	@Success		200	{object}	tUpdTeamResponse
//	@Failure		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/teams/{team_id} [put]
func (s *Server) handlerUserUptTeam(c *gin.Context) {
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

	jBody := tUpdTeamRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.GetTeamByID(c.Request.Context(), uint(teamID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if team.UserID != userID {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	team.Title = jBody.Title
	team.LogoURL = jBody.LogotURL
	team.PhotoURL = jBody.PhotoURL

	team, players, err := s.sport.UpdTeam(c.Request.Context(), team, jBody.Players)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update team", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	playersRes := []tPlayerResponse{}
	for _, p := range *players {
		playersRes = append(playersRes, tPlayerResponse{
			ID:               p.ID,
			FirstName:        p.FirstName,
			SecondName:       p.SecondName,
			LastName:         p.LastName,
			BDay:             formatDate(p.BDay),
			PhotoURL:         p.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(p.PhotoURL),
		})
	}

	c.JSON(http.StatusOK, tUpdTeamResponse{
		ID:               team.ID,
		Title:            team.Title,
		Players:          &playersRes,
		LogoURL:          team.LogoURL,
		LogoExternalURL:  s.getFullUploadURL(team.LogoURL),
		PhotoURL:         team.PhotoURL,
		PhotoExternalURL: s.getFullUploadURL(team.PhotoURL),
		CreatedAt:        formatDateTime(&team.CreatedAt),
	})
}

//	@Summary	Добавить игрока
//	@Schemes
//	@Description	Добавить игрока
//	@Tags			user players
//	@Param			player	body	tNewPlayerRequest	true	"player"
//	@Produce		json
//	@Success		201	{object}	tPlayerResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/players [post]
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

	jBody := tNewPlayerRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !jBody.IsValid() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	player, err := s.sport.NewPlayer(c.Request.Context(), &models.Player{
		FirstName:  jBody.FirstName,
		SecondName: jBody.SecondName,
		LastName:   jBody.LastName,
		PhotoURL:   jBody.PhotoURL,
		UserID:     user.ID,
		BDay:       jBody.BDay.Date(),
	})
	if err != nil {
		s.log.Error("failed create player", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tPlayerResponse{
		ID:               player.ID,
		FirstName:        player.FirstName,
		SecondName:       player.SecondName,
		LastName:         player.LastName,
		PhotoURL:         player.PhotoURL,
		PhotoExternalURL: s.getFullUploadURL(player.PhotoURL),
		BDay:             formatDate(player.BDay),
	})
}

//	@Summary	Добавить/Обновить игроков
//	@Schemes
//	@Description	Добавить/Обновить игроков
//	@Tags			user players
//	@Param			players	body	[]tNewPlayerBatchRequest	true	"players"
//	@Produce		json
//	@Success		201	{object}	tNewPlayerBatchResponse
//	@Failure		400
//	@Failure		409
//	@Failure		500
//	@Router			/user/players/batch [post]
func (s *Server) handlerUserNewPlayerBatch(c *gin.Context) {
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

	jBody := []tNewPlayerBatchRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	batch := []models.Player{}
	for _, p := range jBody {
		if !p.IsValid() {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		batch = append(batch, models.Player{
			ID:         p.ID,
			FirstName:  p.FirstName,
			SecondName: p.SecondName,
			LastName:   p.LastName,
			BDay:       p.BDay.DateTime(),
			UserID:     userID,
		})
	}

	players, err := s.sport.NewPlayerBatch(c.Request.Context(), &batch)
	if err != nil {
		if errors.Is(err, errsport.ErrConflictData) {
			c.Writer.WriteHeader(http.StatusConflict)
			return
		}
		if errors.Is(err, errsport.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		s.log.Error("failed create players", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := tNewPlayerBatchResponse{}
	for _, p := range *players {
		res.Data = append(res.Data, tPlayerBatchResponse{
			ID:         p.ID,
			FirstName:  p.FirstName,
			SecondName: p.SecondName,
			LastName:   p.LastName,
			BDay:       formatDate(p.BDay),
		})
	}

	c.JSON(http.StatusCreated, res)
}

//	@Summary	Все игроки
//	@Schemes
//	@Description	Все игроки
//	@Tags			user players
//	@Produce		json
//	@Param			page	query		int	false	"page number"
//	@Param			limit	query		int	false	"limit size"
//	@Success		200		{object}	tGetPlayersResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/players [get]
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
	count := 0
	if players != nil {
		count = len(*players)
	}
	pg := s.getPagination(c, count)

	res := []tPlayerResponse{}
	if players != nil && pg.TotalRecords > 0 {
		for _, p := range (*players)[pg.StartRow:pg.EndRow] {
			res = append(res, tPlayerResponse{
				ID:               p.ID,
				FirstName:        p.FirstName,
				SecondName:       p.SecondName,
				LastName:         p.LastName,
				PhotoURL:         p.PhotoURL,
				PhotoExternalURL: s.getFullUploadURL(p.PhotoURL),
				BDay:             formatDateTime(p.BDay),
			})
		}
	}

	c.JSON(http.StatusOK, tGetPlayersResponse{
		Pagination: pg,
		Data:       res,
	})
}

//	@Summary	обновить игрока
//	@Schemes
//	@Description	обновить игрока
//	@Tags			user players
//	@Param			player_id	path	int						true	"player id"
//	@Param			id			body	tUpdatePlayerRequest	true	"player"
//	@Produce		json
//	@Success		200	{object}	tPlayerResponse
//	@Failure		204
//	@Failure		400
//	@Failure		500
//	@Router			/user/players/{player_id} [put]
func (s *Server) handlerUserUpdatePlayer(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	playerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdatePlayerRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if !jBody.IsValid() {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	player, err := s.sport.UpdPlayer(c.Request.Context(), &models.Player{
		ID:         uint(playerID),
		FirstName:  jBody.FirstName,
		SecondName: jBody.SecondName,
		LastName:   jBody.LastName,
		PhotoURL:   jBody.PhotoURL,
		UserID:     userID,
		BDay:       jBody.BDay.Date(),
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

	c.JSON(http.StatusOK, tPlayerResponse{
		ID:               player.ID,
		FirstName:        player.FirstName,
		SecondName:       player.SecondName,
		LastName:         player.LastName,
		PhotoURL:         player.PhotoURL,
		PhotoExternalURL: s.getFullUploadURL(player.PhotoURL),
		BDay:             formatDate(player.BDay),
	})
}

//	@Summary	заявки на турнир
//	@Schemes
//	@Description	заявки на турнир
//	@Tags			user tournament
//	@Param			tournament_id	path	int	true	"tournament id"
//	@Produce		json
//	@Success		200	{object}	tGetTournamentApplicationsResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments/{tournament_id}/applications [get]
func (s *Server) handlerGetTournamentApplications(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tournamentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), uint(tournamentID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tournament.UserID != userID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := []tTournamentApplication{}
	for _, a := range tournament.Applications {
		team, err := s.sport.GetTeamByID(c.Request.Context(), a.TeamID)
		if err != nil {
			if errors.Is(err, errstore.ErrNotFoundData) {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if a.Status != models.Draft && a.Status != models.Canceled {
			data = append(data, tTournamentApplication{
				ID:        a.ID,
				TeamID:    a.TeamID,
				TeamTitle: team.Title,
				Status:    string(a.Status),
			})
		}
	}

	c.JSON(http.StatusOK, tGetTournamentApplicationsResponse{Data: data})
}

//	@Summary	заявка турнира
//	@Schemes
//	@Description	заявка турнира
//	@Tags			user tournament
//	@Param			tournament_id	path	int	true	"tournament id"
//	@Param			application_id	path	int	true	"application id"
//	@Produce		json
//	@Success		200	{object}	tGetTorunamentApplicationResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments/{tournament_id}/applications/{application_id} [get]
func (s *Server) handlerGetTournamentApplication(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tournamentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	applicationID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), uint(tournamentID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tournament.UserID != userID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.sport.GetApplicationByID(c.Request.Context(), uint(applicationID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if application.TournamentID != tournament.ID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := s.sport.GetTeamByID(c.Request.Context(), application.TeamID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	players := []tPlayerResponse{}
	for _, p := range application.Players {
		players = append(players, tPlayerResponse{
			ID:               p.ID,
			FirstName:        p.FirstName,
			SecondName:       p.SecondName,
			LastName:         p.LastName,
			PhotoURL:         p.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(p.PhotoURL),
		})
	}

	c.JSON(http.StatusOK, tGetTorunamentApplicationResponse{
		ID:        application.ID,
		TeamID:    application.TeamID,
		TeamTitle: team.Title,
		Status:    string(application.Status),
		Players:   players,
	})
}

//	@Summary	изменить заявку
//	@Schemes
//	@Description	изменить заявку
//	@Tags			user tournament
//	@Param			tournament_id	path	int									true	"tournament id"
//	@Param			application_id	path	int									true	"application id"
//	@param			application		body	tUpdTournamentApplicationRequest	true	"application"
//	@Produce		json
//	@Success		200	{object}	tApplication
//	@Failure		400
//	@Failure		500
//	@Router			/user/tournaments/{tournament_id}/applications/{application_id} [put]
func (s *Server) handlerUpdTournamentApplication(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tournamentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	applicationID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdTournamentApplicationRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var status models.ApplicationStatus
	var ok bool
	if status, ok = applicationTournamentMapStatus[jBody.Status]; !ok {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.sport.UpdApplicationTournament(c.Request.Context(), uint(applicationID), status, uint(tournamentID), userID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, errstore.ErrForbidden) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	team, err := s.sport.GetTeamByID(c.Request.Context(), application.TeamID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tUpdTournamentApplicationResponse{
		ID:        application.ID,
		TeamID:    application.TeamID,
		TeamTitle: team.Title,
		Status:    string(application.Status),
	})
}

//	@Summary	подать заявку
//	@Schemes
//	@Description	подать заявку
//	@Tags			user team
//	@Param			team_id		path	int						true	"team id"
//	@Param			application	body	tNewApplicationRequest	true	"application"
//	@Produce		json
//	@Success		201	{object}	tNewApplicationResponse	"заявка создана"
//	@Failure		400	"не корректный запрос"
//	@Failure		409	"заявка	уже	была создана ранее"
//	@Failure		500
//	@Router			/user/teams/{team_id}/applications [post]
func (s *Server) handlerNewTeamApplication(c *gin.Context) {
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

	jBody := tNewApplicationRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	application, players, err := s.sport.NewApplicationTeam(c.Request.Context(), &jBody.PlayerIDs, jBody.TournamentID, uint(teamID), userID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		if errors.Is(err, errstore.ErrConflictData) {
			c.Writer.WriteHeader(http.StatusConflict)
			return
		}
		s.log.Error("failed create application", zap.Int("team_id", teamID), zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resPlayers := []tPlayerResponse{}
	for _, player := range *players {
		resPlayers = append(resPlayers, tPlayerResponse{
			ID:               player.ID,
			FirstName:        player.FirstName,
			SecondName:       player.SecondName,
			LastName:         player.LastName,
			PhotoURL:         player.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(player.PhotoURL),
		})
	}

	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), application.TournamentID)
	if err != nil {
		s.log.Error("failed get tournament", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tNewApplicationResponse{
		ID:              application.ID,
		TournamentID:    application.TournamentID,
		TournamentTitle: tournament.Title,
		Players:         resPlayers,
		Status:          string(application.Status),
	})
}

//	@Summary	изменить заявку
//	@Schemes
//	@Description	изменить заявку
//	@Tags			user team
//	@Param			team_id			path	int								true	"team id"
//	@Param			application_id	path	int								true	"application	id"
//	@Param			application		body	tUpdApplicationStatusRequest	true	"application status"
//	@Produce		json
//	@Success		200	{object}	tUpdApplicationResponse
//	@Failure		204	"заявка не найдена"
//	@Failure		400	"не найден или не может изменить"
//	@Failure		500
//	@Router			/user/teams/{team_id}/applications/{application_id} [put]
func (s *Server) handlerUpdStatusTeamApplication(c *gin.Context) {
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

	applicationID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bBody, statusCode := s.readBody(c)
	if statusCode > 0 {
		c.Writer.WriteHeader(statusCode)
		return
	}

	jBody := tUpdApplicationStatusRequest{}

	err = json.Unmarshal(bBody, &jBody)
	if err != nil {
		s.log.Debug("failed parse body", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if jBody.Status == nil && jBody.Players == nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var status models.ApplicationStatus
	if jBody.Status != nil {
		var ok bool
		status, ok = applicationMapStatus[*jBody.Status]
		if !ok {
			s.log.Debug("failed parse status", zap.Error(err))
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	application, players, err := s.sport.UpdApplicationTeam(c.Request.Context(), uint(applicationID), jBody.Players, status, uint(teamID), userID)
	if err != nil {
		if errors.Is(err, errstore.ErrForbidden) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		s.log.Error("failed update application", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resPlayers := []tPlayerResponse{}
	for _, player := range *players {
		resPlayers = append(resPlayers, tPlayerResponse{
			ID:               player.ID,
			FirstName:        player.FirstName,
			SecondName:       player.SecondName,
			LastName:         player.LastName,
			PhotoURL:         player.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(player.PhotoURL),
		})
	}
	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), application.TournamentID)
	if err != nil {
		s.log.Error("failed get tournament", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tUpdApplicationResponse{
		ID:              application.ID,
		TournamentID:    application.TournamentID,
		TournamentTitle: tournament.Title,
		Players:         resPlayers,
		Status:          string(application.Status),
	})
}

//	@Summary	заявки команды
//	@Schemes
//	@Description	заявки команды
//	@Tags			user team
//	@Param			team_id	path	int	true	"team id"
//	@Produce		json
//	@Success		200	{object}	tGetApplicationsTeamResponse
//	@Failure		400	"команда не найдена"
//	@Failure		500
//	@Router			/user/teams/{team_id}/applications [get]
func (s *Server) handlerGetTeamApplications(c *gin.Context) {
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

	team, err := s.sport.GetTeamByID(c.Request.Context(), uint(teamID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if team.UserID != userID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	applications, err := s.sport.GetApplicationsTeam(c.Request.Context(), team.ID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.JSON(http.StatusOK, tGetApplicationsTeamResponse{Data: []tApplication{}})
			return
		}
		s.log.Error("failed get applications by team", zap.Uint("team_id", team.ID), zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := []tApplication{}
	for _, a := range *applications {
		data = append(data, tApplication{
			ID:           a.ID,
			TournamentID: a.TournamentID,
			Status:       string(a.Status),
		})
	}

	c.JSON(http.StatusOK, tGetApplicationsTeamResponse{Data: data})
}

//	@Summary	заявка команды
//	@Schemes
//	@Description	заявка команды
//	@Tags			user team
//	@Param			team_id			path	int	true	"team id"
//	@Param			application_id	path	int	true	"application id"
//	@Produce		json
//	@Success		200	{object}	tGetApplicationResponse
//	@Failure		400	"не корректный запрос"
//	@Failure		500
//	@Router			/user/teams/{team_id}/applications/{application_id} [get]
func (s *Server) handlerGetApplication(c *gin.Context) {
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

	team, err := s.sport.GetTeamByID(c.Request.Context(), uint(teamID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if team.UserID != userID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	applicationID, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	application, err := s.sport.GetApplicationByID(c.Request.Context(), uint(applicationID))
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if application.TeamID != team.ID {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resPlayers := []tPlayerResponse{}
	for _, player := range application.Players {
		resPlayers = append(resPlayers, tPlayerResponse{
			ID:               player.ID,
			FirstName:        player.FirstName,
			SecondName:       player.SecondName,
			LastName:         player.LastName,
			PhotoURL:         player.PhotoURL,
			PhotoExternalURL: s.getFullUploadURL(player.PhotoURL),
		})
	}

	tournament, err := s.sport.GetTournamentByID(c.Request.Context(), application.TournamentID)
	if err != nil {
		if errors.Is(err, errstore.ErrNotFoundData) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tGetApplicationResponse{
		ID:              application.ID,
		TournamentID:    application.TournamentID,
		TournamentTitle: tournament.Title,
		Status:          string(application.Status),
		Players:         resPlayers,
	})
}

//	@Summary	загрузка файла
//	@Schemes
//	@Description	загрузка файла
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			file	formData	file	true	"файл"
//	@Success		201		{object}	tHandlerUploadResponse
//	@Failure		400
//	@Failure		500
//	@Router			/user/upload [post]
func (s *Server) handlerUpload(c *gin.Context) {
	userID, err := s.checkAuth(c)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	logoFile, err := c.FormFile("file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		if errors.Is(err, http.ErrMissingBoundary) {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		s.log.Error("failed get file", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.isValidImgExtension(logoFile) {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	dst := fmt.Sprintf("%v/%s", userID, s.genUploadName(logoFile.Filename))
	err = s.saveFile(logoFile, dst)
	if err != nil {
		s.log.Error("failed upload file", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, tHandlerUploadResponse{
		URL:      s.getFullUploadURL(dst),
		Filename: dst,
	})
}
