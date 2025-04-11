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
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)

	// Debug: Log guildId and userId
	utils.ErrorStr("Guild ID:", guildId, "User ID:", userId)

	// format ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	if (err != nil) {
		utils.ErrorStr("Error parsing ticket ID:", err)
		ctx.JSON(400, utils.ErrorStr("Invalid ticket ID"))
		return
	}
	utils.ErrorStr("Ticket ID:", ticketId)

	// get ticket object
	ticket, err := dbclient.Client.Tickets.Get(ctx, ticketId, guildId)
	if err != nil {
		utils.ErrorStr("Error fetching ticket:", err)
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	utils.ErrorStr("Fetched ticket:", ticket)

	// Verify this is a valid ticket and it is closed
	if ticket.UserId == 0 || ticket.Open {
		utils.ErrorStr("Invalid or open ticket:", ticket)
		ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
// ticket.UserId cannot be 0
	if ticket.UserId != userId {
		utils.ErrorStr("User does not own the ticket. Checking permissions...")
		hasPermission, err := utils.HasPermissionToViewTicket(ctx, guildId, userId, ticket)
		if err != nil {
			utils.ErrorStr("Error checking permissions:", err)
			ctx.JSON(err.StatusCode, utils.ErrorJson(err))
			return
		}

		if !hasPermission {
			utils.ErrorStr("User lacks permission to view the transcript")
			ctx.JSON(403, utils.ErrorStr("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	utils.ErrorStr("Fetching transcript from ArchiverClient...")
	transcript, err := utils.ArchiverClient.Get(ctx, guildId, ticketId)
	if err != nil {
		if errors.Is(err, archiverclient.ErrNotFound) {
			utils.ErrorStr("Transcript not found in ArchiverClient")
			ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		} else {
			utils.ErrorStr("Error fetching transcript:", err)
			ctx.JSON(500, utils.ErrorJson(err))
		}
		return
	}
	utils.ErrorStr("Fetched transcript:", transcript)

	// Render
	utils.ErrorStr("Rendering transcript...")
	payload := chatreplica.FromTranscript(transcript, ticketId)
	html, err := chatreplica.Render(payload)
	if err != nil {
		utils.ErrorStr("Error rendering transcript:", err)
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}
	utils.ErrorStr("Transcript rendered successfully")

	ctx.Data(200, "text/html", html)
}
