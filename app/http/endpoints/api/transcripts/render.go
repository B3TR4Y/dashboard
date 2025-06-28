package api

import (
	"errors"
	"strconv"
	"log"
	"github.com/TicketsBot-cloud/archiverclient"
	"github.com/TicketsBot-cloud/dashboard/chatreplica"
	dbclient "github.com/TicketsBot-cloud/dashboard/database"
	"github.com/TicketsBot-cloud/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func GetTranscriptRenderHandler(ctx *gin.Context) {
	log.Println("1")
	guildId := ctx.Keys["guildid"].(uint64)
	userId := ctx.Keys["userid"].(uint64)
	log.Println("2")

	// format ticket ID
	ticketId, err := strconv.Atoi(ctx.Param("ticketId"))
	log.Println("3")
	if err != nil {
		log.Println("4")
		ctx.JSON(400, utils.ErrorStr("Invalid ticket ID"))
		return
	}

	// get ticket object
	ticket, err := dbclient.Client.Tickets.Get(ctx, ticketId, guildId)
	log.Println("5")
	if err != nil {
		log.Println("6")
		ctx.JSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		log.Println("7")
		return
	}

	// Verify this is a valid ticket and it is closed
	log.Println("8")
	if ticket.UserId == 0 || ticket.Open {
		log.Println("9")
		ctx.JSON(404, utils.ErrorStr("Transcript not found"))
		return
	}

	// Verify the user has permissions to be here
	// ticket.UserId cannot be 0
	log.Println("10")
	if ticket.UserId != userId {
		log.Println("11")
		hasPermission, err := utils.HasPermissionToViewTicket(ctx, guildId, userId, ticket)
		log.Println("12")
		if err != nil {
			log.Println("13")
			ctx.JSON(err.StatusCode, utils.ErrorJson(err))
			return
		}

		log.Println("14")
		if !hasPermission {
			log.Println("15")
			ctx.JSON(403, utils.ErrorStr("You do not have permission to view this transcript"))
			return
		}
	}

	// retrieve ticket messages from bucket
	log.Println("16")
	transcript, err := utils.ArchiverClient.Get(ctx, guildId, ticketId)
	log.Println(transcript)
	log.Println(err)
	log.Println(ctx)
	log.Println(guildId)
	log.Println(ticketId)
	log.Println("17")
	if err != nil {
		log.Println("18")
		if err.Error() == "invalid input: magic number mismatch" {
			// Ignore this specific error and proceed
			log.Println("Ignoring error: invalid input: magic number mismatch")
		} else  {
			log.Println("Error: "+ err.Error())

			if errors.Is(err, archiverclient.ErrNotFound) {
				log.Println("19")
				ctx.JSON(404, utils.ErrorStr("Transcript not found"))
			} else {
				log.Println("20")
				ctx.JSON(500, utils.ErrorJson(err))
			}

			return
		}
	}

	// Render
	log.Println("21")
	payload := chatreplica.FromTranscript(transcript, ticketId)
	log.Println("22")
	html, err := chatreplica.Render(payload)
	log.Println("23")
	if err != nil {
		log.Println("24")
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	log.Println("25")
	ctx.Data(200, "text/html", html)
	log.Println("26")
}
