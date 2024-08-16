package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/config"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/domains/app"
	domainSend "github.com/trio-kwek-kwek/GoWhatsappWeb/domains/send"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/internal/rest/helpers"
	pkgError "github.com/trio-kwek-kwek/GoWhatsappWeb/pkg/error"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/pkg/utils"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/pkg/whatsapp"
	"github.com/trio-kwek-kwek/GoWhatsappWeb/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	"os/exec"
)

type serviceSend struct {
	WaCli      *whatsmeow.Client
	appService app.IAppService
}

func NewSendService(waCli *whatsmeow.Client, appService app.IAppService) domainSend.ISendService {
	return &serviceSend{
		WaCli:      waCli,
		appService: appService,
	}
}

func (service serviceSend) SendText(ctx context.Context, request domainSend.MessageRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendMessage(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	// Send message
	msg := &waE2E.Message{Conversation: proto.String(request.Message)}

	// Reply message
	if request.ReplyMessageID != nil && *request.ReplyMessageID != "" {
		participantJID := dataWaRecipient.String()
		if len(*request.ReplyMessageID) < 28 {
			firstDevice, err := service.appService.FirstDevice(ctx)
			if err != nil {
				return response, err
			}
			participantJID = firstDevice.Device
		}

		msg = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text: proto.String(request.Message),
				ContextInfo: &waE2E.ContextInfo{
					StanzaID:    request.ReplyMessageID,
					Participant: proto.String(participantJID),
					QuotedMessage: &waE2E.Message{
						Conversation: proto.String(request.Message),
					},
				},
			},
		}
	}

	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Message sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendImage(ctx context.Context, request domainSend.ImageRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendImage(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	var (
		imagePath      string
		imageThumbnail string
	)

	// Save image to server
	imagePath = fmt.Sprintf("%s/%s", config.PathSendItems, request.Image)

	// Send to WA server
	dataWaCaption := request.Caption
	dataWaImage, err := os.ReadFile(imagePath)
	if err != nil {
		return response, err
	}
	uploadedImage, err := service.WaCli.Upload(context.Background(), dataWaImage, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("failed to upload file: %v", err)
		return response, err
	}

	srcImage, err := imaging.Open(imagePath)
	if err != nil {
		return response, pkgError.InternalServerError(fmt.Sprintf("failed to open image %v", err))
	}

	resizedImage := imaging.Resize(srcImage, 100, 0, imaging.Lanczos)
	imageThumbnail = fmt.Sprintf("%s/thumbnails-%s", config.PathSendItems, request.Image)
	if err = imaging.Save(resizedImage, imageThumbnail); err != nil {
		return response, pkgError.InternalServerError(fmt.Sprintf("failed to save thumbnail %v", err))
	}

	msg := &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
		Caption:       proto.String(dataWaCaption),
		URL:           proto.String(uploadedImage.URL),
		DirectPath:    proto.String(uploadedImage.DirectPath),
		MediaKey:      uploadedImage.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(dataWaImage)),
		FileEncSHA256: uploadedImage.FileEncSHA256,
		FileSHA256:    uploadedImage.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(dataWaImage))),
		ViewOnce:      proto.Bool(request.ViewOnce),
	}}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Message sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendFile(ctx context.Context, request domainSend.FileRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendFile(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	fileBytes := helpers.MultipartFormFileHeaderToBytes(request.File)
	fileMimeType := http.DetectContentType(fileBytes)

	// Send to WA server
	uploadedFile, err := service.WaCli.Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)
	if err != nil {
		fmt.Printf("Failed to upload file: %v", err)
		return response, err
	}

	msg := &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
		URL:           proto.String(uploadedFile.URL),
		Mimetype:      proto.String(fileMimeType),
		Title:         proto.String(request.File.Filename),
		FileSHA256:    uploadedFile.FileSHA256,
		FileLength:    proto.Uint64(uploadedFile.FileLength),
		MediaKey:      uploadedFile.MediaKey,
		FileName:      proto.String(request.File.Filename),
		FileEncSHA256: uploadedFile.FileEncSHA256,
		DirectPath:    proto.String(uploadedFile.DirectPath),
		Caption:       proto.String(request.Caption),
	}}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Document sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendVideo(ctx context.Context, request domainSend.VideoRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendVideo(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	var (
		videoPath    string
		deletedItems []string
	)

	// Save video to server
	videoPath = fmt.Sprintf("%s/%s", config.PathSendItems, request.Video)

	// Check if ffmpeg is installed
	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		return response, pkgError.InternalServerError("ffmpeg not installed")
	}

	// Get thumbnail video with ffmpeg
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	thumbnailVideoPath := fmt.Sprintf("%s/%s", config.PathSendItems, request.Video+".png")
	cmdThumbnail := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01.000", "-vframes", "1", "-f", "image2pipe", "pipe:1")
	cmdThumbnail.Stdout = &outBuffer
	cmdThumbnail.Stderr = &errBuffer
	err = cmdThumbnail.Run()
	if err != nil {
		return response, pkgError.InternalServerError(fmt.Sprintf("failed to create thumbnail %v", errBuffer.String()))
	}

	dataWaVideo, err := os.ReadFile(videoPath)
	if err != nil {
		return response, err
	}
	uploaded, err := service.WaCli.Upload(context.Background(), dataWaVideo, whatsmeow.MediaVideo)
	if err != nil {
		return response, pkgError.InternalServerError(fmt.Sprintf("Failed to upload file: %v", err))
	}
	//remove thumbnail
	deletedItems = append(deletedItems, thumbnailVideoPath)

	msg := &waE2E.Message{VideoMessage: &waE2E.VideoMessage{
		URL:                 proto.String(uploaded.URL),
		Mimetype:            proto.String(http.DetectContentType(dataWaVideo)),
		Caption:             proto.String(request.Caption),
		FileLength:          proto.Uint64(uploaded.FileLength),
		FileSHA256:          uploaded.FileSHA256,
		FileEncSHA256:       uploaded.FileEncSHA256,
		MediaKey:            uploaded.MediaKey,
		DirectPath:          proto.String(uploaded.DirectPath),
		ViewOnce:            proto.Bool(request.ViewOnce),
		JPEGThumbnail:       outBuffer.Bytes(),
		ThumbnailEncSHA256:  outBuffer.Bytes(),
		ThumbnailSHA256:     outBuffer.Bytes(),
		ThumbnailDirectPath: proto.String(uploaded.DirectPath),
	}}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	go func() {
		errDelete := utils.RemoveFile(1, deletedItems...)
		if errDelete != nil {
			logrus.Infof("error when deleting picture: %v", errDelete)
		}
	}()
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Video sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendContact(ctx context.Context, request domainSend.ContactRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendContact(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	msgVCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:;%v;;;\nFN:%v\nTEL;type=CELL;waid=%v:+%v\nEND:VCARD",
		request.ContactName, request.ContactName, request.ContactPhone, request.ContactPhone)
	msg := &waE2E.Message{ContactMessage: &waE2E.ContactMessage{
		DisplayName: proto.String(request.ContactName),
		Vcard:       proto.String(msgVCard),
	}}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Contact sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendLink(ctx context.Context, request domainSend.LinkRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendLink(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	getMetaDataFromURL := utils.GetMetaDataFromURL(request.Link)

	msg := &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{
		Text:         proto.String(fmt.Sprintf("%s\n%s", request.Caption, request.Link)),
		Title:        proto.String(getMetaDataFromURL.Title),
		CanonicalURL: proto.String(request.Link),
		MatchedText:  proto.String(request.Link),
		Description:  proto.String(getMetaDataFromURL.Description),
	}}
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Link sent to %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendLocation(ctx context.Context, request domainSend.LocationRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendLocation(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	// Compose WhatsApp Proto
	msg := &waE2E.Message{
		LocationMessage: &waE2E.LocationMessage{
			DegreesLatitude:  proto.Float64(utils.StrToFloat64(request.Latitude)),
			DegreesLongitude: proto.Float64(utils.StrToFloat64(request.Longitude)),
		},
	}

	// Send WhatsApp Message Proto
	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Send location success %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendAudio(ctx context.Context, request domainSend.AudioRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendAudio(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	autioBytes := helpers.MultipartFormFileHeaderToBytes(request.Audio)
	audioMimeType := http.DetectContentType(autioBytes)

	audioUploaded, err := service.WaCli.Upload(ctx, autioBytes, whatsmeow.MediaAudio)
	if err != nil {
		err = pkgError.WaUploadMediaError(fmt.Sprintf("Failed to upload audio: %v", err))
		return response, err
	}

	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(audioUploaded.URL),
			DirectPath:    proto.String(audioUploaded.DirectPath),
			Mimetype:      proto.String(audioMimeType),
			FileLength:    proto.Uint64(audioUploaded.FileLength),
			FileSHA256:    audioUploaded.FileSHA256,
			FileEncSHA256: audioUploaded.FileEncSHA256,
			MediaKey:      audioUploaded.MediaKey,
		},
	}

	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Send audio success %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}

func (service serviceSend) SendPoll(ctx context.Context, request domainSend.PollRequest) (response domainSend.GenericResponse, err error) {
	err = validations.ValidateSendPoll(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.WaCli, request.Phone)
	if err != nil {
		return response, err
	}

	ts, err := service.WaCli.SendMessage(ctx, dataWaRecipient, service.WaCli.BuildPollCreation(request.Question, request.Options, request.MaxAnswer))
	if err != nil {
		return response, err
	}

	sender, fetchErr := service.appService.FirstDevice(ctx)
	if fetchErr != nil {
		response.Sender = "errorFetch"
	} else {
		response.Sender = sender.Device
	}

	response.TraceCode = request.TraceCode
	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Send poll success %s (server timestamp: %s)", request.Phone, ts.Timestamp.String())
	return response, nil
}
