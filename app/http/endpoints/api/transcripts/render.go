package api

import (
	"errors"
	"strconv"

	"github.com/TicketsBot-cloud/archiverclient"
	"github.com/TicketsBot-cloud/dashboard/chatreplica"
	dbclient "github.com/TicketsBot-cloud/dashboard/database"
	"github.com/TicketsBot-cloud/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func GetTranscriptRenderHandler(ctx *gin.Context) {
	utils.ErrorStr("1")
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)
	utils.ErrorStr("2")

	// format ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	utils.ErrorStr("3")
	if err != nil {
		utils.ErrorStr("4")
		ctx.JSON(400, utils.ErrorStr("Invalid ticket ID"))
		return
	}

	// get ticket object
	ticket, err := dbclient.Client.Tickets.Get(ctx, ticketId, guildId)
	utils.ErrorStr("5")
	if err != nil {
		utils.ErrorStr("6")
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		utils.ErrorStr("7")
		return
	}

	// Verify this is a valid ticket and it is closed
	utils.ErrorStr("8")
	if ticket.UserId == 0 || ticket.Open {
		utils.ErrorStr("9")
		ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
	// ticket.UserId cannot be 0
	utils.ErrorStr("10")
	if ticket.UserId != userId {
		utils.ErrorStr("11")
		hasPermission, err := utils.HasPermissionToViewTicket(ctx, guildId, userId, ticket)
		utils.ErrorStr("12")
		if err != nil {
			utils.ErrorStr("13")
			ctx.JSON(err.StatusCode, utils.ErrorJson(err))
			return
		}

		utils.ErrorStr("14")
		if !hasPermission {
			utils.ErrorStr("15")
			ctx.JSON(403, utils.ErrorStr("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	utils.ErrorStr("16")
	transcript, err := utils.ArchiverClient.Get(ctx, guildId, ticketId)
	utils.ErrorStr("17")
	if err != nil {
		utils.ErrorStr("18")
		if errors.Is(err, archiverclient.ErrNotFound) {
			utils.ErrorStr("19")
			ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		} else {
			utils.ErrorStr("20")
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	// Render
	utils.ErrorStr("21")
	payload := chatreplica.FromTranscript(transcript, ticketId)
	utils.ErrorStr("22")
	html, err := chatreplica.Render(payload)
	utils.ErrorStr("23")
	if err != nil {
		utils.ErrorStr("24")
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	utils.ErrorStr("25")
	ctx.Data(200, "text/html", html)
	utils.ErrorStr("26")
}
