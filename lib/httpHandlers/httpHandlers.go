package httpHandlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Monska85/telegram-gateway/lib/logHelper"
	"github.com/Monska85/telegram-gateway/lib/telegram"
	"github.com/Monska85/telegram-gateway/lib/utils"
)

var logger = logHelper.GetInstance()
var tbot *telegram.TelegramBot

func SetTelegramBot(bot *telegram.TelegramBot) {
	tbot = bot
}

func isPost(r *http.Request) bool {
	return r.Method == "POST"
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "telegram-gateway")
	w.Header().Set("X-Application", "telegram-gateway")
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	var response Response
	response.Success = false

	if !isPost(r) {
		logger.Out(utils.LogError, "Invalid request method", "method", r.Method)
		response.Message = "Invalid request method"
		http.Error(w, utils.GetJsonString(response), http.StatusMethodNotAllowed)
		return
	}

	// Read the entire body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Out(utils.LogError, "Error reading request body", "error", err.Error())
		response.Message = "Error reading request body"
		http.Error(w, utils.GetJsonString(response), http.StatusInternalServerError)
		return
	}

	var request SendMessageRequest
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Out(utils.LogError, "Error parsing request body", "error", err.Error())
		response.Message = "Error parsing request body"
		http.Error(w, utils.GetJsonString(response), http.StatusBadRequest)
		return
	}

	logger.Out(utils.LogDebug, "Received request", "request", request)

	// Send image message if provided
	if request.Image != "" {
		if _, err := tbot.SendImageMessage(int64(request.ChatID), request.Text, request.Image, request.ImageName); err != nil {
			logger.Out(utils.LogError, "Error sending image message")
			response.Message = "Error sending image message"
			http.Error(w, utils.GetJsonString(response), http.StatusInternalServerError)
			return
		}
		logger.Out(utils.LogInfo, "Image message sent successfully")
		response.Success = true
		response.Message = "Image message sent successfully"
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(utils.GetJsonString(response)))
		return
	}

	// If no image is provided, send a regular text message
	if _, err := tbot.SendMessage(int64(request.ChatID), request.Text); err != nil {
		logger.Out(utils.LogError, "Error sending message")
		response.Message = "Error sending message"
		http.Error(w, utils.GetJsonString(response), http.StatusInternalServerError)
		return
	}

	logger.Out(utils.LogInfo, "Message sent successfully")
	response.Success = true
	response.Message = "Message sent successfully"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(utils.GetJsonString(response)))
}
