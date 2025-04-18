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
	fmt.Printf("1")
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)
	fmt.Printf("2")

	// format ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	fmt.Printf("3")
	if err != nil {
		fmt.Printf("4")
		ctx.JSON(400, fmt.Printf("Invalid ticket ID"))
		return
	}

	// get ticket object
	ticket, err := dbclient.Client.Tickets.Get(ctx, ticketId, guildId)
	fmt.Printf("5")
	if err != nil {
		fmt.Printf("6")
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		fmt.Printf("7")
		return
	}

	// Verify this is a valid ticket and it is closed
	fmt.Printf("8")
	if ticket.UserId == 0 || ticket.Open {
		fmt.Printf("9")
		ctx.JSON(404, fmt.Printf("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
	// ticket.UserId cannot be 0
	fmt.Printf("10")
	if ticket.UserId != userId {
		fmt.Printf("11")
		hasPermission, err := utils.HasPermissionToViewTicket(ctx, guildId, userId, ticket)
		fmt.Printf("12")
		if err != nil {
			fmt.Printf("13")
			ctx.JSON(err.StatusCode, utils.ErrorJson(err))
			return
		}

		fmt.Printf("14")
		if !hasPermission {
			fmt.Printf("15")
			ctx.JSON(403, fmt.Printf("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	fmt.Printf("16")
	transcript, err := utils.ArchiverClient.Get(ctx, guildId, ticketId)
	fmt.Printf("17")
	if err != nil {
		fmt.Printf("18")
		if errors.Is(err, archiverclient.ErrNotFound) {
			fmt.Printf("19")
			ctx.JSON(404, fmt.Printf("Transcript not found"))
		} else {
			fmt.Printf("20")
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	// Render
	fmt.Printf("21")
	payload := chatreplica.FromTranscript(transcript, ticketId)
	fmt.Printf("22")
	html, err := chatreplica.Render(payload)
	fmt.Printf("23")
	if err != nil {
		fmt.Printf("24")
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	fmt.Printf("25")
	ctx.Data(200, "text/html", html)
	fmt.Printf("26")
}
